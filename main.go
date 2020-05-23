package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"log"
	"net/http"
	"time"
	"os"
	"strings"
	"text/tabWriter"
)

type Results struct {
	Results []Result `json:"results,omitempty"`
}

type Result struct {
	ClientName string `json:"client_name,omitempty"`
	BillableHours float32 `json:"billable_hours,omitempty"`
}

type Config struct {
	APIKey string
	AccountId string
	Connection string
}

func main() {
	config := &Config{}
	file, err := os.Open("config/config.json") 
	if err != nil {  
		log.Fatal(err)
	}  
	decoder := json.NewDecoder(file) 
	err = decoder.Decode(&config) 
	if err != nil { 
		log.Fatal(err)
	}

	APIKey := config.APIKey
	AccountId := config.AccountId

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter start date (YYMMDD): ")
	start_date, _ := reader.ReadString('\n')
	s := strings.TrimRight(start_date, "\r\n")
	fmt.Print("Enter end date (YYMMDD): ")
	end_date, _ := reader.ReadString('\n')
	e := strings.TrimRight(end_date, "\r\n")

	var url string = "https://api.harvestapp.com/v2/reports/time/projects?from="+s+"&to="+e+""

	client := http.Client {
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+APIKey+"")
	req.Header.Set("Harvest-Account-Id", ""+AccountId+"")
	req.Header.Set("User-Agent", "Harvested Reports")

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	results := &Results{}
	jsonErr := json.Unmarshal(body, results)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	// initialize tabwriter
	w := new(tabwriter.Writer)
	
	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	
	defer w.Flush()

	fmt.Fprintf(w, "\n %s\t%s\t", "Client Name", "Billable Hours")
	fmt.Fprintf(w, "\n %s\t%s\t", "----", "----")
	for _, s := range results.Results {
            fmt.Fprintf(w, "\n %s\t%.2f\t", s.ClientName, s.BillableHours)
    }

}