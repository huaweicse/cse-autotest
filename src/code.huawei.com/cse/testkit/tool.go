package testkit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"math/rand"

	"bytes"

	providerRestApi "code.huawei.com/cse/api/provider/rest"
	"code.huawei.com/cse/common"
	"code.huawei.com/cse/model"
	"code.huawei.com/cse/util"
	. "github.com/onsi/ginkgo"
)

type Param struct {
	Data          map[string]interface{} `json:"data"`
	Keys          []string               `json:"key"`
	DimensionInfo string                 `json:"dimensionInfo"`
	Type          string                 `json:"type"`
}

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
var instanceNames = make(map[string][]string)

func SDKATContext(body func(inputConsumerAddr, inputProviderName, inputProtocol, dimensionInfo, instanceName string, instanceLength int)) {
	Init()

	for consumerType, consumerAddr := range consumers {
		for providerType, providerName := range providers {
			protos := protocols[consumerType]
			if len(protos) == 0 {
				protos = []string{common.ProtocolRest}
			}
			// dimensionInfo
			var dimensionInfo string
			if providerType == common.PlatformGoSDK {
				dimensionInfo = fmt.Sprintf("%s@default#%s", common.ConsumerGoSDK, common.Version30)
			} else if providerType == common.PlatformMesher {
				dimensionInfo = fmt.Sprintf("%s@default#%s", common.ConsumerMesher, common.Version30)
			}
			for _, p := range protos {
				k := fmt.Sprintf("%s|%s|%s|%s|%s", consumerType, consumerAddr, providerType, providerName, p)
				insName, instanceLength := getInstance(k, consumerAddr, providerName, p)
				log.Println(insName, instanceLength)
				Context(fmt.Sprintf("consumer addr: %s, %s -> %s, protocol: %s", consumerAddr, consumerType, providerType, p), func() {
					body(consumerAddr, providerName, p, dimensionInfo, insName, instanceLength)
				})
			}
		}
	}
}
func getInstance(key, consumerAddr, providerName, protocol string) (string, int) {
	if ins, ok := instanceNames[key]; ok {
		lengthIns := len(ins)
		return ins[rand.Intn(lengthIns)%lengthIns], lengthIns
	}
	testUri := fmt.Sprintf("http://%s%s?%s", consumerAddr, providerRestApi.Svc,
		util.FncodeParams([]util.URLParameter{
			{common.ParamProvider: providerName},
			{common.ParamProtocol: protocol},
			{common.ParamTimes: common.CallTimes20Str},
		}))
	log.Println("url=" + testUri)
	instanceList := GetResponseInstanceAliasList(testUri)
	newInsanceList := []string{}
	mTemp := make(map[string]int)
	lengthInstance := len(instanceList)
	for i := 0; i < lengthInstance; i++ {
		if _, ok := mTemp[instanceList[i]]; ok {
			continue
		}
		mTemp[instanceList[i]] = 0
		newInsanceList = append(newInsanceList, instanceList[i])
	}
	if len(newInsanceList) < 2 {
		panic(newInsanceList)
	}
	instanceNames[key] = newInsanceList
	lengthIns := len(newInsanceList)
	return newInsanceList[rand.Intn(lengthIns)%lengthIns], lengthIns
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

func GetResponseInstanceAliasList(u string) []string {
	c := &model.InstanceInfoResponse{}
	resp, err := http.Get(u)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic("status is not ok: " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, c)
	if err != nil {
		panic(err)
	}

	l := make([]string, 0)
	for _, r := range c.Result {
		if r.Provider == nil {
			continue
		}
		if r.Provider.InstanceAlias != "" {
			l = append(l, r.Provider.InstanceAlias)
		} else {
			l = append(l, r.Provider.InstanceName)
		}
	}
	return l
}

func Callcc(url, cctype, dimensionInfo string, items map[string]interface{}, keys []string) {
	p := Param{
		Data:          items,
		Keys:          keys,
		Type:          cctype,
		DimensionInfo: dimensionInfo,
	}
	body, _ := json.Marshal(p)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	if resp == nil {
		panic(resp)
	}
	if resp.StatusCode != http.StatusOK {
		panic("status is not ok: " + resp.Status)
	}

	// read data for resp
	m := make(map[string]interface{})
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if len(body) == 0 {
		panic("body is empty")
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		panic(err)
	}

	if m["Result"] == "" {
		panic("result is nil")
	}
	if m["Result"] != "Success" {
		panic(m["Result"])
	}
}
