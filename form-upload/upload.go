package main

import (
	"bytes"
	"fmt"
	"io"
	"kodo-up/utils"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func createReqBody(data io.Reader, fileName string, upToken string, key string) (string, io.Reader) {
	var buf bytes.Buffer
	bw := multipart.NewWriter(&buf)

	// key（上传后的目录以及文件名）
	part1, _ := bw.CreateFormField("key")
	part1.Write([]byte(key))

	// upToken（上传凭证）
	part2, _ := bw.CreateFormField("token")
	part2.Write([]byte(upToken))

	// file（上传的文件）
	part3, _ := bw.CreateFormFile("file", fileName)
	io.Copy(part3, data)

	bw.Close()
	return bw.FormDataContentType(), &buf
}

// 文件上传
func doUploadFile(addr string, filePath string, upToken string, key string) error {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("文件打开失败")
		return err
	}
	defer f.Close()
	contType, reader := createReqBody(f, filepath.Base(filePath), upToken, key)
	url := fmt.Sprintf("http://%s", addr)
	req, _ := http.NewRequest("POST", url, reader)
	req.Header.Add("Content-Type", contType)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败", err)
		return err
	}
	resp.Body.Close()
	return nil
}

// 字节数组上传
func doUploadBytes(addr string, data io.Reader, upToken string, key string) error {
	contType, reader := createReqBody(data, path.Base(key), upToken, key)
	url := fmt.Sprintf("http://%s", addr)
	req, _ := http.NewRequest("POST", url, reader)
	req.Header.Add("Content-Type", contType)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败", err)
		return err
	}
	resp.Body.Close()
	return nil
}

func main() {
	addr := "up.qiniup.com"
	filePath := "/Users/aanda/Desktop/Images/1.jpg"
	key := "Test-1"
	accessKey := ""
	secretKey := ""
	bucket := "aanda-test"

	upToken := utils.CreateUpToken(accessKey, secretKey, bucket)
	err1 := doUploadFile(addr, filePath, upToken, key)

	key = "Test-2"
	data := []byte("hello, this is qiniu cloud")
	err2 := doUploadBytes(addr, bytes.NewReader(data), upToken, key)

	if err1 != nil || err2 != nil {
		fmt.Println("上传失败")
	} else {
		fmt.Println("上传成功")
	}
}
