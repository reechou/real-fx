package main

import (
	"github.com/reechou/real-fx/config"
	"github.com/reechou/real-fx/servermain"
)

func main() {
	servermain.NewMain(config.NewConfig()).Run()
}
