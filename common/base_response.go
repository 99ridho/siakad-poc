package common

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

type BaseResponseError struct {
	Message   string   `json:"message"`
	Details   []string `json:"details"`
	Timestamp string   `json:"timestamp"`
	Path      string   `json:"path"`
}

type BaseResponse[Data any] struct {
	Status string             `json:"status"`
	Data   *Data              `json:"data,omitempty"`
	Error  *BaseResponseError `json:"error,omitempty"`
}

type PaginationMetadata struct {
	Page         int `json:"page"`
	PageSize     int `json:"page_size"`
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
}

type PaginatedBaseResponse[Data any] struct {
	BaseResponse[Data]
	Paging *PaginationMetadata `json:"paging,omitempty"`
}
