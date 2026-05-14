package grpcx

import (
	"encoding/json"

	"google.golang.org/grpc/encoding"
)

const JSONCodecName = "json"

// JSONCodec 让 gRPC 在内部服务通信时也能直接使用普通 Go 结构体作为消息载体。
type JSONCodec struct{}

func (JSONCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (JSONCodec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (JSONCodec) Name() string {
	return JSONCodecName
}

func init() {
	encoding.RegisterCodec(JSONCodec{})
}
