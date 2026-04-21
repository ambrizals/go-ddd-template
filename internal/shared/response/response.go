package response

type Pagination struct {
	Page        int    `json:"page" example:"1"`
	CurrentPage int    `json:"current_page" example:"1"`
	PerPage     int    `json:"per_page" example:"10"`
	LastPage    int    `json:"last_page" example:"5"`
}

type ErrorItem struct {
	Error   string `json:"error" example:"Validation failed"`
	Message string `json:"message" example:"Email is required"`
}

type APIResponse[T any] struct {
	Payload    *T              `json:"payload"`
	Errors     []ErrorItem     `json:"error"`
	Pagination *Pagination     `json:"pagination,omitempty"`
}

type APIListResponse[T any] struct {
	Payload    []T             `json:"payload"`
	Errors     []ErrorItem     `json:"error"`
	Pagination *Pagination     `json:"pagination,omitempty"`
}

func NewAPIResponse[T any](data *T) APIResponse[T] {
	return APIResponse[T]{
		Payload: data,
		Errors:  []ErrorItem{},
	}
}

func NewAPIListResponse[T any](data []T, page, currentPage, perPage, lastPage int) APIListResponse[T] {
	return APIListResponse[T]{
		Payload: data,
		Errors:  []ErrorItem{},
		Pagination: &Pagination{
			Page:        page,
			CurrentPage: currentPage,
			PerPage:     perPage,
			LastPage:    lastPage,
		},
	}
}

func NewAPIErrorResponse[T any](errorMsg, message string) APIResponse[T] {
	return APIResponse[T]{
		Payload: nil,
		Errors: []ErrorItem{
			{Error: errorMsg, Message: message},
		},
	}
}