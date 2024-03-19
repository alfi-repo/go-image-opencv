package main

import (
	"errors"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"slices"
)

const (
	minResize = 10
	maxResize = 1000
)

var (
	ErrFailedToReadImage      = errors.New("failed to read image")
	ErrFailedToExportImage    = errors.New("failed to export image")
	ErrInvalidResizeDimension = fmt.Errorf("invalid resize dimension. Must be between %d and %d", minResize, maxResize)
	ErrUnknownImageType       = errors.New("unknown image type")
)

type Imaging struct {
	// Image mime type. i.e. image/png, image/jpeg
	Type string
	// OpenCV Mat
	Mat gocv.Mat
}

// Resize image to specified width and height.
// width, height: 10-1000
func (i *Imaging) Resize(width, height int) error {
	if width < minResize || width > maxResize || height < minResize || height > maxResize {
		return ErrInvalidResizeDimension
	}

	gocv.Resize(i.Mat, &i.Mat, image.Point{X: width, Y: height}, 0, 0, gocv.InterpolationDefault)
	return nil
}

// EncodeToJPEG write to JPEG.
// quality: 0-100
func (i *Imaging) EncodeToJPEG(quality int) ([]byte, error) {
	buffer, err := gocv.IMEncodeWithParams(gocv.JPEGFileExt, i.Mat, []int{gocv.IMWriteJpegQuality, quality})
	if err != nil || buffer == nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToExportImage, err)
	}
	defer buffer.Close()

	// Prevent segmentation fault error: https://github.com/hybridgroup/gocv/issues/1005
	newBuffer := make([]byte, buffer.Len())
	copy(newBuffer, buffer.GetBytes())
	return newBuffer, nil
}

// EncodeToPNG write to PNG.
// quality: 0-9
func (i *Imaging) EncodeToPNG(quality int) ([]byte, error) {
	buffer, err := gocv.IMEncodeWithParams(gocv.PNGFileExt, i.Mat, []int{gocv.IMWritePngCompression, quality})
	if err != nil || buffer == nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToExportImage, err)
	}
	defer buffer.Close()

	// Prevent segmentation fault error: https://github.com/hybridgroup/gocv/issues/1005
	newBuffer := make([]byte, buffer.Len())
	copy(newBuffer, buffer.GetBytes())
	return newBuffer, nil
}

// NewImaging creates new Imaging object for processing.
func NewImaging(imageType string, imageFile []byte) (*Imaging, error) {
	if !slices.Contains([]string{"image/png", "image/jpeg"}, imageType) {
		return nil, fmt.Errorf("%w: %s", ErrUnknownImageType, imageType)
	}

	decode, err := gocv.IMDecode(imageFile, gocv.IMReadUnchanged)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToReadImage, err)
	}

	return &Imaging{
		Type: imageType,
		Mat:  decode,
	}, nil
}
