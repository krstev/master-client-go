package main

import (
	"net/http"
	"math/rand"
	"strconv"
	"github.com/ArthurHlt/go-eureka-client/eureka"
	"github.com/hudl/fargo"
	"fmt"
	"bytes"
	"encoding/json"
)

var client = eureka.NewClient([]string{
	"http://127.0.0.1:8761/eureka",
})

type Instance struct {
	HostName string `xml:"application>instance>hostName"`
}

type PrimeNumbers struct {
	PrimeNumbers []int
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	//client := eureka.NewClient([]string{
	//	"http://127.0.0.1:8761/eureka",
	//})
	////instance := eureka.InstanceInfo("test")
	//instance := eureka.NewInstanceInfo("test.com", "test", "69.172.200.235", 80, 30, false) //Create a new instance to register
	//
	////instance.Metadata = &eureka.MetaData{
	////	Map: make(map[string]string),
	////}
	////instance.Metadata.Map["foo"] = "bar"
	//client.RegisterInstance("servicego", instance)
	w.Header().Set("Contetn-Type", "application/json")
	bs := []byte(strconv.Itoa(rand.Intn(100)))
	w.Write(bs)
}

func register() {

	instance := eureka.NewInstanceInfo("127.0.0.1", "servicego", "127.0.0.1", 5555, 30, false)

	client.RegisterInstance("1", instance)
	body, _ := client.GetApplication("SERVICE")

	for _, e := range body.Instances {
		//fmt.Println(e.HostName,":",e.Port.Port)
		//url := e.HomePageUrl + "countVowels?text=qweqwe"
		url := e.HomePageUrl + "primeNumbers?lower=0&upper=10"
		resp, _ := http.Get(url)
		var pn PrimeNumbers
		err := json.NewDecoder(resp.Body).Decode(&pn);
		fmt.Println(url)
		fmt.Println(err)
		fmt.Println(pn.PrimeNumbers)
		fmt.Println(resp.Body)
		fmt.Println(e.StatusPageUrl)
		fmt.Println(e.HomePageUrl)
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		i, _ := strconv.ParseInt(buf.String(), 10, 64)

		//a := buf.String()
		//b := buf.Bytes()

		//c, _ := binary.ReadVarint(buf)
		//i, _ := strconv.ParseInt(a, 10, 64)
		//fmt.Println(i + 25)
		//fmt.Println(b)
		fmt.Println(i)

	}

	//fmt.Println(body.Instances)
}

func makeConnection() fargo.EurekaConnection {
	var c fargo.Config
	c.Eureka.ServiceUrls = []string{"http://127.0.0.1:8761/eureka"}
	c.Eureka.ConnectTimeoutSeconds = 10
	c.Eureka.PollIntervalSeconds = 30
	c.Eureka.Retries = 3
	return fargo.NewConnFromConfig(c)
}

func registerFargo() {
	c := makeConnection()
	//datacentar := eureka.DataCenterInfo{Name: "local", Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo"}
	//instance := fargo.Instance{HostName: "http://127.0.0.1", App: "client", IPAddr: "127.0.0.1", Port: 5555, PortEnabled: true, SecurePort: 443, SecurePortEnabled: false}
	//c.RegisterInstance(&instance)
	//e, _ := fargo.NewConnFromConfigFile("/etc/fargo.gcfg")
	fmt.Println(c.GetApp("MASTER-CLIENT"))
}

func main() {

	//registerFargo()
	register()
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/primeNumbers/", indexHandler)
	//http.HandleFunc("/agg/", newsAggHandler)
	http.ListenAndServe(":5555", nil)

}
