package loadbalance_test

import (
	providerRestApi "code.huawei.com/cse/api/provider/rest"
	"code.huawei.com/cse/common"
	"code.huawei.com/cse/model"
	"code.huawei.com/cse/testkit"
	"code.huawei.com/cse/util"
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"log"
	"net/http"
)

func GetResponceInstanceAliasList(u string) []string {
	c := &model.InstanceInfoResponse{}
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
		if r.Provider.InstanceAlias != "" {
			l = append(l, r.Provider.InstanceAlias)
		} else {
			l = append(l, r.Provider.InstanceName)
		}
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

func test(consumerAddr, providerName, protocol string) {
	instansNames := make(map[string]bool)
	testUri := fmt.Sprintf("http://%s%s?%s", consumerAddr, providerRestApi.Svc,
		util.FncodeParams([]util.URLParameter{
			{common.ParamProvider: providerName},
			{common.ParamProtocol: protocol},
			{common.ParamTimes: common.CallTimes20Str},
		}))
	It("Instance num should more than or equal to 2", func() {
		instansNameList := GetResponceInstanceAliasList(testUri)
		for _, v := range instansNameList {
			instansNames[v] = true
		}
		log.Println("---check instance num, instance list:", instansNameList)
		log.Println("---check instance num, instance num:", len(instansNames))
		Expect(len(instansNames) >= 2).To(BeTrue())
	})
	It("Instance name list should be round robin", func() {
		nameList := GetResponceInstanceAliasList(testUri)
		log.Println("---instance list:", nameList)
		IsRoundRoubin(len(instansNames), nameList...)
	})
}

var _ = Describe("Load balance", func() {
	testkit.SDKATContext(test)
})
