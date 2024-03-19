package main

import (
	"bufio"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"
)

// getImageFromRequest parse image from request and do validate.
// return image type and image bytes
func getImageFromRequest(c echo.Context, imageFormName string, acceptedImageTypes []string) (string, []byte, error) {
	// Check form post.
	formFileHeader, err := c.FormFile(imageFormName)
	if err != nil {
		return "", nil, c.JSON(http.StatusUnprocessableEntity, HTTPErrorResponse{
			Message: ErrFileNotFound.Error(),
		})
	}

	formFileObject, err := formFileHeader.Open()
	if err != nil {
		return "", nil, c.JSON(http.StatusUnprocessableEntity, HTTPErrorResponse{
			Message: ErrOpenFile.Error(),
		})
	}
	defer formFileObject.Close()

	// Check if the uploaded file is not larger than 1MB.
	if formFileHeader.Size > int64(1*1024*1024) {
		return "", nil, c.JSON(http.StatusUnprocessableEntity, HTTPErrorResponse{
			Message: fmt.Sprintf("%s. Submitted file is %d bytes", ErrFileLargerThan1MB.Error(), formFileHeader.Size),
		})
	}

	// Find the content type of the uploaded file.
	contentType, err := DetectContentType(formFileObject)
	if err != nil {
		return "", nil, c.JSON(http.StatusUnprocessableEntity, HTTPErrorResponse{
			Message: ErrOpenFile.Error(),
		})
	}

	// Check if the uploaded file is appropriate type.
	if !slices.Contains(acceptedImageTypes, contentType) {
		return "", nil, c.JSON(http.StatusUnprocessableEntity, HTTPErrorResponse{
			Message: fmt.Sprintf("%s. Submitted file is %s", ErrUnsupportedFile.Error(), contentType),
		})
	}

	// Reset file seek.
	if _, err = formFileObject.Seek(0, io.SeekStart); err != nil {
		log.Printf("handler reset file seek: %#v\n", err)
		return "", nil, c.JSON(http.StatusInternalServerError, HTTPErrorResponse{
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}

	// Read file content.
	byteBuffer := make([]byte, formFileHeader.Size)
	_, err = bufio.NewReader(formFileObject).Read(byteBuffer)
	if err != nil && err != io.EOF {
		log.Printf("handler read file: %#v\n", err)
		return "", nil, c.JSON(http.StatusInternalServerError, HTTPErrorResponse{
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}

	return contentType, byteBuffer, nil
}

// convertHandler Convert image from PNG to JPEG.
func convertHandler(c echo.Context) error {
	imageFormName := "file"
	acceptedImageTypes := []string{"image/png"}
	outputImageJPEGQuality := 95

	// Capture form data.
	imageType, imageFile, _ := getImageFromRequest(c, imageFormName, acceptedImageTypes)
	if imageFile == nil {
		return nil
	}

	// Init image processing.
	imaging, err := NewImaging(imageType, imageFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, HTTPErrorResponse{
			Message: err.Error(),
		})
	}

	// Convert image to JPEG.
	output, err := imaging.EncodeToJPEG(outputImageJPEGQuality)
	if err != nil {
		log.Printf("handler convert image encode: %#v\n", err)
		return c.JSON(http.StatusInternalServerError, HTTPErrorResponse{
			Message: ErrFailedToExportImage.Error(),
		})
	}

	// Write image to response.
	c.Response().Header().Set("Content-Type", "image/jpeg")
	c.Response().WriteHeader(200)
	c.Response().Write(output)
	return nil
}

// resizeHandler Resize image to specified dimensions.
func resizeHandler(c echo.Context) error {
	imageFormName := "file"
	acceptedImageTypes := []string{"image/png", "image/jpeg"}
	outputImageJPEGQuality := 95
	outputImagePNGQuality := 3

	// Capture form data.
	resizeWidth, _ := strconv.Atoi(c.FormValue("width"))
	resizeHeight, _ := strconv.Atoi(c.FormValue("height"))
	imageType, imageFile, _ := getImageFromRequest(c, imageFormName, acceptedImageTypes)
	if imageFile == nil {
		return nil
	}

	// Init image processing.
	imaging, err := NewImaging(imageType, imageFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, HTTPErrorResponse{
			Message: err.Error(),
		})
	}

	// Resize image.
	if err = imaging.Resize(resizeWidth, resizeHeight); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, HTTPErrorResponse{
			Message: err.Error(),
		})
	}

	// Encode image
	var output []byte
	switch imageType {
	case "image/png":
		output, err = imaging.EncodeToPNG(outputImagePNGQuality)
	case "image/jpeg":
		fallthrough
	default:
		output, err = imaging.EncodeToJPEG(outputImageJPEGQuality)
	}
	if err != nil {
		log.Printf("handler convert image encode: %#v\n", err)
		return c.JSON(http.StatusInternalServerError, HTTPErrorResponse{
			Message: ErrImageProcessingFailed.Error(),
		})
	}

	// Write image to response.
	c.Response().Header().Set("Content-Type", imageType)
	c.Response().WriteHeader(200)
	c.Response().Write(output)
	return nil
}

// compressHandler Compress image.
func compressHandler(c echo.Context) error {
	imageFormName := "file"
	acceptedImageTypes := []string{"image/png", "image/jpeg"}
	outputImageJPEGQuality := 80
	outputImagePNGQuality := 9

	// Capture form data.
	imageType, imageFile, _ := getImageFromRequest(c, imageFormName, acceptedImageTypes)
	if imageFile == nil {
		return nil
	}

	// Init image processing.
	imaging, err := NewImaging(imageType, imageFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, HTTPErrorResponse{
			Message: err.Error(),
		})
	}

	// Encode image.
	var output []byte
	switch imageType {
	case "image/png":
		output, err = imaging.EncodeToPNG(outputImagePNGQuality)
	case "image/jpeg":
		fallthrough
	default:
		output, err = imaging.EncodeToJPEG(outputImageJPEGQuality)
	}
	if err != nil {
		log.Printf("handler convert image encode: %#v\n", err)
		return c.JSON(http.StatusInternalServerError, HTTPErrorResponse{
			Message: ErrImageProcessingFailed.Error(),
		})
	}

	// Write image to response.
	c.Response().Header().Set("Content-Type", imageType)
	c.Response().WriteHeader(200)
	c.Response().Write(output)
	return nil
}
