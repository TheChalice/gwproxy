package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gwproxy/externols/github.com/openshift/origin/deploy/api/v1"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
)

type handle struct {
	host string
}
type kinds struct {
	Kind string `json:"kind"`
}

func regurl(r *http.Request) string {

	url := r.URL.String()

	clusterurl := ""

	regws := regexp.MustCompile("access_token")

	if regws.MatchString(url) {

		clusterurl = r.FormValue("Cluster")

	} else {

		clusterurl = r.Header["Cluster"][0]

	}
	return clusterurl
}

func (this *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	cluster := regurl(r)
	//fmt.Printf("cluster\n%+v\n",cluster)
	//fmt.Printf("Header\n %+v\n", r.Header["Origin"][0])
	remote, err := url.Parse("https://" + cluster)

	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)

	proxy.ModifyResponse = rewriteBody

	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	proxy.ServeHTTP(w, r)
}
func rewriteBody(resp *http.Response) (err error) {
	b, err := ioutil.ReadAll(resp.Body) //Read html

	if err != nil {
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		return err
	}

	kinds := kinds{}

	err = json.Unmarshal(b, &kinds)

	if err != nil {
		return err
	}

	if kinds.Kind == "DeploymentConfig" {
		DeploymentConfig := deploy.DeploymentConfig{}

		err = json.Unmarshal(b, &DeploymentConfig)
		//if err!=nil{
		//	return err
		//}
		fmt.Printf("Namespace \n%+v\n", DeploymentConfig.ObjectMeta.Namespace)
	}

	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1) // replace html

	body := ioutil.NopCloser(bytes.NewReader(b))

	resp.Body = body

	resp.ContentLength = int64(len(b))

	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))

	return nil
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
