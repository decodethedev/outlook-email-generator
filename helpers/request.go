package helpers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/decodethedev/email-gen/global"
	"github.com/decodethedev/email-gen/logrus"
	"github.com/tidwall/gjson"
)

func ToBuffer(obj interface{}) *bytes.Buffer {
	a, err := json.Marshal(obj)
	if err != nil {
		panic(err.Error())
	}
	buffer := bytes.NewBuffer(a)
	return buffer
}

func AddHeaders(request *http.Request, headers map[string]string) {
	for headerName, headerValue := range headers {
		request.Header.Add(headerName, headerValue)
	}
}

func AddCookies(request *http.Request, cookies map[string]string) {
	for cookieName, cookieValue := range cookies {
		cookie := http.Cookie{
			Name:  cookieName,
			Value: cookieValue,
		}
		request.AddCookie(&cookie)
	}
}

func ExtractCookieFromRequest(name string, response *http.Response, target *http.Cookie) {
	cookies := response.Cookies()
	for i := range cookies {
		cookie := cookies[i]
		if cookie.Name == name {
			// fmt.Println(cookie)
			*target = *cookie
		}
	}
}

func ExtractHeaderFromRequest(name string, response *http.Response, target *string) {
	headers := response.Header
	for key, value := range headers {
		if key == name {
			// fmt.Println(cookie)
			*target = value[0]
		}
	}
}

func ConvertStringToJson(jsonString string) (map[string]interface{}, error) {
	jsonBytes := []byte(jsonString)
	var jsonRet map[string]interface{}
	err := json.Unmarshal(jsonBytes, &jsonRet)
	if err != nil {
		return nil, err
	}
	return jsonRet, nil
}

func ReadResponseBody(response *http.Response) (string, error) {
	var reader io.ReadCloser
	var err error

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return "", err
		}
		defer reader.Close()
	default:
		reader = response.Body
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	bodyStr := []byte(string(body))
	return string(bodyStr), nil
}

func NewRequest(requestType string, requestUrl string, jsonBody interface{}, headers map[string]string, cookies map[string]string, proxyUrl string) (http.Response, error) {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	useProxy := gjson.Get(global.ConfigJson, "use_proxies").Value().(bool)

	if useProxy && proxyUrl != "" {
		// logrus.Infof("Using %s as proxy.", proxyUrl)
		url, _ := url.Parse("http://" + proxyUrl)
		client = http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(url),
			},
		}
	}

	var requestBody *bytes.Buffer
	switch a := jsonBody.(type) {
	case string:
		bytes.NewBuffer([]byte(a))
		var jsonStr = []byte(jsonBody.(string))
		requestBody = bytes.NewBuffer(jsonStr)
	default:
		requestBody = ToBuffer(jsonBody)
	}

	var requestObject *http.Request
	var err error
	requestObject, err = http.NewRequest(requestType, requestUrl, requestBody)
	if requestType == "GET" {
		requestObject, err = http.NewRequest(requestType, requestUrl, nil)
	}

	if err != nil {
		return http.Response{}, err
	}

	if headers != nil {
		AddHeaders(requestObject, headers) // Automate the process of adding header so you can add headers using json format
	}
	if cookies != nil {
		AddCookies(requestObject, cookies)
	}

	responseObj, err := client.Do(requestObject)

	if err != nil {
		res := http.Response{}
		return res, err
	}

	return *responseObj, nil
}
func CreateMultipartFormData(fieldName, fileName string) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	file := MustOpen(fileName)
	if fw, err = w.CreateFormFile(fieldName, file.Name()); err != nil {
		logrus.Errorf("Error creating writer: %v", err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		logrus.Errorf("Error with io.Copy: %v", err)
	}
	w.Close()
	return b, w
}

func MustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("PWD: ", pwd)
		panic(err)
	}
	return r
}

func NewClientRequest(client *http.Client, requestType string, requestUrl string, jsonBody interface{}, headers map[string]string, cookies map[string]string) (http.Response, error) {
	// client := http.Client{}

	// useProxy := gjson.Get(global.ConfigJson, "use_proxies").Value().(bool)

	// if useProxy && proxyUrl != "" {
	// 	// logrus.Infof("Using %s as proxy.", proxyUrl)
	// 	url, _ := url.Parse("http://" + proxyUrl)
	// 	client = http.Client{
	// 		Transport: &http.Transport{
	// 			Proxy: http.ProxyURL(url),
	// 		},
	// 	}
	// }

	var requestBody *bytes.Buffer
	switch a := jsonBody.(type) {
	case string:
		bytes.NewBuffer([]byte(a))
		var jsonStr = []byte(jsonBody.(string))
		requestBody = bytes.NewBuffer(jsonStr)
	default:
		requestBody = ToBuffer(jsonBody)
	}

	var requestObject *http.Request
	var err error
	requestObject, err = http.NewRequest(requestType, requestUrl, requestBody)
	if requestType == "GET" {
		requestObject, err = http.NewRequest(requestType, requestUrl, nil)
	}

	if err != nil {
		return http.Response{}, err
	}

	if headers != nil {
		AddHeaders(requestObject, headers) // Automate the process of adding header so you can add headers using json format
	}
	if cookies != nil {
		AddCookies(requestObject, cookies)
	}

	responseObj, err := client.Do(requestObject)

	if err != nil {
		res := http.Response{}
		return res, err
	}

	return *responseObj, nil
}
