package config

import (
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-mesh/openlogging"
	"fmt"
)

var (
	GlobalDef *model.GlobalCfg
)

func Init() error {
	GlobalDef = config.GlobalDefinition
	openlogging.Info(fmt.Sprintf("%v", GlobalDef))
	if GlobalDef == nil {
		err := config.Init()
		if err != nil {
			panic(err)
		}
	}
	return nil
}
