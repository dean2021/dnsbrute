package api

import (
	"net/http"
	"github.com/dean2021/dnsbrute/log"
	"regexp"
	"io/ioutil"
)

type virustotal struct{}

func init() {
	registerAPI(virustotal{})
}

// Name 接口名称
func (h virustotal) Name() string {
	return "virustotal"
}

// Query 查询接口
func (h virustotal) Query(domain string) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		url := "https://www.virustotal.com/en/domain/" + domain + "/information/"
		client := http.Client{Timeout: timeout}
		resp, err := client.Get(url)
		if err != nil {
			log.Info("error while fetching www.virustotal.com", err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Info("error while fetching www.virustotal.com:", err)
			return
		}
		reRegexDomain, _ := regexp.Compile(`<a target="_blank" href="/en/domain/(.*?)/information/">`)
		domainUrls := reRegexDomain.FindAllStringSubmatch(string(body), -1)
		for _, domain := range domainUrls {
			if len(domain) == 2 {
				ch <- domain[1]
			}
		}
	}()
	return ch
}
