package request

type TranInReq struct {
	Amount float64 `json:"amount"`
	UserId int64   `json:"userId"`
}

type TranReqWork struct {
	FromAmount float64 `json:"fromAmount"`
	FromUserId int64   `json:"fromUserId"`
	ToAmount   float64 `json:"toAmount"`
	ToUserId   int64   `json:"toUserId"`
}

type TranOutReq struct {
	TranInReq
}
