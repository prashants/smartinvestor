package main

import (
    "fmt"
	"strconv"
	"log"
	"bufio"
    //"net/url"
	"net/http"
	"os"
	"io/ioutil"
	//"encoding/json"
	"github.com/jeffail/gabs"
	//"github.com/wenxiang/go-nestedjson"
	"strings"
)

func main() {

	var closingValue float64

	// Read stock symbols from a file
	symbolFile, err := os.Open("symbols.txt")
	if err != nil {
        log.Fatal(err)
	}
	defer symbolFile.Close()

	scanner := bufio.NewScanner(symbolFile)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Yahoo historic data YQL formats
	// https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.historicaldata%20where%20symbol%20%3D%20%22aapl%22%20and%20startDate%20%3D%20%222014-10-10%22%20and%20endDate%20%3D%20%222014-10-10%22&format=json&diagnostics=true&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback=
	// https://query.yahooapis.com/v1/public/yql?q=select * from yahoo.finance.historicaldata where symbol = "aapl" and startDate = "2014-10-10" and endDate = "2014-10-10"&format=json&diagnostics=true&env=store://datatables.org/alltableswithkeys&callback=
	// https://developer.yahoo.com/yql/console/?q=show%20tables&env=store://datatables.org/alltableswithkeys#h=select+*+from+yahoo.finance.historicaldata+where+symbol+%3D+%22aapl%22+and+startDate+%3D+%222014-10-10%22+and+endDate+%3D+%222014-10-10%22
	// select * from yahoo.finance.historicaldata where symbol = "aapl" and startDate = "2014-10-10" and endDate = "2014-10-10"
	
	// Initialize the URL
	baserurl := "https://query.yahooapis.com/v1/public/yql?"
	var param map[string]string
	param = make(map[string]string)
	param["q"] = "select * from yahoo.finance.historicaldata where symbol = \"aapl\" and startDate = \"2014-10-10\" and endDate = \"2014-10-10\""
	param["format"] = "json"
	param["diagnostics"] = "true"
	param["env"] = "store://datatables.org/alltableswithkeys"
	param["callback"] = ""

	finalURL := baserurl + "q=" + yahooURL(param["q"]) + "&format=" + param["format"] + "&diagnostics=" + param["diagnostics"] + "&env=" + yahooURL(param["env"]) + "&callback=" + param["callback"]
	
	// Make a GET request to Yahoo get the results
    response, err := http.Get(finalURL)
    if err != nil {
        log.Fatal(err)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            log.Fatal(err)
        }

		// Parse the JSON data received from Yahoo
		jsonParsed, err := gabs.ParseJSON([]byte(contents))

		//closingValue = jsonParsed.Search("query", "results", "quote", "Close").Data().(float64)
		//fmt.Printf("%s\n", jsonParsed.Search("query", "results", "quote", "Close").Data().(float64))
		closingValue, err = strconv.ParseFloat(jsonParsed.Search("query", "results", "quote", "Close").Data().(string), 64)
        if err != nil {
            log.Fatal(err)
        }
		fmt.Printf("%s\n", closingValue)
    }
}

// Convert a given url string to a format accepted by Yahoo API
func yahooURL(str string) string {
	// Replacement rules :
	// space => %20
	// = => %3D
	// " => %22
	// : => %3A
	// / => %2F

	str = strings.Replace(str, " ", "%20", -1)
	str = strings.Replace(str, "=", "%3D", -1)
	str = strings.Replace(str, "\"", "%22", -1)
	str = strings.Replace(str, ":", "%3A", -1)
	str = strings.Replace(str, "/", "%2F", -1)

	return str
}