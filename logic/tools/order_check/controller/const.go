package controller

const (
	GodSalesman = "godlike"
	MaxSalesman = 99900
)

const (
	DEFAULT_MAX_WORKER   = 100
	DEFAULT_MAX_CHAN_LEN = 10000
)

const (
	FX_ORDER_SUCCESS    = 1 // FenXiao 订单结算成功,后台结算成功
	FX_ORDER_WAIT       = 2 // 订单等待结算
	FX_ORDER_FAILED     = 3 // 订单失败
	FX_ORDER_SETTLEMENT = 4 // 淘宝结算
)

const (
	TAOBAO_ORDER_SUCCESS    = 1 // 订单成功
	TAOBAO_ORDER_PAY        = 2 // 订单付款
	TAOBAO_ORDER_INVALID    = 3 // 订单失效
	TAOBAO_ORDER_SETTLEMENT = 4 // 订单已结算
)
