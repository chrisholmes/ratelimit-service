package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WhiteList", func() {

	var (
		whiteList *WhiteList
	)

	Describe("has a list of IPs", func() {

		BeforeEach(func() {
      ips := []string{"192.168.1.1"}
			whiteList = NewWhiteList(ips)
		})

		It("reports if does contain an IP", func() {
			ip := "192.168.1.1"
			Expect(whiteList.contains(ip)).To(BeTrue())
		})

		It("reports if doesn't contain an IP", func() {
			ip := "192.168.1.2"
			Expect(whiteList.contains(ip)).To(BeFalse())
		})

	})

})
