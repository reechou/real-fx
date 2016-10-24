package controller

import (
	"sync"
	"time"
	
	"github.com/Sirupsen/logrus"
	"github.com/reechou/real-fx/config"
	"github.com/reechou/real-fx/logic/models"
)

type OrderCheck struct {
	cfg *config.Config
	
	sw *SettlementWorker

	wg   sync.WaitGroup
	stop chan struct{}
	done chan struct{}
}

func NewOrderCheck(cfg *config.Config) *OrderCheck {
	ocw := &OrderCheck{
		cfg:  cfg,
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	ocw.sw = NewSettlementWorker(cfg.WorkerInfo.SWMaxWorker, cfg.WorkerInfo.SWMaxChanLen, cfg)
	return ocw
}

func (ocw *OrderCheck) Stop() {
	close(ocw.stop)
	<-ocw.done
}

func (ocw *OrderCheck) Run() {
	for {
		select {
		case <-time.After(time.Duration(ocw.cfg.WorkerInfo.OrderCheckInterval) * time.Second):
			ocw.runCheck()
		case <-ocw.stop:
			close(ocw.done)
			return
		}
	}
}

func (ocw *OrderCheck) runCheck() {
	err := models.IterateFxWaitOrder(FX_ORDER_SETTLEMENT, ocw.handleOrder)
	if err != nil {
		logrus.Errorf("run check error: %v", err)
	}
}

func (ocw *OrderCheck) handleOrder(idx int, bean interface{}) error {
	order := bean.(*models.FxOrder)
	// check order status
	
	if order.Status == FX_ORDER_SUCCESS {
		// do settlement
		ocw.sw.SettlementOrder(order)
	}
	
	return nil
}
