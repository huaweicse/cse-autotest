package loadbalance_test

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	consumerRestApi "code.huawei.com/cse/api/consumer/rest"
	"code.huawei.com/cse/common"
	"code.huawei.com/cse/model"
	"code.huawei.com/cse/util"
)

var (
	GosdkConsumerAddr  string = os.Getenv(common.EnvConsumerGoSdkAddr)
	MesherConsumerAddr string = os.Getenv(common.EnvConsumerMesherAddr)
)

func GetResponceInstanceAliasList(u string) []string {
	c := &model.InvokeResponce{}
	resp, err := http.Get(u)
	Expect(err).To(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).To(BeNil())

	err = json.Unmarshal(body, c)
	Expect(err).To(BeNil())

	l := make([]string, 0)
	for _, r := range c.Result {
		Expect(r.Provider).NotTo(BeNil())
		if r.Provider == nil {
			continue
		}
		l = append(l, r.Provider.InstanceAlias)
	}
	return l
}

func IsRoundRoubin(cycle int, n ...string) {
	Expect((len(n) - cycle) > 0).To(BeTrue())
	for i := 0; i < len(n)-cycle; i++ {
		Expect(n[i]).NotTo(BeEmpty())
		Expect(n[i+cycle]).To(Equal(n[i]))
		if cycle != 1 {
			Expect(n[i+1]).NotTo(Equal(n[i]))
		}
	}
}

var _ = Describe("Loadbalance", func() {
	if GosdkConsumerAddr == "" {
		GosdkConsumerAddr = common.AddrDefaultGoSDKConsumer
	}
	if MesherConsumerAddr == "" {
		MesherConsumerAddr = common.AddrDefaultMesherConsumer
	}
	gosdk2GosdkUrl := fmt.Sprintf("%s%s?%s", GosdkConsumerAddr, consumerRestApi.TestCommunication,
		util.FncodeParams([]util.URLParameter{
			{common.ParamProvider: common.ProviderGoSDK},
			{common.ParamProtocol: common.ProtocolRest},
			{common.ParamTimes: common.CallTimes20Str},
		}))
	gosdk2MesherUrl := fmt.Sprintf("%s%s?%s", GosdkConsumerAddr, consumerRestApi.TestCommunication,
		util.FncodeParams([]util.URLParameter{
			{common.ParamProvider: common.ProviderMesher},
			{common.ParamProtocol: common.ProtocolRest},
			{common.ParamTimes: common.CallTimes20Str},
		}))
	mesher2GosdkUrl := fmt.Sprintf("%s%s?%s", MesherConsumerAddr, consumerRestApi.TestCommunication,
		util.FncodeParams([]util.URLParameter{
			{common.ParamProvider: common.ProviderGoSDK},
			{common.ParamProtocol: common.ProtocolRest},
			{common.ParamTimes: common.CallTimes20Str},
		}))
	mesher2MesherUrl := fmt.Sprintf("%s%s?%s", MesherConsumerAddr, consumerRestApi.TestCommunication,
		util.FncodeParams([]util.URLParameter{
			{common.ParamProvider: common.ProviderMesher},
			{common.ParamProtocol: common.ProtocolRest},
			{common.ParamTimes: common.CallTimes20Str},
		}))

	Context(fmt.Sprintf("Default, received instanceAlias list should be [%s]", "RoundRobin"), func() {
		It(common.GoSdk2GoSdk, func() {
			nameList := GetResponceInstanceAliasList(gosdk2GosdkUrl)
			log.Printf("%s, instanceAlias list: %s", common.GoSdk2GoSdk, nameList)
			IsRoundRoubin(common.InstanceNumDefault, nameList...)
		})
		It(common.GoSdk2Mesher, func() {
			nameList := GetResponceInstanceAliasList(gosdk2MesherUrl)
			log.Printf("%s, instanceAlias list: %s", common.GoSdk2Mesher, nameList)
			IsRoundRoubin(common.InstanceNumDefault, nameList...)
		})
		It(common.Mesher2GoSdk, func() {
			nameList := GetResponceInstanceAliasList(mesher2GosdkUrl)
			log.Printf("%s, instanceAlias list: %s", common.Mesher2GoSdk, nameList)
			IsRoundRoubin(common.InstanceNumDefault, nameList...)
		})
		It(common.Mesher2Mesher, func() {
			nameList := GetResponceInstanceAliasList(mesher2MesherUrl)
			log.Printf("%s, instanceAlias list: %s", common.Mesher2Mesher, nameList)
			IsRoundRoubin(common.InstanceNumDefault, nameList...)
		})
	})
})
