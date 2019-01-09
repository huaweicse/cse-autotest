package util

import (
	"os"
	"path/filepath"
	"sync"
)

const (
	//ChassisConfDir is constant of type string
	ChassisConfDir = "CHASSIS_CONF_DIR"
	//ChassisHome is constant of type string
	ChassisHome       = "CHASSIS_HOME"
	PAAS_PROJECT_NAME = "PAAS_PROJECT_NAME"
)

const Definition = "microservice.yaml"

const Auth string = "auth.yaml"

var configDir string
var homeDir string
var once sync.Once

func initDir() {
	if h := os.Getenv(ChassisHome); h != "" {
		homeDir = h
	} else {
		wd, err := GetWorkDir()
		if err != nil {
			panic(err)
		}
		homeDir = wd
	}

	// set conf dir, CHASSIS_CONF_DIR has highest priority
	if confDir := os.Getenv(ChassisConfDir); confDir != "" {
		configDir = confDir
	} else {
		// CHASSIS_HOME has second most high priority
		configDir = filepath.Join(homeDir, "conf")
	}
}

//GetWorkDir is a function used to get the working directory
func GetWorkDir() (string, error) {
	wd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return wd, nil
}

//ChassisHomeDir is function used to get the home directory of chassis
func ChassisHomeDir() string {
	once.Do(initDir)
	return homeDir
}

//GetConfDir is a function used to get the configuration directory
func GetConfDir() string {
	once.Do(initDir)
	return configDir
}

//GetDefinition is a function used to join .yaml file name with configuration path
func GetDefinition() string {
	return filepath.Join(GetConfDir(), Definition)
}

func GetAuth() string {
	return filepath.Join(GetConfDir(), Auth)
}
