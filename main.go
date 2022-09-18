package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func GoCounter(user_url string, out_chan chan map[string]int) {
	if _, err := url.ParseRequestURI(user_url); err != nil {
		out_chan <- map[string]int{user_url: 0}
		return
	}
	resp, err := http.Get(user_url)
	if err != nil {
		out_chan <- map[string]int{user_url: 0}
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		out_chan <- map[string]int{user_url: 0}
		return
	}
	out_chan <- map[string]int{user_url: strings.Count(string(body), "Go")}
}

func main() {
	output := make(chan map[string]int)
	var path_to_file string
	flag.StringVar(&path_to_file, "path", "./test_urls.txt", "")
	flag.Parse()

	f, err := os.Open(path_to_file)
	if err != nil {
		fmt.Println("error open file with name: " + path_to_file)
		os.Exit(-1)
	}
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for _, line := range lines {
		go GoCounter(line, output)
	}

	total_count := 0
	for val := range output {
		for k, v := range val {
			fmt.Printf("Count of Go for %s = %d\n", k, v)
			total_count += v
		}
	}
	fmt.Printf("Total count: %d", total_count)
}
