package main

import (
	"net/http"
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"net/url"
	"net/http/httputil"
)

func getPort() string {
	port := os.Getenv("PORT")
	return ":"+ port
}
func logSetup() {
	aURL := os.Getenv("A_CONDITION_URL")
	bURL := os.Getenv("B_CONDITION_URL")
	defultURL := os.Getenv("DEFAULT_CONDITION_URL")
	fmt.Printf("redirect to URL a : %s\n", aURL)
	fmt.Printf("redirect tor URL b: %s\n", bURL)
	fmt.Printf("redirect to default URL: %s\n", defultURL)
	fmt.Printf("server is running on %s\n", getPort())

}

type requestPayloadStruct struct{
	ProxyCondition string `json:"proxy_condition"`
}

func requestBodyDecoder(req *http.Request) *json.Decoder {
	body, err  := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v", err)
		panic(err)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(body)))
}

func parseReq(req *http.Request) requestPayloadStruct{
	decoder := requestBodyDecoder(req)
	var payload requestPayloadStruct
	err := decoder.Decode(&payload)
	if err != nil {
		panic(err)
	}
	return payload

}

func getProxyURL(payload string) string {
	switch payload {
	case "A":
		return os.Getenv("A_CONDITION_URL")
	case "B":
		return os.Getenv("B_CONDITION_URL")
	default:
		return os.Getenv("DEFAULT_CONDITION_URL")
	}
}

// serves a reverse proxy from a given url
func serveReverseProxy(u string, w http.ResponseWriter, req *http.Request) {
	// parse the request
	url, _ := url.Parse(u)

	// update the header ro allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host
	fmt.Println("url", url)
	// create new reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ServeHTTP(w, req)
}

// Log the typeform payload and redirect url
func logRequestPayload(requestionPayload requestPayloadStruct, proxyURL string) {
	fmt.Printf("proxy_condition: %s, proxy_url: %s\n", requestionPayload.ProxyCondition, proxyURL)
}


func handleRequestAndRedirect(w http.ResponseWriter, req *http.Request) {
		body := parseReq(req)
		url := getProxyURL(body.ProxyCondition)
		logRequestPayload(body, url)
		serveReverseProxy(url, w, req)

}
func reverseServer1(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("reverse server 1"))
	return
}

func reverseServer2(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("reverse server 2"))
}

func reverseServer3(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("reverse server 3"))
}



func main() {
	logSetup()

	http.HandleFunc("/a", reverseServer1)
	http.HandleFunc("/b", reverseServer2)
	http.HandleFunc("/default", reverseServer3)
	http.HandleFunc("/", handleRequestAndRedirect)
	err := http.ListenAndServe(getPort(), nil)
	if err != nil {
		panic(err)
	}
	
	
}