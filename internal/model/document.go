package model

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime/multipart"
)

type Document struct {
	Id                  int          `json:"id"`
	Type                string       `json:"type"`
	FriendlyDescription string       `json:"friendlyDescription"`
	CreatedDate         string       `json:"createdDate"`
	Direction           string       `json:"direction"`
	Filename            string       `json:"filename"`
	CreatedBy           User         `json:"createdBy"`
	ReceivedDateTime    string       `json:"receivedDateTime"`
	Note                DocumentNote `json:"note"`
	ReformattedTime     string
	Infected            bool `json:"infected"`
}

type DocumentNote struct {
	Description string `json:"description"`
	Name        string `json:"name"`
}

func EncodeFileToBase64(file multipart.File) (string, error) {
	defer func() {
		_ = file.Close()
	}()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}

	err := file.Close()
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
