package models

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

type FxAccountLimit struct {
	ID              int64  `xorm:"pk autoincr"`
	UnionId         string `xorm:"not null default '' varchar(128) unique"`
	CheckIn         int64  `xorm:"not null default 0 int"`
	WithdrawalTimes int64  `xorm:"not null default 0 int"`
	CreatedAt       int64  `xorm:"not null default 0 int"`
	UpdatedAt       int64  `xorm:"not null default 0 int"`
}

func CreateFxAccountLimt(info *FxAccountLimit) error {
	if info.UnionId == "" {
		return fmt.Errorf("fx account limit union_id[%s] cannot be nil.", info.UnionId)
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now
	_, err := x.Insert(info)
	if err != nil {
		logrus.Errorf("create fx account limit error: %v", err)
		return err
	}
	logrus.Infof("fx account limit union_id[%s] create success.", info.UnionId)

	return nil
}
