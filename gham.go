package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [callsign]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func getJSON(url string) interface{} {
	resp, err := http.Get(url)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err.Error())
	}

	var i interface{}
	if err := json.Unmarshal(body, &i); err != nil {
		panic(err)
	}

	return i
}

func printJSON(d interface{}) {
	m := map[string]string{
		"lastUpdate":   "Last Update",
		"licName":      "License Name",
		"frn":          "Federal Registration Number (FRN)",
		"callsign":     "Call Sign",
		"categoryDesc": "Category Description",
		"statusDesc":   "License Status",
		"expiredDate":  "Date Expires",
		"licenseID":    "License ID",
		"licDetailURL": "FCC URL",
	}

	var s = make([]string, len(m))

	a := d.(map[string]interface{})

	for k, v := range a {
		switch vv := v.(type) {
		case string:
			if x, ok := m[k]; ok {
				s = append(s, "\x1b[32;1m"+x+"\x1b[0m"+": "+vv+"\n")
			}
		case map[string]interface{}:
			printJSON(vv)
		case []interface{}:
			for _, vvv := range vv {
				printJSON(vvv)
			}
		}
	}

	sort.Strings(s)
	fmt.Println(strings.Join(s, ""))
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("missing callsign.")
		os.Exit(1)
	}

	var url = "http://data.fcc.gov/api/license-view/basicSearch/getLicenses?searchValue=%S&format=json"
	url = strings.Replace(url, "%S", strings.ToUpper(args[0]), 1)

	printJSON(getJSON(url))

	os.Exit(0)
}
