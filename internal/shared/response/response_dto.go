package response

type PaginationDTO struct {
	Page        int `json:"page" example:"1"`
	CurrentPage int `json:"current_page" example:"1"`
	PerPage     int `json:"per_page" example:"10"`
	LastPage    int `json:"last_page" example:"5"`
}

type ErrorItemDTO struct {
	Error   string `json:"error" example:"Validation failed"`
	Message string `json:"message" example:"Email is required"`
}

type APIResponseDTO struct {
	Payload    interface{}  `json:"payload"`
	Errors     []ErrorItemDTO `json:"error"`
	Pagination *PaginationDTO `json:"pagination,omitempty"`
}