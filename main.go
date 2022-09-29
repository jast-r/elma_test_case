package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
)

func GoCounter(user_url string, total_count *int, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()
	go_count := 0
	if _, err := url.Parse(user_url); err != nil {
		fmt.Printf("Count of \"Go\" in %s = %d\n", user_url, go_count)
		return
	}
	resp, err := http.Get(user_url)
	if err != nil {
		fmt.Printf("Count of \"Go\" in %s = %d\n", user_url, go_count)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Count of \"Go\" in %s = %d\n", user_url, go_count)
		return
	}
	go_count = strings.Count(string(body), "Go")
	mutex.Lock()
	*total_count += go_count
	mutex.Unlock()
	fmt.Printf("Count of \"Go\" in %s = %d\n", user_url, go_count)
}

func main() {
	var wg sync.WaitGroup
	var path_to_file string
	var mutex sync.Mutex
	total_count := 0

	flag.StringVar(&path_to_file, "path", "./urls.txt", "")
	flag.Parse()

	f, err := os.Open(path_to_file)
	if err != nil {
		fmt.Println("error open file with name: " + path_to_file)
		os.Exit(-1)
	}
	defer f.Close()
	var urls []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	for _, url := range urls {
		// больше 7, т.к. основной поток тоже считается за горутину
		if runtime.NumGoroutine() > 7 {
			wg.Wait()
		}
		wg.Add(1)
		go GoCounter(url, &total_count, &wg, &mutex)
	}
	wg.Wait()
	fmt.Printf("Total count of \"Go\": %d\n", total_count)
}
