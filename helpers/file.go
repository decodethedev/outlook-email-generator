package helpers

import (
	"bufio"
	"os"

	"github.com/decodethedev/email-gen/global"
	"github.com/decodethedev/email-gen/logrus"
	"github.com/tidwall/gjson"
)

func LoadFiles() {

	proxiesFile := gjson.Get(global.ConfigJson, "proxies_file_name").Value().(string)

	// Read proxies
	{

		proxiesFile, err := os.Open(proxiesFile)
		if err != nil {
			logrus.Error(err.Error())
			os.Exit(0)
		}
		defer proxiesFile.Close()

		proxiesScanner := bufio.NewScanner(proxiesFile)

		for proxiesScanner.Scan() {
			proxy := global.Proxy{
				ProxyString: proxiesScanner.Text(),
			}
			global.ProxyQueue = append(global.ProxyQueue, proxy)
		}

		if len(global.ProxyQueue) == 0 {
			logrus.Error("[ERROR] No proxy found.")
			os.Exit(0)
		}

		logrus.Successf("[File] Got %d proxies file parsed and added to queue.", len(global.ProxyQueue))
	}
}

func SaveMailAccount(email string, password string, authKey string) error {
	// Save to file

	file1, _ := os.OpenFile("./done/emails.txt", os.O_APPEND|os.O_WRONLY, 0600)
	defer file1.Close()

	if _, err := file1.WriteString(email + ":" + password + "\n"); err != nil {
		return err
	}

	global.AccountsCreated += 1

	return nil
}
