package models

import (
	"github.com/Sirupsen/logrus"
)

type TaobaoOrderReal struct {
	ID         int64   `xorm:"id"`
	GoodsId    string  `xorm:"goodsId"`
	OrderId    string  `xorm:"orderId"`
	PayMoney   float32 `xorm:"payMoney"`
	GoodsState int64   `xorm:"goodsState"`
}

func GetTaobaoOrder(info *TaobaoOrderReal) (bool, error) {
	has, err := x.Where("orderId = ?", info.OrderId).And("goodsId = ?", info.GoodsId).And("payMoney = ?", info.PayMoney).Get(info)
	if err != nil {
		logrus.Errorf("get taobao order error: %v", err)
		return false, err
	}
	if !has {
		logrus.Errorf("get taobao order: no this order.")
		return false, nil
	}

	return true, nil
}
