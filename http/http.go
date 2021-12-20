package http

import (
	"io/ioutil"
	"net/http"
	"crypto/tls"
	"strings"
	"time"
	"github.com/rdlucklib/rdluck_tools/log"
)

var Log *log.Log

func init() {
	Log = log.Init("20060102.http")
}

func HttpGet(url string) string {
	sTime := time.Now()
	res, err := http.Get(url)
	eTime := time.Now()
	t := eTime.Sub(sTime).Nanoseconds() / 1000000
	if err != nil {
		return err.Error()
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err.Error()
	}
	Log.Println(t, "ms", "GET", "URL", url, "RESPONSE", string(result))
	return string(result)
}

// func HttpPost(url string, params url.Values) ([]byte, error) {
// 	body := ioutil.NopCloser(strings.NewReader(params.Encode())) //把form数据编下码
// 	client := &http.Client{}
// 	req, err := http.NewRequest("POST", url, body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8") //这个一定要加，不加form的值post不过去
// 	resp, err := client.Do(req)                                                       //发送
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close() //一定要关闭resp.Body
// 	return ioutil.ReadAll(resp.Body)
// }

func HttpPost(url, postData string, params ...string) ([]byte, error) {
	contentType := "application/x-www-form-urlencoded;charset=utf-8"
	if len(params) > 0 && params[0] != "" {
		contentType = params[0]
	}
	sTime := time.Now()
	resp, err := http.Post(url,
		contentType,
		strings.NewReader(postData))
	eTime := time.Now()
	t := eTime.Sub(sTime).Nanoseconds() / 1000000
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	Log.Println(t, "ms", "POST", "URL", url, "DATA", postData, "RESPONSE", string(b))
	return b, err
}

// func HttpPost(url string, params url.Values) ([]byte, error) {
// 	http.DefaultClient.Timeout = 60 * time.Second
// 	resp, err := http.PostForm(url, params)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	return ioutil.ReadAll(resp.Body)
// }

func HttpsPost(url, postData string, params ...string) ([]byte, error) {
	body := ioutil.NopCloser(strings.NewReader(postData))
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	contentType := "application/x-www-form-urlencoded;charset=utf-8"
	if len(params) > 0 && params[0] != "" {
		contentType = params[0]
	}
	req.Header.Set("Content-Type", contentType)
	sTime := time.Now()
	resp, err := client.Do(req)
	eTime := time.Now()
	t := eTime.Sub(sTime).Nanoseconds() / 1000000
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	Log.Println(t, "ms", "POST", "URL", url, "DATA", postData, "RESPONSE", string(b))
	return b, err
}

func Post(url, postData string, params ...string) ([]byte, error) {
	if strings.HasPrefix(url, "https://") {
		return HttpsPost(url, postData, params...)
	} else {
		return HttpPost(url, postData, params...)
	}
}

func Get(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}
