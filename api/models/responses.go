package models

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseId struct {
	Message int64 `json:"message"`
}

type ResponseOK struct {
	Message string `json:"message"`
}