package config

import (
	"bytes"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"code.huawei.com/cse/common"
	"code.huawei.com/cse/model"
	"code.huawei.com/cse/util"
	"github.com/go-mesh/openlogging"
)

const _DefaultRegion string = "cn-north-5"

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
	err := LoadConfig()
	if err != nil {
		return err
	}
	def := GlobalDef
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
func LoadConfig() error {
	cfg := &model.GlobalCfg{}
	buf := bytes.NewBuffer([]byte{})

	err := loadConfig(util.GetAuth(), buf)
	if err != nil {
		openlogging.GetLogger().Errorf("read auth.yaml file failed ,err : %v", err)
	}
	err = loadConfig(util.GetDefinition(), buf)
	if err != nil {
		openlogging.GetLogger().Errorf("read microservice.yaml file failed , err : %v", err)
	}
	if buf == nil || len(buf.Bytes()) == 0 {
		cfg = DefaultDef
	} else {
		err = yaml.Unmarshal(buf.Bytes(), cfg)
		if err != nil {
			openlogging.GetLogger().Errorf("unmarshal to cache failed please check yaml file, err : %v", err)
			return err
		}
	}

	GlobalDef = cfg
	return nil
}

// loafConfig load config
func loadConfig(filePath string, buf *bytes.Buffer) error {
	globalDefFileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		openlogging.GetLogger().Errorf("read config file failed:%s", filePath)
		return err
	}
	globalDefFileContent = bytes.Replace(globalDefFileContent, []byte("---"), []byte(""), -1)
	buf.Write(globalDefFileContent)
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

	envProject := getProject()

	if envProject != "" {
		newDef.ProjectId = envProject
	} else if oldDef.ProjectId != "" {
		newDef.ProjectId = oldDef.ProjectId
	} else {
		newDef.ProjectId = "default"
	}
	if oldDef.Region != "" {
		newDef.Region = oldDef.Region
	} else {
		newDef.Region = _DefaultRegion
	}
	if oldDef.ConfigCenterUrl != "" {
		newDef.ConfigCenterUrl = oldDef.ConfigCenterUrl
	} else {
		newDef.ConfigCenterUrl = "http://127.0.0.1:3000"
	}
	if oldDef.Cse.Credentials.SecretKey != "" {
		newDef.Cse.Credentials.SecretKey = oldDef.Cse.Credentials.SecretKey
	}
	if oldDef.Cse.Credentials.AccessKey != "" {
		newDef.Cse.Credentials.AccessKey = oldDef.Cse.Credentials.AccessKey
	}
	if oldDef.Cse.Credentials.AkskCustomCipher != "" {
		newDef.Cse.Credentials.AkskCustomCipher = oldDef.Cse.Credentials.AkskCustomCipher
	}
	if oldDef.Cse.Credentials.Project.AkskCustomCipher != "" {
		newDef.Cse.Credentials.Project.AkskCustomCipher = oldDef.Cse.Credentials.Project.AkskCustomCipher
	}

	return &newDef
}

func getProject() string {
	s, ok := os.LookupEnv(util.PAAS_PROJECT_NAME)
	if ok {
		return s
	}
	return ""
}
