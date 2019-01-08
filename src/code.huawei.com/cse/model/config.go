package model

/****************配置项****************/
type GlobalCfg struct {
	AppID           string    `yaml:"appid"`
	Name            string    `yaml:"name"`
	Version         string    `yaml:"version"`
	Region          string    `yaml:"region"`
	ProjectId       string    `yaml:"projectId"`
	Http            Protocol  `yaml:"http"`
	ConfigCenterUrl string    `yaml:"configCenterUrl"`
	Cse             CseStruct `yaml:"cse"`
}

// Protocol protocol structure
type Protocol struct {
	Listen string `yaml:"listenAddress"`
}

// CredentialStruct aksk信息
type CredentialStruct struct {
	AccessKey        string        `yaml:"accessKey"`
	SecretKey        string        `yaml:"secretKey"`
	AkskCustomCipher string        `yaml:"akskCustomCipher"`
	Project          ProjectStruct `yaml:"project"`
}

type ProjectStruct struct {
	AkskCustomCipher string `yaml:"akskCustomCipher"`
}

//CseStruct 设置注册中心SC的地址，要开哪些传输协议， 调用链信息等
type CseStruct struct {
	Credentials CredentialStruct `yaml:"credentials"`
}
