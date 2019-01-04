package service

import (
	"github.com/go-chassis/go-chassis/core"
)

var RestInvoker *core.RestInvoker

func init() {
	RestInvoker = core.NewRestInvoker()
}
