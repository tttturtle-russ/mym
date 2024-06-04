package archlinux

import (
	"Mym/consts"
	"Mym/mirrors"
	"Mym/util"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

const ArchMirrorBaseUrl = "https://archlinux.org/mirrorlist/?"
const ArchConfig = "/etc/pacman.d/mirrorlist"

//consts ArchCountries = []string{
//	"Australia",
//	"Austria",
//	"Azerbaijan",
//	"Bangladesh",
//	"Belarus",
//	"Belgium",
//	"Bosnia and Herzegovina",
//	"Brazil",
//	"Bulgaria",
//	"Colombia",
//	"Canada",
//}

type Arch struct {
	GenerateDate time.Time
	Table        []mirrors.MirrorSite
	Protocols    []consts.Protocol
	Country      string
	IPVersion    []consts.IPVersion
	PreFetch     bool
	Url          string
}

func NewArch(options ...mirrors.Option) Arch {
	arch := Arch{}
	for _, option := range options {
		option(&arch)
	}
	arch.Url = ArchMirrorBaseUrl
	for _, protocol := range arch.Protocols {
		if !util.IsValidProtocol(protocol) {
			log.Fatalf("Error: invalid protocol")
		}
		arch.Url += fmt.Sprintf("protocol=%s&", protocol)
	}
	if len(arch.Protocols) == 0 {
		arch.Protocols = append(arch.Protocols, consts.HTTP, consts.HTTPS)
	}
	for _, ipVersion := range arch.IPVersion {
		if !util.IsValidIPVersion(ipVersion) {
			log.Fatalf("Error: invalid IP version")
		}
		arch.Url += fmt.Sprintf("ip_version=%d&", ipVersion)
	}
	if len(arch.IPVersion) == 0 {
		arch.IPVersion = append(arch.IPVersion, consts.IPV4)
	}
	arch.Url += "country=" + arch.Country
	if arch.PreFetch {
		arch.List()
	}
	return arch
}

func WithPreFetch() mirrors.Option {
	return func(m mirrors.Mirror) {
		m.(*Arch).PreFetch = true
	}
}

func WithProtocols(protocols ...consts.Protocol) mirrors.Option {
	return func(m mirrors.Mirror) {
		m.(*Arch).Protocols = protocols
	}
}

func WithCountry(country string) mirrors.Option {
	return func(m mirrors.Mirror) {
		m.(*Arch).Country = country
	}
}

func WithIPVersion(ipVersion ...consts.IPVersion) mirrors.Option {
	return func(m mirrors.Mirror) {
		m.(*Arch).IPVersion = ipVersion
	}
}

func (a *Arch) getGenerateDate(mirrorlist []string) int {
	for index, line := range mirrorlist {
		if strings.HasPrefix(line, "## Generated") {
			a.GenerateDate, _ = time.Parse("2006-01-02", line[len(line)-10:])
			// +3 to skip the "##" and blank line
			return index + 3
		}
	}
	return -1
}

func (a *Arch) List() []mirrors.MirrorSite {
	var result []mirrors.MirrorSite
	var mirrorlist []string

	if a.Table != nil {
		return a.Table
	}
	response, err := resty.New().
		R().
		Get(a.Url)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	mirrorlist = strings.Split(response.String(), "\n")

	index := a.getGenerateDate(mirrorlist)
	if index == -1 {
		log.Fatalf("Error: can't find generate date")
	}

	for _, line := range mirrorlist[index:] {
		if strings.HasPrefix(line, "#Server") {
			url := line[10:]
			site := mirrors.MirrorSite{
				Url:      url,
				Location: a.Country,
				Protocol: a.Protocols,
			}
			result = append(result, site)
		}
	}
	a.Table = result
	return result
}

func (a *Arch) Timeout() ([]mirrors.MirrorSite, []float64) {
	var timeout2site = make(map[float64]mirrors.MirrorSite)
	var timeoutSlice []float64
	var result []mirrors.MirrorSite

	client := resty.New().
		SetTimeout(consts.DefaultTimeout).
		R()
	for _, site := range a.Table {
		resp, err := client.Get(site.Url)
		if err != nil {
			if !strings.Contains(err.Error(), "Client.Timeout exceeded") {
				log.Fatalf("Error: %v", err)
			}
		}
		timeout := resp.Time().Seconds()
		timeoutSlice = append(timeoutSlice, timeout)
		timeout2site[timeout] = site
	}
	sort.Sort(sort.Float64Slice(timeoutSlice))
	for _, timeout := range timeoutSlice {
		result = append(result, timeout2site[timeout])
	}
	return result, timeoutSlice
}

func (a *Arch) WriteToConfig() error {
	//var root bool
	var err error

	//if runtime.GOOS != "windows" && os.Getuid() == 0 {
	//	root = true
	//}
	//if !root {
	//	err = util.AskForRoot()
	//}
	//if err != nil {
	//	log.Println(err.Error())
	//	log.Fatalln("Failed to ask for sudo")
	//}

	err = os.Rename(ArchConfig, ArchConfig+".bak")
	if err != nil {
		log.Println(err.Error())
		log.Fatalln("Failed to rename config file")
	}
	oldconfig, err := os.Open(ArchConfig + ".bak")
	if err != nil {
		log.Fatalln("Failed to open file")
	}
	defer oldconfig.Close()
	config, err := os.Create(ArchConfig)
	if err != nil {
		log.Fatalln("Failed to create new config file")
	}
	defer config.Close()
	_, err = io.Copy(config, oldconfig)
	return err
}

func (a *Arch) IsDeployed() bool {
	_, err := os.Stat(ArchConfig + ".bak")
	return !errors.Is(err, os.ErrNotExist)
}
