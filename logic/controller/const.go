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
	FX_ORDER_SUCCESS    = 1 // 订单成功
	FX_ORDER_FAIL       = 2 // 订单失败
	FX_ORDER_WAIT       = 3 // 订单等待结算
	FX_ORDER_SETTLEMENT = 4 // 订单结算成功
)

const (
	WITHDRAWAL_WAITING = iota
	WITHDRAWAL_DONE
)
