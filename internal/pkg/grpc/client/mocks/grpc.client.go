package mocks

import (
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type GrpcClientMock struct {
	mock.Mock
}

func (_m *GrpcClientMock) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *GrpcClientMock) GetGrpcConnection() *grpc.ClientConn {
	ret := _m.Called()

	var r0 *grpc.ClientConn
	if rf, ok := ret.Get(0).(func() *grpc.ClientConn); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*grpc.ClientConn)
		}
	}

	return r0
}

type mockNewGrpcClientConstructor interface {
	mock.TestingT
	Cleanup(func())
}

func newGrpcClient(t mockNewGrpcClientConstructor) *GrpcClientMock {
	m := &GrpcClientMock{}
	m.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectationsForObjects(t) })

	return m
}
