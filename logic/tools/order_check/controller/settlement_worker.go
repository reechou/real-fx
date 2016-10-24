package controller

import (
	//"fmt"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/real-fx/config"
	"github.com/reechou/real-fx/logic/models"
)

type SettlementWorker struct {
	orderChanList []chan *models.FxOrder

	cfg *config.Config

	wg   sync.WaitGroup
	stop chan struct{}
}

func NewSettlementWorker(maxWorker, maxChanLen int64, cfg *config.Config) *SettlementWorker {
	sw := &SettlementWorker{
		cfg:  cfg,
		stop: make(chan struct{}),
	}
	for i := 0; i < maxWorker; i++ {
		orderChan := make(chan *models.FxOrder, maxChanLen)
		sw.orderChanList = append(sw.orderChanList, orderChan)
		sw.wg.Add(1)
		go sw.runWorker(orderChan, sw.stop)
	}
	return sw
}

func (sw *SettlementWorker) Close() {
	close(sw.stop)
	sw.wg.Wait()
}

func (sw *SettlementWorker) SettlementOrder(order *models.FxOrder) {
	idx := order.ID % len(sw.orderChanList)
	select {
	case sw.orderChanList[idx] <- order:
	case time.After(5 * time.Second):
		logrus.Errorf("settlement into order channel timeout.")
	}
}

func (sw *SettlementWorker) runWorker(orderChan chan *models.FxOrder, stop chan struct{}) {
	for {
		select {
		case order := <-orderChan:
			sw.do(order)
		case <-stop:
			sw.wg.Done()
			return
		}
	}
}

func (sw *SettlementWorker) do(order *models.FxOrder) {
	if order.Status != FX_ORDER_SUCCESS {
		logrus.Errorf("order[%v] cannot be settlement.", order)
		return
	}

	// check status
	checkOrder := &models.FxOrder{
		OrderId: order.OrderId,
	}
	_, err := models.GetFxOrderInfo(checkOrder)
	if err != nil {
		logrus.Errorf("get fx order[%v] status error: %v", order, err)
		return err
	}
	if checkOrder.Status != FX_ORDER_SUCCESS {
		logrus.Errorf("order[%v] cannot be settlement, order status: %d", order, checkOrder.Status)
		return
	}

	var levelReturns []float32
	for i := 0; i < len(sw.cfg.SettlementCommission.LevelPer); i++ {
		lReturn := order.ReturnMoney * float32(sw.cfg.SettlementCommission.LevelPer[i]/100)
		levelReturns = append(levelReturns, lReturn)
	}

	settlementFxOrder := &models.SettlementFxOrderInfo{
		Status:        FX_ORDER_SETTLEMENT,
		Order:         order,
		OrderAddMoney: levelReturns[0],
	}
	err = models.SettlementOwnerFxOrder(settlementFxOrder)
	if err != nil {
		logrus.Errorf("do settlement order[%v] settlement owner order error: %v", order, err)
		return err
	}
	logrus.Infof("order_id[%s] settlement for owner[%s] with return_money[%f] success", order.OrderId, order.UnionId, levelReturns[0])

	now := time.Now().Unix()

	//month := fmt.Sprintf(time.Now().Format("200601"))
	//err = sw.updateFxAccountMonth(month, order.UnionId, levelReturns[0])
	//if err != nil {
	//	logrus.Errorf("do settlement order[%v] update fx account month owner order error: %v", order, err)
	//	return err
	//}

	var recordList []models.FxOrderSettlementRecord
	recordList = append(recordList, models.FxOrderSettlementRecord{
		UnionId:     order.UnionId,
		OrderId:     order.OrderId,
		ReturnMoney: levelReturns[0],
		SourceId:    order.UnionId,
		Level:       0,
		CreatedAt:   now,
		UpdatedAt:   now,
	})

	unionId := order.UnionId
	for i := 1; i < len(levelReturns); i++ {
		// get upper
		fxAccount := &models.FxAccount{
			UnionId: unionId,
		}
		has, err := models.GetFxAccount(fxAccount)
		if err != nil {
			logrus.Errorf("do settlement order[%v] in level[%d] get fx account from union_id[%d] error: %v",
				order, i, unionId, err)
			return
		}
		if !has {
			logrus.Debugf("do settlement no this account[%s]", unionId)
			break
		}
		// add return money
		unionId = fxAccount.Superior
		fxAccount.UnionId = unionId
		err = models.AddFxAccountMoney(levelReturns[i], fxAccount)
		if err != nil {
			logrus.Errorf("do settlement order[%v] in level[%d] add money in fx account from union_id[%d] error: %v",
				order, i, unionId, err)
			return
		}
		logrus.Infof("order_id[%s] settlement for upper user[%s][level-%d] with return_money[%f] success", order.OrderId, unionId, i, levelReturns[i])

		//err = sw.updateFxAccountMonth(month, unionId, levelReturns[i])
		//if err != nil {
		//	logrus.Errorf("do settlement order[%v] update fx account month union_id[%s][level-%d] order error: %v", order, unionId, i, err)
		//	return err
		//}

		recordList = append(recordList, models.FxOrderSettlementRecord{
			UnionId:     unionId,
			OrderId:     order.OrderId,
			ReturnMoney: levelReturns[i],
			SourceId:    order.UnionId,
			Level:       int64(i),
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	// insert history
	err = models.CreateFxOrderSettlementRecordList(recordList)
	if err != nil {
		logrus.Errorf("create fx order[%d] settlement record list error: %v", order, err)
	}
}

func (sw *SettlementWorker) updateFxAccountMonth(month, unionId string, returnMoney float32) error {
	fxAccountMonth := &models.FxAccountMonthAchievement{
		UnionId:              unionId,
		Month:                month,
		ThisMonthAchievement: returnMoney,
	}
	has, err := models.GetFxAccountMonthAchievement(fxAccountMonth)
	if err != nil {
		logrus.Errorf("get fx account month[%s] achievement error: %v", month, err)
		return
	}
	if !has {
		err = models.CreateFxAccountMonthAchievement(fxAccountMonth)
		if err != nil {
			logrus.Errorf("create fx account month[%s] achievement error: %v", month, err)
			return
		}
	} else {
		err = models.UpdateFxAccountMonthAchievement(returnMoney, fxAccountMonth)
		if err != nil {
			logrus.Errorf("update fx account union_id[%s] month[%s] achievement error: %v", unionId, month, err)
			return
		}
	}
	return nil
}
