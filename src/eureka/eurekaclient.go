package eureka

import (
	"fmt"
	"go-microservice-eureka/src/github.com/eriklupander/goeureka/util"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var instanceId string

func Register(port int) {
	instanceId = util.GetUUID()
	dir, _ := os.Getwd()
	data, _ := ioutil.ReadFile(dir + "/templates/newtemplate.json")

	tpl := string(data)
	tpl = strings.Replace(tpl, "${ipAddress}", util.GetLocalIP(), -1)
	tpl = strings.Replace(tpl, "${port}", strconv.Itoa(port), -1)
	tpl = strings.Replace(tpl, "${instanceId}", instanceId, -1)

	// Register.
	registerAction := HttpAction{
		Url:         "http://127.0.0.1:8761/eureka/apps/go-client",
		Method:      "POST",
		ContentType: "application/json",
		Body:        tpl,
	}
	var result bool
	for {
		result = DoHttpRequest(registerAction)
		if result {
			break
		} else {
			time.Sleep(time.Second * 5)
		}
	}
}

func StartHeartbeat() {
	for {
		time.Sleep(time.Second * 30)
		heartbeat()
	}
}

func heartbeat() {
	heartbeatAction := HttpAction{
		Url:    "http://127.0.0.1:8761/eureka/apps/go-client/" + util.GetLocalIP() + ":go-client:" + instanceId,
		Method: "PUT",
	}
	DoHttpRequest(heartbeatAction)
}

func Deregister() {
	fmt.Println("Trying to deregister application...")
	// Deregister
	deregisterAction := HttpAction{
		Url:    "http://127.0.0.1:8761/eureka/apps/go-client/" + util.GetLocalIP() + ":go-client:" + instanceId,
		Method: "DELETE",
	}
	DoHttpRequest(deregisterAction)
	fmt.Println("Deregistered application, exiting. Check Eureka...")
}
