package circuit_test

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Timeout", func() {
	Describe("timeout", func() {
		Context("http server side has 10 secs process duration, client side timeout is 5 secs", func() {
			It("should be timeout", func() {
			})
		})

		Context("grpc server side has 10 secs process duration, client side timeout is 5 secs", func() {
			It("should be timeout", func() {
			})
		})

	})
	Describe("circuit", func() {
		Context("http server side has 10 secs process duration, client side timeout is 5 secs", func() {
			It("should open", func() {
			})
		})
		Context("after open. try to access other API", func() {
			It("should success", func() {
			})
		})
		Context("after sleep time", func() {
			It("should close", func() {
			})
		})
	})
})
