package main

import (
	"bufio"
	"fmt"
	"mime/multipart"
	"net/http"
)

type HTTPErrorResponse struct {
	Message string `json:"message"`
}

// DetectContentType detects the content type of the uploaded file.
func DetectContentType(file multipart.File) (string, error) {
	peekBuff := bufio.NewReader(file)
	peekSniff, err := peekBuff.Peek(512)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrOpenFile, err)
	}

	return http.DetectContentType(peekSniff), nil
}
