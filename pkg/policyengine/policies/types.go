package policies

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/cilium/cilium/pkg/spanstat"
	"github.com/cilium/cilium/pkg/u8proto"
)

type noOpStats struct{}

func (n noOpStats) WaitingForPolicyRepository() *spanstat.SpanStat {
	return new(spanstat.SpanStat)
}

func (n noOpStats) PolicyCalculation() *spanstat.SpanStat {
	return new(spanstat.SpanStat)
}

type fakeEndpoint struct{}

func (f *fakeEndpoint) GetID() uint64 {
	return 0
}

func (f *fakeEndpoint) GetNamedPort(ingress bool, name string, proto u8proto.U8proto) uint16 {
	return 0
}

func (f *fakeEndpoint) PolicyDebug(fields logrus.Fields, msg string) {
	fmt.Printf("%s: %v\n", msg, fields)
}

func (f *fakeEndpoint) IsHost() bool {
	return false
}

// MapStateSize is used for the initial golang map size
func (f *fakeEndpoint) MapStateSize() int {
	return 0
}
