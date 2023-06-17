package global

import (
	"math/rand"

	"github.com/zenthangplus/goccm"
)

var ConfigJson string

var ProxyQueue []Proxy

var Random rand.Source
var RandomObj *rand.Rand

var AccountsCreated int
var Errors int

var GoroutinesHandler goccm.ConcurrencyManager

var Emails []HotmailEmail
