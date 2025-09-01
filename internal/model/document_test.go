package model

import (
	"bytes"
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"strings"
	"testing"
)

func TestEncodeFileToBase64(t *testing.T) {
	var buff bytes.Buffer

	formWriter := multipart.NewWriter(io.Writer(&buff))
	file, _ := formWriter.CreateFormFile("document-upload", "data.txt")
	_, _ = io.Copy(file, strings.NewReader("test-string"))

	_ = formWriter.Close()

	buffReader := bytes.NewReader(buff.Bytes())
	formReader := multipart.NewReader(buffReader, formWriter.Boundary())
	multipartForm, _ := formReader.ReadForm(1 << 20)

	multipartFiles := multipartForm.File["document-upload"]
	multipartFile, _ := multipartFiles[0].Open()

	base64File, err := EncodeFileToBase64(multipartFile)

	assert.Nil(t, err)
	assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("test-string")), base64File)
}
