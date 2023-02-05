package utils

import (
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

func CreateUpToken(accessKey string, secretKey string, bucket string) string {
	// 上传策略
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	// 上传凭证
	upToken := putPolicy.UploadToken(mac)
	return upToken
}
