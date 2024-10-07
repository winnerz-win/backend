package jawss3

import (
	"bytes"
	"fmt"
	"jtools/jparallel"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type JAwsS3Config struct {
	AwsAccessKey       string `yaml:"aws_access_key"`
	AwsSecretAccessKey string `yaml:"aws_secret_access_key"`
	AwsToken           string `yaml:"aws_token"`
	S3Region           string `yaml:"s3_region"`
	S3Name             string `yaml:"s3_name"`
	CdnPath            string `yaml:"cdn_path"`
}

type S3UploadParams []S3UploadParam
type S3UploadParam struct {
	S3Endpoint   string `json:"s3_endpoint"`
	FilePath     string `json:"file_path"`
	IsFileRemove bool   `json:"is_file_remove"`
}

type S3DeleteParams []S3DeleteParam
type S3DeleteParam struct {
	S3Path string `json:"s3_path"`
}

var isInit bool

var (
	client_credential *credentials.Credentials
	client_config     *aws.Config
	client            *s3.S3
)

var config JAwsS3Config

func GetS3Region() string {
	return config.S3Region
}
func GetS3Name() string {
	return config.S3Name
}
func GetCdnPath() string {
	return config.CdnPath
}

func InitS3(_config JAwsS3Config) error {
	isInit = false

	config = _config

	client_credential = credentials.NewStaticCredentials(config.AwsAccessKey, config.AwsSecretAccessKey, config.AwsToken)
	client_config = aws.NewConfig().WithRegion(config.S3Region).WithCredentials(client_credential)
	se, err := session.NewSession()
	if err != nil {
		return err
	}

	client = s3.New(se, client_config)
	isInit = true

	return nil
}

func S3Upload(param S3UploadParam) (*s3.PutObjectOutput, error) {
	if !isInit {
		return nil, fmt.Errorf("[S3Upload Fail] not call InitS3 - %v", param)
	}

	f, err := os.Open(param.FilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		f.Close()

		if param.IsFileRemove {
			os.Remove(param.FilePath)
		}
	}()

	fInfo, _ := f.Stat()
	if err != nil {
		return nil, err
	}

	size := fInfo.Size()
	buffer := make([]byte, size)
	f.Read(buffer)
	fBytes := bytes.NewReader(buffer)
	fType := http.DetectContentType(buffer)

	putParams := &s3.PutObjectInput{
		Bucket:        aws.String(config.S3Name),
		Key:           aws.String(param.S3Endpoint),
		Body:          fBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fType),
	}
	res, err := client.PutObject(putParams)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func S3Uploads(params S3UploadParams) error {
	_, pErrors, _ := jparallel.Foreach(
		params,
		func(i int, param S3UploadParam) (string, error) {
			_, err := S3Upload(param)
			if err != nil {
				return "", err
			}
			return "", nil
		},
		200,
	)
	return pErrors.Error()
}

func S3Delete(param S3DeleteParam) (*s3.DeleteObjectOutput, error) {
	if !isInit {
		return nil, fmt.Errorf("[S3Delete Fail] not call InitS3 - %v", param)
	}

	delParam := &s3.DeleteObjectInput{
		Bucket: aws.String(config.S3Name),
		Key:    aws.String(param.S3Path),
	}
	res, err := client.DeleteObject(delParam)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func S3Deletes(params S3DeleteParams) error {
	_, pErrors, _ := jparallel.Foreach(
		params,
		func(i int, param S3DeleteParam) (string, error) {
			_, err := S3Delete(param)
			if err != nil {
				return "", err
			}
			return "", nil
		},
		200,
	)

	for _, pErr := range pErrors {
		if pErr.Err != nil {
			return pErr.Err
		}
	}

	return nil
}
