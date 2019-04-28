package model

import "encoding/json"

type ErrResponse struct {
	Code    int32  `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (er *ErrResponse) Error() string {
	j, _ := json.Marshal(er)
	return string(j)
}

func NewErrResponse(code int32, t string, message string) error {
	er := ErrResponse{
		Code:    code,
		Type:    t,
		Message: message,
	}
	return &er
}
