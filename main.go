package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gwproxy/externols/github.com/openshift/origin/build/api/v1"
	"gwproxy/externols/github.com/openshift/origin/deploy/api/v1"
	"gwproxy/externols/github.com/openshift/origin/image/api/v1"
	"gwproxy/externols/github.com/openshift/origin/project/api/v1"
	"gwproxy/externols/github.com/openshift/origin/route/api/v1"
	"gwproxy/externols/github.com/openshift/origin/user/api/v1"
	"gwproxy/externols/k8s.io/kubernetes/pkg/api/v1"
	"gwproxy/externols/k8s.io/kubernetes/pkg/apis/apps/v1beta1"
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
type Everything interface {
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
	//fmt.Printf("kind\n%+v\n",kind)
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
func formatkind(kind string, body []byte) []byte {

	var structtype interface{}

	switch kind {

	case "ProjectList":
		structtype = ocproject.ProjectList{}

	case "User":
		structtype = ocuser.User{}

	case "BuildList":
		structtype = build.BuildList{}

	case "BuildConfigList":
		structtype = build.BuildConfigList{}

	case "BuildConfig":
		structtype = build.BuildConfig{}

	case "ImageStreamList":
		structtype = image.ImageStreamList{}

	case "ImageStreamTag":
		structtype = image.ImageStreamTag{}

	case "ImageStream":
		structtype = image.ImageStream{}

	case "ImageStreamImage":
		structtype = image.ImageStreamImage{}

	case "DeploymentConfig":
		structtype = deploy.DeploymentConfig{}

	case "DeploymentConfigList":
		structtype = deploy.DeploymentConfigList{}

	case "DeploymentList":
		structtype = k8sapis.DeploymentList{}

	case "ReplicationControllerList":
		structtype = k8sv1.ReplicationControllerList{}

	case "PodList":
		structtype = k8sv1.PodList{}

	case "Pod":
		structtype = k8sv1.Pod{}

	case "ServiceList":
		structtype = k8sv1.ServiceList{}

	case "Service":
		structtype = k8sv1.Service{}

	case "EndpointsList":
		structtype = k8sv1.EndpointsList{}

	case "RouteList":
		structtype = route.RouteList{}

	case "Route":
		structtype = route.Route{}

	case "SecretList":
		structtype = k8sv1.SecretList{}

	case "Secret":
		structtype = k8sv1.Secret{}

	case "EventList":
		structtype = k8sv1.EventList{}

	case "ConfigMapList":
		structtype = k8sv1.ConfigMapList{}

	case "ConfigMap":
		structtype = k8sv1.ConfigMap{}

	case "PersistentVolumeClaimList":
		structtype = k8sv1.PersistentVolumeClaimList{}

	case "PersistentVolumeClaim":
		structtype = k8sv1.PersistentVolumeClaim{}

	default:
		fmt.Printf("%+v\n", kind)
		structtype = "missmatch"

	}
	if structtype == "missmatch" {

		return body

	}
	err := json.Unmarshal(body, &structtype)

	if err != nil {
		panic(err)
	}

	body, err = json.Marshal(structtype)

	return body
}

func rewriteBody(resp *http.Response) (err error) {
	//fmt.Printf("%+v\n", resp.StatusCode)
	if resp.StatusCode == 101 {
		return nil
	}

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

	b = formatkind(kinds.Kind, b)

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
