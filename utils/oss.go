package utils

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

const (
	OSSDriver   = "oss"
	LocalDriver = "local"
)

type UploadConf struct {
	Filter          string
	FileMaxSize     int64 // 上传文件的最大尺寸 单位：Mb  0:不限大小
	Driver          string
	UploadConfLocal *UploadConfLocal
	UploadConfOSS   *UploadConfOSS
}

var UploadConfig *UploadConf

func InitUploadConf() *UploadConf {
	var conf UploadConf
	conf.UploadConfLocal = InitUploadConfLocal()
	conf.UploadConfOSS = InitUploadConfOSS()
	conf.Filter = ".jpg,.jpeg,.png,.bmp"
	conf.FileMaxSize = 50
	conf.Driver = OSSDriver
	UploadConfig = &conf
	return &conf
}

type UploadConfLocal struct {
	Host string // 根域名
}

func InitUploadConfLocal() *UploadConfLocal {
	var conf UploadConfLocal
	conf.Host = viper.GetString("serve.host")
	return &conf
}

type UploadConfOSS struct {
	Host            string // 根域名
	Endpoint        string // 节点
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

func InitUploadConfOSS() *UploadConfOSS {
	var conf UploadConfOSS
	conf.Endpoint = viper.GetString("oss.endpoint")
	conf.BucketName = viper.GetString("oss.bucketname")
	conf.Host = "https://" + conf.BucketName + "." + conf.Endpoint
	conf.AccessKeyID = viper.GetString("oss.accessKeyID")
	conf.AccessKeySecret = viper.GetString("oss.accessKeySecret")
	return &conf
}

type UploadResult struct {
	Host         string
	RelativePath string
}

// UploadFile 上传单个文件到本地目录
func UploadFile(ctx *gin.Context, formName string) (result UploadResult, err error) {
	file, err := ctx.FormFile(formName)
	if err != nil {
		err = errors.New("获取文件数据失败")
		return
	}

	return actionUploadFile(ctx, file)
}

func UploadMultipartFile(ctx *gin.Context, formName string) (resultSlice []UploadResult, err error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		err = errors.New(fmt.Sprintf("获取图片失败， err:%v\n", err))
		return
	}
	files := form.File[formName]

	for _, file := range files {
		result, err := actionUploadFile(ctx, file)
		if err != nil {
			err = errors.New("上传失败： [error]:" + err.Error())
			return nil, err
		}
		resultSlice = append(resultSlice, result)
	}

	return
}

func actionUploadFile(ctx *gin.Context, file *multipart.FileHeader) (result UploadResult, err error) {
	uuid, err := GetUUID()
	file.Filename = path.Join(uuid + RandomCode(5) + ".png")

	// 限制上传文件大小
	if UploadConfig.FileMaxSize != 0 {
		if file.Size > UploadConfig.FileMaxSize*1024*1024 {
			fmt.Println(file.Size)
			err = errors.New(fmt.Sprintf("上传文件尺寸过大，只能上传小于：%dMb的文件\n", UploadConfig.FileMaxSize))
			return
		}
	}
	// 限制上传文件的格式
	suffix := strings.ToLower(path.Ext(file.Filename))
	if suffix == "" {
		err = errors.New(fmt.Sprintf("上传文件格式错误，只能上传格式为%s的文件\n", UploadConfig.Filter))
		return
	}

	if UploadConfig.Filter != "" {
		filterSlice := strings.Split(UploadConfig.Filter, ",")
		flag := false
		for _, v := range filterSlice {
			if v == suffix {
				flag = true
				break
			}
		}
		if !flag {
			err = errors.New(fmt.Sprintf("上传文件格式错误，只能上传格式为%s的文件\n", UploadConfig.Filter))
			return
		}
	}
	switch UploadConfig.Driver {
	case LocalDriver:
		return uploadFileToLocal(ctx, file)
	case OSSDriver:
		return uploadFileToOSS(ctx, file)
	default:
		return result, errors.New(fmt.Sprintf("未选择上传引擎，请前往配置文件进行配置\n"))
	}
}

func uploadFileToLocal(ctx *gin.Context, file *multipart.FileHeader) (result UploadResult, err error) {
	result.Host = UploadConfig.UploadConfLocal.Host
	filepath := path.Join("static/uploadfile", time.Now().Format("20060102"))
	result.RelativePath = path.Join(filepath, file.Filename)
	// 文件目录是否存在， 不存在创建
	if _, err := os.Stat(filepath); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(filepath, os.ModePerm)
		}
	}
	err = ctx.SaveUploadedFile(file, result.RelativePath)
	if err != nil {
		{
			err = errors.New(fmt.Sprintf("上传失败，%v", err))
			return
		}
	}
	return
}

func uploadFileToOSS(ctx *gin.Context, file *multipart.FileHeader) (result UploadResult, err error) {
	result.Host = UploadConfig.UploadConfOSS.Host

	localResult, err := uploadFileToLocal(ctx, file)
	if err != nil {
		return
	}

	// oss 连接
	ossPath := path.Join("upload", time.Now().Format("20060102"))
	result.RelativePath = path.Join(ossPath, file.Filename)
	fmt.Sprintf("endpoint %s id %s secret %s", UploadConfig.UploadConfOSS.Endpoint, UploadConfig.UploadConfOSS.AccessKeyID, UploadConfig.UploadConfOSS.AccessKeySecret)
	client, err := oss.New(UploadConfig.UploadConfOSS.Endpoint, UploadConfig.UploadConfOSS.AccessKeyID, UploadConfig.UploadConfOSS.AccessKeySecret)
	if err != nil {
		os.Remove(localResult.RelativePath)
		err = errors.New("文件上传云端失败1")
	}
	bucket, err := client.Bucket(UploadConfig.UploadConfOSS.BucketName)
	if err != nil {
		os.Remove(localResult.RelativePath)
		err = errors.New("文件上传云端失败2")
	}

	// 上传云端
	err = bucket.PutObjectFromFile(result.RelativePath, localResult.RelativePath)
	if err != nil {
		os.Remove(localResult.RelativePath)
		err = errors.New("文件上传云端失败3")
	}
	return
}
