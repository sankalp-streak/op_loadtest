package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type ress struct {
	mu      sync.Mutex
	suc     float32
	fail    float32
	time    time.Duration
	Total   float32
	TimeOut float32
}

func shuffleStrings(input []string) []string {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator with the current time

	shuffled := make([]string, len(input))
	copy(shuffled, input)

	// Fisher-Yates shuffle algorithm
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}

func flowTest(curls []string, common []string, timeout int) ress {
	curls = append(curls, shuffleStrings(common)[:200]...)
	mp := map[string]int{}
	for i:=0;i<len(curls);i++{
		mp[curls[i][:55]]++
	}
	for x,y := range mp{
		fmt.Println(x,y)
	}
	r := ress{suc: 0, fail: 0, Total: 0}
	for i := 0; i < len(curls); i++ {
		st, _ := RunCurl(curls[i], timeout)
		if st == -1 {
			continue
		}
		if st == -2 {
			r.TimeOut++
		}
		r.mu.Lock()
		r.Total++
		if st >= 200 && st <= 299 {
			r.suc++
		} else if st > 299 {
			r.fail++
		}
		r.mu.Unlock()
		// Check if it's the 10th iteration
		if i%10 == 0 {
			time.Sleep(500 * time.Millisecond)
		}
	}
	fmt.Println("Failsss", r.fail)
	return r
}

func flowTester(users int, timeout int) ress {
	curls := getCurls("curls/loggedout.sh")
	curls = append(curls, getCurls("curls/login.sh")...)
	common := getCurls("curls/common.sh")

	r := ress{suc: 0, fail: 0, Total: 0}
	wg := sync.WaitGroup{}

	for i := 0; i < users; i++ {
		wg.Add(1)
		go func() {
			x := flowTest(curls, common, timeout)
			r.mu.Lock()
			r.suc += x.suc
			r.fail += x.fail
			r.Total += x.Total
			r.TimeOut += x.TimeOut
			r.time += x.time
			r.mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	return r
}

func multiSeatch(users int, timeout int) ress {
	wg := sync.WaitGroup{}
	curls := getCurls("curls/multisearch.sh")
	l := len(curls)
	fmt.Println(len(curls))

	r := ress{suc: 0, fail: 0, Total: 0}
	for i := 0; i < users; i++ {
		randomIndex := rand.Intn(l)
		pick := curls[randomIndex] // get the value from the slice
		wg.Add(1)
		go func() {
			rSleep := rand.Intn(6) - 1
			time.Sleep(time.Duration(rSleep) * time.Second)
			st, t := RunCurl(pick, timeout)
			r.mu.Lock()
			r.Total++
			r.time += t
			if st == -2 {
				r.TimeOut++
			}
			if st == -1 {
				r.Total--
			} else if st >= 200 && st <= 299 {
				r.suc++
			} else if st > 299 {
				r.fail++
			}
			r.mu.Unlock()
			wg.Done()
		}()

	}
	wg.Wait()
	return r
}

func stech(users int, timeout int) ress {
	wg := sync.WaitGroup{}
	curls := getCurls("curls/stech.sh")
	l := len(curls)
	fmt.Println(len(curls))

	r := ress{suc: 0, fail: 0, Total: 0}
	for i := 0; i < users; i++ {
		randomIndex := rand.Intn(l)
		pick := curls[randomIndex] // get the value from the slice
		wg.Add(1)
		go func() {
			st, t := RunCurl(pick, timeout)
			r.mu.Lock()
			r.Total++
			r.time += t
			if st == -2 {
				r.TimeOut++
			}
			if st == -1 {
				r.Total--
			} else if st >= 200 && st <= 299 {
				r.suc++
			} else if st > 299 {
				r.fail++
			}
			r.mu.Unlock()
			wg.Done()
		}()

	}
	wg.Wait()
	return r
}

func filterCoreAppStrings(input []string, substring string) []string {
	var filtered []string

	for _, str := range input {
		// Check if the substring is present in the current string
		if strings.Contains(str, substring) {
			filtered = append(filtered, str)
		}
	}

	return filtered
}

func coreAppLoadTest(file string, users int, timeout int) ress {
	curls := filterCoreAppStrings(getCurls(file), "https://api-op.streak.tech/")
	l := len(curls)
	fmt.Println("Total Curls", len(curls))

	r := ress{suc: 0, fail: 0, Total: 0}
	wg := sync.WaitGroup{}
	for i := 0; i < users; i++ {
		randomIndex := rand.Intn(l)
		pick := curls[randomIndex]
		wg.Add(1)
		go func() {
			rSleep := rand.Intn(6) - 1
			time.Sleep(time.Duration(rSleep) * time.Second)
			st, t := RunCurl(pick, timeout)
			r.mu.Lock()
			r.time += t
			r.Total++
			if st == -2 {
				r.TimeOut++
			}
			if st == -1 {
				r.Total--
			} else if st >= 200 && st <= 299 {
				r.suc++
			} else if st > 299 {
				r.fail++
			}
			r.mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	return r
}

func stable5(file string, totalCalls int, timeout int) ress{

	curls := filterCoreAppStrings(getCurls(file), "https://api-op.streak.tech/")
	l := len(curls)
	fmt.Println("Total Curls", len(curls))

	r := ress{suc: 0, fail: 0, Total: 0}
	wg := sync.WaitGroup{}

	duration := 5 * time.Minute
	interval := duration / time.Duration(totalCalls)

	// Create a ticker to trigger the function at regular intervals
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Loop to call the function 1000 times
	for i := 0; i < totalCalls; i++ {
		randomIndex := rand.Intn(l)
		pick := curls[randomIndex]
		wg.Add(1)
		go func() {
			st, t := RunCurl(pick, timeout)
			r.mu.Lock()
			r.time += t
			r.Total++
			if st == -2 {
				r.TimeOut++
			}
			if st == -1 {
				r.Total--
			} else if st >= 200 && st <= 299 {
				r.suc++
			} else if st > 299 {
				r.fail++
			}
			r.mu.Unlock()
			wg.Done()
		}()
		<-ticker.C // Wait for the next tick
	}
	wg.Wait()
	return r
}

func stable1(file string, totalCalls int, timeout int) ress{

	curls := filterCoreAppStrings(getCurls(file), "https://api-op.streak.tech/")
	l := len(curls)
	fmt.Println("Total Curls", len(curls))

	r := ress{suc: 0, fail: 0, Total: 0}
	wg := sync.WaitGroup{}

	duration := 1 * time.Minute
	interval := duration / time.Duration(totalCalls)

	// Create a ticker to trigger the function at regular intervals
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Loop to call the function 1000 times
	for i := 0; i < totalCalls; i++ {
		randomIndex := rand.Intn(l)
		pick := curls[randomIndex]
		wg.Add(1)
		go func() {
			st, t := RunCurl(pick, timeout)
			r.mu.Lock()
			r.time += t
			r.Total++
			if st == -2 {
				r.TimeOut++
			}
			if st == -1 {
				r.Total--
			} else if st >= 200 && st <= 299 {
				r.suc++
			} else if st > 299 {
				r.fail++
			}
			r.mu.Unlock()
			wg.Done()
		}()
		<-ticker.C // Wait for the next tick
	}
	wg.Wait()
	return r
}

func finalLoadTest(users int, TestName string, timeout int) ress {
	var r ress
	if TestName == "multisearch" {
		r = multiSeatch(users, timeout)
	} else if TestName == "flow" {
		r = flowTester(users, timeout)
	} else if TestName == "core" {
		r = coreAppLoadTest("curls/common.sh", users, timeout)
	} else if TestName == "core1" {
		r = coreAppLoadTest("curls/marketplace.sh", users, timeout)
	} else if TestName == "core2" {
		r = coreAppLoadTest("curls/loginapi.sh", users, timeout)
	} else if TestName == "stech" {
		r = stech(users, timeout)
	}else if TestName == "stable5" {
		r = stable5("curls/common.sh", users, timeout)
	}else if TestName == "stable1" {
		r = stable1("curls/common.sh", users, timeout)
	}else if TestName == "stable5login" {
		r = stable5("curls/loginapi.sh", users, timeout)
	}else if TestName == "stable1login" {
		r = stable1("curls/loginapi.sh", users, timeout)
	}
	
	return r
}
