package constant

type DataResponse struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
}

type ListResponse struct {
	StatusCode int `json:"status_code"`
	// Pagination *helper.Pagination `json:"_pagination"`
	Data    interface{} `json:"data"`
	Message string
}

type StatusResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
