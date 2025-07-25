// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	amqp "github.com/streadway/amqp"
	mock "github.com/stretchr/testify/mock"
)

// IRabbitMQClient is an autogenerated mock type for the IRabbitMQClient type
type IRabbitMQClient struct {
	mock.Mock
}

type IRabbitMQClient_Expecter struct {
	mock *mock.Mock
}

func (_m *IRabbitMQClient) EXPECT() *IRabbitMQClient_Expecter {
	return &IRabbitMQClient_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with no fields
func (_m *IRabbitMQClient) Close() {
	_m.Called()
}

// IRabbitMQClient_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type IRabbitMQClient_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *IRabbitMQClient_Expecter) Close() *IRabbitMQClient_Close_Call {
	return &IRabbitMQClient_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *IRabbitMQClient_Close_Call) Run(run func()) *IRabbitMQClient_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *IRabbitMQClient_Close_Call) Return() *IRabbitMQClient_Close_Call {
	_c.Call.Return()
	return _c
}

func (_c *IRabbitMQClient_Close_Call) RunAndReturn(run func()) *IRabbitMQClient_Close_Call {
	_c.Run(run)
	return _c
}

// Consume provides a mock function with given fields: queueName, handler
func (_m *IRabbitMQClient) Consume(queueName string, handler func(amqp.Delivery)) error {
	ret := _m.Called(queueName, handler)

	if len(ret) == 0 {
		panic("no return value specified for Consume")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, func(amqp.Delivery)) error); ok {
		r0 = rf(queueName, handler)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IRabbitMQClient_Consume_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Consume'
type IRabbitMQClient_Consume_Call struct {
	*mock.Call
}

// Consume is a helper method to define mock.On call
//   - queueName string
//   - handler func(amqp.Delivery)
func (_e *IRabbitMQClient_Expecter) Consume(queueName interface{}, handler interface{}) *IRabbitMQClient_Consume_Call {
	return &IRabbitMQClient_Consume_Call{Call: _e.mock.On("Consume", queueName, handler)}
}

func (_c *IRabbitMQClient_Consume_Call) Run(run func(queueName string, handler func(amqp.Delivery))) *IRabbitMQClient_Consume_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(func(amqp.Delivery)))
	})
	return _c
}

func (_c *IRabbitMQClient_Consume_Call) Return(_a0 error) *IRabbitMQClient_Consume_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IRabbitMQClient_Consume_Call) RunAndReturn(run func(string, func(amqp.Delivery)) error) *IRabbitMQClient_Consume_Call {
	_c.Call.Return(run)
	return _c
}

// Publish provides a mock function with given fields: queue, body
func (_m *IRabbitMQClient) Publish(queue string, body []byte) error {
	ret := _m.Called(queue, body)

	if len(ret) == 0 {
		panic("no return value specified for Publish")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte) error); ok {
		r0 = rf(queue, body)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IRabbitMQClient_Publish_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Publish'
type IRabbitMQClient_Publish_Call struct {
	*mock.Call
}

// Publish is a helper method to define mock.On call
//   - queue string
//   - body []byte
func (_e *IRabbitMQClient_Expecter) Publish(queue interface{}, body interface{}) *IRabbitMQClient_Publish_Call {
	return &IRabbitMQClient_Publish_Call{Call: _e.mock.On("Publish", queue, body)}
}

func (_c *IRabbitMQClient_Publish_Call) Run(run func(queue string, body []byte)) *IRabbitMQClient_Publish_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].([]byte))
	})
	return _c
}

func (_c *IRabbitMQClient_Publish_Call) Return(_a0 error) *IRabbitMQClient_Publish_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IRabbitMQClient_Publish_Call) RunAndReturn(run func(string, []byte) error) *IRabbitMQClient_Publish_Call {
	_c.Call.Return(run)
	return _c
}

// NewIRabbitMQClient creates a new instance of IRabbitMQClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIRabbitMQClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *IRabbitMQClient {
	mock := &IRabbitMQClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
