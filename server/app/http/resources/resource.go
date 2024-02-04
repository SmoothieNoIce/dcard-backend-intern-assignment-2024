package resources

type JSONResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginateResouce struct {
	TotalPages int         `json:"total_pages"`
	Count      int         `json:"count"`
	Results    interface{} `json:"results"`
}
