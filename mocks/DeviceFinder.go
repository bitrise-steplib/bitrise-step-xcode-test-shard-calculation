// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	destination "github.com/bitrise-io/go-xcode/v2/destination"
	mock "github.com/stretchr/testify/mock"
)

// DeviceFinder is an autogenerated mock type for the DeviceFinder type
type DeviceFinder struct {
	mock.Mock
}

// FindDevice provides a mock function with given fields: _a0
func (_m *DeviceFinder) FindDevice(_a0 destination.Simulator) (destination.Device, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for FindDevice")
	}

	var r0 destination.Device
	var r1 error
	if rf, ok := ret.Get(0).(func(destination.Simulator) (destination.Device, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(destination.Simulator) destination.Device); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(destination.Device)
	}

	if rf, ok := ret.Get(1).(func(destination.Simulator) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListDevices provides a mock function with no fields
func (_m *DeviceFinder) ListDevices() (*destination.DeviceList, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ListDevices")
	}

	var r0 *destination.DeviceList
	var r1 error
	if rf, ok := ret.Get(0).(func() (*destination.DeviceList, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *destination.DeviceList); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*destination.DeviceList)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewDeviceFinder creates a new instance of DeviceFinder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDeviceFinder(t interface {
	mock.TestingT
	Cleanup(func())
}) *DeviceFinder {
	mock := &DeviceFinder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
