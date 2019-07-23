package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

type handle struct {
	host string
}

func regurl(r *http.Request) string {

	url :=r.URL.String()

	clusterurl :=""

	regws := regexp.MustCompile("access_token")

	if regws.MatchString(url) {

		clusterurl=r.FormValue("Cluster")

	}else {

		clusterurl=r.Header["Cluster"][0]

	}
	return clusterurl
}

func (this *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	cluster:=regurl(r)

	//fmt.Printf("cluster\n%+v\n",cluster)

	//fmt.Printf("Header\n %+v\n", r.Header["Origin"][0])

	remote, err := url.Parse("https://"+cluster)

	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)

	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	proxy.ServeHTTP(w, r)
}

func startServer() {
	//被代理的服务器host和port
	h := &handle{}
	log.Print("Service starting in 10001")

	err := http.ListenAndServe(":10001", h)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
func main() {
	startServer()
}
