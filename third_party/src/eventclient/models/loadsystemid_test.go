package models_test

import (
	. "eventclient/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Loadsystemid", func() {
	sysConfig := LoadSystemIDFromPath("./systemIdTest.json", 0)

	Describe("Loadsystemid", func() {
		Context("Loadsystemid", func() {
			It("should be success to Loadsystemid", func() {

				Expect(sysConfig).ShouldNot(BeNil())
				err := sysConfig.WriteConfigToFile()
				Expect(err).ShouldNot(HaveOccurred())

			})
		})
	})
})
