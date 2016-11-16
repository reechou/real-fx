package controller

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/real-fx/logic/tools/order_check/config"
	"github.com/reechou/real-fx/logic/tools/order_check/fx_models"
)

type SettlementWorker struct {
	orderChanList []chan *fx_models.FxOrder

	cfg *config.Config

	wg   sync.WaitGroup
	stop chan struct{}
}

func NewSettlementWorker(maxWorker, maxChanLen int, cfg *config.Config) *SettlementWorker {
	sw := &SettlementWorker{
		cfg:  cfg,
		stop: make(chan struct{}),
	}
	for i := 0; i < maxWorker; i++ {
		orderChan := make(chan *fx_models.FxOrder, maxChanLen)
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

func (sw *SettlementWorker) SettlementOrder(order *fx_models.FxOrder) {
	idx := int(order.ID) % len(sw.orderChanList)
	select {
	case sw.orderChanList[idx] <- order:
	case <-time.After(5 * time.Second):
		logrus.Errorf("settlement into order channel timeout.")
	}
}

func (sw *SettlementWorker) runWorker(orderChan chan *fx_models.FxOrder, stop chan struct{}) {
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

func (sw *SettlementWorker) do(order *fx_models.FxOrder) {
	if order.Status != FX_ORDER_WAIT {
		logrus.Errorf("order[%v] cannot be settlement.", order)
		return
	}

	// check status
	checkOrder := &fx_models.FxOrder{
		OrderId: order.OrderId,
	}
	_, err := fx_models.GetFxOrderInfo(checkOrder)
	if err != nil {
		logrus.Errorf("get fx order[%v] status error: %v", order, err)
		return
	}
	if checkOrder.Status != FX_ORDER_WAIT {
		logrus.Errorf("order[%v] cannot be settlement, order status: %d", order, checkOrder.Status)
		return
	}

	var levelReturns []float32
	for i := 0; i < len(sw.cfg.SettlementCommission.LevelPer); i++ {
		lReturn := order.ReturnMoney * float32(sw.cfg.SettlementCommission.LevelPer[i]) / 100.0 * float32(sw.cfg.Score.EnlargeScale)
		levelReturns = append(levelReturns, lReturn)
	}

	settlementFxOrder := &fx_models.SettlementFxOrderInfo{
		Status:        FX_ORDER_SUCCESS,
		Order:         order,
		OrderAddMoney: levelReturns[0],
	}
	err = fx_models.SettlementOwnerFxOrder(settlementFxOrder)
	if err != nil {
		logrus.Errorf("do settlement order[%v] settlement owner order error: %v", order, err)
		return
	}
	logrus.Infof("order_id[%s] settlement for owner[%s] with return_money[%f] success", order.OrderId, order.UnionId, levelReturns[0])

	now := time.Now().Unix()

	//month := fmt.Sprintf(time.Now().Format("200601"))
	//err = sw.updateFxAccountMonth(month, order.UnionId, levelReturns[0])
	//if err != nil {
	//	logrus.Errorf("do settlement order[%v] update fx account month owner order error: %v", order, err)
	//	return err
	//}

	fxAccount := &fx_models.FxAccount{
		UnionId: order.UnionId,
	}
	has, err := fx_models.GetFxAccount(fxAccount)
	if err != nil {
		logrus.Errorf("do settlement order[%v] in level[0] get fx account from union_id[%d] error: %v",
			order, order.UnionId, err)
		return
	}
	if !has {
		logrus.Errorf("do settlement no this owner account[%s]", order.UnionId)
		return
	}

	//var recordList []fx_models.FxOrderSettlementRecord
	//recordList = append(recordList, fx_models.FxOrderSettlementRecord{
	//	AccountId:   fxAccount.ID,
	//	UnionId:     order.UnionId,
	//	OrderId:     order.OrderId,
	//	ReturnMoney: levelReturns[0],
	//	SourceId:    order.UnionId,
	//	Level:       0,
	//	CreatedAt:   now,
	//	UpdatedAt:   now,
	//})

	var historyList []fx_models.FxAccountHistory
	historyList = append(historyList, fx_models.FxAccountHistory{
		AccountId:  fxAccount.ID,
		UnionId:    fxAccount.UnionId,
		Score:      levelReturns[0],
		ChangeType: int64(FX_HISTORY_TYPE_ORDER_0),
		ChangeDesc: FxHistoryDescs[FX_HISTORY_TYPE_ORDER_0],
		CreatedAt:  now,
	})

	unionId := fxAccount.Superior
	for i := 1; i < len(levelReturns); i++ {
		// get upper
		fxAccount := &fx_models.FxAccount{
			UnionId: unionId,
		}
		has, err := fx_models.GetFxAccount(fxAccount)
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
		err = fx_models.AddFxAccountMoney(levelReturns[i], fxAccount)
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

		//recordList = append(recordList, fx_models.FxOrderSettlementRecord{
		//	AccountId:   fxAccount.ID,
		//	UnionId:     unionId,
		//	OrderId:     order.OrderId,
		//	ReturnMoney: levelReturns[i],
		//	SourceId:    order.UnionId,
		//	Level:       int64(i),
		//	CreatedAt:   now,
		//	UpdatedAt:   now,
		//})
		historyList = append(historyList, fx_models.FxAccountHistory{
			AccountId:  fxAccount.ID,
			UnionId:    fxAccount.UnionId,
			Score:      levelReturns[i],
			ChangeType: int64(FX_HISTORY_TYPE_ORDER_0 + i),
			ChangeDesc: FxHistoryDescs[FX_HISTORY_TYPE_ORDER_0+i],
			CreatedAt:  now,
		})
		unionId = fxAccount.Superior
	}

	// insert history
	//err = fx_models.CreateFxOrderSettlementRecordList(recordList)
	//if err != nil {
	//	logrus.Errorf("create fx order[%d] settlement record list error: %v", order, err)
	//}
	fx_models.CreateFxAccountHistoryList(historyList)
}

func (sw *SettlementWorker) updateFxAccountMonth(month, unionId string, returnMoney float32) error {
	fxAccountMonth := &fx_models.FxAccountMonthAchievement{
		UnionId:              unionId,
		Month:                month,
		ThisMonthAchievement: returnMoney,
	}
	has, err := fx_models.GetFxAccountMonthAchievement(fxAccountMonth)
	if err != nil {
		logrus.Errorf("get fx account month[%s] achievement error: %v", month, err)
		return err
	}
	if !has {
		err = fx_models.CreateFxAccountMonthAchievement(fxAccountMonth)
		if err != nil {
			logrus.Errorf("create fx account month[%s] achievement error: %v", month, err)
			return err
		}
	} else {
		err = fx_models.UpdateFxAccountMonthAchievement(returnMoney, fxAccountMonth)
		if err != nil {
			logrus.Errorf("update fx account union_id[%s] month[%s] achievement error: %v", unionId, month, err)
			return err
		}
	}
	return nil
}
