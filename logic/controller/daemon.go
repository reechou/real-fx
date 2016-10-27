package controller

import (
	"github.com/reechou/real-fx/config"
)

type Daemon struct {
	cfg *config.Config

	cww *CashWithdrawalWorker
}

func NewDaemon(cfg *config.Config) *Daemon {
	d := &Daemon{
		cfg: cfg,
	}
	d.cww = NewCashWithdrawalWorker(DEFAULT_MAX_WORKER, DEFAULT_MAX_CHAN_LEN, d.cfg)
	return d
}
