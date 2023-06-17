package helpers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/decodethedev/email-gen/global"
	"github.com/decodethedev/email-gen/logrus"
)

func UpdateJSONConfig() error {
	fileBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}

	global.ConfigJson = string(fileBytes)
	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[global.RandomObj.Intn(len(letterBytes))]
	}
	return string(b)
}

func NormalizePhone(phone string, number string) string {
	numberWithoutPlus := number[1:]
	phoneWithoutCountryCode := strings.Replace(phone, numberWithoutPlus, "", 1)

	finalPhone := ""

	count := 0
	for _, char := range phoneWithoutCountryCode {
		if count > 2 {
			count = 0
			finalPhone += " "
		}
		finalPhone += string(char)
		count++
	}

	finalPhone = (number + " " + finalPhone)

	logrus.Printf("Normalized phone: %s", finalPhone)

	return finalPhone
}

func GetProxyTimezone(client *http.Client) (string, error) {
	response, err := NewClientRequest(client, "GET", "http://worldtimeapi.org/api/ip", nil, nil, nil)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", errors.New("failed to get current timezone")
	}

	type Timezone struct {
		Timezone string `json:"timezone"`
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	bodyStr := []byte(string(body))
	timezone := new(Timezone)

	errr := json.Unmarshal(bodyStr, &timezone)
	if errr != nil {
		return "", errr
	}

	// logrus.Successf("Timezone: %s", timezone.Timezone)

	return timezone.Timezone, nil
}

func GetRandomAvatarPath() (string, string) {
	files, err := ioutil.ReadDir("./bin/avatars/")
	if err != nil {
		logrus.Print(err.Error())
	}

	randomFile := files[global.RandomObj.Intn(len(files))]

	// read and encode by base64

	return "./bin/avatars/" + randomFile.Name(), randomFile.Name()
}
