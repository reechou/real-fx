package fx_models

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

type FxAccount struct {
	ID             int64   `xorm:"pk autoincr"`
	UnionId        string  `xorm:"not null default '' varchar(128) unique"`
	Phone          string  `xorm:"not null default '' varchar(16)"`
	Name           string  `xorm:"not null default '' varchar(128)"`
	CanWithdrawals float32 `xorm:"not null default 0.000 decimal(10,3)"`
	AllScore       float32 `xorm:"not null default 0.000 decimal(10,3) index"`
	Ticket         string  `xorm:"not null default '' varchar(128)"`
	Superior       string  `xorm:"not null default '' varchar(128) index"`
	Status         int64   `xorm:"not null default 0 int"`
	CreatedAt      int64   `xorm:"not null default 0 int"`
	UpdatedAt      int64   `xorm:"not null default 0 int"`
}

type FxAccountMonthAchievement struct {
	ID                   int64   `xorm:"pk autoincr"`
	UnionId              string  `xorm:"not null default '' varchar(128) unique(uni_user_month)"`
	Month                string  `xorm:"not null default '' varchar(16) unique(uni_user_month)"`
	ThisMonthAchievement float32 `xorm:"not null default 0.000 decimal(9,3)"`
	CreatedAt            int64   `xorm:"not null default 0 int"`
	UpdatedAt            int64   `xorm:"not null default 0 int"`
}

func AddFxAccountMoney(allAdd float32, info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	var err error
	_, err = x.Exec("update fx_account set can_withdrawals=can_withdrawals+?, all_score=all_score+?, updated_at=? where union_id=?",
		allAdd, allAdd, info.UpdatedAt, info.UnionId)
	return err
}

func GetFxAccount(info *FxAccount) (bool, error) {
	has, err := x.Where("union_id = ?", info.UnionId).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		logrus.Debugf("cannot find fx account from wx_unionid[%s]", info.UnionId)
		return false, nil
	}
	return true, nil
}

func CreateFxAccountMonthAchievement(info *FxAccountMonthAchievement) error {
	if info.UnionId == "" || info.Month == "" {
		return fmt.Errorf("argvs cannot be nil")
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		logrus.Errorf("create fx account month achievement error: %v", err)
		return err
	}
	logrus.Infof("create fx account month achievement from wx_unionid[%s] [%s] success.", info.UnionId, info.Month)

	return nil
}

func UpdateFxAccountMonthAchievement(addMoney float32, info *FxAccountMonthAchievement) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Exec("update fx_account_month_achievement set this_month_achievement=this_month_achievement+?, updated_at=? where union_id=? and month=?",
		addMoney, info.UpdatedAt, info.UnionId, info.Month)
	return err
}

func GetFxAccountMonthAchievement(info *FxAccountMonthAchievement) (bool, error) {
	has, err := x.Where("union_id = ?", info.UnionId).And("month = ?", info.Month).Get(info)
	if err != nil {
		logrus.Errorf("get union_id[%s] month[%s] fx account month achievement error: %v", info.UnionId, info.Month, err)
		return false, err
	}
	if !has {
		return false, nil
	}
	return true, nil
}
