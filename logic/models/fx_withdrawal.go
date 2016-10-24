package models

import (
	"time"

	"github.com/Sirupsen/logrus"
)

type WithdrawalRecord struct {
	ID              int64   `xorm:"pk autoincr"`
	AccountId       int64   `xorm:"not null default 0 int index"`
	UnionId         string  `xorm:"not null default '' varchar(128) index"`
	WithdrawalMoney float32 `xorm:"not null default 0.000 decimal(9,3)"`
	Balance         float32 `xorm:"not null default 0.000 decimal(9,3)"`
	Status          int64   `xorm:"not null default 0 int"`
	CreatedAt       int64   `xorm:"not null default 0 int"`
	UpdatedAt       int64   `xorm:"not null default 0 int"`
}

func CreateWithdrawalRecord(info *WithdrawalRecord) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	_, err := x.Insert(info)
	if err != nil {
		logrus.Errorf("create fx withdrawal record union_id[%s] error: %v", info.UnionId, err)
		return err
	}
	logrus.Infof("create fx withdrawal record union_id[%s] withdrawal[%f] balance[%f] create success.", info.UnionId, info.WithdrawalMoney, info.Balance)
	return nil
}

func UpdateWithdrawalRecordStatus(info *WithdrawalRecord) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Cols("status", "updated_at").Update(info, &WithdrawalRecord{ID: info.ID})
	return err
}

func GetWithdrawalRecordListCount(unionId string) (int64, error) {
	count, err := x.Where("union_id = ?", unionId).Count(&WithdrawalRecord{})
	if err != nil {
		logrus.Errorf("union_id[%s] get withdrawal record list count error: %v", unionId, err)
		return 0, err
	}
	return count, nil
}

func GetWithdrawalRecordListCountById(accountId int64) (int64, error) {
	count, err := x.Where("account_id = ?", accountId).Count(&WithdrawalRecord{})
	if err != nil {
		logrus.Errorf("account_id[%d] get withdrawal record list count error: %v", accountId, err)
		return 0, err
	}
	return count, nil
}

func GetWithdrawalRecordList(unionId string, offset, num int64) ([]WithdrawalRecord, error) {
	var list []WithdrawalRecord
	err := x.Where("union_id = ?", unionId).Limit(int(num), int(offset)).Find(&list)
	if err != nil {
		logrus.Errorf("union_id[%s] get withdrawal record list error: %v", unionId, err)
		return nil, err
	}
	return list, nil
}

func GetWithdrawalRecordListById(accountId int64, offset, num int64) ([]WithdrawalRecord, error) {
	var list []WithdrawalRecord
	err := x.Where("account_id = ?", accountId).Limit(int(num), int(offset)).Find(&list)
	if err != nil {
		logrus.Errorf("account_id[%d] get withdrawal record list error: %v", accountId, err)
		return nil, err
	}
	return list, nil
}
