package kcli_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestKCLI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KCLI")
}
