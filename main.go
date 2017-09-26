package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	DEFAULT_PORT     = "8080"
	CF_FORWARDED_URL = "X-Cf-Forwarded-Url"
	DEFAULT_LIMIT    = 10
)

var (
	whiteList *WhiteList
)

func main() {
	log.SetOutput(os.Stdout)

	list := os.Getenv("WHITE_LIST")
  rawWhiteList := strings.Split(list, ",")
	sort.Strings(rawWhiteList)
  whiteList = NewWhiteList(rawWhiteList)

	log.Printf("whiteList is [%s]\n", strings.Join(whiteList.whiteList, ", "))

	http.Handle("/", newProxy())
	log.Fatal(http.ListenAndServe(":"+getPort(), nil))
}

func newProxy() http.Handler {
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			forwardedURL := req.Header.Get(CF_FORWARDED_URL)

			url, err := url.Parse(forwardedURL)
			if err != nil {
				log.Fatalln(err.Error())
			}

			req.URL = url
			req.Host = url.Host
		},
		Transport: newWhiteListedRoundTripper(),
	}
	return proxy
}

func getPort() string {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}
	return port
}

func skipSslValidation() bool {
	var skipSslValidation bool
	var err error
	if skipSslValidation, err = strconv.ParseBool(os.Getenv("SKIP_SSL_VALIDATION")); err != nil {
		skipSslValidation = true
	}
	return skipSslValidation
}
func getEnv(env string, defaultValue int) int {
	var (
		v      string
		config int
	)
	if v = os.Getenv(env); len(v) == 0 {
		return defaultValue
	}

	config, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return config
}

type WhiteListedRoundTripper struct {
	whiteList *WhiteList
	transport http.RoundTripper
}

func newWhiteListedRoundTripper() *WhiteListedRoundTripper {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSslValidation()},
	}
	return &WhiteListedRoundTripper{
		whiteList: whiteList,
		transport: tr,
	}
}

func (r *WhiteListedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var err error
	var res *http.Response

	remoteIP := strings.Split(req.RemoteAddr, ":")[0]

	log.Printf("request from [%s]\n", remoteIP)
	if !r.whiteList.contains(remoteIP) {
		resp := &http.Response{
			StatusCode: 401,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Doesn't match whitelist")),
		}
		return resp, nil
	}

	res, err = r.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	return res, err
}
