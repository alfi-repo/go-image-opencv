package main

import "errors"

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOpenFile              = errors.New("failed to open file")
	ErrFileNotFound          = errors.New("file not found")
	ErrFileLargerThan1MB     = errors.New("file larger than 1MB")
	ErrImageProcessingFailed = errors.New("image processing failed")
)
