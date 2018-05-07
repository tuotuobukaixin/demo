package datacollector_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDatacollector(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Datacollector Suite")
}
