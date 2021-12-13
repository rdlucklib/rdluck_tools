package upload

import (
	"net/url"
	"os"
	"mime/multipart"
	"path/filepath"
	"fmt"
	"bytes"
	"io"
	"net/http"
)

//文件上传
//newfileUploadRequest(上传目标路径, "上传参数", "上传文件参数名", "被上传文件路径")

func NewfileUploadRequest(uri string, params url.Values, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	for key, val := range params {
		_ = writer.WriteField(key, val[0])
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", uri, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	fmt.Println("Content-Type", writer.FormDataContentType())
	return request, err
}
