package core

type ResponseError struct {
	ErrMsg string `json:"error_msg"`
	ErrCode int `json:"err_code"`
}

const ErrorNoResponse = "%s 微信服务器无反馈"