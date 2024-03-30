package domain

type Response struct {
	Code    int16       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func PackRes(data interface{}) Response {
	return Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}
