package gcrypto_test

import (
	"errors"
	. "gcrypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGcrypto(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gcrypto Suite")
}

var _ = Describe("Gcrypto", func() {
	var (
		algo string
		init InitFunc
	)

	BeforeEach(func() {
		algo = "aes"
		init = func() (Engine, error) {
			return nil, nil
		}
	})

	Context("Register a algorithm", func() {
		It("Register algorithm success", func() {
			Expect(Register(algo, init)).To(Succeed())
		})

		It("Register algorithm failed", func() {
			Expect(Register(algo, init)).To(Equal(errors.New("engine already registered " + algo)))
		})
	})

	Context("create a new crypto engine", func() {
		It("Create crypto engine success", func() {
			_, err := New(algo)
			Expect(err).To(Succeed())
		})

		It("Create crypto engine failed", func() {
			_, err := New("des")
			Expect(err).To(Equal(errors.New("no such algorithm des")))
		})
	})
})
