package api

import (
	"net/http"
	"regexp"
	"io/ioutil"
	"strings"
	"fmt"
	"github.com/dean2021/dnsbrute/log"
)

type netcraft struct{}

func init() {
	registerAPI(netcraft{})
}

// Name 接口名称
func (h netcraft) Name() string {
	return "netcraft"
}

func (h netcraft) parser(url string) {
	url = strings.Replace(url, " ", "+", -1)
	url = fmt.Sprintf("http://searchdns.netcraft.com%s", url)
	client := http.Client{Timeout: timeout}
	resp, err := client.Get(url)
	if err != nil {
		log.Info("error while fetching searchdns.netcraft.com:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Info("error while fetching searchdns.netcraft.com:", err)
		return
	}
	// 获取域名
	// 获取下一页
	reRegexDomain, _ := regexp.Compile(`<a href="http://(.*?)/" rel="nofollow">`)
	domainUrls := reRegexDomain.FindAllStringSubmatch(string(body), -1)
	for _, domain := range domainUrls {
		ch <- domain[1]
	}
	// 获取下一页
	reRegexText, _ := regexp.Compile(`<A href="(.*?)"><b>Next page</b></a>`)
	nextUrls := reRegexText.FindStringSubmatch(string(body))
	if len(nextUrls) == 2 {
		nextUrl := nextUrls[1]
		h.parser(nextUrl)
	}
}

var ch = make(chan string)
// Query 查询接口
func (h netcraft) Query(domain string) <-chan string {
	go func() {
		defer close(ch)
		url := fmt.Sprintf("/?restriction=site+contains&position=limited&host=%s", domain)
		h.parser(url)
	}()
	return ch
}
