package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/decodethedev/email-gen/global"
	"github.com/decodethedev/email-gen/helpers"
	"github.com/decodethedev/email-gen/logrus"
	"github.com/jaswdr/faker"
	"github.com/playwright-community/playwright-go"
	"github.com/tidwall/gjson"
)

// func MakeSimpliedConfig(config string, amsc string, emailType bool) string {

// 	simplifiedConfig := make(map[string]interface{}, 0)
// 	parseMap := map[string]string{
// 		"hpgid":  "hpgid",
// 		"scid":   "scid",
// 		"uiflvr": "uiflvr",
// 		"tcxt":   "clientTelemetry.tcxt",
// 		"hip":    "WLXAccount.hip",
// 		"uaid":   "uaid",
// 		"canary": "apiCanary",
// 	}

// 	variable := ""
// 	for key, value := range parseMap {
// 		variable = ""
// 		js := gjson.Get(config, value).String()

// 		variable = js

// 		if strings.Contains(variable, "{") {

// 			json.Unmarshal([]byte(variable), &variable)
// 		}

// 		simplifiedConfig[key] = variable
// 	}

// 	simplifiedConfig["amsc"] = amsc

// 	fake := faker.New()
// 	person := fake.Person()
// 	name := person.Name()

// 	firstName := strings.Split(name, " ")[0]
// 	simplifiedConfig["firstName"] = strings.ReplaceAll(firstName, ".", "")

// 	lastName := strings.Split(name, " ")[1]
// 	simplifiedConfig["lastName"] = lastName

// 	email := ""

// 	if emailType {
// 		email = strings.ToLower(firstName + lastName + helpers.RandStringBytes(3) + "@hotmail.com")
// 	} else {
// 		email = strings.ToLower(firstName + lastName + helpers.RandStringBytes(3) + "@outlook.com")
// 	}
// 	simplifiedConfig["email"] = email

// 	data, _ := json.Marshal(simplifiedConfig)

// 	return fmt.Sprintf("%s", data)
// }

func GetConfig(client *http.Client, emailType bool) (string, error) {

	res, err := helpers.NewClientRequest(client, "GET", "https://signup.live.com/signup?lic=1", nil, nil, nil)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	defer res.Body.Close()

	amsc := strings.Split(res.Header.Get("Set-Cookie"), ";")[0]

	bodyStr, err := helpers.ReadResponseBody(&res)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("GetConfig: %s", bodyStr)
	}

	// fmt.Println(bodyStr)

	re := regexp.MustCompile(`(?m)var t0=(.*?);`)
	match := re.FindStringSubmatch(bodyStr)

	config := match[0]
	config = strings.ReplaceAll(config, "var t0=", "")
	config = strings.ReplaceAll(config, ";", "")

	var s map[string]interface{}

	err = json.Unmarshal([]byte(config), &s)
	if err != nil {
		log.Fatal(err)
	}

	e, _ := json.Marshal(s)
	config = string(e)

	simplifiedConfig := make(map[string]interface{}, 0)

	simplifiedConfig["hpgid"] = s["hpgid"]
	simplifiedConfig["scid"] = s["scid"]
	simplifiedConfig["uiflvr"] = s["uiflvr"]

	simplifiedConfig["tcxt"] = s["clientTelemetry"].(map[string]interface{})["tcxt"].(string)

	str, _ := json.Marshal(s["WLXAccount"].(map[string]interface{})["hip"])
	simplifiedConfig["hip"] = string(str)

	simplifiedConfig["uaid"] = s["uaid"]
	simplifiedConfig["canary"] = s["apiCanary"]

	simplifiedConfig["amsc"] = amsc

	fake := faker.New()
	person := fake.Person()
	name := person.Name()

	firstName := strings.Split(name, " ")[0]
	simplifiedConfig["firstName"] = strings.ReplaceAll(firstName, ".", "")

	lastName := strings.Split(name, " ")[1]
	simplifiedConfig["lastName"] = lastName

	email := ""

	if emailType {
		email = strings.ToLower(firstName + lastName + helpers.RandStringBytes(3) + "@hotmail.com")
	} else {
		email = strings.ToLower(firstName + lastName + helpers.RandStringBytes(3) + "@outlook.com")
	}
	simplifiedConfig["email"] = email

	data, _ := json.Marshal(simplifiedConfig)
	strData := string(data)

	// {
	// 	var data = strings.NewReader(`{"pageApiId":200639,"clientDetails":[],"country":"US","userAction":"","source":"PageView","clientTelemetryData":{"category":"PageLoad","pageName":"200639","eventInfo":{"timestamp":1666171020952,"enforcementSessionToken":null,"perceivedPlt":4139,"networkLatency":2678,"appVersion":null,"networkType":null,"precaching":null,"bundleVersion":null,"deviceYear":null,"isMaster":null,"bundleHits":null,"bundleMisses":null}},"cxhFunctionRes":null,"netId":null,"uiflvr":1001,"uaid":"` + gjson.Get(strData, "uaid").String() + `","scid":100118,"hpgid":200639}`)
	// 	req, err := http.NewRequest("POST", "https://signup.live.com/API/ReportClientEvent?lcid=1033&wa=wsignin1.0&rpsnv=13&ct=1666171012&rver=7.0.6737.0&wp=MBI_SSL&wreply=https%3a%2f%2foutlook.live.com%2fowa%2f%3fRpsCsrfState%3d52dd0e10-7f66-b2c7-1391-150fe3bea1a5&id=292841&CBCXT=out&lw=1&fl=dob%2cflname%2cwld&cobrandid=90015&aadredir=0&lic=1&uaid=81df3a4e8b66417cbc3bd0754b9f6247", data)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	req.Header.Set("authority", "signup.live.com")
	// 	req.Header.Set("accept", "application/json")
	// 	req.Header.Set("accept-language", "en-GB,en-US;q=0.9,en;q=0.8")
	// 	req.Header.Set("canary", gjson.Get(strData, "canary").String())
	// 	req.Header.Set("content-type", "application/json")
	// 	req.Header.Set("cookie", "amsc="+gjson.Get(strData, "amsc").String())
	// 	req.Header.Set("hpgid", "200639")
	// 	req.Header.Set("origin", "https://signup.live.com")
	// 	req.Header.Set("referer", "https://signup.live.com/signup")
	// 	req.Header.Set("scid", "100118")
	// 	req.Header.Set("sec-ch-ua", `"Chromium";v="106", "Google Chrome";v="106", "Not;A=Brand";v="99"`)
	// 	req.Header.Set("sec-ch-ua-mobile", "?0")
	// 	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	// 	req.Header.Set("sec-fetch-dest", "empty")
	// 	req.Header.Set("sec-fetch-mode", "cors")
	// 	req.Header.Set("sec-fetch-site", "same-origin")
	// 	req.Header.Set("tcxt", gjson.Get(strData, "tcxt").String())
	// 	req.Header.Set("uaid", gjson.Get(strData, "uaid").String())
	// 	req.Header.Set("uiflvr", "1001")
	// 	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	// 	req.Header.Set("x-ms-apitransport", "xhr")
	// 	req.Header.Set("x-ms-apiversion", "2")
	// 	resp, err := client.Do(req)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	defer resp.Body.Close()
	// 	bodyText, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	if resp.StatusCode != 200 {
	// 		return "", errors.New("Error: " + string(bodyText))
	// 	}

	// 	simplifiedConfig["canary"] = gjson.Get(string(bodyText), "apiCanary").String()
	// 	simplifiedConfig["tcxt"] = gjson.Get(string(bodyText), "telemetryContext").String()
	// }

	// data, _ = json.Marshal(simplifiedConfig)
	// strData = string(data)

	return strData, nil

	// pw, err := playwright.Run()
	// if err != nil {
	// 	logrus.Errorf("[CONFIG] Could not launch playwright: %s", err.Error())
	// 	return nil, "", err
	// }

	// proxyEn := gjson.Get(global.ConfigJson, "use_proxies").Bool()

	// launchOptions := playwright.BrowserTypeLaunchOptions{
	// 	Headless: playwright.Bool(true),
	// }
	// if proxyEn {

	// 	url, err := url.Parse("http://" + proxy)
	// 	if err != nil {
	// 		logrus.Errorf("[CONFIG] Could not parse proxy: %s", err.Error())
	// 		return nil, "", err
	// 	}

	// 	username := url.User.Username()
	// 	password, _ := url.User.Password()

	// 	proxt := playwright.BrowserTypeLaunchOptionsProxy{
	// 		Server:   &url.Host,
	// 		Username: &username,
	// 		Password: &password,
	// 	}

	// 	launchOptions = playwright.BrowserTypeLaunchOptions{
	// 		Headless: playwright.Bool(true),
	// 		Proxy:    &proxt,
	// 	}

	// }

	// browser, err := pw.Firefox.Launch(launchOptions)
	// if err != nil {
	// 	logrus.Errorf("[CONFIG] Could not launch browser: %v", err)
	// 	return nil, "", err
	// }

	// viewportSettings := playwright.BrowserNewContextOptionsViewport{
	// 	Width:  playwright.Int(50),
	// 	Height: playwright.Int(50),
	// }

	// browserContext, err := browser.NewContext(playwright.BrowserNewContextOptions{
	// 	Viewport: &viewportSettings,
	// })

	// if err != nil {
	// 	logrus.Errorf("[CONFIG] Could not create context: %v", err)
	// 	return nil, "", err
	// }

	// page, err := browserContext.NewPage()
	// if err != nil {
	// 	logrus.Errorf("[CONFIG] Could not create page: %v", err)
	// 	return nil, "", err
	// }

	// excluded_resource_types := []string{"font", "stylesheet", "image"}
	// amsc := ""

	// handler := func(route playwright.Route, request playwright.Request) {

	// 	// if strings.Contains(request.URL(), "https://signup.live.com/signup") {
	// 	// 	content, _ := ioutil.ReadFile("./bin/login.html")

	// 	// 	route.Fulfill(playwright.RouteFulfillOptions{
	// 	// 		Body: playwright.String(string(content)),
	// 	// 	})
	// 	// }

	// 	allHeaders := request.Headers()

	// 	for _, excluded_resource_type := range excluded_resource_types {
	// 		if request.ResourceType() == excluded_resource_type {
	// 			route.Abort()
	// 			return
	// 		}
	// 	}

	// 	for key, value := range allHeaders {
	// 		if strings.Contains(key, "cookie") && strings.Contains(value, "amsc") {
	// 			amsc = value
	// 			break
	// 		}
	// 	}

	// 	route.Continue()
	// }

	// page.Route("**/*", handler)
	// page.Goto("https://signup.live.com/signup", playwright.PageGotoOptions{WaitUntil: playwright.WaitUntilStateNetworkidle})

	// page.Reload(playwright.PageReloadOptions{WaitUntil: playwright.WaitUntilStateNetworkidle})

	// value, err := page.Evaluate("$Config")
	// if err != nil {
	// 	logrus.Errorf("[CONFIG] Could not evaluate: %v", err)
	// 	return nil, "", err
	// }

	// browser.Close()

	// return value, amsc, nil

}

func Login(email string, proxy string) (string, error) {
	pw, err := playwright.Run()
	if err != nil {
		logrus.Errorf("[LOGIN] Could not launch playwright: %s", err.Error())
		return "", err
	}

	proxyEn := gjson.Get(global.ConfigJson, "use_proxies").Bool()

	launchOptions := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	}
	if proxyEn {
		logrus.Printf("[LOGIN] Using proxy: %s", proxy)
		url, err := url.Parse("http://" + proxy)
		if err != nil {
			logrus.Errorf("[LOGIN] Could not parse proxy: %s", err.Error())
			return "", err
		}

		username := url.User.Username()
		password, _ := url.User.Password()

		proxt := playwright.BrowserTypeLaunchOptionsProxy{
			Server:   &url.Host,
			Username: &username,
			Password: &password,
		}

		launchOptions = playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(false),
			Proxy:    &proxt,
		}

	}

	browser, err := pw.Firefox.Launch(launchOptions)
	if err != nil {
		logrus.Errorf("[LOGIN] Could not launch browser: %v", err)
		return "", err
	}

	viewportSettings := playwright.BrowserNewContextOptionsViewport{
		Width:  playwright.Int(570),
		Height: playwright.Int(570),
	}

	browserContext, err := browser.NewContext(playwright.BrowserNewContextOptions{
		Viewport: &viewportSettings,
	})

	if err != nil {
		logrus.Errorf("[LOGIN] Could not create context: %v", err)
		browser.Close()
		return "", err
	}

	page, err := browserContext.NewPage()
	if err != nil {
		logrus.Errorf("[LOGIN] Could not create page: %v", err)
		browser.Close()
		return "", err
	}

	excluded_resource_types := []string{"image", "font"}
	authKey := ""

	handler := func(route playwright.Route, request playwright.Request) {
		// if strings.Contains(request.URL(), "https://login.live.com/login.srf") {
		// 	content, _ := ioutil.ReadFile("./bin/login.html")

		// 	route.Fulfill(playwright.RouteFulfillOptions{
		// 		Body: playwright.String(string(content)),
		// 	})
		// }

		for _, excluded_resource_type := range excluded_resource_types {
			if request.ResourceType() == excluded_resource_type {
				route.Abort()
				return
			}
		}

		if authKey == "" {
			allHeaders := request.Headers()
			for key, value := range allHeaders {
				if strings.Contains(key, "cookie") && strings.Contains(value, "RPSSecAuth") {
					splittedValue := strings.Split(value, ";")
					for _, v := range splittedValue {
						if strings.Contains(v, "RPSSecAuth") || strings.Contains(v, "RPSAuth") {
							splittedValue2 := strings.Split(v, "=")
							authKey = splittedValue2[1]
						}
					}
				}
			}
		}

		route.Continue()

	}

	page.Route("**/*", handler)
	page.Goto("https://login.live.com/login.srf?id=292841&aadredir=1&CBCXT=out&lw=1&fl=dob%2cflname%2cwld&cobrandid=90015")

	options := playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateVisible}

	page.WaitForSelector("input#i0116", options)

	page.Type("input#i0116", "breh")
	for i := 0; i < 3; i++ {
		page.Fill("input#i0116", email)
	}

	page.Click("input#idSIButton9")

	page.WaitForSelector("input#i0118", options)

	page.Type("input#i0118", "breh")
	for i := 0; i < 3; i++ {
		page.Fill("input#i0118", "@GeneratorPassword123")
	}

	page.Click("input#idSIButton9")
	page.WaitForSelector("input#idSIButton9", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})

	page.WaitForSelector("input#idSIButton9", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateVisible})
	page.Click("input#idSIButton9")

	getTitle := func() string {
		title := page.URL()
		return title
	}

	counts := 0

	for !(strings.Contains(getTitle(), "https://outlook.live.com/mail/0/")) {
		if counts > 50 {
			logrus.Errorf("[LOGIN] Could not login: %v", err)

			browser.Close()
			return "", errors.New("could not login, timeout 5 seconds")
		}

		time.Sleep(100 * time.Millisecond)
		counts++
	}

	page.Goto("https://outlook.live.com/mail/0/options/mail/accounts/popImap", playwright.PageGotoOptions{WaitUntil: playwright.WaitUntilStateLoad})

	page.WaitForSelector("span.ms-ChoiceFieldLabel", options)

	time.Sleep(500 * time.Millisecond)

	for i := 0; i < 3; i++ {
		entries, err := page.QuerySelectorAll("span.ms-ChoiceFieldLabel")
		if err != nil {
			logrus.Errorf("[LOGIN] Could not query selector: %v", err)

			browser.Close()
			return "", err
		}

		for _, entry := range entries {
			text, err := entry.InnerText()
			if err != nil {
				logrus.Errorf("[LOGIN] Could not get inner text: %v", err)

				browser.Close()
				return "", err
			}

			if text == "Yes" {
				entry.Click()
			}
		}
	}

	page.WaitForSelector("span.ms-Button-label", options)

	done := false

	resHandler := func(response playwright.Response) {
		if strings.Contains(response.URL(), "https://outlook.live.com/owa/0/service.svc?action=SetConsumerMailbox&app=Mail") {
			done = true
		}
	}

	page.On("response", resHandler)

	{
		entries, err := page.QuerySelectorAll("span.ms-Button-label")
		if err != nil {
			logrus.Errorf("[LOGIN] Could not query selector: %v", err)

			browser.Close()
			return "", err
		}

		for _, entry := range entries {
			text, err := entry.InnerText()
			if err != nil {
				logrus.Errorf("[LOGIN] Could not get inner text: %v", err)

				browser.Close()
				return "", err
			}

			if text == "Save" {
				entry.Click()
			}
		}
	}

	counts = 0

	for !done {
		if counts > 50 {
			logrus.Errorf("[LOGIN] Could not save: %v", err)
			browser.Close()
			return "", errors.New("could not save, timeout 5 seconds")
		}
		time.Sleep(100 * time.Millisecond)
		counts++
	}

	browser.Close()

	return authKey, nil
}
