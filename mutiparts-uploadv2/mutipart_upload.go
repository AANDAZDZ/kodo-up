package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kodo-up/utils"
	"net/http"
	"os"
	"strconv"
)

func Transform(resp *http.Response) map[string]interface{} {
	var result map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(string(body)), &result)
	return result
}

func getBody(filePath string) []byte {
	file, _ := os.Open(filePath)
	defer file.Close()
	body, _ := ioutil.ReadAll(file)
	return body
}

func main() {
	host := "up.qiniup.com"
	bucket := "aanda-test"
	key := "Test-11"
	accessKey := ""
	secretKey := ""
	token := utils.CreateUpToken(accessKey, secretKey, bucket)
	filePath := "/Users/aanda/Desktop/Images/1.jpg"

	body := getBody(filePath)

	// init
	mutiPartUploadInput := MutiPartUploadInput{
		Host:      host,
		Bucket:    bucket,
		Key:       key,
		KeyEncode: true,
		Token:     token,
	}
	resp1, _ := mutiPartUploadInput.initiateMultipartUpload()
	mutiPartUploadInput.UploadId = fmt.Sprintf("%v", Transform(resp1)["uploadId"])

	// up
	uploadPartInput := UploadPartInput{
		MutiPartUploadInput: mutiPartUploadInput,
		PartNum:             1,
		Body:                body,
	}
	resp2, err := uploadPartInput.uploadPart()
	eTag := fmt.Sprintf("%v", Transform(resp2)["etag"])
	if err == nil {
		fmt.Println("上传分片：" + strconv.Itoa(uploadPartInput.PartNum) + " 成功, " + "etag: " + eTag)
	}

	// complete
	completeMultipart := CompleteMultipartInput{
		Parts: []Part{{ETag: eTag, PartNum: uploadPartInput.PartNum}},
	}
	_, err1 := completeMultipart.completeMultipartUpload(&mutiPartUploadInput)
	if err1 == nil {
		fmt.Println("合并分片成功")
	}

	// list
	listPartsInput := ListPartsInput{
		MutiPartUploadInput: mutiPartUploadInput,
		MaxParts:            5,
		Offset:              0,
	}
	resp3, err2 := listPartsInput.ListParts()
	if err2 == nil {
		fmt.Println(Transform(resp3)["parts"])
	}

	// abort
	abortMultipartInput := AbortMultipartInput{
		MutiPartUploadInput: mutiPartUploadInput,
	}
	err3 := abortMultipartInput.abortMultipartUpload()
	if err3 == nil {
		fmt.Println("中止成功")
	}
}
