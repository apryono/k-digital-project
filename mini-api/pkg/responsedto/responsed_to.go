package responsedto

import "github.com/k-digital-project/mini-api/usecase/viewmodel"

func ErrorResponse(message interface{}) viewmodel.ResponseErrorVM {
	err := []interface{}{message}
	res := viewmodel.ResponseErrorVM{Messages: err}

	return res
}

func SuccessResponse(data interface{}, meta interface{}) viewmodel.ResponseSuccessVM {
	return viewmodel.ResponseSuccessVM{
		Data: data,
		Meta: meta,
	}
}
