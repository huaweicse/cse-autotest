package config

import (
	"code.huawei.com/cse/common"
	"io/ioutil"

	"code.huawei.com/cse/model"
	"code.huawei.com/cse/util"

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
	"os"
)

var (
	GlobalDef  *model.GlobalCfg
	DefaultDef = &model.GlobalCfg{
		AppID:   "sdkat",
		Name:    "sdkat_provider_mesher",
		Version: "1.0",
		Http: model.Protocol{
			Listen: "127.0.0.1:9090",
		},
	}
)

func Init() error {
	def := &model.GlobalCfg{}
	globalDefFilePath := util.GetDefinition()
	globalDefFileContent, err := ioutil.ReadFile(globalDefFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			glog.Infof("config file %s not exist, use default conf", globalDefFilePath)
			GlobalDef = DefaultDef
			return nil
		}
		return err
	}
	err = yaml.Unmarshal(globalDefFileContent, def)
	if err != nil {
		return err
	}
	if n := os.Getenv(common.EnvServiceName); n != "" {
		def.Name = n
	}
	if n := os.Getenv(common.EnvVersion); n != "" {
		def.Version = n
	}
	if n := os.Getenv(common.EnvAppId); n != "" {
		def.Name = n
	}
	GlobalDef = CheckGlobalDef(def)
	return nil
}

func CheckGlobalDef(oldDef *model.GlobalCfg) *model.GlobalCfg {
	if oldDef == nil {
		return DefaultDef
	}

	newDef := *DefaultDef
	if oldDef.AppID != "" {
		newDef.AppID = oldDef.AppID
	}
	if oldDef.Name != "" {
		newDef.Name = oldDef.Name
	}
	if oldDef.Version != "" {
		newDef.Version = oldDef.Version
	}
	if oldDef.Http.Listen != "" {
		newDef.Http.Listen = oldDef.Http.Listen
	}
	return &newDef
}
