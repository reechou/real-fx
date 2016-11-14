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

const (
	WITHDRAWAL_WAITING = iota
	WITHDRAWAL_DONE
)

const (
	FX_HISTORY_TYPE_SIGN = iota
	FX_HISTORY_TYPE_INVITE
	FX_HISTORY_TYPE_ORDER_0
	FX_HISTORY_TYPE_ORDER_1
	FX_HISTORY_TYPE_ORDER_2
)

var (
	FxHistoryDescs = []string{"每日签到", "邀请下线", "订单返积分", "一级分销积分", "二级分销积分"}
)
