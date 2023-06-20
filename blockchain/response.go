package main

type ResponseSuccess struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func NewResponseSuccess(data interface{}) *ResponseSuccess {
	return &ResponseSuccess{
		Status: "success",
		Data:   data,
	}
}

const (
	ERROR_PUSH_BAD_DATA = "10001"
	ERROR_BAD_LIMIT     = "10002"
	ERROR_BAD_OFFSET    = "10003"
	ERROR_REPOSITORY    = "10004"
)

type ResponseFail struct {
	Status    string      `json:"status"`
	ErrorCode string      `json:"error_code"`
	Message   interface{} `json:"data"`
}

func NewResponseFail(errorCode string, message interface{}) *ResponseFail {
	return &ResponseFail{Status: "fail", ErrorCode: errorCode, Message: message}
}
