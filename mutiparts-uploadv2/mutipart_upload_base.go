package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type MutiPartUploadInput struct {
	Host      string
	Token     string
	Bucket    string
	Key       string
	KeyEncode bool
	UploadId  string
}

type UploadPartInput struct {
	MutiPartUploadInput
	PartNum int
	Body    []byte
}

type CompleteMultipartInput struct {
	Parts      []Part            `json:"parts"`
	Fname      string            `json:"fname"`
	MimeType   string            `json:"mimeType"`
	Metadata   map[string]string `json:"metadata"`
	CustomVars map[string]string `json:"customVars"`
}

type Part struct {
	ETag    string `json:"etag"`
	PartNum int    `json:"partNumber"`
}

type AbortMultipartInput struct {
	MutiPartUploadInput
}

type ListPartsInput struct {
	MutiPartUploadInput
	MaxParts int
	Offset   int
}

func (m *MutiPartUploadInput) initiateMultipartUpload() (*http.Response, error) {
	var key string

	if m.KeyEncode {
		key = base64.URLEncoding.EncodeToString([]byte(m.Key))
	} else {
		key = m.Key
	}

	url := fmt.Sprintf("http://%s/buckets/%s/objects/%s/uploads", m.Host, m.Bucket, key)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", "UpToken "+m.Token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("init request fail", err)
	}
	return resp, err
}

func (u *UploadPartInput) uploadPart() (*http.Response, error) {
	var key string

	if u.KeyEncode {
		key = base64.URLEncoding.EncodeToString([]byte(u.Key))
	} else {
		key = u.Key
	}

	url := fmt.Sprintf("http://%s/buckets/%s/objects/%s/uploads/%s/%d", u.Host, u.Bucket, key, u.UploadId, u.PartNum)

	reader := bytes.NewReader(u.Body)
	req, _ := http.NewRequest("PUT", url, reader)
	req.Header.Add("Authorization", "UpToken "+u.Token)
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Content-Length", strconv.Itoa(len(u.Body)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("upload request fail")
	}
	return resp, err
}

func (c *CompleteMultipartInput) completeMultipartUpload(m *MutiPartUploadInput) (*http.Response, error) {
	var key string

	if m.KeyEncode {
		key = base64.URLEncoding.EncodeToString([]byte(m.Key))
	} else {
		key = m.Key
	}

	url := fmt.Sprintf("http://%s/buckets/%s/objects/%s/uploads/%s", m.Host, m.Bucket, key, m.UploadId)

	bytesData, _ := json.Marshal(c)
	reader := bytes.NewReader(bytesData)
	req, _ := http.NewRequest("POST", url, reader)

	req.Header.Add("Authorization", "UpToken "+m.Token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("complete request fail")
	}
	return resp, err
}

func (a *AbortMultipartInput) abortMultipartUpload() error {
	var key string

	if a.KeyEncode {
		key = base64.URLEncoding.EncodeToString([]byte(a.Key))
	} else {
		key = a.Key
	}

	url := fmt.Sprintf("http://%s/buckets/%s/objects/%s/uploads/%s", a.Host, a.Bucket, key, a.UploadId)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", "UpToken "+a.Token)

	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		fmt.Println("abort request fail")
	}
	return err
}

func (l *ListPartsInput) ListParts() (*http.Response, error) {
	var key string

	if l.KeyEncode {
		key = base64.URLEncoding.EncodeToString([]byte(l.Key))
	} else {
		key = l.Key
	}

	url := fmt.Sprintf("http://%s/buckets/%s/objects/%s/uploads/%s?max-parts=%d&part-number-marker=%d", l.Host, l.Bucket, key, l.UploadId, l.MaxParts, l.Offset)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "UpToken "+l.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("list request fail")
	}
	return resp, err
}
