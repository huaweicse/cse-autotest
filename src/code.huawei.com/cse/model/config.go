package model

/****************配置项****************/
type GlobalCfg struct {
	AppID   string   `yaml:"appid"`
	Name    string   `yaml:"name"`
	Version string   `yaml:"version"`
	Http    Protocol `yaml:"http"`
}

// Protocol protocol structure
type Protocol struct {
	Listen string `yaml:"listenAddress"`
}
