package api

import (
	"net/http"
	"github.com/dean2021/dnsbrute/log"
	"io/ioutil"
	"encoding/json"
)

type sublist3r struct{}

func init() {
	registerAPI(sublist3r{})
}

// Name 接口名称
func (h sublist3r) Name() string {
	return "sublist3r"
}

// Query 查询接口
func (h sublist3r) Query(domain string) <-chan string {
	ch := make(chan string)

	type Domain []interface{}
	go func() {
		defer close(ch)
		var domains Domain
		url := "https://api.sublist3r.com/search.php?domain=" + domain
		client := http.Client{Timeout: timeout}
		resp, err := client.Get(url)
		if err != nil {
			log.Info("error while fetching api.sublist3r.com:", err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Info("error while fetching api.sublist3r.com:", err)
			return
		}
		json.Unmarshal(body, &domains)
		for _, subdomain := range domains {
			ch <- subdomain.(string)
		}
	}()

	return ch
}
