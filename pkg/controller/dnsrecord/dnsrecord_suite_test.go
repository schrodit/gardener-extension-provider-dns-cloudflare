package dnsrecord_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHealthz(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DNSRecord Actuator Suite")
}
