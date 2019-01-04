package communication_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	consumerRestApi "code.huawei.com/cse/api/consumer/rest"
	"code.huawei.com/cse/common"
	"code.huawei.com/cse/model"
	"code.huawei.com/cse/util"
	"encoding/json"
)

var (
	GosdkConsumerAddr  string = os.Getenv(common.EnvConsumerGoSdkAddr)
	MesherConsumerAddr string = os.Getenv(common.EnvConsumerMesherAddr)
)

func CheckReponseProvider(u, key string) {
	c := &model.InvokeResponce{}
	resp, err := http.Get(u)
	Expect(err).To(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).To(BeNil())

	err = json.Unmarshal(body, c)
	Expect(err).To(BeNil())

	for _, r := range c.Result {
		Expect(r.StatusCode).To(Equal(http.StatusOK))
		Expect(r.Error).To(BeEmpty())
		k := strings.Join([]string{r.Provider.MicroService.ServiceName, r.Provider.MicroService.Version, r.Provider.MicroService.Application}, "/")
		Expect(k).To(Equal(key))
	}
}

var _ = Describe("Communication", func() {
	if GosdkConsumerAddr == "" {
		GosdkConsumerAddr = common.AddrDefaultGoSDKConsumer
	}
	if MesherConsumerAddr == "" {
		MesherConsumerAddr = common.AddrDefaultMesherConsumer
	}
	gosdkProviderKey := strings.Join([]string{common.ProviderGoSDK, common.Version30, common.AppSDKAT}, "/")
	mesherProviderKey := strings.Join([]string{common.ProviderMesher, common.Version30, common.AppSDKAT}, "/")
	gosdk2GosdkUrl := fmt.Sprintf("%s%s?%s", GosdkConsumerAddr, consumerRestApi.TestCommunication,
		util.FncodeParams([]util.URLParameter{
			{common.ParamProvider: common.ProviderGoSDK},
			{common.ParamProtocol: common.ProtocolRest},
		}))
	gosdk2MesherUrl := fmt.Sprintf("%s%s?%s", GosdkConsumerAddr, consumerRestApi.TestCommunication,
		util.FncodeParams([]util.URLParameter{
			{common.ParamProvider: common.ProviderMesher},
			{common.ParamProtocol: common.ProtocolRest},
		}))
	mesher2GosdkUrl := fmt.Sprintf("%s%s?%s", MesherConsumerAddr, consumerRestApi.TestCommunication,
		util.FncodeParams([]util.URLParameter{
			{common.ParamProvider: common.ProviderGoSDK},
			{common.ParamProtocol: common.ProtocolRest},
		}))
	mesher2MesherUrl := fmt.Sprintf("%s%s?%s", MesherConsumerAddr, consumerRestApi.TestCommunication,
		util.FncodeParams([]util.URLParameter{
			{common.ParamProvider: common.ProviderMesher},
			{common.ParamProtocol: common.ProtocolRest},
		}))

	Context("Should get responce from target provider", func() {
		It(common.GoSdk2GoSdk, func() {
			CheckReponseProvider(gosdk2GosdkUrl, gosdkProviderKey)
		})
		It(common.GoSdk2Mesher, func() {
			CheckReponseProvider(gosdk2MesherUrl, mesherProviderKey)
		})
		It(common.Mesher2GoSdk, func() {
			CheckReponseProvider(mesher2GosdkUrl, gosdkProviderKey)
		})
		It(common.Mesher2Mesher, func() {
			CheckReponseProvider(mesher2MesherUrl, mesherProviderKey)
		})
	})
})
