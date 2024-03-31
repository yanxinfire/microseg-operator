package types

const Success = "success"

type SuccessResponse struct {
	APIVersion string      `json:"apiVersion"`
	Data       interface{} `json:"data"`
}

type FailResponse struct {
	APIVersion string  `json:"apiVersion"`
	Error      RespErr `json:"error"`
}

type RespErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ListRespData struct {
	Items        interface{} `json:"items"`
	ItemsPerPage int         `json:"itemsPerPage"`
	TotalItems   int         `json:"totalItems"`
}

type SingleRespData struct {
	Item interface{} `json:"item"`
}
