package main

import (
	"os"
	"testing"
)

func TestImaging_EncodeToJPEG(t *testing.T) {
	imageFile, err := os.ReadFile("testdata/SampleSmallJPGImage.jpg")
	if err != nil {
		t.Fatalf("failed to open sample image file: %v", err)
	}

	imaging, err := NewImaging("image/jpeg", imageFile)
	if err != nil {
		t.Fatalf("failed to imaging sample image file: %v", err)
	}

	type args struct {
		quality int
	}
	tests := []struct {
		name    string
		fields  *Imaging
		args    args
		wantErr bool
	}{
		{name: "Should pass encode to JPG", fields: imaging, args: args{quality: 70}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := tt.fields

			_, err = i.EncodeToJPEG(tt.args.quality)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeToJPEG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestImaging_EncodeToPNG(t *testing.T) {
	imageFile, err := os.ReadFile("testdata/SampleSmallJPGImage.jpg")
	if err != nil {
		t.Fatalf("failed to open sample image file: %v", err)
	}

	imaging, err := NewImaging("image/jpeg", imageFile)
	if err != nil {
		t.Fatalf("failed to imaging sample image file: %v", err)
	}

	type args struct {
		quality int
	}
	tests := []struct {
		name    string
		fields  *Imaging
		args    args
		wantErr bool
	}{
		{name: "Should pass encode to PNG", fields: imaging, args: args{quality: 9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := tt.fields
			_, err = i.EncodeToPNG(tt.args.quality)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeToPNG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestImaging_Resize(t *testing.T) {
	imageFile, err := os.ReadFile("testdata/SampleSmallJPGImage.jpg")
	if err != nil {
		t.Fatalf("failed to open sample image file: %v", err)
	}

	imaging, err := NewImaging("image/jpeg", imageFile)
	if err != nil {
		t.Fatalf("failed to imaging sample image file: %v", err)
	}

	type args struct {
		width  int
		height int
	}
	tests := []struct {
		name    string
		fields  *Imaging
		args    args
		wantErr bool
	}{
		{name: "Should pass resize image", fields: imaging, args: args{width: 100, height: 100}},
		{name: "Width too small expect error", fields: imaging, args: args{width: 9, height: 100}, wantErr: true},
		{name: "Height too small expect error", fields: imaging, args: args{width: 100, height: 9}, wantErr: true},
		{name: "Dimension too small expect error", fields: imaging, args: args{width: 9, height: 9}, wantErr: true},
		{name: "Width too large expect error", fields: imaging, args: args{width: 1001, height: 100}, wantErr: true},
		{name: "Height too large expect error", fields: imaging, args: args{width: 100, height: 1001}, wantErr: true},
		{name: "Dimension too large expect error", fields: imaging, args: args{width: 1001, height: 1001}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := tt.fields
			if err = i.Resize(tt.args.width, tt.args.height); (err != nil) != tt.wantErr {
				t.Errorf("Resize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewImaging(t *testing.T) {
	imageFile, err := os.ReadFile("testdata/SampleSmallJPGImage.jpg")
	if err != nil {
		t.Fatalf("failed to open sample image file: %v", err)
	}

	type args struct {
		imageType string
		imageFile []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Should pass NewImaging", args: args{"image/jpeg", imageFile}},
		{name: "Unsupported image type expect error", args: args{"image/gif", imageFile}, wantErr: true},
		{name: "Empty image file expect error", args: args{"image/png", nil}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = NewImaging(tt.args.imageType, tt.args.imageFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewImaging() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
