package response

type ErrorResponse struct {
	Status 	int `json:"status"`
	Code 	string `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Data interface{} `json:"data"`
}