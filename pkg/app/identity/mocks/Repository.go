// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/karta0898098/iam/pkg/app/identity/entity"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

type Repository_Expecter struct {
	mock *mock.Mock
}

func (_m *Repository) EXPECT() *Repository_Expecter {
	return &Repository_Expecter{mock: &_m.Mock}
}

// FindUserByUsername provides a mock function with given fields: ctx, username
func (_m *Repository) FindUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	ret := _m.Called(ctx, username)

	var r0 *entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*entity.User, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.User); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_FindUserByUsername_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindUserByUsername'
type Repository_FindUserByUsername_Call struct {
	*mock.Call
}

// FindUserByUsername is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
func (_e *Repository_Expecter) FindUserByUsername(ctx interface{}, username interface{}) *Repository_FindUserByUsername_Call {
	return &Repository_FindUserByUsername_Call{Call: _e.mock.On("FindUserByUsername", ctx, username)}
}

func (_c *Repository_FindUserByUsername_Call) Run(run func(ctx context.Context, username string)) *Repository_FindUserByUsername_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Repository_FindUserByUsername_Call) Return(profile *entity.User, err error) *Repository_FindUserByUsername_Call {
	_c.Call.Return(profile, err)
	return _c
}

func (_c *Repository_FindUserByUsername_Call) RunAndReturn(run func(context.Context, string) (*entity.User, error)) *Repository_FindUserByUsername_Call {
	_c.Call.Return(run)
	return _c
}

// StoreSession provides a mock function with given fields: ctx, session
func (_m *Repository) StoreSession(ctx context.Context, session *entity.Session) error {
	ret := _m.Called(ctx, session)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.Session) error); ok {
		r0 = rf(ctx, session)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_StoreSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StoreSession'
type Repository_StoreSession_Call struct {
	*mock.Call
}

// StoreSession is a helper method to define mock.On call
//   - ctx context.Context
//   - session *entity.Session
func (_e *Repository_Expecter) StoreSession(ctx interface{}, session interface{}) *Repository_StoreSession_Call {
	return &Repository_StoreSession_Call{Call: _e.mock.On("StoreSession", ctx, session)}
}

func (_c *Repository_StoreSession_Call) Run(run func(ctx context.Context, session *entity.Session)) *Repository_StoreSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.Session))
	})
	return _c
}

func (_c *Repository_StoreSession_Call) Return(err error) *Repository_StoreSession_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *Repository_StoreSession_Call) RunAndReturn(run func(context.Context, *entity.Session) error) *Repository_StoreSession_Call {
	_c.Call.Return(run)
	return _c
}

// StoreUser provides a mock function with given fields: ctx, user
func (_m *Repository) StoreUser(ctx context.Context, user *entity.User) error {
	ret := _m.Called(ctx, user)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_StoreUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StoreUser'
type Repository_StoreUser_Call struct {
	*mock.Call
}

// StoreUser is a helper method to define mock.On call
//   - ctx context.Context
//   - user *entity.User
func (_e *Repository_Expecter) StoreUser(ctx interface{}, user interface{}) *Repository_StoreUser_Call {
	return &Repository_StoreUser_Call{Call: _e.mock.On("StoreUser", ctx, user)}
}

func (_c *Repository_StoreUser_Call) Run(run func(ctx context.Context, user *entity.User)) *Repository_StoreUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.User))
	})
	return _c
}

func (_c *Repository_StoreUser_Call) Return(err error) *Repository_StoreUser_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *Repository_StoreUser_Call) RunAndReturn(run func(context.Context, *entity.User) error) *Repository_StoreUser_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
