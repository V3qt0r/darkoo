package api

type Response struct {
	Code 			int 		`json:"code" yaml:"code" example:"500"`
	Message			string		`json:"message" yaml:"message"`
	Details 		interface{} `json:"details" yaml:"details"`
}


func NewResponse(code int, message string, details interface{}) *Response {
	return &Response{
		Code:		code,
		Message: 	message,
		Details: 	details,
	}
}