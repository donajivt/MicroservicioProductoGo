package dtos

type ResponseDto struct {
	IsSuccess bool        `json:"is_success"`
	Result    interface{} `json:"result"`
	Message   string      `json:"message"`
}
