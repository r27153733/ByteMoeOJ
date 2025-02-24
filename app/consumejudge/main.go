package main

// SIGUSR1 toggle the pause/resume consumption
import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"
	"log"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/IBM/sarama"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/r27153733/ByteMoeOJ/app/consumejudge/grpccompressor"
	"github.com/r27153733/ByteMoeOJ/app/consumejudge/logsarama"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/helper/compile"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/helper/mq"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/model"
	"github.com/r27153733/ByteMoeOJ/app/moejudge/pb"
	"github.com/r27153733/ByteMoeOJ/app/wasmexecclient/pb/wasm"
	"github.com/r27153733/ByteMoeOJ/app/wasmexecclient/wasmexecutor"
	"github.com/r27153733/ByteMoeOJ/lib/cgroup"
	"github.com/r27153733/ByteMoeOJ/lib/stringu"
	"github.com/r27153733/fastgozero/core/conf"
	"github.com/r27153733/fastgozero/core/service"
	"github.com/r27153733/fastgozero/core/stores/sqlx"
	"github.com/r27153733/fastgozero/zrpc"
	"github.com/valyala/bytebufferpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"
)

type KafkaConfig struct {
	Addr string

	Brokers  []string
	Version  string
	Group    string
	Topics   []string
	Assignor string `json:",default=range"`
	Oldest   bool   `json:",default=true"`
	Verbose  bool   `json:",default=true"`
}

type Config struct {
	service.ServiceConf
	WasmExecRPC    zrpc.RpcClientConf
	Kafka          KafkaConfig
	DataSourceName string
}

var c Config

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func init() {
	flag.Parse()
	conf.MustLoad(*configFile, &c)
	c.MustSetUp()

	db, err := sql.Open("pgx", c.DataSourceName)
	if err != nil {
		panic(err)
	}

	sqlxDB := sqlx.NewSqlConnFromDB(db)
	judgeModel = model.NewJudgeModel(sqlxDB)
	judgeCaseModel = model.NewJudgeCaseModel(sqlxDB)
	problemLangModel = model.NewProblemLangModel(sqlxDB)
	problemDataModel = model.NewProblemDataModel(sqlxDB)
	producer = mq.NewSyncProducer(&mq.KafkaConfig{
		Addr:          c.Kafka.Addr,
		Brokers:       c.Kafka.Brokers,
		Version:       c.Kafka.Version,
		Verbose:       false,
		CertFile:      "",
		KeyFile:       "",
		CaFile:        "",
		TLSSkipVerify: false,
	})

	encoding.RegisterCompressor(grpccompressor.ZstdCompressor{})
	wasmExecutor = wasmexecutor.NewWasmExecutor(
		zrpc.MustNewClient(c.WasmExecRPC,
			zrpc.WithDialOption(
				grpc.WithDefaultCallOptions(
					grpc.UseCompressor(grpccompressor.ZstdName),
					// 减一防止溢出
					grpc.MaxCallRecvMsgSize(math.MaxInt-1)),
			),
		),
	)
	zrpc.DontLogClientContentForMethod(wasm.WasmExecutor_Execute_FullMethodName)
}

var (
	judgeModel       model.JudgeModel
	judgeCaseModel   model.JudgeCaseModel
	problemLangModel model.ProblemLangModel
	problemDataModel model.ProblemDataModel
	wasmExecutor     wasmexecutor.WasmExecutor
	producer         sarama.SyncProducer
)

func main() {
	keepRunning := true
	log.Println("Starting a new Sarama consumer")

	if c.Kafka.Verbose {
		sarama.Logger = logsarama.LogX{}
	}

	version, err := sarama.ParseKafkaVersion(c.Kafka.Version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()
	config.Version = version
	config.ChannelBufferSize = cgroup.AvailableCPUs() * 2
	config.Consumer.MaxProcessingTime = time.Second * 10
	//config.Consumer.Offsets.AutoCommit.Enable = false
	switch c.Kafka.Assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", c.Kafka.Assignor)
	}

	if c.Kafka.Oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	/**
	 * Setup a new Sarama consumer group
	 */
	consumer := Consumer{
		ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(c.Kafka.Brokers, c.Kafka.Group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, c.Kafka.Topics, &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}

			ctx := context.Background()

			var id uuid.UUID
			copy(id[:], message.Value)
			judge, err := judgeModel.FindOne(context.Background(), id)
			if err != nil {
				in := new(pb.JudgeReq)
				err = proto.Unmarshal(message.Value[16:], in)
				if err != nil {
					return err
				}
				judge = &model.Judge{
					Id:         id,
					UserId:     pb.ToUUID(in.UserId),
					ProblemId:  pb.ToUUID(in.ProblemId),
					Status:     0,
					Code:       stringu.B2S(in.Code),
					Lang:       int16(in.Lang),
					TimeUsed:   0,
					MemoryUsed: 0,
				}
			}
			err = ConsumeJudge(ctx, judge)
			if err != nil {
				_, _, err = producer.SendMessage(&sarama.ProducerMessage{
					Topic: message.Topic,
					Key:   sarama.ByteEncoder(message.Key),
					Value: sarama.ByteEncoder(message.Value),
				})
				if err != nil {
					return err
				}
			}
			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

const (
	Waiting uint8 = iota
	Processing
	Compiling
	Running
	Judging
	Accepted
	PresentationError
	MemoryLimitExceeded
	TimeLimitExceeded
	OutputLimitExceeded
	RuntimeError
	CompilationError
	WrongAnswer
)

func ConsumeJudge(ctx context.Context, judge *model.Judge) error {
	switch uint8(judge.Status) {
	case Waiting:
		return ConsumeJudgeProcessing(ctx, judge)
	case Processing, Compiling, Running, Judging:
		return ConsumeJudgeCompiling(ctx, judge)
	case Accepted, MemoryLimitExceeded, TimeLimitExceeded, OutputLimitExceeded, RuntimeError, CompilationError, PresentationError, WrongAnswer:
		return nil
	default:
		panic("unhandled default case")
	}
}

func ConsumeJudgeProcessing(ctx context.Context, judge *model.Judge) error {
	judge.Status = int16(Processing)
	_, err := judgeModel.Upsert(ctx, judge)
	if err != nil {
		return err
	}
	return ConsumeJudgeCompiling(ctx, judge)
}

func ConsumeJudgeCompiling(ctx context.Context, judge *model.Judge) error {
	judge.Status = int16(Compiling)
	_ = judgeModel.Update(ctx, judge)

	var buf bytebufferpool.ByteBuffer
	err := compile.CompileWASM(stringu.S2B(judge.Code), pb.LangType(judge.Lang), &buf)
	if err != nil {
		if err == compile.Err {
			judge.Status = int16(CompilationError)
			err = judgeModel.Update(ctx, judge)
		}
		if err != nil {
			return err
		}
	}

	return ConsumeJudgeRunning(ctx, judge, buf.Bytes())
}

const langAll = 66

func ConsumeJudgeRunning(ctx context.Context, judge *model.Judge, wasmBinary []byte) error {
	judge.Status = int16(Running)
	_ = judgeModel.Update(ctx, judge)

	langCtx, err := problemLangModel.FindOneByProblemIdLang(ctx, judge.ProblemId, judge.Lang)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			langCtx, err = problemLangModel.FindOneByProblemIdLang(ctx, judge.ProblemId, langAll)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	datas, err := problemDataModel.FindJudgeDataHash(ctx, judge.ProblemId)
	if err != nil {
		return err
	}

	req := &wasmexecutor.WasmExecutionRequest{
		WasmBinaryArr:  [][]byte{wasmBinary},
		Inputs:         make([]*wasm.WasmExecutionInput, len(datas)),
		Compression:    0,
		ReturnHashOnly: true,
	}

	for i := 0; i < len(datas); i++ {
		req.Inputs[i] = &wasm.WasmExecutionInput{
			Stdin:       stringu.S2B(datas[i].Input),
			Args:        nil,
			Envs:        nil,
			MemoryLimit: uint64(langCtx.MemoryLimit),
			FuelLimit:   uint64(langCtx.TimeLimit),
			StdoutLimit: uint64(datas[i].OutputLen) * 2,
			StderrLimit: uint64(datas[i].OutputLen) * 2,
		}
	}

	res, err := wasmExecutor.Execute(context.Background(), req)
	if err != nil {
		return err
	}

	return ConsumeJudgeJudging(ctx, judge, res.Outputs, datas)
}

func ConsumeJudgeJudging(ctx context.Context, judge *model.Judge, outputs []*wasmexecutor.WasmExecutionOutput, datas []model.ProblemData) error {
	judge.Status = int16(Judging)
	_ = judgeModel.Update(ctx, judge)

	cases := make([]model.JudgeCase, len(outputs))
	for i := 0; i < len(outputs); i++ {
		status := uint8(outputs[i].Status)

		if uint8(outputs[i].Status) == Accepted {
			if !memoryEqual[uint64, int64](outputs[i].StdoutHash, datas[i].OutputHash) {
				if !memoryEqual[uint64, int64](outputs[i].StdoutTokenStreamHash, datas[i].OutputTokenHash) {
					status = WrongAnswer
				} else {
					status = PresentationError
				}
			}
		}

		cases[i] = model.JudgeCase{
			Id:            uuid.UUID{},
			JudgeId:       judge.Id,
			ProblemDataId: datas[i].Id,
			Status:        int16(status),
			TimeUsed:      int64(outputs[i].FuelConsumed),
			MemoryUsed:    int64(outputs[i].MemoryUsed),
			Reason:        stringu.B2S(outputs[i].Stderr),
		}

		judge.MemoryUsed = max(cases[i].MemoryUsed, judge.MemoryUsed)
		judge.TimeUsed += cases[i].TimeUsed
		judge.Status = max(cases[i].Status, judge.Status)
	}

	_, err := judgeCaseModel.BatchInsert(ctx, cases)
	if err != nil {
		return err
	}

	return ConsumeJudgeEnd(ctx, judge)
}

func ConsumeJudgeEnd(ctx context.Context, judge *model.Judge) error {
	err := judgeModel.Update(ctx, judge)
	if err != nil {
		return err
	}
	return nil
}

func memoryEqual[A comparable, B comparable](a A, b B) bool {
	if unsafe.Sizeof(a) != unsafe.Sizeof(b) { // 常量，编译时优化。
		panic("BUG!")
	}
	return a == *(*A)(unsafe.Pointer(&b))
}
