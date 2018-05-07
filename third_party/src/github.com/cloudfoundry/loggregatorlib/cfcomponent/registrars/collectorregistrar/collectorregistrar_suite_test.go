package collectorregistrar_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCollectorregistrar(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Collectorregistrar Suite")
}
