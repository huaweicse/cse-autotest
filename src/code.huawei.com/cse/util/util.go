package util

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"code.huawei.com/cse/common"
	"net/http"
	"time"
)

// URLParameter maintains the list of parameters to be added in URL
type URLParameter map[string]string

func FormatURL(scheme, format string, v ...interface{}) string {
	return fmt.Sprintf("%s://%s", scheme, fmt.Sprintf(format, v...))
}

func FormateHttpUrl(format string, v ...interface{}) string {
	return FormatURL(common.SchemeHttp, format, v...)
}

func FncodeParams(params []URLParameter) string {
	encoded := []string{}
	for _, param := range params {
		for k, v := range param {
			encoded = append(encoded, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
		}
	}
	return strings.Join(encoded, "&")
}

var myHostName string

func HostName() string {
	if myHostName != "" {
		return myHostName
	}
	if n, err := os.Hostname(); err != nil {
		myHostName = "unknown"
	} else {
		myHostName = n
	}
	return myHostName
}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

func IsContentTypeJson(h http.Header) bool {
	if h == nil {
		return false
	}
	t := h[common.HeaderContentType]
	if t == nil || len(t) == 0 {
		return false
	}
	for _, v := range t {
		if i := strings.Index(v, common.ContentTypeJson); i != -1 {
			return true
		}
	}
	return false
}
