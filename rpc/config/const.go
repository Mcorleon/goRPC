package config

import "grpcClient/rpc/protocol"

const (
	Protocol_MsgVersion = 1
	NET_TRANS_PROTOCOL="tcp"
	CompressType=0
	SerializeType=protocol.Gob
)
