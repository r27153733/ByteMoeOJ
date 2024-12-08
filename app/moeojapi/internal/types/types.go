// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3

package types

type Example struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type GetProblemReq struct {
	ID string `form:"id"`
}

type GetUserReq struct {
	ID string `form:"id"`
}

type JudgeReq struct {
	ProblemID string `json:"problem_id"`
	Code      string `json:"code"`
	Lang      uint8  `json:"lang"` // 编程语言
}

type JudgeResp struct {
	Result JudgeResult `json:"result"`
}

type JudgeResult struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Status     uint8  `json:"status"` // 判题状态，如 0：等待、1：通过、2：错误、3：超时、4：运行时错误等
	Code       string `json:"code"`
	Lang       uint8  `json:"lang"`                 // 编程语言
	Message    string `json:"message,optional"`     // 结果描述
	TimeUsed   uint64 `json:"time_used,optional"`   // 用时
	MemoryUsed uint64 `json:"memory_used,optional"` // 使用内存
}

type JudgeStatusReq struct {
	JudgeID string `form:"judge_id""`
}

type JudgeStatusResp struct {
	JudgeID string `json:"judge_id"`
}

type ListJudgeReq struct {
	ID        string `form:"id,optional"`
	UserID    string `form:"user_id,optional"`    // 可选，根据用户 ID 查询
	ProblemID string `form:"problem_id,optional"` // 可选，根据题目 ID 查询
	GroupID   string `form:"group_id,optional"`
	Status    uint8  `form:"status,optional"`    // 可选，根据判题状态查询
	Page      int    `form:"page,optional"`      // 分页页码
	PageSize  int    `form:"page_size,optional"` // 分页大小
}

type ListJudgeResp struct {
	List     []JudgeResult `json:"list"`               // 判题结果列表
	Total    int           `json:"total"`              // 总记录数
	Page     int           `json:"page,optional"`      // 当前页码
	PageSize int           `json:"page_size,optional"` // 当前页大小
}

type ListProblemReq struct {
	GroupID      string `form:"group_id,optional"`
	SearchString string `form:"contains_string"`
	Page         int    `form:"page,optional"`
	PageSize     int    `form:"page_size,optional"`
}

type ListProblemResp struct {
	List     []Problem `json:"list"`
	Total    int       `json:"total,optional"`
	Page     int       `json:"page,optional"`
	PageSize int       `json:"page_size,optional"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Problem struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Example    []Example `json:"example"`
	InitCode   string    `json:"init_code"`
	Lang       []uint8   `json:"lang"`
	Status     uint8     `json:"status"`
	Difficulty uint8     `json:"difficulty"`
	PassRate   uint8     `json:"pass_rate"`
}

type TryReq struct {
	ProblemID string `json:"problem_id"`
	Code      string `json:"code"`
	Lang      uint8  `json:"lang"` // 编程语言
}

type TryResp struct {
	Output string `json:"output"`
}

type UserInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserInfoToken struct {
	UserInfo
	Token string `json:"token"`
}

type WasmReq struct {
	ProblemID string `json:"problem_id"`
	Code      string `json:"code"`
	Lang      uint8  `json:"lang"` // 编程语言
}

type WasmResp struct {
	WasmBinary []byte `json:"wasm_binary"`
}
