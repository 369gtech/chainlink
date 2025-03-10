// Code generated by mockery v2.10.1. DO NOT EDIT.

package mocks

import (
	types "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	mock "github.com/stretchr/testify/mock"
)

// HeadListener is an autogenerated mock type for the HeadListener type
type HeadListener struct {
	mock.Mock
}

// Connected provides a mock function with given fields:
func (_m *HeadListener) Connected() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ListenForNewHeads provides a mock function with given fields: handleNewHead, done
func (_m *HeadListener) ListenForNewHeads(handleNewHead types.NewHeadHandler, done func()) {
	_m.Called(handleNewHead, done)
}

// ReceivingHeads provides a mock function with given fields:
func (_m *HeadListener) ReceivingHeads() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
