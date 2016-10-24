package main

import (
	"github.com/reechou/real-fx/config"
	"github.com/reechou/real-fx/logic/tools/order_check/controller"
)

func main() {
	controller.NewOrderCheck(config.NewConfig()).Run()
}
