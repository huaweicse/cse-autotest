package testkit

import (
	"code.huawei.com/cse/common"
	"fmt"
	. "github.com/onsi/ginkgo"
	"log"
	"os"
	"strings"
	"sync"
)

// k: consumer type, v: protocol
var protocols = map[string][]string{
	common.PlatformGoSDK:  {common.ProtocolRest},
	common.PlatformMesher: {common.ProtocolRest},
}

// k: consumer type, v:consumer addr
var consumers map[string]string
var defaultConsumers = map[string]string{
	common.PlatformGoSDK:  "127.0.0.1:8000",
	common.PlatformMesher: "127.0.0.1:8080",
}

//provider type
var providers map[string]string
var defaultProviders = map[string]string{
	common.PlatformGoSDK:  common.ProviderGoSDK,
	common.PlatformMesher: common.ProviderMesher,
}
var once sync.Once

func SDKATContext(body func(inputConsumerAddr, inputProviderName, inputProtocol string)) {
	Init()
	for consumerType, consumerAddr := range consumers {
		for providerType, providerName := range providers {
			protos := protocols[consumerType]
			if len(protos) == 0 {
				protos = []string{common.ProtocolRest}
			}
			for _, p := range protos {
				Context(fmt.Sprintf("consumer addr: %s, %s -> %s, protocol: %s", consumerAddr, consumerType, providerType, p), func() {
					body(consumerAddr, providerName, p)
				})
			}
		}
	}
}

func initConsumerAndProvider() {
	consumers = make(map[string]string)
	//如果设置了consumer地址，则使用设置的consumer地址
	isConsumerSet := false
	if v := os.Getenv(common.EnvConsumerGoSdkAddr); v != "" {
		consumers[common.PlatformGoSDK] = v
		isConsumerSet = true
	}
	if v := os.Getenv(common.EnvConsumerMesherAddr); v != "" {
		consumers[common.PlatformMesher] = v
		isConsumerSet = true
	}
	//未设置consumer，用默认地址
	if !isConsumerSet {
		consumers = defaultConsumers
	}

	providers = make(map[string]string)
	// format "gosdk,mesher"
	if v := os.Getenv("SDKAT_PROVIDER"); v != "" {
		//设置了provider
		parts := strings.Split(v, ",")
		for _, p := range parts {
			providers[p] = defaultProviders[p]
		}
	} else {
		//未设置provider，则与consumer信息对其
		for p := range consumers {
			providers[p] = defaultProviders[p]
		}
	}
	log.Println("---consumers: ", consumers)
	log.Println("---providers: ", providers)
	log.Println("---protocols: ", protocols)
}

func Init() {
	once.Do(initConsumerAndProvider)
}
