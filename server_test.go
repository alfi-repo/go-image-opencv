package main

import (
	"mime/multipart"
	"os"
	"testing"
)

func TestDetectContentType(t *testing.T) {
	imageFile, err := os.Open("testdata/SampleSmallJPGImage.jpg")
	if err != nil {
		t.Fatalf("failed to open sample image file: %v", err)
	}
	defer imageFile.Close()

	emptyImageFile, _ := os.Open("testdata/missing_file.jpg")

	type args struct {
		file multipart.File
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "Should pass and return image/jpeg", args: args{file: imageFile}, want: "image/jpeg"},
		{name: "Empty file expect error", args: args{file: emptyImageFile}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectContentType(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectContentType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DetectContentType() got = %v, want %v", got, tt.want)
			}
		})
	}
}
