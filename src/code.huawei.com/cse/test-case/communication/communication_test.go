package communication_test

import (
	"code.huawei.com/cse/testkit"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"encoding/json"

	providerRestApi "code.huawei.com/cse/api/provider/rest"
	"code.huawei.com/cse/common"
	"code.huawei.com/cse/model"
	"code.huawei.com/cse/util"
)

func CheckReponseProvider(u, key string) {
	c := &model.InstanceInfoResponse{}
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

func test(consumerAddr, providerName, protocol, dimensionInfo string) {
	providerKey := strings.Join([]string{providerName, common.Version30, common.AppSDKAT}, "/")
	testUri := fmt.Sprintf("http://%s%s?%s", consumerAddr, providerRestApi.Svc,
		util.FncodeParams([]util.URLParameter{
			{common.ParamProvider: providerName},
			{common.ParamProtocol: protocol},
		}))
	It("Response should contains expected provider info", func() {
		CheckReponseProvider(testUri, providerKey)
	})
}

var _ = Describe("Communication", func() {
	testkit.SDKATContext(test)
})
