package aggretastic_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAggretastic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Aggretastic Suite")
}
