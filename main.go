package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getIP() string {
	url := "https://api.ipify.org?format=text"
	fmt.Printf("Getting IP address from ipify\n")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("My IP is:%s\n", ip)
	return string(ip)
}

// get geolocation with ip address
func getGeolocation(ip_address string) {
	fmt.Print("Getting geolocation\n")
	url := "https://ipwhois.app/json/" + ip_address
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Geolocation:%s\n", body)
}

func main() {
	my_ip := getIP()
	getGeolocation(my_ip)
}
