package viewmodel

type ResponseErrorVM struct {
	Messages interface{} `json:"messages"`
}

type ResponseSuccessVM struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}
