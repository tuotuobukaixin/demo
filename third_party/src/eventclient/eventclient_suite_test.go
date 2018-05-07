package eventclient_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEventclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Eventclient Suite")
}
