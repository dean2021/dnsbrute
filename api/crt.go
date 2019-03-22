package api

import (
	"net/http"
	"github.com/dean2021/dnsbrute/log"
	"regexp"
	"io/ioutil"
)

type crt struct{}

func init() {
	registerAPI(crt{})
}

// Name 接口名称
func (h crt) Name() string {
	return "crt"
}

// Query 查询接口
func (h crt) Query(domain string) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		url := "https://crt.sh/?q=%25." + domain
		client := http.Client{Timeout: timeout}
		resp, err := client.Get(url)
		if err != nil {
			log.Info("error while fetching crt.sh:", err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Info("error while fetching crt.sh:", err)
			return
		}
		// 获取域名
		reRegexDomain, _ := regexp.Compile(`<TD>([a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)+)</TD>`)
		domainUrls := reRegexDomain.FindAllStringSubmatch(string(body), -1)
		for _, domain := range domainUrls {
			if len(domain) == 2 {
				ch <- domain[1]
			}
		}
	}()
	return ch
}
