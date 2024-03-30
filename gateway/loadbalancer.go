package gateway

import (
	"github.com/zhaommmmomo/zim/common/config"
)

type LoadBalance int

const (
	roundRobin LoadBalance = iota
)

type (
	loadBalancer interface {
		register(*reactor)
		next() *reactor
		index(uint8) *reactor
		iterate(func(uint8, *reactor) bool)
		len() uint8
	}

	baseLoadBalancer struct {
		reactors []*reactor
		size     uint8
	}

	roundRobinLoadBalancer struct {
		baseLoadBalancer
		nextIndex uint8
	}
)

func newLoadBalancer() loadBalancer {
	lb := config.GetGatewayEpollLoadBalancer()
	switch LoadBalance(lb) {
	case roundRobin:
		return new(roundRobinLoadBalancer)
	}
	return new(roundRobinLoadBalancer)
}

func (lb *baseLoadBalancer) register(r *reactor) {
	r.idx = lb.size
	lb.reactors = append(lb.reactors, r)
	lb.size++
}

func (lb *baseLoadBalancer) index(i uint8) *reactor {
	if i > lb.size {
		return nil
	}
	return lb.reactors[i]
}

func (lb *baseLoadBalancer) iterate(f func(uint8, *reactor) bool) {
	for i, r := range lb.reactors {
		if !f(uint8(i), r) {
			break
		}
	}
}

func (lb *baseLoadBalancer) len() uint8 {
	return lb.size
}

func (lb *roundRobinLoadBalancer) next() *reactor {
	if lb.nextIndex >= lb.size {
		lb.nextIndex = 0
	}
	r := lb.reactors[lb.nextIndex]
	lb.nextIndex++
	return r
}
