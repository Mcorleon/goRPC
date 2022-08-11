package client

import (
	"context"
	"errors"
)
type ClientProxy interface {
	Call(context.Context, string, interface{}, ...interface{}) (interface{}, error)
}

type RPCClientProxy struct {
	option Option
}

func NewClientProxy(option Option) ClientProxy {
	cp := &RPCClientProxy{
		option:  option,
	}
	return cp
}

func (cp *RPCClientProxy) Call(ctx context.Context, servicePath string, stub interface{}, params ...interface{}) (interface{}, error) {
	service, err := NewService(servicePath)
	if err != nil {
		return nil, err
	}
	client := NewClient(cp.option)
	addr := service.SelectAddr()
	err = client.Connect(addr) //TODO 长连接管理
	if err != nil {
		return nil, err
	}
	retries := cp.option.Retries
	for retries > 0 {
		retries--
		return client.Invoke(ctx, service, stub, params...)
	}
	return nil, errors.New("error")
}