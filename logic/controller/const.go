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
	WITHDRAWAL_DONE    = 1    // 提现完成
	WITHDRAWAL_WAITING = 2    // 审核中
	WITHDRAWAL_FAIL    = 3    // 提现失败
)

// 积分历史记录类型
const (
	FX_HISTORY_TYPE_SIGN = iota // 签到
	FX_HISTORY_TYPE_INVITE      // 邀请
	FX_HISTORY_TYPE_ORDER_0     // 订单主
	FX_HISTORY_TYPE_ORDER_1     // 1级分销
	FX_HISTORY_TYPE_ORDER_2     // 2级分销
	FX_HISTORY_TYPE_WITHDRAWAL  // 提现
	FX_HISTORY_TYPE_SCORE_MALL  // 积分商城
)

var (
	FxHistoryDescs = []string{
		"每日签到",
		"邀请下线",
		"订单返积分",
		"一级分销积分",
		"二级分销积分",
		"提现",
		"积分商城",
	}
)
