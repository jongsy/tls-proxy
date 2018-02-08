package main

import (
	"net/http/httputil"
	"net/http"
	"fmt"
	"log"
	"strconv"
	"flag"
	"bytes"
	"io/ioutil"
)

func NewReverseProxy(host string) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = "https"
		req.URL.Host = host
		req.Host = host

		fmt.Printf("%+v\n", req)

	}
	return &httputil.ReverseProxy{Director: director}
}

func rewriteBody(resp *http.Response) (err error) {
	b, err := ioutil.ReadAll(resp.Body) //Read html
	if err != nil {
		return  err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}
	//b = bytes.Replace(b, []byte("denied"), []byte("schmerver"), -1) // replace html
	s := string(b[:])
	fmt.Println(s)
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return nil
}

func main() {
	portFlag := flag.Int("port", 8087, "port. default: 8087")
	hostFlag := flag.String("host", "google.com", "host. default: google.com")
	flag.Parse()
	port := ":" + strconv.Itoa(*portFlag)
	host := *hostFlag

	proxy := NewReverseProxy(host)
	proxy.ModifyResponse = rewriteBody
	log.Fatal(http.ListenAndServe(port, proxy))
}