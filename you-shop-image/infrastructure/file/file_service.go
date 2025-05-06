package file

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/TechwizsonORG/image-service/config/model"
	"github.com/TechwizsonORG/image-service/err"
)

type FileService struct {
	s3ProxyConfig model.S3ProxyConfig
}

func NewTypeService(s3ProxyConfig model.S3ProxyConfig) *FileService {
	return &FileService{
		s3ProxyConfig: s3ProxyConfig,
	}
}

func (f *FileService) getBase64BasicAuth() string {
	authEncode := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", f.s3ProxyConfig.Username, f.s3ProxyConfig.Password)))
	return fmt.Sprintf("Basic %s", authEncode)
}
func (f *FileService) SaveToRemoteServer(filename string, extension string, file multipart.File) (string, error) {
	fullname := fmt.Sprintf("%s.%s", filename, extension)
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	part, partErr := writer.CreateFormFile("file", fullname)
	if partErr != nil {
		return "", partErr
	}

	_, copyErr := io.Copy(part, file)
	if copyErr != nil {
		return "", copyErr
	}

	closeErr := writer.Close()
	if closeErr != nil {
		return "", closeErr
	}

	imgDest := fmt.Sprintf("%s:%d/%s/", f.s3ProxyConfig.Host, f.s3ProxyConfig.Port, f.s3ProxyConfig.Folder)
	req, newReqErr := http.NewRequest(http.MethodPut, imgDest, buf)
	if newReqErr != nil {
		return "", newReqErr
	}
	req.Header.Set("Authorization", f.getBase64BasicAuth())
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	res, resErr := client.Do(req)
	if resErr != nil {
		return "", resErr
	}
	defer res.Body.Close()
	return fmt.Sprintf("%s%s", imgDest, fullname), nil
}

func (f *FileService) GetFile(filepath string) ([]byte, *err.AppError) {
	req, newReqErr := http.NewRequest(http.MethodGet, filepath, nil)
	if newReqErr != nil {
		return make([]byte, 0), nil
	}
	req.Header.Set("Authorization", f.getBase64BasicAuth())

	client := &http.Client{}

	res, resErr := client.Do(req)
	if resErr != nil {
		return make([]byte, 0), nil
	}
	defer res.Body.Close()

	var buf bytes.Buffer
	io.Copy(&buf, res.Body)
	return buf.Bytes(), nil
}

func (f *FileService) DeleteFile(filepath string) *err.AppError {
	defaultErr := err.NewAppError(500, "Failed when deleting file", "Failed when deleting file", nil)
	req, newReqErr := http.NewRequest(http.MethodDelete, filepath, nil)

	if newReqErr != nil {
		return defaultErr
	}
	req.Header.Set("Authorization", f.getBase64BasicAuth())

	client := &http.Client{}

	res, resErr := client.Do(req)
	if resErr != nil {
		return defaultErr
	}
	defer res.Body.Close()
	return nil
}
