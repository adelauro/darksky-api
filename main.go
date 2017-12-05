package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

//GeoIP freegeoip struct
//{"ip":"127.0.0.1","country_code":"","country_name":"","region_code":"","region_name":"","city":"","zip_code":"","time_zone":"","latitude":0,"longitude":0,"metro_code":0}
type GeoIP struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"cit"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MetroCode   int64   `json:"metro_code"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	strs := strings.Split(r.RemoteAddr, ":")
	log.Println("GeoIP: " + strs[0])
	resp, err := http.Get("http://freegeoip.net/json/" + strs[0])
	if err != nil || resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		w.Write([]byte("Error GET freegeoip.net"))
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		w.Write([]byte("freegeoip.net read response"))
		return
	}
	defer r.Body.Close()
	g := &GeoIP{}
	if err = json.Unmarshal(body, g); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		w.Write([]byte("freegeoip.net content"))
		return
	}
	log.Println(g)
	k, ok := os.LookupEnv("DARKSKY_KEY")
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("DARKSKY_KEY not set")
		return
	}
	strs = []string{}
	for i := 0; i > -7; i--{
		t := time.Now().AddDate(0, 0, i)
		url := "https://api.darksky.net/forecast/" + k + "/" + strings.Join([]string{
			strconv.FormatFloat(g.Latitude, 'f', -1, 64),
			strconv.FormatFloat(g.Longitude, 'f', -1, 64),
			strconv.Itoa(int(t.Unix()))}, ",")
		url += "/?exclude=currently,minutely,hourly,alerts,flags"
		log.Println(url)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			w.Write([]byte("Error GET darksky"))
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			w.Write([]byte("darksky read response"))
			return
		}
		defer r.Body.Close()
		strs = append(strs, string(body))
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	ret := "[" + strings.Join(strs, ",") + "]"
	w.Write([]byte(ret))
}

func main() {
	_, ok := os.LookupEnv("DARKSKY_KEY")
	if !ok {
		fmt.Println("DARKSKY_KEY not set exiting...")
		os.Exit(1)
	}
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", Handler).Methods("GET", "OPTIONS")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":80", r))
}
