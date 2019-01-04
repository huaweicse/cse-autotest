package communication_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCommunication(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Communication Suite")
}
