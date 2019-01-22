package loadbalance_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"bytes"

	"time"

	providerRestApi "code.huawei.com/cse/api/provider/rest"
	"code.huawei.com/cse/common"
	"code.huawei.com/cse/testkit"
	"code.huawei.com/cse/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func IsRoundRoubin(cycle int, n ...string) bool {
	Expect((len(n) - cycle) > 0).To(BeTrue())
	Expect(cycle > 1).To(BeTrue())
	length := len(n) - cycle
	for i := 0; i < length; i++ {
		Expect(n[i]).NotTo(BeEmpty())
		if !Expect(n[i+cycle]).To(Equal(n[i])) {
			return false
		}
	}
	return true
}
func isStrategySessionStickiness(cycle int, n ...string) bool {
	Expect((len(n) - cycle) > 0).To(BeTrue())
	Expect(cycle > 1).To(BeTrue())
	length := len(n) - 1
	for i := 0; i < length; i++ {
		Expect(n[i]).NotTo(BeEmpty())
		if !Expect(n[i+1]).To(Equal(n[i])) {
			return false
		}
	}
	return true
}
func isWeightedResponse(cycle int, n ...string) bool {
	Expect((len(n) - cycle) > 0).To(BeTrue())
	Expect(cycle > 1).To(BeTrue())
	l := len(n) - 1
	ma := make(map[string]int)
	for i := 0; i < l; i++ {
		Expect(n[i]).NotTo(BeEmpty())
		if _, ok := ma[n[i]]; !ok {
			ma[n[i]] = 1
		}
		ma[n[i]] += 1
	}
	for _, v := range ma {
		if v/len(n)*100 >= 80 {
			return true
		}
	}

	return false
}
func isRand(cycle int, n ...string) {
	Expect((len(n) - cycle) > 0).To(BeTrue())
	Expect(cycle > 1).To(BeTrue())
	Expect(isStrategySessionStickiness(cycle, n...)).To(BeFalse())
	Expect(IsRoundRoubin(cycle, n...)).To(BeFalse())
	Expect(isWeightedResponse(cycle, n...)).To(BeFalse())
}
func test(consumerAddr, providerName, protocol, dimensionInfo string) {
	instansNames := make(map[string]bool)
	var instanceName string
	testUri := fmt.Sprintf("http://%s%s?%s", consumerAddr, providerRestApi.Svc,
		util.FncodeParams([]util.URLParameter{
			{common.ParamProvider: providerName},
			{common.ParamProtocol: protocol},
			{common.ParamTimes: common.CallTimes20Str},
		}))
	It("Instance num should more than or equal to 2", func() {
		instansNameList := testkit.GetResponceInstanceAliasList(testUri)
		for k, v := range instansNameList {
			if k == 0 || instanceName == "" {
				instanceName = v
			}
			if _, ok := instansNames[v]; !ok {
				instansNames[v] = true
			}
		}
		log.Println("---check instance num, instance list:", instansNameList)
		log.Println("---check instance num, instance num:", len(instansNames))
		Expect(len(instansNames) >= 2).To(BeTrue())
	})
	instansNum := len(instansNames)

	// RoundRobin test
	testLB("Instance name list should be round robin", dimensionInfo,
		consumerAddr, testUri, "RoundRobin", "", providerName, instanceName, instansNum)
	// session stickiness test
	testLB("Instance name list should be session stickiness", dimensionInfo, consumerAddr,
		testUri, "SessionStickiness", "10", providerName, instanceName, instansNum)

	// WeightedResponse test
	testLB("Instance name list should be WeightedResponse", dimensionInfo, consumerAddr,
		testUri, "WeightedResponse", "", providerName, instanceName, instansNum)

	// Rand test
	testLB("Instance name list should be Random", dimensionInfo, consumerAddr, testUri,
		"Random", "", providerName, instanceName, instansNum)
}
func testLB(text, dimensionInfo, consumerAddr, testUri,
	t, st, providerName, instanceName string, instansNameNums int) {
	m := map[string]interface{}{
		"cse.loadbalance.strategy.name":                                 t,
		"cse.loadbalance.SessionStickinessRule.sessionTimeoutInSeconds": st,
	}
	callcc(fmt.Sprintf("http://%s%s", consumerAddr, providerRestApi.ConfigCenterAdd),
		"add", dimensionInfo, m, nil)
	time.Sleep(3 * time.Second)
	if t == "WeightedResponse" {
		curl := fmt.Sprintf("http://%s%s?%s", consumerAddr, providerRestApi.Svc,
			util.FncodeParams([]util.URLParameter{
				{common.ParamProvider: providerName},
				{common.ParamProtocol: fmt.Sprintf("delayInstance/%s/1000", instanceName)},
				{common.ParamTimes: common.CallTimes10Str},
			}))
		testkit.GetResponceInstanceAliasList(curl)
		time.Sleep(30 * time.Second)
	}
	It(text, func() {
		nameList := testkit.GetResponceInstanceAliasList(testUri)

		log.Println(fmt.Sprintf("type:%s,---instance list:%v", t, nameList))
		switch t {
		case "RoundRobin":
			IsRoundRoubin(instansNameNums, nameList...)
		case "SessionStickiness":
			isStrategySessionStickiness(instansNameNums, nameList...)
		case "WeightedResponse":
			isWeightedResponse(instansNameNums, nameList...)
		case "Random":
			isRand(instansNameNums, nameList...)
		default:
		}
	})

}

var _ = Describe("Load balance", func() {
	testkit.SDKATContext(test)
})

func callcc(url, cctype, dimensionInfo string, items map[string]interface{}, keys []string) {
	p := Param{
		Data:          items,
		Keys:          keys,
		Type:          cctype,
		DimensionInfo: dimensionInfo,
	}
	body, _ := json.Marshal(p)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	Expect(err).To(BeNil())
	Expect(resp).NotTo(BeNil())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	// read data for resp
	m := make(map[string]interface{})
	body, err = ioutil.ReadAll(resp.Body)
	Expect(err).To(BeNil())
	Expect(len(body)).NotTo(Equal(0))
	err = json.Unmarshal(body, &m)
	Expect(err).To(BeNil())

	Expect(m["Result"]).NotTo(BeEmpty())
	Expect(m["Result"]).To(Equal("Success"))
}

type Param struct {
	Data          map[string]interface{} `json:"data"`
	Keys          []string               `json:"key"`
	DimensionInfo string                 `json:"dimensionInfo"`
	Type          string                 `json:"type"`
}
