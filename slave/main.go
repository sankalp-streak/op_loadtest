package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"optest/loadtestPB"
	"os"
	"strconv"
	"time"

	"github.com/StreakAI/prism/prismPB"
	"github.com/StreakAI/prism/slave"
	"google.golang.org/protobuf/proto"
)

var logger *log.Logger
var symbols []string

func setupLogger(port string) {
	for {
		ct := time.Now()
		y := strconv.FormatInt(int64(ct.Year()), 10)
		m := strconv.FormatInt(int64(ct.Month()), 10)
		d := strconv.FormatInt(int64(ct.Day()), 10)
		fn := "screener_slave_" + port + "_" + y + "_" + m + "_" + d + ".log"
		of, _ := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		logger = log.New(of, "", log.Ldate|log.Ltime)
		logger.Println("New log file", fn)
		time.Sleep(1 * time.Hour)
	}
}

var SLAVE_IP string
var SLAVE_PORT string
var MASTER_PATH string
var MASTER_REDIS_AUTH string

func initializeConfig() {
	c, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("Error reading config file", err)
		panic(err)
	}
	var data map[string]interface{}
	er := json.Unmarshal(c, &data)
	if er != nil {
		log.Println("Erorr unmarshalling config json", er)
		panic(er)
	}
	SLAVE_IP = data["slave_ip"].(string)
	SLAVE_PORT = data["slave_port"].(string)
	MASTER_PATH = data["master_path"].(string)
	MASTER_REDIS_AUTH = data["master_redis_auth"].(string)
}

func getIpAddress() string {
	log.Print("getting ip address")
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/public-ipv4")
	if err != nil {
		log.Print("err getting ip address from ipinfo : ", err)
		// return private ip
		return getPrivateIp()
	}
	if resp.StatusCode != http.StatusOK {
		return getPrivateIp()
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("err reading response from ipinfo : ", err)
	}
	return string(bodyBytes)
}

func getPrivateIp() string {
	log.Print("getting private ip address")
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/local-ipv4")
	if err != nil {
		log.Fatal("err getting ip address from ipinfo : ", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Error response", resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("err reading response from ipinfo : ", err)
	}
	return string(bodyBytes)
}

func resolveData(r *prismPB.TaskRequest) []loadtestPB.Request {
	req := make([]loadtestPB.Request, len(r.Data))
	for i, v := range r.Data {
		err := proto.Unmarshal(v.Value, &req[i])
		if err != nil {
			log.Fatal("could not deserialize request")
		}
	}
	return req
}

func processAlerts(tid string, req *prismPB.TaskRequest) (bool, []proto.Message) {
	logger.Println("Received slave task", tid)
	data := resolveData(req)
	users := len(data)
	resp := make([]proto.Message, users)
	fmt.Println("TOTAL HITS PLANNED:", users)
	r := finalLoadTest(users, data[0].TestName)
	fmt.Println(r)
	resp[0] = &loadtestPB.Response{Success: r.suc, Fails: r.fail, Total: r.Total}
	fmt.Println("Results===============>", resp[0])
	return true, resp
}

func main() {
	initializeConfig()
	go setupLogger(SLAVE_PORT)
	time.Sleep(2 * time.Second)

	SLAVE_IP = getIpAddress() // "localhost" //
	log.Println(SLAVE_IP)
	logger.Println("Starting screener lambda slave")
	logger.Println("Starting slave in", SLAVE_IP, SLAVE_PORT)
	s, err := slave.Init(SLAVE_IP, SLAVE_PORT)
	if err != nil {
		logger.Println("Error init slave", err)
		return
	}

	s.SetupLogger("stdout", make(map[string]string, 0))
	s.SetupRedisChannel(MASTER_PATH, MASTER_REDIS_AUTH, SLAVE_IP, SLAVE_PORT, 100, "loadtest")

	s.SetCallback(processAlerts)
	s.Serve()
}
