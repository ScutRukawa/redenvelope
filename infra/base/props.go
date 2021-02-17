package base

import (
	"fmt"
	"redenvelope/infra"
	"sync"

	"github.com/tietang/props/kvs"
)

var props kvs.ConfigSource

//Props x
func Props() kvs.ConfigSource {
	return props
}

//PropsStarter x
type PropsStarter struct {
	infra.BaseStarter
}

func (p *PropsStarter) Init(ctx infra.StarterContext) {
	props = ctx.Props()
	fmt.Println("props 初始化")
	GetSystemAccount()
}

//SystemAccount x
type SystemAccount struct {
	AccountNo   string
	AccountName string
	UserID      string
	UserName    string
}

var systemAccount *SystemAccount
var systemAccountOnce sync.Once

//GetSystemAccount x
func GetSystemAccount() *SystemAccount {
	systemAccountOnce.Do(func() {
		systemAccount = new(SystemAccount)
		err := kvs.Unmarshal(Props(), systemAccount, "system.account")
		if err != nil {
			panic(err)
		}
	})
	return systemAccount
}

func GetEnvelopeActivityLink() string {
	link := Props().GetDefault("envelope.link", "/v1/envelope/link")
	return link
}

func GetEnvelopeDomain() string {
	domain := Props().GetDefault("envelope.domain", "http://localhost")
	return domain
}
