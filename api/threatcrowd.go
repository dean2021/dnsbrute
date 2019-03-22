package api

import (
	"net/http"
	"github.com/dean2021/dnsbrute/log"
	"io/ioutil"
	"encoding/json"
)

type threatcrowd struct{}

func init() {
	registerAPI(threatcrowd{})
}

// Name 接口名称
func (h threatcrowd) Name() string {
	return "threatcrowd"
}

// Query 查询接口
func (h threatcrowd) Query(domain string) <-chan string {
	ch := make(chan string)
	type ThreatCrowd struct {
		Emails     []string      `json:"emails"`
		Hashes     []string      `json:"hashes"`
		Permalink  string        `json:"permalink"`
		References []interface{} `json:"references"`
		Resolutions []struct {
			IPAddress    string `json:"ip_address"`
			LastResolved string `json:"last_resolved"`
		} `json:"resolutions"`
		ResponseCode string   `json:"response_code"`
		Subdomains   []string `json:"subdomains"`
		Votes        int      `json:"votes"`
	}
	go func() {
		defer close(ch)
		var domains ThreatCrowd
		url := "https://www.threatcrowd.org/searchApi/v2/domain/report/?domain=" + domain
		client := http.Client{Timeout: timeout}
		resp, err := client.Get(url)
		if err != nil {
			log.Info("error while fetching www.threatcrowd.org:", err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Info("error while fetching www.threatcrowd.org:", err)
			return
		}
		json.Unmarshal(body, &domains)
		for _, subdomain := range domains.Subdomains {
			ch <- subdomain
		}
	}()

	return ch
}
