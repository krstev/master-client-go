package service

import (
	"bytes"
	"encoding/json"
	"github.com/ArthurHlt/go-eureka-client/eureka"
	"master-client-go/src/model"
	"math/rand"
	"net/http"
	"strconv"
	"io/ioutil"
	"fmt"
	"time"
)

var v1ApiVersion = "v1"
var v2ApiVersion = "v2"

func GetInstaces(appId string) []eureka.InstanceInfo {
	var client = eureka.NewClient([]string{
		"http://127.0.0.1:8761/eureka",
	})
	body, _ := client.GetApplication(appId)
	return body.Instances
}

func getRandomServiceInstance() string {
	instances := GetInstaces(serviceAppID)
	instance := instances[rand.Intn(len(instances))]
	return instance.HomePageUrl
}

func GetPrimes(ch chan []int, instaceUrl string, lower, upper int) {
	defer wg.Done()
	//time.Sleep(time.Duration(rand.Intn(10)) * time.Microsecond)
	url := instaceUrl + "primeNumbers?lower=" + strconv.Itoa(lower) + "&upper=" + strconv.Itoa(upper)
	resp, _ := http.Get(url)
	var pn model.PrimeNumbers

	err := json.NewDecoder(resp.Body).Decode(&pn)
	if err != nil {
		panic(err)
	}
	ch <- pn.PrimeNumbers
}

func getVowels(ch chan int64, instaceUrl, text string) {
	defer wg.Done()
	url := instaceUrl + "countVowels?text=" + text
	resp, _ := http.Get(url)

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	i, _ := strconv.ParseInt(buf.String(), 10, 64)

	ch <- i
}

func searchWebReplica(c chan string, search string, replicaNumber int) {
	channel := make(chan string, replicaNumber)
	for i := 0; i < replicaNumber; i++ {
		go searchWeb(channel, search, v2ApiVersion)
	}
	c <- <-channel
}

func searchImageReplica(c chan string, search string, replicaNumber int) {
	channel := make(chan string, replicaNumber)
	for i := 0; i < replicaNumber; i++ {
		go searchImage(channel, search, v2ApiVersion)
	}
	c <- <-channel
}

func searchVideoReplica(c chan string, search string, replicaNumber int) {
	newC := make(chan string, replicaNumber)
	for i := 0; i < replicaNumber; i++ {
		go searchVideo(newC, search, v2ApiVersion)
	}
	c <- <-newC
}

func searchWeb(c chan string, search string, apiVersion string) {
	googleSearch(c, search, "web", apiVersion)
}

func searchImage(c chan string, search string, apiVersion string) {
	googleSearch(c, search, "image", apiVersion)
}

func searchVideo(c chan string, search string, apiVersion string) {
	googleSearch(c, search, "video", apiVersion)
}

func googleSearch(c chan string, search string, media string, apiVersion string) {

	url := getRandomServiceInstance() + "/google/" + apiVersion + "/" + media + "?search" + search
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	c <- string(bodyBytes)
}

func getV1InstanceUrls() []string {
	var result []string
	instances := GetInstaces(serviceAppID)

	for _, i := range instances {
		result = append(result, i.HomePageUrl)
	}
	return result
}

func getV2InstanceUrls() [] string {
	var result []string
	url := "http://127.0.0.1:8000/api/v1/namespaces/default/endpoints/" + minikubeService
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	//res
	//err := json.NewDecoder(resp.Body).Decode(&request)
	return result
}

func PrimeNumbers(w http.ResponseWriter, r *http.Request, urls []string){
	startTime := logStart()
	defer r.Body.Close()
	var request model.PrimeRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		panic(err)
	}
	numOfInstance := len(urls)

	ch := make(chan []int, numOfInstance)
	wg.Add(numOfInstance)

	if numOfInstance == 0 {
		panic("No available instances")
	}

	scope := (request.Limit / numOfInstance) + 1

	start := 0
	end := scope

	for _, e := range urls {
		if end > request.Limit {
			end = request.Limit
		}
		go GetPrimes(ch, e, start, end)
		start = end + 1
		end = end + scope
		if end > request.Limit {
			end = request.Limit
		}
	}

	wg.Wait()
	close(ch)
	var results []int
	for result := range ch {
		results = append(results, result...)
	}
	if err := json.NewEncoder(w).Encode(results); err != nil {
		panic(err)
	}
	logEnd(startTime)
}

func CountVowels(w http.ResponseWriter, r *http.Request,urls []string){
	startTime := logStart()

	defer r.Body.Close()
	var request model.VowelsRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		panic(err)
	}

	text := request.Text
	numOfInstance := len(urls)

	if numOfInstance == 0 {
		panic("No available instances")

	}

	ch := make(chan int64, numOfInstance)
	wg.Add(numOfInstance)

	length := len(request.Text)
	scope := (length / numOfInstance) + 1

	start := 0
	end := scope

	for _, e := range urls {
		t := text[start:end]
		go getVowels(ch, e, string(t))
		start = end
		end = end + scope
		if end > length {
			end = length
		}
	}

	wg.Wait()
	close(ch)
	var results int64
	for result := range ch {
		results = results + result
	}
	if err := json.NewEncoder(w).Encode(results); err != nil {
		panic(err)
	}
	logEnd(startTime)
}

func logStart() time.Time {
	startTime := time.Now()
	fmt.Println("Start time : ", startTime)
	return startTime
}

func logEnd(startTime time.Time) {
	endTime := time.Now()
	fmt.Println("End time : ", endTime)
	fmt.Println("Total : ", endTime.Sub(startTime).String())
}

