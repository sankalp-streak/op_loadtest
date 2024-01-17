package main

import (
	"fmt"
	"log"
	"optest/loadtestPB"
	"os"
	"strconv"
	"time"

	"github.com/StreakAI/prism/master"
	"github.com/StreakAI/prism/prismPB"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var mstr *master.Master
var logger *log.Logger

func setupLogger() {
	for {
		ct := time.Now()
		y := strconv.FormatInt(int64(ct.Year()), 10)
		m := strconv.FormatInt(int64(ct.Month()), 10)
		d := strconv.FormatInt(int64(ct.Day()), 10)
		fn := "load" + y + "_" + m + "_" + d + ".log"
		of, _ := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		logger = log.New(of, "", log.Ldate|log.Ltime)
		logger.Println("New log file", fn)
		time.Sleep(1 * time.Hour)
	}
}

func initializeTaskMaster() {
	mstr, _ = master.Init()
	mstr.SetupRedisChannel("rch.streak.tech:6379", "b626aaa7a231aabf6aa3df5fc5caa847c202dc3bd0d361a1291bb2f855afe8ba7727d3e4273e9d2a4ad5020c8695bbaf9b8d051c5bf43b868795a130d482acae", "loadtest")
}

func resolveData(r *prismPB.TaskRequest) []loadtestPB.Response {
	req := make([]loadtestPB.Response, len(r.Data))
	for i, v := range r.Data {
		err := proto.Unmarshal(v.Value, &req[i])
		if err != nil {
			log.Fatal("could not deserialize request")
		}
	}
	return req
}

func main() {
	go setupLogger()
	time.Sleep(2 * time.Second)
	logger.Println("Starting screener lambda server")
	initializeTaskMaster()
	users, _ := strconv.Atoi(os.Args[1])
	fmt.Println("Number of Users:", users)
	tasks := make([]protoreflect.ProtoMessage, users)
	for i := 0; i < users; i++ {
		tasks[i] = &loadtestPB.Request{TestName: os.Args[2]}
	}

	key := "LOAD_TEST"
	logger.Println("Sending to prism master", key, len(tasks))

	time.Sleep(35 * time.Second)
	startTime := time.Now()
	resp := mstr.Do(key, tasks, &loadtestPB.Response{})
	var Success float32 = 0
	var Total float32 = 0
	var Fails float32 = 0
	for _, x := range resp {
		Success += x.(*loadtestPB.Response).Success
		Total += x.(*loadtestPB.Response).Total
		Fails += x.(*loadtestPB.Response).Fails
	}
	percent := (float64(Fails) / float64(Total)) * 100.0
	fmt.Println("Result for=> Users:", users,"Total Request:",Total, "Success:", Success, "Fails:", Fails, "Fail Percent:", percent)

	elapsedTime := time.Since(startTime)
	fmt.Printf("Master took %f Seconds\n", elapsedTime.Seconds())

}
