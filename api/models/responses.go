package models

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseOK struct {
	Message interface{} `json:"message"`
}