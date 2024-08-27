// Code generated by mockery v2.44.2. DO NOT EDIT.

package order

import (
	context "context"

	lks "github.com/daemondxx/lks_back/internal/api/lks"
	mock "github.com/stretchr/testify/mock"
)

// MockLksAPI is an autogenerated mock type for the LksAPI type
type MockLksAPI struct {
	mock.Mock
}

type MockLksAPI_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLksAPI) EXPECT() *MockLksAPI_Expecter {
	return &MockLksAPI_Expecter{mock: &_m.Mock}
}

// GetActualDuty provides a mock function with given fields: ctx, p
func (_m *MockLksAPI) GetActualDuty(ctx context.Context, p lks.AuthPayload) ([]lks.CurrentDuty, error) {
	ret := _m.Called(ctx, p)

	if len(ret) == 0 {
		panic("no return value specified for GetActualDuty")
	}

	var r0 []lks.CurrentDuty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, lks.AuthPayload) ([]lks.CurrentDuty, error)); ok {
		return rf(ctx, p)
	}
	if rf, ok := ret.Get(0).(func(context.Context, lks.AuthPayload) []lks.CurrentDuty); ok {
		r0 = rf(ctx, p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]lks.CurrentDuty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, lks.AuthPayload) error); ok {
		r1 = rf(ctx, p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockLksAPI_GetActualDuty_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetActualDuty'
type MockLksAPI_GetActualDuty_Call struct {
	*mock.Call
}

// GetActualDuty is a helper method to define mock.On call
//   - ctx context.Context
//   - p lks.AuthPayload
func (_e *MockLksAPI_Expecter) GetActualDuty(ctx interface{}, p interface{}) *MockLksAPI_GetActualDuty_Call {
	return &MockLksAPI_GetActualDuty_Call{Call: _e.mock.On("GetActualDuty", ctx, p)}
}

func (_c *MockLksAPI_GetActualDuty_Call) Run(run func(ctx context.Context, p lks.AuthPayload)) *MockLksAPI_GetActualDuty_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(lks.AuthPayload))
	})
	return _c
}

func (_c *MockLksAPI_GetActualDuty_Call) Return(_a0 []lks.CurrentDuty, _a1 error) *MockLksAPI_GetActualDuty_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockLksAPI_GetActualDuty_Call) RunAndReturn(run func(context.Context, lks.AuthPayload) ([]lks.CurrentDuty, error)) *MockLksAPI_GetActualDuty_Call {
	_c.Call.Return(run)
	return _c
}

// GetArchiveDuty provides a mock function with given fields: ctx, p, month, year
func (_m *MockLksAPI) GetArchiveDuty(ctx context.Context, p lks.AuthPayload, month int, year int) ([]lks.ArchiveDuty, error) {
	ret := _m.Called(ctx, p, month, year)

	if len(ret) == 0 {
		panic("no return value specified for GetArchiveDuty")
	}

	var r0 []lks.ArchiveDuty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, lks.AuthPayload, int, int) ([]lks.ArchiveDuty, error)); ok {
		return rf(ctx, p, month, year)
	}
	if rf, ok := ret.Get(0).(func(context.Context, lks.AuthPayload, int, int) []lks.ArchiveDuty); ok {
		r0 = rf(ctx, p, month, year)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]lks.ArchiveDuty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, lks.AuthPayload, int, int) error); ok {
		r1 = rf(ctx, p, month, year)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockLksAPI_GetArchiveDuty_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetArchiveDuty'
type MockLksAPI_GetArchiveDuty_Call struct {
	*mock.Call
}

// GetArchiveDuty is a helper method to define mock.On call
//   - ctx context.Context
//   - p lks.AuthPayload
//   - month int
//   - year int
func (_e *MockLksAPI_Expecter) GetArchiveDuty(ctx interface{}, p interface{}, month interface{}, year interface{}) *MockLksAPI_GetArchiveDuty_Call {
	return &MockLksAPI_GetArchiveDuty_Call{Call: _e.mock.On("GetArchiveDuty", ctx, p, month, year)}
}

func (_c *MockLksAPI_GetArchiveDuty_Call) Run(run func(ctx context.Context, p lks.AuthPayload, month int, year int)) *MockLksAPI_GetArchiveDuty_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(lks.AuthPayload), args[2].(int), args[3].(int))
	})
	return _c
}

func (_c *MockLksAPI_GetArchiveDuty_Call) Return(_a0 []lks.ArchiveDuty, _a1 error) *MockLksAPI_GetArchiveDuty_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockLksAPI_GetArchiveDuty_Call) RunAndReturn(run func(context.Context, lks.AuthPayload, int, int) ([]lks.ArchiveDuty, error)) *MockLksAPI_GetArchiveDuty_Call {
	_c.Call.Return(run)
	return _c
}

// GetPerspectiveDuty provides a mock function with given fields: ctx, p, month, year
func (_m *MockLksAPI) GetPerspectiveDuty(ctx context.Context, p lks.AuthPayload, month int, year int) ([]lks.PerspectiveDuty, error) {
	ret := _m.Called(ctx, p, month, year)

	if len(ret) == 0 {
		panic("no return value specified for GetPerspectiveDuty")
	}

	var r0 []lks.PerspectiveDuty
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, lks.AuthPayload, int, int) ([]lks.PerspectiveDuty, error)); ok {
		return rf(ctx, p, month, year)
	}
	if rf, ok := ret.Get(0).(func(context.Context, lks.AuthPayload, int, int) []lks.PerspectiveDuty); ok {
		r0 = rf(ctx, p, month, year)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]lks.PerspectiveDuty)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, lks.AuthPayload, int, int) error); ok {
		r1 = rf(ctx, p, month, year)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockLksAPI_GetPerspectiveDuty_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPerspectiveDuty'
type MockLksAPI_GetPerspectiveDuty_Call struct {
	*mock.Call
}

// GetPerspectiveDuty is a helper method to define mock.On call
//   - ctx context.Context
//   - p lks.AuthPayload
//   - month int
//   - year int
func (_e *MockLksAPI_Expecter) GetPerspectiveDuty(ctx interface{}, p interface{}, month interface{}, year interface{}) *MockLksAPI_GetPerspectiveDuty_Call {
	return &MockLksAPI_GetPerspectiveDuty_Call{Call: _e.mock.On("GetPerspectiveDuty", ctx, p, month, year)}
}

func (_c *MockLksAPI_GetPerspectiveDuty_Call) Run(run func(ctx context.Context, p lks.AuthPayload, month int, year int)) *MockLksAPI_GetPerspectiveDuty_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(lks.AuthPayload), args[2].(int), args[3].(int))
	})
	return _c
}

func (_c *MockLksAPI_GetPerspectiveDuty_Call) Return(_a0 []lks.PerspectiveDuty, _a1 error) *MockLksAPI_GetPerspectiveDuty_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockLksAPI_GetPerspectiveDuty_Call) RunAndReturn(run func(context.Context, lks.AuthPayload, int, int) ([]lks.PerspectiveDuty, error)) *MockLksAPI_GetPerspectiveDuty_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockLksAPI creates a new instance of MockLksAPI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockLksAPI(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockLksAPI {
	mock := &MockLksAPI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}