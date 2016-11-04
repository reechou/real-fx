package models

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/reechou/real-fx/logic/tools/order_check/config"
)

var x *xorm.Engine

func InitDB(cfg *config.Config) {
	var err error
	x, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		cfg.TaobaoDBInfo.User,
		cfg.TaobaoDBInfo.Pass,
		cfg.TaobaoDBInfo.Host,
		cfg.TaobaoDBInfo.DBName))
	if err != nil {
		logrus.Fatalf("Fail to init new engine: %v", err)
	}
	x.SetLogger(nil)
	x.SetMapper(core.GonicMapper{})
	x.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
}
