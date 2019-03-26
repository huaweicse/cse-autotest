package loadbalance_test

import (
	"fmt"
	"log"

	"time"

	providerRestApi "code.huawei.com/cse/api/provider/rest"
	"code.huawei.com/cse/common"
	"code.huawei.com/cse/testkit"
	"code.huawei.com/cse/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func IsRoundRobin(cycle int, n ...string) bool {
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
	Expect(IsRoundRobin(cycle, n...)).To(BeFalse())
	Expect(isWeightedResponse(cycle, n...)).To(BeFalse())
}
func test(consumerAddr, providerName, protocol, dimensionInfo, instanceName string, instancesLength int) {

	testUri := fmt.Sprintf("http://%s%s?%s", consumerAddr, providerRestApi.Svc,
		util.FncodeParams([]util.URLParameter{
			{common.ParamProvider: providerName},
			{common.ParamProtocol: protocol},
			{common.ParamTimes: common.CallTimes20Str},
		}))

	// RoundRobin test
	testLB("Instance name list should be round robin", dimensionInfo,
		consumerAddr, testUri, "RoundRobin", "", providerName, instanceName, instancesLength)
	// session stickiness test
	testLB("Instance name list should be session stickiness", dimensionInfo, consumerAddr,
		testUri, "SessionStickiness", "10", providerName, instanceName, instancesLength)

	// WeightedResponse test
	testLB("Instance name list should be WeightedResponse", dimensionInfo, consumerAddr,
		testUri, "WeightedResponse", "", providerName, instanceName, instancesLength)

	// Rand test
	testLB("Instance name list should be Random", dimensionInfo, consumerAddr, testUri,
		"Random", "", providerName, instanceName, instancesLength)
}
func testLB(text, dimensionInfo, consumerAddr, testUri,
	t, st, providerName, instanceName string, instancesLength int) {
	m := map[string]interface{}{
		"cse.loadbalance.strategy.name":                                 t,
		"cse.loadbalance.SessionStickinessRule.sessionTimeoutInSeconds": st,
	}
	testkit.Callcc(fmt.Sprintf("http://%s%s", consumerAddr, providerRestApi.ConfigCenterAdd),
		"add", dimensionInfo, m, nil)
	time.Sleep(3 * time.Second)
	if t == "WeightedResponse" {
		curl := fmt.Sprintf("http://%s%s?%s", consumerAddr, providerRestApi.Svc,
			util.FncodeParams([]util.URLParameter{
				{common.ParamProvider: providerName},
				{common.ParamProtocol: fmt.Sprintf("delayInstance/%s/1000", instanceName)},
				{common.ParamTimes: common.CallTimes10Str},
			}))
		testkit.GetResponseInstanceAliasList(curl)
		time.Sleep(30 * time.Second)
	}
	It(text, func() {
		nameList := testkit.GetResponseInstanceAliasList(testUri)

		log.Println(fmt.Sprintf("type:%s,---instance list:%v", t, nameList))
		switch t {
		case "RoundRobin":
			IsRoundRobin(instancesLength, nameList...)
		case "SessionStickiness":
			isStrategySessionStickiness(instancesLength, nameList...)
		case "WeightedResponse":
			isWeightedResponse(instancesLength, nameList...)
		case "Random":
			isRand(instancesLength, nameList...)
		default:
		}
	})

}

var _ = Describe("Load balance", func() {
	testkit.SDKATContext(test)
})
