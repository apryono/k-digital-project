package viewmodel

type ResponseErrorVM struct {
	Messages interface{} `json:"messages"`
}

type ResponsesSuccessVM struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}
