package util

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestProxy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoFibonacciHeap Suite")
}

var _ = Describe("Test initialization", func() {
	Context("Register suite setup and teardown function", func() {
		BeforeSuite(func() {
		})

		AfterSuite(func() {
		})
	})
})
