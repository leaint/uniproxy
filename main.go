package main

import (
	"crypto/tls"
	"flag"
	"io"
	"net/http"
	"strings"
)

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
}

var client = &http.Client{
	Transport: tr,
	// Timeout:   20 * time.Second,
}

var servAddr = "127.0.0.22:8080"

type MyHandler struct {
}

func (h *MyHandler) ServeHTTP(ores http.ResponseWriter, oreq *http.Request) {

	method := oreq.Method

	if method == http.MethodOptions {

		ores.Header().Add("Access-Control-Allow-Origin", "*")
		ores.Header().Add("Access-Control-Allow-Headers", "X-Cookie, X-Referer, X-Origin")
		ores.WriteHeader(http.StatusNoContent)
		return
	}

	requri := strings.TrimLeft(oreq.URL.RequestURI(), "/")
	req, _ := http.NewRequest(method, requri, oreq.Body)
	req.Header = oreq.Header.Clone()
	// req.Proto = Proto

	req.Header.Set("Cookie", req.Header.Get("X-Cookie"))
	req.Header.Set("Referer", req.Header.Get("X-Referer"))
	req.Header.Set("Origin", req.Header.Get("X-Origin"))

	req.Header.Del("X-Referer")
	req.Header.Del("X-Cookie")
	req.Header.Del("X-Origin")

	res, e := client.Do(req)

	if e != nil {

		ores.WriteHeader(http.StatusInternalServerError)
		ores.Write([]byte(e.Error()))
		return
	}

	//将Cookie暴露在客户端可访问的响应头中
	xcookies := ""
	for _, v := range res.Cookies() {
		xcookies += v.Name + "=" + v.Value + ";"
	}
	if xcookies != "" {
		ores.Header().Add("X-Cookie", xcookies)
	}

	for k, v := range res.Header {
		if k == "Set-Cookie" {
			continue
		}
		for _, vi := range v {
			ores.Header().Set(k, vi)
		}
	}
	ores.Header().Set("Access-Control-Allow-Origin", "*")
	ores.Header().Set("Access-Control-Allow-Headers", "X-Cookie, X-Referer, X-Origin")
	ores.Header().Set("Access-Control-Expose-Headers", "*")
	ores.WriteHeader(res.StatusCode)

	defer res.Body.Close()

	io.Copy(ores, res.Body)

}

func main() {

	hoststr := flag.String("bind", servAddr, "ip:port")
	flag.Parse()

	handler := &MyHandler{}
	http.ListenAndServe(*hoststr, handler)
}
