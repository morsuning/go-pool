package gopool_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGopool(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoPool Suite")
}
