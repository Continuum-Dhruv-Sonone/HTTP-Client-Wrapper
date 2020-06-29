package main

import (
	"bufio"
	"fmt"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/clientWrapper/client"
)

func main() {

	h := client.CircuitBreakerConfig{
		CircuitConfig: hystrix.CommandConfig{
			Timeout:                1,
			MaxConcurrentRequests:  10,
			ErrorPercentThreshold:  20,
			RequestVolumeThreshold: 5,
			SleepWindow:            5,
		},
		BaseURL: "http://localhost:9090/hello",
	}
	err := client.Register([]client.CircuitBreakerConfig{h})
	if err != nil {
		fmt.Println(err)
	}

	for index := 0; index < 100; index++ {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:9090", nil)

		//Creating a new client
		cl := client.New()
		if err != nil {
			panic(err)
		}

		resp, err := cl.Do(req)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()

		fmt.Println("Response status:", resp.Status)

		scanner := bufio.NewScanner(resp.Body)
		for i := 0; scanner.Scan() && i < 5; i++ {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}

}
