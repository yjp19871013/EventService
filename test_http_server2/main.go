package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

func main() {
	http.HandleFunc("/", index)

	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		panic(err)
	}
}
