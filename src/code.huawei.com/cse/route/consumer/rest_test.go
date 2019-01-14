package consumer

import _ "github.com/huaweicse/auth/adaptor/gochassis"

import (
	"testing"

	"io/ioutil"

	"strings"

	"os"

	"code.huawei.com/cse/config"
	"code.huawei.com/cse/util"
	"github.com/go-mesh/openlogging"
	"github.com/stretchr/testify/assert"
)

var addParam Param

func init() {
	os.Setenv(util.ChassisHome, "/usr/cse-auth/test")
	os.Setenv(util.PAAS_PROJECT_NAME, "default")
	config.Init()
	addParam = Param{
		Data: map[string]interface{}{
			"cse.loadbalance.backoff.minMs": "300",
		},
		DimensionInfo: "ROUTERClient@default#0.0.1",
		Keys:          []string{"cse.loadbalance.backoff.minMs"},
	}
}
func TestServer_ConfigCenter(t *testing.T) {
	// test add api
	openlogging.GetLogger().Info("test start")

	resp, err := ccAdd(addParam)
	assert.Nil(t, err)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.True(t, strings.Contains(string(body), "Success"))

	resp, err = ccDelete(addParam)
	assert.Nil(t, err)
	body, _ = ioutil.ReadAll(resp.Body)
	assert.True(t, strings.Contains(string(body), "Success"))
}
