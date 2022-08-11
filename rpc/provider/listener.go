package provider

import (
	"fmt"
	"grpcClient/rpc/config"
	"grpcClient/rpc/protocol"
	"io"
	"log"
	"net"
)

type Listener interface {
	Run()
	SetHandler(string, Handler)
	Close()
}
type RPCListener struct {
	ServiceIp   string
	ServicePort int
	Handlers    map[string]Handler
	nl          net.Listener
}

func NewRPCListener(serviceIp string, servicePort int) *RPCListener {
	return &RPCListener{ServiceIp: serviceIp,
		ServicePort: servicePort,
		Handlers:    make(map[string]Handler)}
}
func (l *RPCListener) Run() {
	addr := fmt.Sprintf("%s:%d", l.ServiceIp, l.ServicePort)
	nl, err := net.Listen(config.NET_TRANS_PROTOCOL, addr) //tcp
	if err != nil {
		panic(err)
	}
	l.nl = nl
	for {
		conn, err := l.nl.Accept()
		if err != nil {
			continue
		}
		go l.handleConn(conn)
	}
}
func (l *RPCListener) Close() {
	if l.nl != nil {
		l.nl.Close()
	}
}

func (l *RPCListener) SetHandler(name string, handler Handler) {
	if _, ok := l.Handlers[name]; ok {
		log.Printf("%s is registered!\n", name)
		return
	}
	l.Handlers[name] = handler
}

func (l *RPCListener) handleConn(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("server %s catch panic err:%s\n", conn.RemoteAddr(), err)
		}
	}()
	for {
		msg, err := l.receiveData(conn)
		if err != nil || msg == nil {
			return
		}
		handler, ok := l.Handlers[msg.ServiceClass]
		if !ok {
			return
		}
		result, err := handler.Handle(msg.ServiceMethod, msg.Payload)
		err = l.sendData(conn, result)
		if err != nil {
			return
		}
	}
}

func (l *RPCListener) receiveData(conn net.Conn) (*protocol.RPCMsg, error) {
	msg, err := protocol.Read(conn)
	if err != nil {
		if err != io.EOF { //close
			return nil, err
		}
	}
	return msg, nil
}
func (l *RPCListener) sendData(conn net.Conn, payload []interface{}) error {
	header := protocol.Header([protocol.HEADER_LEN]byte{})
	header[0] = protocol.MagicNumber
	header.SetVersion(config.Protocol_MsgVersion)
	header.SetMsgType(protocol.Response)
	header.SetCompressType(config.CompressType)
	header.SetSerializeType(config.SerializeType)
	resMsg := protocol.RPCMsg{}
	resMsg.Payload = payload
	return  resMsg.Send(conn,&header)
}