package loadbalance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLoadbalance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Load balancing Suite")
}
