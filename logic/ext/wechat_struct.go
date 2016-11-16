package ext

const (
	WECHAT_RESPONSE_OK = 1000
)

type WithdrawalReq struct {
	OpenId      string `json:"openId"`
	TotalAmount int64  `json:"total_amount"`
	MchBillno   string `json:"mch_billno"`
}

type WechatResponse struct {
	Code int64       `json:"state"`
	Msg  string      `json:"message,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
