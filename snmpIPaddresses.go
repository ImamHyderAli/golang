package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
)

type Device struct {
	IP string `json:"ip,omitempty"`
}

var devices []Device
var ipAddress []string

func GetDeviceEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range devices {
		if item.IP == params["ip"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Device{})
}

func GetdevicesEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(devices)
}

func main() {
	router := mux.NewRouter()
	out, _ := exec.Command("/bin/sh", "-c", "echo tcs123| sudo -S nmap -sU 10.132.32.198-255 -p 161 --open -oG - | awk '/161\\/open.*/{print $2}'").Output()
	str1 := string(out)
	strs := strings.Split(str1, "\n")
	ipAddress = strs
	for index := range strs {
		devices = append(devices, Device{IP: strs[index]})
	}
	fmt.Println("Serving at http://localhost:12345/devices ...")
	fmt.Println(ipAddress)
	router.HandleFunc("/devices", GetdevicesEndpoint).Methods("GET")
	router.HandleFunc("/devices/{ip}", GetDeviceEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":12345", router))
}
