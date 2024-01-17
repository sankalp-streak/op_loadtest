package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type ress struct {
	mu    sync.Mutex
	suc   float32
	fail  float32
	Total float32
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

func flowTest(curls []string, common []string) ress {
	curls = append(curls, shuffleStrings(common)[:200]...)
	r := ress{suc: 0, fail: 0, Total: 0}
	for i := 0; i < len(curls); i++ {
		st := RunCurl(curls[i])
		if st == -1 {
			continue
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

func flowTester(users int) ress {
	curls := getCurls("curls/loggedout.sh")
	curls = append(curls, getCurls("curls/login.sh")...)
	common := getCurls("curls/common.sh")

	r := ress{suc: 0, fail: 0, Total: 0}
	wg := sync.WaitGroup{}

	for i := 0; i < users; i++ {
		wg.Add(1)
		go func() {
			x := flowTest(curls, common)
			r.mu.Lock()
			r.suc += x.suc
			r.fail += x.fail
			r.Total += x.Total
			r.mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	return r
}

func multiSeatch(users int) ress {
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
			st := RunCurl(pick)
			r.mu.Lock()
			r.Total++
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

func coreAppLoadTest(users int) ress {
	curls := filterCoreAppStrings(getCurls("curls/common.sh"), "https://api-op.streak.tech/")
	l := len(curls)
	fmt.Println("Total Curls", len(curls))

	r := ress{suc: 0, fail: 0, Total: 0}
	wg := sync.WaitGroup{}
	for i := 0; i < users; i++ {
		randomIndex := rand.Intn(l)
		pick := curls[randomIndex]
		wg.Add(1)
		go func() {
			st := RunCurl(pick)
			r.mu.Lock()
			r.Total++
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

func finalLoadTest(users int, TestName string) ress {
	var r ress
	if TestName == "multisearch" {
		r = multiSeatch(users)
	} else if TestName == "flow" {
		r = flowTester(users)
	} else if TestName == "core" {
		r = coreAppLoadTest(users)
	}
	return r
}
