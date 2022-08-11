package client

import (
	"context"
	"errors"
	"fmt"
	"grpcClient/rpc/config"
	"grpcClient/rpc/protocol"
	"log"
	"net"
	"reflect"
	"time"
)

type Client interface {
	Connect(string) error
	Invoke(context.Context, *Service, interface{}, ...interface{}) (interface{}, error)
	Close()
}

type Option struct {
	Retries           int
	ConnectionTimeout time.Duration
	SerializeType     protocol.SerializeType
	CompressType      protocol.CompressType
}
var DefaultOption = Option{
	Retries:           3,
	ConnectionTimeout: 5 * time.Second,
	SerializeType:     protocol.Gob,
	CompressType:      protocol.None,
}
type RPCClient struct {
	conn   net.Conn
	option Option
}
func NewClient(option Option) Client {
	return &RPCClient{option: option}
}

func (cli *RPCClient) Connect(addr string) error {
	conn, err := net.DialTimeout(config.NET_TRANS_PROTOCOL, addr, cli.option.ConnectionTimeout)
	if err != nil {
		return err
	}
	cli.conn = conn
	return nil
}
func (cli *RPCClient) Invoke(ctx context.Context, service *Service, stub interface{}, params ...interface{}) (interface{}, error) {
	cli.makeCall(service, stub)
	return cli.wrapCall(ctx, stub, params...)
}
func (cli *RPCClient) Close() {
	if cli.conn != nil {
		cli.conn.Close()
	}
}

func (cli *RPCClient) makeCall(service *Service, methodPtr interface{}) {
	container := reflect.ValueOf(methodPtr).Elem()

	handler := func(req []reflect.Value) []reflect.Value {
		numOut := container.Type().NumOut()
		errorHandler := func(err error) []reflect.Value {
			outArgs := make([]reflect.Value, numOut)
			for i := 0; i < len(outArgs)-1; i++ {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}
		inArgs := make([]interface{}, 0, len(req))
		for _, arg := range req {
			inArgs = append(inArgs, arg.Interface())
		}

		msg := protocol.RPCMsg{}
		msg.ServiceClass = service.Class
		msg.ServiceMethod = service.Method
		msg.Payload = inArgs
		header := protocol.Header([protocol.HEADER_LEN]byte{})
		header[0] = protocol.MagicNumber
		header.SetVersion(config.Protocol_MsgVersion)
		header.SetMsgType(protocol.Request)
		header.SetCompressType(config.CompressType)
		header.SetSerializeType(config.SerializeType)
		err := msg.Send(cli.conn,&header)
		if err != nil {
			log.Printf("send err:%v\n", err)
			return errorHandler(err)
		}
		respMsg, err := protocol.Read(cli.conn)
		if err != nil {
			return errorHandler(err)
		}
		respDecode:=respMsg.Payload
		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut {
				if respDecode[i] == nil {
					outArgs[i] = reflect.Zero(container.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(respDecode[i])
				}
			} else {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
		}
		return outArgs
	}
	container.Set(reflect.MakeFunc(container.Type(), handler))
}

func (cli *RPCClient) wrapCall(ctx context.Context, stub interface{}, params ...interface{}) (interface{}, error) {
	f := reflect.ValueOf(stub).Elem()
	if len(params) != f.Type().NumIn() {
		return nil, errors.New(fmt.Sprintf("params not adapted: %d-%d", len(params), f.Type().NumIn()))
	}
	in := make([]reflect.Value, len(params))
	for idx, param := range params {
		in[idx] = reflect.ValueOf(param)
	}
	result := f.Call(in)
	return result, nil
}