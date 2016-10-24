package logic

import (
	"time"

	"github.com/reechou/real-fx/config"
	"github.com/reechou/real-fx/server"
)

type LogicThinking interface {
	InitRouter(s *server.Server)
	Shutdown(timeout time.Duration)
}

func NewLogic(cfg *config.Config) (LogicThinking, error) {
	l := NewModuleLogic(cfg)
	return l, nil
}
