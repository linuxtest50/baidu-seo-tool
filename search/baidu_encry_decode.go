package search

import (
	"errors"
	"github.com/kevin-zx/http-util"
	"net/http"
	"strings"
	"time"
)

func DecodeBaiduEncURL(baiduUrl string) string {
	response, _, err := GetWebconAndRealUrlFromBaiduUrl(baiduUrl)
	if err != nil {
		return ""
	}
	return response.Request.URL.String()
}

func GetWebconAndRealUrlFromBaiduUrl(baiduUrl string) (response *http.Response, webCon string, err error) {
	response, err = http_util.SendRequest(baiduUrl, map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36"}, "GET", nil, 10*time.Second)
	if err != nil {
		return nil, "", err
	}
	//fmt.Println(baiduUrl)
	realURL := ""
	if response.StatusCode == 200 {
		webCon, _ = http_util.ReadContentFromResponse(response, "")
		if strings.Contains(webCon, "window.opener&&window.opener.bds&&window.opener.bds.pdc&&window.opener.bds.pdc.sendLinkLog") {
			part1 := strings.Split(webCon, "window.location.replace(\"")
			if len(part1) < 2 {
				return nil, "", errors.New("can't go real page")
			} else {
				realURL = strings.Split(part1[1], "\")},timeout")[0]
			}
		} else if strings.Contains(webCon, `JSON.parse(localStorage.getItem("tc_time_log")`) {
			ps := strings.Split(webCon, "\n")
			for _, p := range ps {
				if strings.Contains(p, "window.location.replace(") && strings.Contains(p, ")") {
					start := strings.Index(p, `("`)
					end := strings.LastIndex(p, `")`)
					if end > start+1 && (start > 0 && end > 0) {
						realURL = p[start+2 : end]
					}
				}
			}
		} else {
			return
		}

		if realURL == "" {
			return
		} else {
			response, err = http_util.GetWebResponseFromUrlWithHeader(realURL, map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36"})
			if err != nil {
				return
			}
			webCon, err = http_util.ReadContentFromResponse(response, "")
			if err != nil {
				return
			}
			return
		}

	}
	return nil, "", errors.New("wrong status")
}
