package lifopool_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLifopool(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LifoPool 测试套件")
}
