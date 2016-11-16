package controller

import (
	"github.com/reechou/real-fx/config"
	"github.com/reechou/real-fx/logic/ext"
)

type Daemon struct {
	cfg *config.Config

	we  *ext.WechatExt
	cww *CashWithdrawalWorker
}

func NewDaemon(cfg *config.Config) *Daemon {
	d := &Daemon{
		cfg: cfg,
	}
	d.cww = NewCashWithdrawalWorker(DEFAULT_MAX_WORKER, DEFAULT_MAX_CHAN_LEN, d.cfg)
	d.we = ext.NewWechatExt(d.cfg)
	return d
}
