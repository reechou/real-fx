package fx_models

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-xorm/xorm"
)

type FxOrder struct {
	ID          int64   `xorm:"pk autoincr"`
	AccountId   int64   `xorm:"not null default 0 int index"`
	UnionId     string  `xorm:"not null default '' varchar(128) index"`
	OrderId     string  `xorm:"not null default '' varchar(128) unique"`
	OrderName   string  `xorm:"not null default '' varchar(128)"`
	Price       float32 `xorm:"not null default 0.000 decimal(10,3)"`
	ReturnMoney float32 `xorm:"not null default 0.000 decimal(9,3)"`
	Status      int64   `xorm:"not null default 0 int index"`
	CreatedAt   int64   `xorm:"not null default 0 int"`
	UpdatedAt   int64   `xorm:"not null default 0 int"`
}

type FxOrderSettlementRecord struct {
	ID          int64   `xorm:"pk autoincr"`
	AccountId   int64   `xorm:"not null default 0 int index"`
	UnionId     string  `xorm:"not null default '' varchar(128) index"`
	OrderId     string  `xorm:"not null default 0 int index"`
	ReturnMoney float32 `xorm:"not null default 0.000 decimal(9,3)"`
	SourceId    string  `xorm:"not null default '' varchar(128)"`
	Level       int64   `xorm:"not null default 0 int"`
	CreatedAt   int64   `xorm:"not null default 0 int"`
	UpdatedAt   int64   `xorm:"not null default 0 int"`
}

func GetFxOrderInfo(info *FxOrder) (bool, error) {
	has, err := x.Where("order_id = ?", info.OrderId).Get(info)
	if err != nil {
		logrus.Errorf("get fx order[%s] error: %v", info.OrderId, err)
		return false, fmt.Errorf("get fx order[%s] error: %v", info.OrderId, err)
	}
	if !has {
		logrus.Errorf("get fx order[%s] has no this order.", info.OrderId)
		return false, nil
	}
	return true, nil
}

func UpdateFxOrderStatus(info *FxOrder) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Id(info.ID).Cols("status", "updated_at").Update(info)
	return err
}

func IterateFxWaitOrder(status int64, f xorm.IterFunc) error {
	logrus.Debugf("interate fx waiting orders...")
	err := x.Where("status = ?", status).Iterate(&FxOrder{}, f)
	if err != nil {
		logrus.Errorf("iterate fx order error: %v", err)
		return err
	}
	return nil
}

func CreateFxOrderSettlementRecordList(list []FxOrderSettlementRecord) error {
	if len(list) == 0 {
		return nil
	}
	_, err := x.Insert(&list)
	if err != nil {
		logrus.Errorf("create fx order settlement record list error: %v", err)
		return err
	}
	return nil
}