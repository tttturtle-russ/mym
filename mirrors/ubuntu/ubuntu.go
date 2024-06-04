package ubuntu

import (
	"Mym/consts"
	"Mym/mirrors"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"time"
)

const UbuntuConfig = "/etc/apt/sources.list"
const UbuntuBaseUrl = "http://mirrors.ubuntu.com/"

type UbuntuVersion string

const (
	V1404 UbuntuVersion = "14.04 LTS"
	V1604 UbuntuVersion = "16.04 LTS"
	V1804 UbuntuVersion = "18.04 LTS"
	V2004 UbuntuVersion = "20.04 LTS"
	V2204 UbuntuVersion = "22.04 LTS"
	V2210 UbuntuVersion = "22.10"
	V2304 UbuntuVersion = "23.04"
)

type Ubuntu struct {
	GeneratedDate time.Time
	Version       UbuntuVersion
	Protocols     []consts.Protocol
	Country       string
	IPVersion     []consts.IPVersion
	PreFetch      bool
	Url           string
}

func NewUbuntu(options ...mirrors.Option) Ubuntu {
	u := Ubuntu{}
	for _, op := range options {
		op(&u)
	}
	get, err := resty.New().R().Get(UbuntuBaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(get)
	return u
}

func (u *Ubuntu) Timeout() ([]mirrors.MirrorSite, []float64) {
	//TODO implement me
	panic("implement me")
}

func (u *Ubuntu) List() []mirrors.MirrorSite {
	//TODO implement me
	panic("implement me")
}

func (u *Ubuntu) WriteToConfig() error {
	//TODO implement me
	panic("implement me")
}

func (u *Ubuntu) IsDeployed() bool {
	_, err := os.Stat(UbuntuConfig + ".bak")
	return !errors.Is(err, os.ErrNotExist)
}
