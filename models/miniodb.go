package models

import (
	"github.com/githubchry/gomdb/drivers"
	"bytes"
	"github.com/minio/minio-go"
	"log"
	"net/url"
	"time"
)

//创建存储桶
//location := "us-east-1"
//bucketNameArr := [...]string{"music", "photo"}
func MakeBucket(bucketNameArr []string, location string) error {
	//range遍历数组
	for _, bucketName := range bucketNameArr {
		// 创建存储桶。
		err := drivers.MinioDbConn.MakeBucket(bucketName, location)
		if err != nil {
			// 检查存储桶是否已经存在。
			exists, err := drivers.MinioDbConn.BucketExists(bucketName)
			if err == nil && exists {
				log.Printf("We already own %s\n", bucketName)
			} else {
				log.Print(err)
				return err
			}
		} else {
			log.Printf("Successfully created %s\n", bucketName)
		}
	}

	return nil
}

//删除存储桶
func RemoveBucket(bucketNameArr []string) error {
	//range遍历数组
	for _, bucketName := range bucketNameArr {
		// 创建存储桶。
		err := drivers.MinioDbConn.RemoveBucket(bucketName)
		if err != nil {
			// 检查存储桶是否已经存在。
			exists, err := drivers.MinioDbConn.BucketExists(bucketName)
			if err == nil && !exists {
				log.Printf("%s does not exist!\n", bucketName)
			} else {
				log.Print(err)
				return err
			}
		} else {
			log.Printf("Successfully remove %s\n", bucketName)
		}
	}

	return nil
}


// 获取上传url
func PreUpload(bucketName string, fileName string) string {

	presignedURL, err := drivers.MinioDbConn.PresignedPutObject(bucketName, fileName, time.Second * 24 * 60 * 60)
	if err != nil {
		log.Println(err)
	}
	//log.Println("Successfully generated presigned URL", presignedURL)
	return presignedURL.String()
}

// 获取下载url
func PreDownload(bucketName string, fileName string) string {
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+fileName)

	presignedURL, err := drivers.MinioDbConn.PresignedGetObject(bucketName, fileName, time.Second * 60 * 2, reqParams)
	if err != nil {
		log.Println(err)
	}
	//log.Println("Successfully generated presigned URL", presignedURL)
	return presignedURL.String()
}

// 上传本地文件 //"image/png"   "image/jpg"
func FPutObject(bucketName, objName, filePath, contentType string) error {
	//从路径获取文件名: filepath.Base(files)   获取文件后缀path.Ext(files)
	uploadInfo, err := drivers.MinioDbConn.FPutObject(bucketName, objName, filePath, minio.PutObjectOptions{ContentType: contentType});
	if err != nil {
		return err
	}
	log.Println("Successfully uploaded file object: ", uploadInfo)
	return nil
}


// 下载文件到本地
func FGetObject(bucketName, objName, filePath string) error {
	err := drivers.MinioDbConn.FGetObject(bucketName, objName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}


// 上传内存文件 //"image/png"   "image/jpg"
func PutObject(bucketName, objName, contentType string, obj []byte) error {
	reader := bytes.NewReader(obj)
	uploadInfo, err := drivers.MinioDbConn.PutObject(bucketName, objName, reader, reader.Size(), minio.PutObjectOptions{ContentType: contentType});
	if err != nil {
		return err
	}
	log.Println("Successfully uploaded object: ", uploadInfo)
	return nil
}


// 下载文件到内存
func GetObject(bucketName, objName string) ([]byte, error) {
	object, err := drivers.MinioDbConn.GetObject(bucketName, objName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	stat, _ := object.Stat()
	obj := make([]byte, stat.Size)
	object.Read(obj)

	return obj, nil
}

// 删除文件对象
func RemoveObject(bucketName, objName string) error {
	err := drivers.MinioDbConn.RemoveObject(bucketName, objName)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// 删除多个文件对象
func RemoveObjects(bucketName string, objectsCh <-chan string) {

	for rErr := range drivers.MinioDbConn.RemoveObjects(bucketName, objectsCh) {
		log.Println("Error detected during deletion: ", rErr)
	}
}
