package fx_models

import (
	"time"

	"github.com/Sirupsen/logrus"
)

type SettlementFxOrderInfo struct {
	Status        int64
	Order         *FxOrder
	OrderAddMoney float32
}

func SettlementOwnerFxOrder(info *SettlementFxOrderInfo) error {
	now := time.Now().Unix()
	_, err := x.Exec("update fx_account fa, fx_order fo set fa.can_withdrawals=fa.can_withdrawals+?, fa.updated_at=?, fo.status=?, fo.updated_at=? where fa.union_id=? and fo.order_id=?",
		info.OrderAddMoney, now, info.Status, now, info.Order.UnionId, info.Order.OrderId)
	if err != nil {
		logrus.Errorf("settlement owner[%s] fx order[%s] error: %v", info.Order.UnionId, info.Order.OrderId, err)
		return err
	}
	return nil
}
