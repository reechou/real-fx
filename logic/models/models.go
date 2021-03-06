package models

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/reechou/real-fx/config"
)

var x *xorm.Engine

func InitDB(cfg *config.Config) {
	var err error
	x, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		cfg.DBInfo.User,
		cfg.DBInfo.Pass,
		cfg.DBInfo.Host,
		cfg.DBInfo.DBName))
	if err != nil {
		logrus.Fatalf("Fail to init new engine: %v", err)
	}
	x.SetLogger(nil)
	x.SetMapper(core.GonicMapper{})
	x.TZLocation, _ = time.LoadLocation("Asia/Shanghai")

	if err = x.Sync2(new(FxAccount),
		new(FxAccountMonthAchievement),
		new(FxAccountFollow),
		new(FxTeam),
		new(FxTeamMember),
		new(FxOrder),
		new(FxOrderSettlementRecord),
		new(FxOrderWaitSettlementRecord),
		new(WithdrawalRecord),
		new(FxAccountHistory),
		new(WithdrawalRecordError)); err != nil {
		logrus.Fatalf("Fail to sync database: %v", err)
	}
}
