package main

import (
	"bufio"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/decodethedev/email-gen/global"
	"github.com/decodethedev/email-gen/helpers"
	"github.com/decodethedev/email-gen/logrus"
	"github.com/decodethedev/email-gen/ms"
	"github.com/playwright-community/playwright-go"
	"github.com/tidwall/gjson"
	"github.com/zenthangplus/goccm"
)

func threadHandler(emailType bool) {
	username, _, err := ms.CreateMSAccount(emailType)
	if err != nil {
		logrus.Errorf("[ACCOUNT] Failed to create account %s (%s)", username, err.Error())
		global.Errors++
	}
	global.GoroutinesHandler.Done()
}

func main() {
	err := playwright.Install()
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())

	global.Random = rand.NewSource(time.Now().UnixNano())
	global.RandomObj = rand.New(global.Random)

	helpers.UpdateJSONConfig()
	helpers.LoadFiles()

	accounts := 0
	reader := bufio.NewReader(os.Stdin)
	{
		logrus.Print("How many accounts do you want to create? ")
		str, err := reader.ReadString('\n')
		if err != nil {
			logrus.Errorf("[ERROR] Error while reading input: %s", err)
			return
		}

		accounts, err = strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(str, "\r", ""), "\n", ""))
		if err != nil {
			logrus.Errorf("[ERROR] Error while converting input to integer: %s", err)
			return
		}
	}

	emailType := false
	{
		logrus.Print("hotmail.com or outlook.com? (1/0): ")
		str, err := reader.ReadString('\n')
		if err != nil {
			logrus.Errorf("[ERROR] Error while reading input: %s", err)
			return
		}

		if strings.Contains(str, "1") {
			emailType = true
		}
	}

	threads := int(gjson.Get(global.ConfigJson, "threads").Int())
	delay := int(gjson.Get(global.ConfigJson, "delay_per_thread").Int())

	global.GoroutinesHandler = goccm.New(threads)
	count := 0

	global.Errors = 0

	for {
		accountsLeft := accounts - count

		if global.AccountsCreated >= accounts {
			break
		}

		if global.Errors > 0 {
			for i := 0; i < global.Errors; i++ {
				global.GoroutinesHandler.Wait()

				go threadHandler(emailType)

				global.Errors--
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
		}

		if accountsLeft > 0 {
			global.GoroutinesHandler.Wait()

			go threadHandler(emailType)

			time.Sleep(time.Duration(delay) * time.Millisecond)

			count += 1

		}

	}

	global.GoroutinesHandler.WaitAllDone()

	logrus.Successf("[SUCCESS] Created %d accounts", global.AccountsCreated)

	logrus.Print("Enter to exit.")
	reader2 := bufio.NewReader(os.Stdin)
	reader2.ReadString('\n')
	os.Exit(0)

}

