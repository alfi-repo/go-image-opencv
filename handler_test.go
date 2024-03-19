package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name             string
		testFile         string
		wantStatusCode   int
		wantResponseType string
	}{
		{name: "Should pass and return image/jpeg", testFile: "SamplePNGImage_100kbmb.png", wantStatusCode: http.StatusOK, wantResponseType: "image/jpeg"},
		{name: "Sent JPEG file expect error", testFile: "SampleJPGImage_100kbmb.jpg", wantStatusCode: http.StatusUnprocessableEntity, wantResponseType: "application/json"},
		{name: "Sent large PNG file expect error", testFile: "SamplePNGImage_2mbmb.png", wantStatusCode: http.StatusUnprocessableEntity, wantResponseType: "application/json"},
		{name: "Sent empty file data expect error", testFile: "", wantStatusCode: http.StatusUnprocessableEntity, wantResponseType: "application/json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				osFile *os.File
				err    error
			)

			if tt.testFile != "" {
				osFile, err = os.Open("testdata/" + tt.testFile)
				if err != nil {
					t.Fatalf("failed to open sample image file: %v", err)
				}
				defer osFile.Close()
			}

			e := echo.New()
			bodyBuffer := new(bytes.Buffer)
			mw := multipart.NewWriter(bodyBuffer)
			part, err := mw.CreateFormFile("file", "image.jpg")
			if err != nil {
				t.Fatalf("failed to create mw writer: %v", err)
			}
			io.Copy(part, osFile)
			mw.Close()

			req := httptest.NewRequest(http.MethodPost, "/v1/convert", bodyBuffer)
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, mw.FormDataContentType())
			c := e.NewContext(req, rec)

			if assert.NoError(t, convertHandler(c)) {
				assert.Equal(t, tt.wantStatusCode, rec.Code)
				assert.Contains(t, rec.Result().Header.Get("Content-Type"), tt.wantResponseType)
			}
		})
	}
}

func TestResize(t *testing.T) {
	type args struct {
		width  int
		height int
	}
	tests := []struct {
		name             string
		testFile         string
		args             args
		wantStatusCode   int
		wantResponseType string
	}{
		{
			name:             "Should pass and return image/jpeg",
			testFile:         "SampleJPGImage_100kbmb.jpg",
			args:             args{width: 100, height: 100},
			wantStatusCode:   http.StatusOK,
			wantResponseType: "image/jpeg",
		},
		{
			name:             "Should pass and return image/png",
			testFile:         "SamplePNGImage_100kbmb.png",
			args:             args{width: 100, height: 100},
			wantStatusCode:   http.StatusOK,
			wantResponseType: "image/png",
		},
		{
			name:             "Dimension too small expect error",
			testFile:         "SamplePNGImage_100kbmb.png",
			args:             args{width: 1, height: 1},
			wantStatusCode:   http.StatusUnprocessableEntity,
			wantResponseType: "application/json",
		},
		{
			name:             "Dimension too large expect error",
			testFile:         "SamplePNGImage_100kbmb.png",
			args:             args{width: 10000, height: 10000},
			wantStatusCode:   http.StatusUnprocessableEntity,
			wantResponseType: "application/json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osFile, err := os.Open("testdata/" + tt.testFile)
			if err != nil {
				t.Fatalf("failed to open sample image file: %v", err)
			}
			defer osFile.Close()

			e := echo.New()
			bodyBuffer := new(bytes.Buffer)
			mw := multipart.NewWriter(bodyBuffer)
			mw.WriteField("width", strconv.Itoa(tt.args.width))
			mw.WriteField("height", strconv.Itoa(tt.args.height))
			part, err := mw.CreateFormFile("file", tt.testFile)
			if err != nil {
				t.Fatalf("failed to create mw writer: %v", err)
			}
			io.Copy(part, osFile)
			mw.Close()

			req := httptest.NewRequest(http.MethodPost, "/v1/resize", bodyBuffer)
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, mw.FormDataContentType())
			c := e.NewContext(req, rec)

			if assert.NoError(t, resizeHandler(c)) {
				assert.Equal(t, tt.wantStatusCode, rec.Code)
				assert.Contains(t, rec.Result().Header.Get("Content-Type"), tt.wantResponseType)
			}
		})
	}
}

func TestCompress(t *testing.T) {
	tests := []struct {
		name             string
		testFile         string
		wantStatusCode   int
		wantResponseType string
	}{
		{
			name:             "Should pass and return image/jpeg",
			testFile:         "SampleJPGImage_100kbmb.jpg",
			wantStatusCode:   http.StatusOK,
			wantResponseType: "image/jpeg",
		},
		{
			name:             "Should pass and return image/png",
			testFile:         "SamplePNGImage_100kbmb.png",
			wantStatusCode:   http.StatusOK,
			wantResponseType: "image/png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osFile, err := os.Open("testdata/" + tt.testFile)
			if err != nil {
				t.Fatalf("failed to open sample image file: %v", err)
			}
			defer osFile.Close()

			e := echo.New()
			bodyBuffer := new(bytes.Buffer)
			mw := multipart.NewWriter(bodyBuffer)
			part, err := mw.CreateFormFile("file", tt.testFile)
			if err != nil {
				t.Fatalf("failed to create mw writer: %v", err)
			}
			io.Copy(part, osFile)
			mw.Close()

			req := httptest.NewRequest(http.MethodPost, "/v1/compress", bodyBuffer)
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, mw.FormDataContentType())
			c := e.NewContext(req, rec)

			if assert.NoError(t, compressHandler(c)) {
				assert.Equal(t, tt.wantStatusCode, rec.Code)
				assert.Contains(t, rec.Result().Header.Get("Content-Type"), tt.wantResponseType)
			}
		})
	}
}
