package xjson

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin/binding"
)

type jsoniterBinding struct{}

func (jsoniterBinding) Name() string {
	return "json"
}

func (jsoniterBinding) Bind(req *http.Request, obj any) error {
	if req.Body == nil {
		return errors.New("invalid request")
	}
	return JSON.NewDecoder(req.Body).Decode(obj)
}

// 🔥 必须实现
func (jsoniterBinding) BindBody(body []byte, obj any) error {
	return JSON.Unmarshal(body, obj)
}

func InitGin() {
	binding.JSON = &jsoniterBinding{}
}
