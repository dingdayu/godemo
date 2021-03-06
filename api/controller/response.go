package controller

import "fmt"

var (
	SuccessResponse = Response{Code: 10000, Message: "请求成功"} //success default

	// 系统响应    100xxx
	ErrInternalServer = Response{Code: 10001, Message: "系统错误"}
	ErrMissParams     = Response{Code: 10002, Message: "缺少参数"}
	ErrFailParams     = Response{Code: 10003, Message: "参数格式错误"}
	ErrNotExist       = Response{Code: 10004, Message: "数据不存在"}
	ErrDefault        = Response{Code: 10005, Message: "操作失败"}
	ErrDataPermission = Response{Code: 10006, Message: "没有此数据权限"}

	// Auth 认证登陆 响应 101xx
	ErrAuthForbidden    = Response{Code: 10100, Message: "登录信息不存"}
	ErrAuthUserNotFound = Response{Code: 10101, Message: "登陆用户不存在"}
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e Response) Error() string {
	return fmt.Sprintf(e.Message)
}

func NewErrResponse(message string) Response {
	s := ErrInternalServer
	s.Message = message
	return s
}

func NewSucResponse(data interface{}) Response {
	s := SuccessResponse
	s.Data = data
	return s
}
func NewErrParamsResponse(message string) Response {
	s := ErrMissParams
	s.Message = message
	return s
}
