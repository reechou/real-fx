package models

import (
	"fmt"
	"strconv"
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
	ReturnMoney float32 `xorm:"not null default 0.000 decimal(9,3)" json:"-"`
	Status      int64   `xorm:"not null default 0 int index"`
	CreatedAt   int64   `xorm:"not null default 0 int index"`
	UpdatedAt   int64   `xorm:"not null default 0 int"`
}

type FxOrderWaitSettlementRecord struct {
	ID          int64   `xorm:"pk autoincr"`
	AccountId   int64   `xorm:"not null default 0 int index"`
	UnionId     string  `xorm:"not null default '' varchar(128) index"`
	OrderId     string  `xorm:"not null default '' varchar(128) index"`
	ReturnMoney float32 `xorm:"not null default 0.000 decimal(9,3)"`
	Level       int64   `xorm:"not null default 0 int index"`
	CreatedAt   int64   `xorm:"not null default 0 int index"`
}

type FxOrderSettlementRecord struct {
	ID          int64   `xorm:"pk autoincr"`
	AccountId   int64   `xorm:"not null default 0 int index"`
	UnionId     string  `xorm:"not null default '' varchar(128) index"`
	OrderId     string  `xorm:"not null default 0 int index"`
	ReturnMoney float32 `xorm:"not null default 0.000 decimal(9,3)"`
	SourceId    string  `xorm:"not null default '' varchar(128)"`
	Level       int64   `xorm:"not null default 0 int index"`
	CreatedAt   int64   `xorm:"not null default 0 int index"`
	UpdatedAt   int64   `xorm:"not null default 0 int"`
}

func CreateFxOrder(info *FxOrder) error {
	if info.UnionId == "" || info.OrderId == "" {
		return fmt.Errorf("fx order union_id[%s] order_id[%s] cannot be nil.", info.UnionId, info.OrderId)
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now
	_, err := x.Insert(info)
	if err != nil {
		logrus.Errorf("create fx order error: %v", err)
		return err
	}
	logrus.Infof("fx order union_id[%s] order_id[%s] create success.", info.UnionId, info.OrderId)

	return nil
}

func UpdateFxOrderStatus(info *FxOrder) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Id(info.ID).Cols("status", "updated_at").Update(info)
	return err
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

func GetFxOrderListCount(unionId string) (int64, error) {
	count, err := x.Where("union_id = ?", unionId).Count(&FxOrder{})
	if err != nil {
		logrus.Errorf("union_id[%s] get fx order list count error: %v", unionId, err)
		return 0, err
	}
	return count, nil
}

func GetFxOrderListCountById(accountId int64) (int64, error) {
	count, err := x.Where("account_id = ?", accountId).Count(&FxOrder{})
	if err != nil {
		logrus.Errorf("account_id[%d] get fx order list count error: %v", accountId, err)
		return 0, err
	}
	return count, nil
}

func GetFxOrderList(unionId string, offset, num, status int64) ([]FxOrder, error) {
	var fxOrderList []FxOrder
	err := x.Where("union_id = ?", unionId).And("status = ?", status).Desc("created_at").Limit(int(num), int(offset)).Find(&fxOrderList)
	if err != nil {
		logrus.Errorf("union_id[%s] get fx order list error: %v", unionId, err)
		return nil, err
	}
	return fxOrderList, nil
}

func GetFxOrderListById(accountId int64, offset, num, status int64) ([]FxOrder, error) {
	var fxOrderList []FxOrder
	err := x.Where("account_id = ?", accountId).And("status = ?", status).Desc("created_at").Limit(int(num), int(offset)).Find(&fxOrderList)
	if err != nil {
		logrus.Errorf("account_id[%d] get fx order list error: %v", accountId, err)
		return nil, err
	}
	return fxOrderList, nil
}

func IterateFxWaitOrder(status int64, f xorm.IterFunc) error {
	err := x.Where("status = ?", status).Iterate(&FxOrder{}, f)
	if err != nil {
		logrus.Errorf("iterate fx order error: %v", err)
		return err
	}
	return nil
}

func CreateFxOrderSettlementRecord(info *FxOrderSettlementRecord) error {
	if info.UnionId == "" || info.OrderId == "" || info.SourceId == "" {
		return fmt.Errorf("fx order settlement record union_id[%s] order_id[%s] source_id[%s] cannot be nil.", info.UnionId, info.OrderId, info.SourceId)
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now
	_, err := x.Insert(info)
	if err != nil {
		logrus.Errorf("create fx order settlement record error: %v", err)
		return err
	}
	logrus.Infof("fx order settlement record union_id[%s] order_id[%s] source_id[%s] create success.", info.UnionId, info.OrderId, info.SourceId)
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

func GetFxOrderSettlementRecordListCount(unionId string) (int64, error) {
	count, err := x.Where("union_id = ?", unionId).Count(&FxOrderSettlementRecord{})
	if err != nil {
		logrus.Errorf("union_id[%s] get fx order settlement record list count error: %v", unionId, err)
		return 0, err
	}
	return count, nil
}

func GetFxOrderSettlementRecordListCountById(accountId int64) (int64, error) {
	count, err := x.Where("account_id = ?", accountId).Count(&FxOrderSettlementRecord{})
	if err != nil {
		logrus.Errorf("account_id[%d] get fx order settlement record list count error: %v", accountId, err)
		return 0, err
	}
	return count, nil
}

func GetFxOrderSettlementRecordList(unionId string, offset, num int64) ([]FxOrderSettlementRecord, error) {
	var fxOrderSMRecordList []FxOrderSettlementRecord
	err := x.Where("union_id = ?", unionId).Desc("created_at").Limit(int(num), int(offset)).Find(&fxOrderSMRecordList)
	if err != nil {
		logrus.Errorf("union_id[%s] get fx order settlement record list error: %v", unionId, err)
		return nil, err
	}
	return fxOrderSMRecordList, nil
}

func GetFxOrderSettlementRecordListByid(accountId int64, offset, num int64) ([]FxOrderSettlementRecord, error) {
	var fxOrderSMRecordList []FxOrderSettlementRecord
	err := x.Where("account_id = ?", accountId).Desc("created_at").Limit(int(num), int(offset)).Find(&fxOrderSMRecordList)
	if err != nil {
		logrus.Errorf("accunt_id[%d] get fx order settlement record list error: %v", accountId, err)
		return nil, err
	}
	return fxOrderSMRecordList, nil
}

func GetFxOrderSettlementRecordFromOrderId(info *FxOrderSettlementRecord) (bool, error) {
	has, err := x.Where("order_id = ?", info.OrderId).Get(info)
	if err != nil {
		logrus.Errorf("get fx order sr from order id[%s] error: %v", info.OrderId, err)
		return has, err
	}
	return has, nil
}

func CreateFxOrderWaitSettlementRecordList(list []FxOrderWaitSettlementRecord) error {
	if len(list) == 0 {
		return nil
	}
	_, err := x.Insert(&list)
	if err != nil {
		logrus.Errorf("create fx order wait settlement record list error: %v", err)
		return err
	}
	return nil
}

func GetFxOrderWaitSettlementRecordListSumById(accountId int64, status int64) (float32, error) {
	//sum, err := x.Where("union_id = ?", unionId).Sum(&FxOrderWaitSettlementRecord{}, "return_money")
	//sum, err := x.Table("fx_order_wait_settlement_record").Select("fx_order_wait_settlement_record.return_money").
	//	Join("LEFT", "fx_order", "fx_order_wait_settlement_record.order_id = fx_order.order_id").Where("fx_order.status = ?", status).
	//	And("fx_order_wait_settlement_record.union_id = ?", unionId).Sum()
	results, err := x.Query("select sum(r.return_money) as sum_money from fx_order_wait_settlement_record as r left join fx_order as o on r.order_id = o.order_id where o.status = ? and r.account_id = ?", status, accountId)
	if err != nil {
		logrus.Errorf("account_id[%d] get fx order wait settlement record list sum(return_money) error: %v", accountId, err)
		return 0, err
	}
	if len(results) == 0 {
		return 0, nil
	}
	sumMoney := results[0]["sum_money"]
	sum, err := strconv.ParseFloat(string(sumMoney), 64)
	if err != nil {
		logrus.Errorf("account_id[%d] get fx order wait settlement record list strconv sum(return_money) error: %v", accountId, err)
		return 0, err
	}
	return float32(sum), nil
}

func GetFxOrderWaitSettlementRecordListCount(unionId string) (int64, error) {
	count, err := x.Where("union_id = ?", unionId).And("level = 0").Count(&FxOrderWaitSettlementRecord{})
	if err != nil {
		logrus.Errorf("union_id[%s] get fx order wait settlement record list count error: %v", unionId, err)
		return 0, err
	}
	return count, nil
}

func GetFxOrderWaitSettlementRecordListCountById(accountId int64, status int) (int64, error) {
	count, err := x.Table("fx_order_wait_settlement_record").Join("LEFT", "fx_order", "fx_order_wait_settlement_record.order_id = fx_order.order_id").
		Where("fx_order_wait_settlement_record.account_id = ?", accountId).
		And("fx_order_wait_settlement_record.level = 0").
		And("fx_order.status = ?", status).
		Count(&FxOrderWaitSettlementRecord{})
	//count, err := x.Where("account_id = ?", accountId).And("level = 0").Count(&FxOrderWaitSettlementRecord{})
	if err != nil {
		logrus.Errorf("account_id[%d] get fx order wait settlement record list count error: %v", accountId, err)
		return 0, err
	}
	return count, nil
}

func GetFxOrderWaitSettlementRecordList(unionId string, offset, num int64) ([]FxOrderWaitSettlementRecord, error) {
	var fxOrderWSMRecordList []FxOrderWaitSettlementRecord
	err := x.Where("union_id = ?", unionId).And("level = 0").Desc("created_at").Limit(int(num), int(offset)).Find(&fxOrderWSMRecordList)
	if err != nil {
		logrus.Errorf("union_id[%s] get fx order wait settlement record list error: %v", unionId, err)
		return nil, err
	}
	return fxOrderWSMRecordList, nil
}

func GetFxOrderWaitSettlementRecordListById(accountId int64, offset, num int64, status int) ([]FxOrderWaitSettlementRecord, error) {
	var fxOrderWSMRecordList []FxOrderWaitSettlementRecord
	err := x.Table("fx_order_wait_settlement_record").Select("fx_order_wait_settlement_record.*").
		Join("LEFT", "fx_order", "fx_order_wait_settlement_record.order_id = fx_order.order_id").Where("fx_order_wait_settlement_record.account_id = ?", accountId).
		And("fx_order_wait_settlement_record.level = 0").And("fx_order.status = ?", status).Find(&fxOrderWSMRecordList)
	//err := x.Where("account_id = ?", accountId).And("level = 0").Desc("created_at").Limit(int(num), int(offset)).Find(&fxOrderWSMRecordList)
	if err != nil {
		logrus.Errorf("account_id[%d] get fx order wait settlement record list error: %v", accountId, err)
		return nil, err
	}
	return fxOrderWSMRecordList, nil
}
