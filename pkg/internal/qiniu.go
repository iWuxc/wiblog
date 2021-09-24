// Package internal provides ...
package internal

import (
	"context"
	"errors"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	"path/filepath"
	"wiblog/pkg/conf"
)

type UploadParams struct {
	Name string
	Size int64
	Data io.Reader

	Config conf.Qiniu
}

// QiniuUpload 七牛云上传
func QiniuUpload(params UploadParams) (string, error) {
	if params.Config.AccessKey == "" || params.Config.SecretKey == "" {
		return "", errors.New("qiniu config error")
	}
	key := completeQiniuKey(params.Name)

	//设置上传策略
	putPolicy := &storage.PutPolicy{
		Scope:      params.Config.Bucket,
		Expires:    3600,
		InsertOnly: 1,
	}
	mac := qbox.NewMac(params.Config.AccessKey, params.Config.SecretKey)
	//上传token
	uploadToken := putPolicy.UploadToken(mac)
	//上传配置
	cfg := &storage.Config{
		UseHTTPS: false,
		Zone:     &storage.ZoneHuabei,
	}
	//uploader
	uploader := storage.NewFormUploader(cfg)

	ret := storage.PutRet{}
	putExtra := &storage.PutExtra{}

	err := uploader.Put(context.Background(), &ret, uploadToken, key,
		params.Data, params.Size, putExtra)
	if err != nil {
		return "", err
	}
	url := "https://" + params.Config.Domain + "/" + key
	return url, nil
}

type DeleteParams struct {
	Name string
	Days int

	Config conf.Qiniu
}

func QiniuDelete(params DeleteParams) error {

	key := completeQiniuKey(params.Name)

	mac := qbox.NewMac(params.Config.AccessKey, params.Config.SecretKey)

	//上传配置
	cfg := &storage.Config{
		UseHTTPS: false,
		Zone:     &storage.ZoneHuabei,
	}

	// manager
	bucketManager := storage.NewBucketManager(mac, cfg)

	if params.Days > 0 {
		return bucketManager.DeleteAfterDays(params.Config.Bucket, key, params.Days)
	}
	return bucketManager.Delete(params.Config.Bucket, key)

}

// completeQiniuKey 修复路径
func completeQiniuKey(name string) string {
	ext := filepath.Ext(name)

	switch ext {
	case ".bmp", ".png", ".jpg",
		".gif", ".ico", ".jpeg":

		name = "blog/img/" + name
	case ".mov", ".mp4":
		name = "blog/video/" + name
	case ".go", ".js", ".css",
		".cpp", ".php", ".rb",
		".java", ".py", ".sql",
		".lua", ".html", ".sh",
		".xml", ".cs":

		name = "blog/code/" + name
	case ".txt", ".md", ".ini",
		".yaml", ".yml", ".doc",
		".ppt", ".pdf":

		name = "blog/document/" + name
	case ".zip", ".rar", ".tar",
		".gz":

		name = "blog/archive/" + name
	default:
		name = "blog/other/" + name
	}
	return name
}
