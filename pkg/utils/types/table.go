package types

type TableData struct {
	Data       interface{} `json:"data"`
	PageSize   int         `json:"page_size"`
	PageNo     int         `json:"page_no"`
	TotalPage  int         `json:"total_page"`
	TotalCount int         `json:"total_count"`
}
