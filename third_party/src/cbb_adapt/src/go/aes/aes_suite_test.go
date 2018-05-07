package aes_test

import (
	"gcrypto"
	. "gcrypto/aes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Aes Suite")
}

var _ = Describe("Aes", func() {
	var (
		data    string
		eng     gcrypto.Engine
		err     error
		encdata string
		decdata string
	)

	BeforeEach(func() {
		data = "rnd-mushroom.huawei.com"
	})

	Describe("Init aes angine", func() {
		It("Init aes angine success", func() {
			eng, err = Init()
			Expect(err).To(Succeed())
		})
	})

	Describe("Encrypt and Decrypt", func() {
		It("Encrypt", func() {
			encdata, err = eng.Encrypt(0, data)
			Expect(err).To(Succeed())
		})

		It("Decrypt", func() {
			decdata, err = eng.Decrypt(0, encdata)
			Expect(err).To(Succeed())
			Expect(decdata).To(Equal(data))
		})
	})
})
