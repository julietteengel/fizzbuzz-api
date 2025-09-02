package model

type FizzBuzzRequest struct {
	Int1  int    `json:"int1" validate:"required,min=1"`
	Int2  int    `json:"int2" validate:"required,min=1"`
	Limit int    `json:"limit" validate:"required,min=1,max=10000"`
	Str1  string `json:"str1" validate:"required,min=1,max=100"`
	Str2  string `json:"str2" validate:"required,min=1,max=100"`
}

type FizzBuzzResponse struct {
	Result []string `json:"result"`
	Count  int      `json:"count"`
}

type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}