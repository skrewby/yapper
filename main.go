package main

import (
	"github.com/skrewby/yapper/controller"
	"github.com/skrewby/yapper/utils"
)

func main() {
	env := utils.GetEnvironmentVariables()
	c := controller.NewController(env)
	c.StartServer()
}
