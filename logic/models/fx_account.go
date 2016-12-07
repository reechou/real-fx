package models

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	AccountStatusFollow = iota
	AccountStatusUnfollow
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
	SignTime       int64   `xorm:"not null default 0 int index"`
	Status         int64   `xorm:"not null default 0 int"`
	CreatedAt      int64   `xorm:"not null default 0 int index"`
	UpdatedAt      int64   `xorm:"not null default 0 int"`
}

type FxAccountStatus struct {
	ID        int64  `xorm:"pk autoincr"`
	UnionId   string `xorm:"not null default '' varchar(128) unique"`
	CreatedAt int64  `xorm:"not null default 0 int"`
	UpdatedAt int64  `xorm:"not null default 0 int"`
}

type FxAccountMonthAchievement struct {
	ID                   int64   `xorm:"pk autoincr"`
	UnionId              string  `xorm:"not null default '' varchar(128) unique(uni_user_month)"`
	Month                string  `xorm:"not null default '' varchar(16) unique(uni_user_month)"`
	ThisMonthAchievement float32 `xorm:"not null default 0.000 decimal(9,3)"`
	CreatedAt            int64   `xorm:"not null default 0 int"`
	UpdatedAt            int64   `xorm:"not null default 0 int"`
}

type FxAccountFollow struct {
	ID        int64  `xorm:"pk autoincr"`
	UnionId   string `xorm:"not null default '' varchar(128)"`
	WXAccount string `xorm:"not null default '' varchar(64) unique(wx_account)"`
	OpenId    string `xorm:"not null default '' varchar(128) unique(wx_account)"`
	Status    int64  `xorm:"not null default 0 int"`
	CreatedAt int64  `xorm:"not null default 0 int"`
	UpdatedAt int64  `xorm:"not null default 0 int"`
}

func CreateFxAccount(info *FxAccount) (err error) {
	if info.UnionId == "" {
		return fmt.Errorf("wx union_id[%s] cannot be nil.", info.UnionId)
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err = x.Insert(info)
	if err != nil {
		logrus.Errorf("create fx account error: %v", err)
		return err
	}
	logrus.Infof("create fx account from wx_unionid[%s] success.", info.UnionId)

	return
}

func UpdateFxAccountBaseInfo(info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Cols("phone", "name", "updated_at").Update(info, &FxAccount{UnionId: info.UnionId})
	return err
}

func UpdateFxAccountStatus(info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Cols("status", "updated_at").Update(info, &FxAccount{UnionId: info.UnionId})
	return err
}

func UpdateFxAccountSalesman(info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Cols("ticket", "phone", "updated_at").Update(info, &FxAccount{UnionId: info.UnionId})
	return err
}

func UpdateFxAccountSignTime(allAdd float32, info *FxAccount) (int64, error) {
	now := time.Now().Unix()
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	dayZero := t.Unix() - 8*3600
	result, err := x.Exec("update fx_account set can_withdrawals=can_withdrawals+?, all_score=all_score+?, updated_at=?, sign_time=? where union_id=? and sign_time < ?",
		allAdd, allAdd, now, now, info.UnionId, dayZero)
	if err != nil {
		logrus.Errorf("update fx_account sign time error: %v", err)
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		logrus.Errorf("get affected error: %v", err)
		return 0, err
	}
	return affected, nil
}

func AddFxAccountMoney(allAdd float32, info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	var err error
	_, err = x.Exec("update fx_account set can_withdrawals=can_withdrawals+?, all_score=all_score+?, updated_at=? where union_id=?",
		allAdd, allAdd, info.UpdatedAt, info.UnionId)
	if err != nil {
		return err
	}
	logrus.Infof("fx account[%s] add money[%f] success.", info.UnionId, allAdd)
	return nil
}

func MinusFxAccountMoney(allMinus float32, info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	var err error
	_, err = x.Exec("update fx_account set can_withdrawals=can_withdrawals-?, updated_at=? where union_id=? and can_withdrawals >= ?",
		allMinus, info.UpdatedAt, info.UnionId, allMinus)
	if err != nil {
		return err
	}
	logrus.Infof("fx account[%s] minus money[%f] success.", info.UnionId, allMinus)
	return nil
}

func GetFxAccount(info *FxAccount) (bool, error) {
	has, err := x.Where("union_id = ?", info.UnionId).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		logrus.Errorf("cannot find fx account from wx_unionid[%s]", info.UnionId)
		return false, nil
	}
	return true, nil
}

func GetFxAccountById(info *FxAccount) (bool, error) {
	has, err := x.Where("id = ?", info.UnionId).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		logrus.Errorf("cannot find fx account from account_id[%s]", info.ID)
		return false, nil
	}
	return true, nil
}

func GetLowerPeopleCount(unionId string) (int64, error) {
	count, err := x.Where("superior = ?", unionId).Count(&FxAccount{})
	if err != nil {
		logrus.Errorf("union_id[%s] get lower peoples list count error: %v", unionId, err)
		return 0, err
	}
	return count, nil
}

func GetLowerPeople(unionId string, offset, num int64) ([]FxAccount, error) {
	var lowerPeoples []FxAccount
	err := x.Where("superior = ?", unionId).Desc("created_at").Limit(int(num), int(offset)).Find(&lowerPeoples)
	if err != nil {
		logrus.Errorf("union_id[%s] lower peoples list error: %v", unionId, err)
		return nil, err
	}
	return lowerPeoples, nil
}

func CreateFxAccountFollow(info *FxAccountFollow) error {
	if info.UnionId == "" || info.WXAccount == "" || info.OpenId == "" {
		return fmt.Errorf("argvs cannot be nil.")
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		logrus.Errorf("create fx account follow error: %v", err)
		return err
	}
	logrus.Infof("create fx account follow from wx_unionid[%s] [%s-%s] success.", info.UnionId, info.WXAccount, info.OpenId)

	return nil
}

func UpdateFxAccountFollowStatus(info *FxAccountFollow) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Cols("status", "updated_at").Update(info, &FxAccountFollow{WXAccount: info.WXAccount, OpenId: info.OpenId})
	return err
}

func GetFxAccountFollow(info *FxAccountFollow) error {
	has, err := x.Where("wx_account = ?", info.WXAccount).And("open_id = ?", info.OpenId).Get(info)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("cannot find fx account from wx_account[%s] open_id[%s]", info.WXAccount, info.OpenId)
	}
	return nil
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
