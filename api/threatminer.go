package api

import (
	"net/http"
	"github.com/dean2021/dnsbrute/log"
	"regexp"
	"io/ioutil"
)

type threatminer struct{}

func init() {
	registerAPI(threatminer{})
}

// Name 接口名称
func (h threatminer) Name() string {
	return "threatminer"
}

// Query 查询接口
func (h threatminer) Query(domain string) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		url := "https://www.threatminer.org/getData.php?e=subdomains_container&q=" + domain + "&t=0&rt=10&p=1"
		client := http.Client{Timeout: timeout}
		resp, err := client.Get(url)
		if err != nil {
			log.Info("error while fetching www.threatminer.org:", err)
			return
		}
		defer resp.Body.Close()
		reRegexText, _ := regexp.Compile(`href="domain\.php\?q=(.*?)"`)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Info("error while fetching www.threatminer.org:", err)
			return
		}
		domainArr := reRegexText.FindAllStringSubmatch(string(body), -1)
		for _, v := range domainArr {
			ch <- v[1]
		}
	}()

	return ch
}
