package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Continent string

const (
	NA Continent = "North America (NA)"
	SA           = "South America"
	A            = "Arctic"
	AA           = "Antarctic"
	AF           = "Africa"
	AS           = "Asia"
	AU           = "Australia"
	EU           = "Europe"
)

type IPLocation struct {
	As           string
	City         string
	Country      string
	CountryCode  string
	ISP          string
	Latitude     float32 `json:"lat"`
	Longitude    float32 `json:"lon"`
	Organization string  `json:"org"`
	IP           string  `json:"query"`
	Region       string
	RegionName   string
	Status       string
	TimeZone     string
	Zip          string
	Message      string
}

func printLocale(locale IPLocation) {

	// fmt.Println("locale: ", locale)
	fmt.Println()
	fmt.Println("As: ", locale.As)
	fmt.Println("City: ", locale.City)
	fmt.Println("Country: ", locale.Country)
	fmt.Println("CountryCode: ", locale.CountryCode)
	fmt.Println("IP: ", locale.IP)
	fmt.Println("ISP: ", locale.ISP)
	fmt.Println("Latitude: ", locale.Latitude)
	fmt.Println("Longitude: ", locale.Longitude)
	fmt.Println("Organization: ", locale.Organization)
	fmt.Println("Region: ", locale.Region)
	fmt.Println("Status: ", locale.Status)
	fmt.Println("TimeZone: ", locale.TimeZone)
	fmt.Println("Zip: ", locale.Zip)
	fmt.Println()
}

func readIP(args []string) (string, string) {

	var ip string
	if len(os.Args) < 2 {
		fmt.Println("Please enter the IP you want to check!")
		reader := bufio.NewReader(os.Stdin)
		ip, _ = reader.ReadString('\n')

		// Split away the DOS "\r\n" or linux "\n" line ending
		ip = strings.TrimSuffix(ip, "\n")
		ip = strings.TrimSuffix(ip, "\r")

		fmt.Println("Did you know that you can provide the ip as a command line argument?")
	} else {
		ip = os.Args[1]
	}

	// Validate that input is an IPv4
	if validateIP(ip) {
		return ip, ""
	} else {
		return "", "Not an IPv4."
	}
}

func validateIP(ip string) bool {

	// fmt.Println("IP: ", ip)
	if ip == "" {
		return false
	}

	address := strings.Split(ip, ".")
	if len(address) > 4 {
		return false
	}
	for _, v := range address {
		i, err := strconv.Atoi(v)
		if err != nil || i > 255 || v == "" {
			return false
		}
	}
	return true
}

func locateIP(ip string) IPLocation {

	// I use a different endpoint because IPvoid is responding with a rendered HTML page
	endpoint := "http://ip-api.com/json/"

	// Append url query which is the ip
	url := endpoint + ip

	// An http client to send the response
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	// Request the API
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	// Read response body as bytes
	localeBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parsing json-bytes into the corresponding struct
	locale := IPLocation{}
	err = json.Unmarshal(localeBytes, &locale)
	if err != nil {
		log.Fatal(err)
	}

	return locale
}

func main() {
	fmt.Println("Hei Vlad!")
	ip, err := readIP(os.Args)
	if err != "" {
		fmt.Println("Please provide an IPv4.")
	} else {
		locale := locateIP(ip)
		printLocale(locale)
	}
}
