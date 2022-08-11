package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	SPLIT_LEN = 4
)

type RPCMsg struct {
	ServiceClass  string
	ServiceMethod string
	Payload       []interface{}
}

//func (msg *RPCMsg) Send(writer io.Writer) error {
//	//send header
//	_, err := writer.Write(msg.Header[:])
//	if err != nil {
//		return err
//	}
//
//	//write body total len :4 byte
//	dataLen := SPLIT_LEN + len(msg.ServiceClass) + SPLIT_LEN + len(msg.ServiceMethod) + SPLIT_LEN + len(msg.Payload)
//	err = binary.Write(writer, binary.BigEndian, uint32(dataLen)) //4
//	if err != nil {
//		return err
//	}
//
//	//write service class len :4 byte
//	err = binary.Write(writer, binary.BigEndian, uint32(len(msg.ServiceClass)))
//	if err != nil {
//		return err
//	}
//
//	//write service class
//	err = binary.Write(writer, binary.BigEndian, util.StringToByte(msg.ServiceClass))
//	if err != nil {
//		return err
//	}
//
//	//write service method len :4 byte
//	err = binary.Write(writer, binary.BigEndian, uint32(len(msg.ServiceMethod)))
//	if err != nil {
//		return err
//	}
//
//	//write service method
//	err = binary.Write(writer, binary.BigEndian, util.StringToByte(msg.ServiceMethod))
//	if err != nil {
//		return err
//	}
//
//	//write payload len :4 byte
//	err = binary.Write(writer, binary.BigEndian, uint32(len(msg.Payload)))
//	if err != nil {
//		return err
//	}
//
//	//write payload
//	//err = binary.Write(writer, binary.BigEndian, msg.Payload)
//	_, err = writer.Write(msg.Payload)
//	if err != nil {
//		return err
//	}
//	return nil
//
//}

func (msg *RPCMsg) Send(writer io.Writer, header *Header) error {
	coder := Codecs[header.SerializeType()]
	_, err := writer.Write(header[:])
	if err != nil {
		return err
	}
	data, err := coder.Encode(msg)
	if err != nil {
		return err
	}
	dataLen := len(data)
	err = binary.Write(writer, binary.BigEndian, uint32(dataLen)) //4
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	return nil

}

func Read(r io.Reader) (*RPCMsg, error) {
	msg := RPCMsg{}
	err := msg.Decode(r)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

//func (msg *RPCMsg) Decode(r io.Reader) error {
//	//read header
//	_, err := io.ReadFull(r, msg.Header[:])
//	if !msg.Header.CheckMagicNumber() { //magicNumber
//		return fmt.Errorf("magic number error: %v", msg.Header[0])
//	}
//
//	//total body len
//	headerByte := make([]byte, 4)
//	_, err = io.ReadFull(r, headerByte)
//	if err != nil {
//		return err
//	}
//	bodyLen := binary.BigEndian.Uint32(headerByte)
//
//	//read all body
//	data := make([]byte, bodyLen)
//	_, err = io.ReadFull(r, data)
//
//	//service class len
//	start := 0
//	end := start + SPLIT_LEN
//	classLen := binary.BigEndian.Uint32(data[start:end]) //0,4
//
//	//service class
//	start = end
//	end = start + int(classLen)
//	msg.ServiceClass = util.ByteToString(data[start:end]) //4,x
//
//	//service method len
//	start = end
//	end = start + SPLIT_LEN
//	methodLen := binary.BigEndian.Uint32(data[start:end]) //x,x+4
//
//	//service method
//	start = end
//	end = start + int(methodLen)
//	msg.ServiceMethod = util.ByteToString(data[start:end]) //x+4, x+4+y
//
//	//payload len
//	start = end
//	end = start + SPLIT_LEN
//	binary.BigEndian.Uint32(data[start:end]) //x+4+y, x+y+8 payloadLen
//
//	//payload
//	start = end
//	msg.Payload = data[start:]
//	return nil
//
//}

func (msg *RPCMsg) Decode(r io.Reader) error {
	//read header
	header := Header([HEADER_LEN]byte{})
	_, err := io.ReadFull(r, header[:])

	if err != nil {
		return err
	}
	if !header.CheckMagicNumber() { //magicNumber
		return fmt.Errorf("magic number error: %v", header[0])
	}
	dataByte := make([]byte, 4)
	_, err = io.ReadFull(r, dataByte)
	if err != nil {
		return err
	}
	dataLen := binary.BigEndian.Uint32(dataByte)
	//read all body
	data := make([]byte, dataLen)
	_, err = io.ReadFull(r, data)

	coder := Codecs[header.SerializeType()]
	err = coder.Decode(data, msg)
	return nil

}
