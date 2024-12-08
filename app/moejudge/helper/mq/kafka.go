package mq

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/IBM/sarama"
	"log"
	"os"
)

type KafkaProducer struct {
	Topic    string
	Producer sarama.SyncProducer
}

func NewKafkaProducer(topic string, p sarama.SyncProducer) *KafkaProducer {
	return &KafkaProducer{
		Topic:    topic,
		Producer: p,
	}
}

func (k *KafkaProducer) Send(key string, message []byte) error {
	_, _, err := k.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: k.Topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(message),
	})
	return err
}

// Config holds the configuration options for the Kafka client and server settings.
type KafkaConfig struct {
	// Addr is the address the server binds to.
	Addr string `json:"addr"`

	// Brokers is the comma-separated list of Kafka brokers to connect to.
	Brokers []string `json:"brokers"`

	// Version specifies the Kafka cluster version.
	Version string `json:"version"`

	// Verbose enables Sarama logging for debugging.
	Verbose bool `json:"verbose,optional"`

	// CertFile is the optional certificate file used for client authentication.
	CertFile string `json:"certificate,optional"`

	// KeyFile is the optional key file used for client authentication.
	KeyFile string `json:"key,optional"`

	// CaFile is the optional certificate authority file for TLS client authentication.
	CaFile string `json:"ca,optional"`

	// TLSSkipVerify indicates whether to skip the verification of the TLS server certificate.
	TLSSkipVerify bool `json:"tls-skip-verify,optional"`
}

func NewSyncProducer(c *KafkaConfig) sarama.SyncProducer {
	version, err := sarama.ParseKafkaVersion(c.Version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}
	if c.Verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	config.Version = version
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	tlsConfig := createTlsConfiguration(c)
	if tlsConfig != nil {
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Enable = true
	}

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewSyncProducer(c.Brokers, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	return producer
}

func createTlsConfiguration(c *KafkaConfig) (t *tls.Config) {
	if c.CertFile != "" && c.KeyFile != "" && c.CaFile != "" {
		cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
		if err != nil {
			log.Fatal(err)
		}

		caCert, err := os.ReadFile(c.CaFile)
		if err != nil {
			log.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: c.TLSSkipVerify,
		}
	}
	// will be nil by default if nothing is provided
	return t
}
