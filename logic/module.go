package logic

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/real-fx/config"
	"github.com/reechou/real-fx/logic/controller"
	"github.com/reechou/real-fx/logic/fx"
	"github.com/reechou/real-fx/logic/models"
	"github.com/reechou/real-fx/logic/wx"
	"github.com/reechou/real-fx/server"
	"github.com/reechou/real-fx/utils"
)

type ModuleLogic struct {
	cfg    *config.Config
	daemon *controller.Daemon
}

func NewModuleLogic(cfg *config.Config) *ModuleLogic {
	d := controller.NewDaemon(cfg)
	models.InitDB(cfg)
	return &ModuleLogic{
		cfg:    cfg,
		daemon: d,
	}
}

func (ml *ModuleLogic) InitRouter(s *server.Server) {
	s.InitRouter(utils.IsDebugEnabled(),
		wx.NewRouter(ml.daemon),
		fx.NewRouter(ml.daemon))
}

func (ml *ModuleLogic) Shutdown(timeout time.Duration) {
	ch := make(chan struct{})
	go func() {
		ml.close()
		close(ch)
	}()
	select {
	case <-ch:
		logrus.Debug("Clean logic shutdown succeeded")
	case <-time.After(timeout * time.Second):
		logrus.Error("Force shutdown server.")
	}
}

func (ml *ModuleLogic) close() {

}
