package ms

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/decodethedev/email-gen/global"
	"github.com/decodethedev/email-gen/helpers"
	"github.com/decodethedev/email-gen/logrus"
	"github.com/decodethedev/email-gen/utils"
	"github.com/tidwall/gjson"
)

func createAccount(simplifiedConfig string, client *http.Client) (string, error) {

	headers := map[string]string{
		"authority":         "signup.live.com",
		"accept":            "application/json",
		"accept-language":   "en-US,en;q=0.6",
		"canary":            gjson.Get(simplifiedConfig, "canary").String(),
		"hpgid":             gjson.Get(simplifiedConfig, "hpgid").String(),
		"origin":            "https://signup.live.com",
		"referer":           "https://signup.live.com/signup",
		"scid":              gjson.Get(simplifiedConfig, "scid").String(),
		"sec-fetch-dest":    "empty",
		"sec-fetch-mode":    "cors",
		"sec-fetch-site":    "same-origin",
		"sec-gpc":           "1",
		"tcxt":              gjson.Get(simplifiedConfig, "tcxt").String(),
		"uaid":              gjson.Get(simplifiedConfig, "uaid").String(),
		"uiflvr":            gjson.Get(simplifiedConfig, "uiflvr").String(),
		"user-agent":        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
		"x-ms-apitransport": "xhr",
		"x-ms-apiversion":   "2",
		"cookie":            gjson.Get(simplifiedConfig, "amsc").String(),
	}

	bodyStr := fmt.Sprintf(`{"RequestTimeStamp": "2022-11-18T06:28:55.706Z", "MemberName": "%s", "CheckAvailStateMap": ["%s:undefined"], "EvictionWarningShown": [], "UpgradeFlowToken": "", "FirstName": "%s", "LastName": "%s", "MemberNameChangeCount": 1, "MemberNameAvailableCount": 1, "MemberNameUnavailableCount": 0, "CipherValue": "KAO4R/Q2Mz2hxRbVeFCrJ9dKyYGiYIcc2Umr45vSmj7OGUAs0RDVyxT1+UZhyDd7z4w3WL4dWV8ymmaJ3sOhDChX0vBdVyWR2zangBmSF+kjGeNtsKY3RcBfy6Z80YFTxBvgaJMXVn3trU0f4VH8R15BEC/F8uzvIO2rmmHtaezhlxsVLITVAr57tYrp/+I9MpkSO02yX4dWuUBNq3+1Rz8Hti3x4JIsKmDBpGHDEXrfzHi5lxVGbzTBqwWzgf/BE3nI2+gvwRoN34EE9W6TBoZcOW9e4PBIv4xyL/urbrNaEqvkIlA82X/EP54BlVragndttt4ZT6bxt6P1TU4XXw==", "SKI": "AF99E0B5CB4A0FCD26625571F926665CC11334CE", "BirthDate": "04:04:2000", "Country": "US", "IsOptOutEmailDefault": false, "IsOptOutEmailShown": true, "IsOptOutEmail": false, "LW": true, "SiteId": "292841", "IsRDM": 0, "ReturnUrl": null, "SignupReturnUrl": null, "uiflvr": %s, "uaid": "%s", "SuggestedAccountType": "OUTLOOK", "SuggestionType": "Locked", "HFId": "%s", "encAttemptToken": "", "dfpRequestId": "", "scid": %s, "hpgid": 200650}`, gjson.Get(simplifiedConfig, "email").String(),
		gjson.Get(simplifiedConfig, "email").String(),
		gjson.Get(simplifiedConfig, "firstName").String(),
		gjson.Get(simplifiedConfig, "lastName").String(),
		gjson.Get(simplifiedConfig, "uiflvr").String(),
		gjson.Get(simplifiedConfig, "uaid").String(),
		gjson.Get(gjson.Get(simplifiedConfig, "hip").String(), "fid").String(),
		gjson.Get(simplifiedConfig, "scid").String())

	response, err := helpers.NewClientRequest(client, "POST", "https://signup.live.com/API/CreateAccount",
		bodyStr, headers, nil)

	if err != nil {
		logrus.Errorf("Error creating account: %s", err.Error())
		return "", err
	}

	// Get body
	body, err := helpers.ReadResponseBody(&response)
	if err != nil {
		logrus.Errorf("Error reading response body: %s", err.Error())
		return "", err
	}

	return body, nil
}

func createAccountCaptcha(simplifiedConfig string, client *http.Client, enc string, df string, captchaToken string) (string, error) {

	headers := map[string]string{
		"authority":         "signup.live.com",
		"accept":            "application/json",
		"accept-language":   "en-US,en;q=0.6",
		"canary":            gjson.Get(simplifiedConfig, "canary").String(),
		"hpgid":             gjson.Get(simplifiedConfig, "hpgid").String(),
		"origin":            "https://signup.live.com",
		"referer":           "https://signup.live.com/signup",
		"scid":              gjson.Get(simplifiedConfig, "scid").String(),
		"sec-fetch-dest":    "empty",
		"sec-fetch-mode":    "cors",
		"sec-fetch-site":    "same-origin",
		"sec-gpc":           "1",
		"tcxt":              gjson.Get(simplifiedConfig, "tcxt").String(),
		"uaid":              gjson.Get(simplifiedConfig, "uaid").String(),
		"uiflvr":            gjson.Get(simplifiedConfig, "uiflvr").String(),
		"user-agent":        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
		"x-ms-apitransport": "xhr",
		"x-ms-apiversion":   "2",
		"cookie":            gjson.Get(simplifiedConfig, "amsc").String(),
	}

	bodyStr := fmt.Sprintf(`{
"RequestTimeStamp": "2022-11-18T07:21:43.132Z",
"MemberName": "%s",
"CheckAvailStateMap": ["%s:undefined"],
"EvictionWarningShown": [],
"UpgradeFlowToken": {},
"FirstName": "%s",
"LastName": "%s",
"MemberNameChangeCount": 1,
"MemberNameAvailableCount": 1,
"MemberNameUnavailableCount": 0,
"CipherValue": "qbJ3iA4uAym8J2eDczXrfq9IeRT4xPLEm9XybpUwBlP34x6kH02WWk5Hv0QuVwewLjL7zyiMaMR4pyF1FcbeNR/zniHKAyxiPCkiwedFj+/qE9PyMU78pp5NJYrNKaaA3LsPKVXWq7EKYekWJ7xkkVzsTiNu7nslZFBS9C2oZ2ZAjPbUtLx1SbGaXEXHyjAGimTAIcuj0hye5BwrFf7o2Sd8yYOCbf2ccdzUfDMMxQLVs3bJ1uNmcSY14JENw9NS19/cTPgtQjpuMdNpQVdv9e7MsB1guLhuQs5WPygPHDMdK3/vl6FMjInT4Zmp+nH5J1oZZHpR3xZwLas9ttmsoA==",
"SKI": "AF99E0B5CB4A0FCD26625571F926665CC11334CE",
"BirthDate": "18:09:2000",
"Country": "US",
"IsOptOutEmailDefault": false,
"IsOptOutEmailShown": true,
"IsOptOutEmail": false,
"LW": true,
"SiteId": "292841",
"IsRDM": 0,
"WReply": null,
"ReturnUrl": null,
"SignupReturnUrl": "GAYASF",
"uiflvr": %s,
"uaid": "%s",
"SuggestedAccountType": "EASI",
"SuggestionType": "Prefer",
"HFId": "%s",
"HType": "enforcement",
"HSol": "%s",
"HPId": "%s",
"encAttemptToken": "%s",
"dfpRequestId": "%s",
"scid": %s,
"hpgid": 201040
}`, gjson.Get(simplifiedConfig, "email").String(),
		gjson.Get(simplifiedConfig, "email").String(),
		gjson.Get(simplifiedConfig, "firstName").String(),
		gjson.Get(simplifiedConfig, "lastName").String(),
		gjson.Get(simplifiedConfig, "uiflvr").String(),
		gjson.Get(simplifiedConfig, "uaid").String(),
		gjson.Get(gjson.Get(simplifiedConfig, "hip").String(), "fid").String(),
		captchaToken,
		gjson.Get(gjson.Get(simplifiedConfig, "hip").String(), "enforcement.pid").String(),
		enc,
		df,
		gjson.Get(simplifiedConfig, "scid").String())

	bodyStr = strings.ReplaceAll(bodyStr, "GAYASF", "https://login.live.com/login.srf%3fcobrandid%3d90015%26id%3d292841%26cobrandid%3d90015%26id%3d292841%26contextid%3dED557D44DA295E3B%26opid%3d626703912739EC56%26mkt%3dEN-US%26lc%3d1033%26bk%3d1668750792%26uaid%3d68a91b8a23ee4724931dbbb3aff9c94b")

	response, err := helpers.NewClientRequest(client, "POST", "https://signup.live.com/API/CreateAccount",
		bodyStr, headers, nil)

	if err != nil {
		logrus.Errorf("Error creating account (captcha): %s", err)
		return "", err
	}

	// Get body
	body, err := helpers.ReadResponseBody(&response)
	if err != nil {
		logrus.Errorf("Error reading response body (captcha): %s", err)
		return "", err
	}

	return body, nil
}

func getCaptcha(token string, client *http.Client, simplifiedConfig string) (string, error) {

	headers := map[string]string{
		"content-type": "application/json",
	}

	res, err := helpers.NewClientRequest(client, "POST", "https://api.anycaptcha.com/createTask", fmt.Sprintf(`
	{
		"clientKey": "%s",
		"task": {
			"type": "FunCaptchaTaskProxyless",
			"websitePublicKey": "%s",
			"websiteURL": "%s"
		}
	}`, token, gjson.Get(gjson.Get(simplifiedConfig, "hip").String(), "enforcement.pid"), gjson.Get(gjson.Get(simplifiedConfig, "hip").String(), "enforcement.url")), headers, nil)
		
	if err != nil {
		logrus.Errorf("Error creating captcha task: %s", err)
		return "", err
	}

	body, err := helpers.ReadResponseBody(&res)
	if err != nil {
		logrus.Errorf("Error reading response body (captcha): %s", err)
		return "", err
	}

	status := gjson.Get(body, "errorId").String()
	if status != "0" {
		return "", errors.New("Error creating captcha task: " + body)
	}

	taskID := gjson.Get(body, "taskId").String()

	// logrus.Successf("Captcha task created with ID: %s", taskID)

	result := ""

	for result == "" {
		response, err := helpers.NewClientRequest(client, "POST", "https://api.anycaptcha.com/getTaskResult", fmt.Sprintf(`{
			"clientKey": "%s",
			"taskId": %s
		}`, token, taskID), headers, nil)

		if err != nil {
			logrus.Errorf("Error getting captcha task result: %s", err)
			return "", err
		}

		bodyStr2, err := helpers.ReadResponseBody(&response)
		if err != nil {
			logrus.Errorf("Error reading response body (captcha): %s", err)
			return "", err
		}

		status = gjson.Get(bodyStr2, "errorId").String()
		if status != "0" {
			return "", errors.New("Error getting captcha task result: " + bodyStr2)
		}

		result = gjson.Get(bodyStr2, "solution.token").String()
	}

	return result, nil
}

// func parseImagePath(src string) string {
// 	// Regex: \"hipChallengeUrl\":\"(.*)\",\"imagePath\"

// 	re := regexp.MustCompile(`\"hipChallengeUrl\":\"(.*)\",\"imagePath\"`)
// 	match := re.FindStringSubmatch(src)

// 	if len(match) == 0 {
// 		return ""
// 	}

// 	return match[1]
// }

// func getImageCaptcha(token string, simplifiedConfig string, client *http.Client) string {
// 	headers := map[string]string{
// 		"content-type": "application/json",
// 	}

// 	// Get image captcha
// 	response, _ := helpers.NewClientRequest(client, "GET", gjson.Get(gjson.Get(simplifiedConfig, "hip").String(), "url").String(), nil, nil, nil)

// 	// Get body
// 	body, _ := helpers.ReadResponseBody(&response)

// 	// Get image
// 	imagePath := fmt.Sprintf(`"%s"`, parseImagePath(body))
// 	imagePath, _ = strconv.Unquote(imagePath)

// 	// Get image and encode base64

// 	imgBdy := ""
// 	{
// 		res, _ := helpers.NewClientRequest(client, "GET", imagePath, map[string]string{
// 			"Accept":          "*/*",
// 			"Accept-Language": "en-US,en;q=0.6",
// 			"Connection":      "keep-alive",
// 			// "Cookie": "logonLatency=LGN01=638020949733023967; mkt=en-US; mkt1=en-US; amsc=3EF0G2OCnEx/VefWohrhdlDzzllO+DMz9mMAi2vRyHLxH68T/Rw4nFxfg7B5X5aZFfWp71nvl7E9e1tMnEj9+trwU9MsCYm3rkfqoDa/ielHRY/Z1f5rgsvojcWxFry5+n052J6U3J2zFw5YSLHo9iCwO5xIiqpXqRgsN78qvbj1zDW1RW0AFFmnbX21hAfn++mI8vSYeW4lY2HztQ2ZQyZqLdA2j2SO8KhE58lgDAJMU3wXzn7qAgZHXeLWf+CFf8ZongtvJgYgGKgNRaO93g4sfL0cBReckM7vaCyKYM6tPH9DYf3nlWDtOanqj+e2:2:3c; fptctx2=taBcrIH61PuCVH7eNCyH0FWPWMZs3CpAZMKmhMiLe%252bHdXFXNzpMdAlKdbbEGnY%252fmLpekSdWzWZyQE6eoaUnAAzzRQivxyq5CaqBUqyvCg2KpoBhyhkyRciybG33V1bcM9VthmLls0H00%252f%252fGxalX3nUsCB%252fl03kOoW5bLbwFJWSSpiCFvb8tgd1TMLprVRfnG5wQSi34RdToGAzLMWHVIlANSwx6gd2KHlF1Ad8DzozVGpwwZJVjmaS02E3gTnbquqle0DKv0ZZeC5OH6SDagtpMLdwC3gDzRbuMnU1CIBnw%253d; MUID=ff6e8b3cb0c640aca74b656f985b24e3",
// 			"Referer":        "https://signup.live.com/",
// 			"Sec-Fetch-Dest": "script",
// 			"Sec-Fetch-Mode": "no-cors",
// 			"Sec-Fetch-Site": "same-site",
// 			"Sec-GPC":        "1",
// 			"User-Agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
// 		}, nil, nil)
// 		imgBdy, _ = helpers.ReadResponseBody(&res)

// 		imgBdy = base64.StdEncoding.EncodeToString([]byte(imgBdy))
// 		fmt.Println(imgBdy)
// 	}

// 	// Send image to any captcha

// 	res, _ := helpers.NewClientRequest(client, "POST", "https://api.anycaptcha.com/createTask", fmt.Sprintf(`{
// 		"clientKey": "%s",
// 		"task": {
// 			"type": "ImageToTextTask",
// 			"body": "%s",
// 			"subType: "MICROSOFT"
// 		}
// 	}`, token, imgBdy), headers, nil)

// 	bodyStr, _ := helpers.ReadResponseBody(&res)

// 	if gjson.Get(bodyStr, "errorId").String() != "0" {
// 		fmt.Println(bodyStr)
// 		return ""
// 	}

// 	taskID := gjson.Get(bodyStr, "taskId").String()

// 	result := ""
// 	for result == "" {
// 		response, _ := helpers.NewClientRequest(client, "POST", "https://api.anycaptcha.com/getTaskResult", fmt.Sprintf(`{
// 			"clientKey": "%s",
// 			"taskId": %s
// 		}`, token, taskID), headers, nil)

// 		bodyStr2, _ := helpers.ReadResponseBody(&response)
// 		if gjson.Get(bodyStr2, "errorId").String() != "0" {
// 			fmt.Println(bodyStr2)
// 			print("Breh")
// 			return ""
// 		}

// 		if gjson.Get(bodyStr2, "status").String() == "processing" {
// 			continue
// 		}

// 		result = gjson.Get(bodyStr2, "solution.text").String()
// 	}

// 	return result
// } IMAGE CAPTCHA

func CreateMSAccount(emailType bool) (username string, password string, err error) {
	start := time.Now()
	client := http.Client{}
	proxyEn := gjson.Get(global.ConfigJson, "use_proxies").Bool()

	prxyStr := ""

	if proxyEn {

		randomInt := global.RandomObj.Intn(len(global.ProxyQueue))
		prxyStr = global.ProxyQueue[randomInt].ProxyString

		url, err := url.Parse("http://" + prxyStr)
		if err != nil {
			fmt.Println(err.Error())
			return "", "", err
		}

		proxt := http.Transport{
			Proxy: http.ProxyURL(url),
		}

		client = http.Client{
			Transport: &proxt,
		}
	}

	// Install depps

	simplifiedConfig, err := utils.GetConfig(&client, emailType)
	if err != nil {
		return "", "", err
	}

	email := gjson.Get(simplifiedConfig, "email").String()

	logrus.Infof("[ACCOUNT] Got simplified config for %s", email)

	data, err := createAccount(simplifiedConfig, &client)
	if err != nil {
		return "", "", err
	}

	captcha := false
	if strings.Contains(data, "hipEnforcement") {
		captcha = true
	} else {
		if strings.Contains(data, "hip") {
			captcha = false
		} else {
			logrus.Errorf("[ACCOUNT] Failed to create account for %s", email)
			fmt.Println(data)
			return "", "", errors.New("unknown error")
		}
	}

	encAttemptToken := (gjson.Get(gjson.Get(data, "error.data").String(), "encAttemptToken"))
	dfpRequestId := (gjson.Get(gjson.Get(data, "error.data").String(), "dfpRequestId"))
	
	authKey := ""

	if captcha {
		token := ""
		for token == "" {
			token, err = getCaptcha(gjson.Get(global.ConfigJson, "any_captcha_token").String(), &client, simplifiedConfig)
			if err != nil {
				token = ""
			}
		}

		logrus.Infof("[ACCOUNT] Got captcha token %s", email)

		data, err = createAccountCaptcha(simplifiedConfig, &client, encAttemptToken.String(), dfpRequestId.String(), token)
		if err != nil {
			return "", "", err
		}

		logrus.Infof("[ACCOUNT] Responded with captcha token %s", email)

		eemail := gjson.Get(data, "signinName").String()

		// // fmt.Println(data) POP ENABLE SECTION IS BROKEN

		// if eemail == "" {
		// 	logrus.Errorf("[ACCOUNT] Failed to create account for %s", email)
		// 	fmt.Println(data)
		// 	return "", "", errors.New("Unknown error")
		// }

		// authKey, err = utils.Login(eemail, prxyStr)

		// count := 0
		// for err != nil {
		// 	if count == 3 {
		// 		return "", "", err
		// 	}
		// 	authKey, err = utils.Login(eemail, prxyStr)
		// 	count++
		// }

		logrus.Infof("[ACCOUNT] Logged in %s", eemail)

	} else {
		return "", "", errors.New("no captcha")
	}

	logrus.Successf("Done creating account: %s", email)

	helpers.SaveMailAccount(email, "@GeneratorPassword123", authKey)
	t := time.Now()
	elapsed := t.Sub(start)

	logrus.Successf("[%s] Time taken: %s", email, elapsed)

	return email, "@GeneratorPassword123", nil
}
