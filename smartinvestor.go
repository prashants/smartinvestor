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

// Globals
var baserurl string
var param map[string]string

func main() {

	var closingValue float64

	// Initialize Global varibles
	// Yahoo historic data YQL formats
	// https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.historicaldata%20where%20symbol%20%3D%20%22aapl%22%20and%20startDate%20%3D%20%222014-10-10%22%20and%20endDate%20%3D%20%222014-10-10%22&format=json&diagnostics=true&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback=
	// https://query.yahooapis.com/v1/public/yql?q=select * from yahoo.finance.historicaldata where symbol = "aapl" and startDate = "2014-10-10" and endDate = "2014-10-10"&format=json&diagnostics=true&env=store://datatables.org/alltableswithkeys&callback=
	// https://developer.yahoo.com/yql/console/?q=show%20tables&env=store://datatables.org/alltableswithkeys#h=select+*+from+yahoo.finance.historicaldata+where+symbol+%3D+%22aapl%22+and+startDate+%3D+%222014-10-10%22+and+endDate+%3D+%222014-10-10%22
	// select * from yahoo.finance.historicaldata where symbol = "aapl" and startDate = "2014-10-10" and endDate = "2014-10-10"

	baserurl = "https://query.yahooapis.com/v1/public/yql?"
	param = make(map[string]string)
	param["format"] = "json"
	param["diagnostics"] = "true"
	param["env"] = "store://datatables.org/alltableswithkeys"
	param["callback"] = ""

	// Read stock symbols from a file
	symbolFile, err := os.Open("symbols.txt")
	if err != nil {
        log.Fatal(err)
	}
	defer symbolFile.Close()

	// Output file to save results
	opFile, err := os.OpenFile("output.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
        log.Fatal(err)
	}
	defer opFile.Close()

	// Setup reader and writer buffers
	scanner := bufio.NewScanner(symbolFile)
	writer := bufio.NewWriter(opFile)

	// Flush buffers
    defer func() {
        if err = writer.Flush(); err != nil {
            log.Fatal(err)
        }
    }()

	// Write output csv headers
	_, err = fmt.Fprintf(writer, "%s,%s\n", "SCRIPT", "2014-10-10")
	if err != nil {
		log.Fatal(err)
	}
	// Flush data
	if err = writer.Flush(); err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		fmt.Printf("Processing : %s\n", scanner.Text())

		// Read the closing value
		closingValue = getQuotes(scanner.Text(), "2014-10-10", "2014-10-10")

		// Write the data to output file
		_, err = fmt.Fprintf(writer, "%s,%f\n", scanner.Text(), closingValue)
		if err != nil {
			log.Fatal(err)
		}

		// Flush data
        if err = writer.Flush(); err != nil {
            log.Fatal(err)
        }
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
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

// Fetch quotes from Yahoo
func getQuotes(symbol string, fromdate string, todate string) float64 {

	// Yahoo historic data YQL formats
	// https://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.historicaldata%20where%20symbol%20%3D%20%22aapl%22%20and%20startDate%20%3D%20%222014-10-10%22%20and%20endDate%20%3D%20%222014-10-10%22&format=json&diagnostics=true&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback=
	// https://query.yahooapis.com/v1/public/yql?q=select * from yahoo.finance.historicaldata where symbol = "aapl" and startDate = "2014-10-10" and endDate = "2014-10-10"&format=json&diagnostics=true&env=store://datatables.org/alltableswithkeys&callback=
	// https://developer.yahoo.com/yql/console/?q=show%20tables&env=store://datatables.org/alltableswithkeys#h=select+*+from+yahoo.finance.historicaldata+where+symbol+%3D+%22aapl%22+and+startDate+%3D+%222014-10-10%22+and+endDate+%3D+%222014-10-10%22
	// select * from yahoo.finance.historicaldata where symbol = "aapl" and startDate = "2014-10-10" and endDate = "2014-10-10"

	var closingValue float64
	var closingValueStr string
	var closingValueInterface interface{} 

	// Setup query parameters
	param["q"] = "select * from yahoo.finance.historicaldata where symbol = \"" + symbol + "\" and startDate = \"" + fromdate + "\" and endDate = \"" + todate + "\""

	finalURL := baserurl + "q=" + yahooURL(param["q"]) + "&format=" + param["format"] + "&diagnostics=" + param["diagnostics"] + "&env=" + yahooURL(param["env"]) + "&callback=" + param["callback"]
	
	// Make a GET request to Yahoo get the results
    response, err := http.Get(finalURL)
    if err != nil {
        log.Fatal(err)
    }
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the JSON data received from Yahoo
	jsonParsed, err := gabs.ParseJSON([]byte(contents))
	if err != nil {
		log.Fatal(err)
	}

	closingValueInterface = jsonParsed.Search("query", "results", "quote", "Close").Data()
	switch v := closingValueInterface.(type) {
		case string:
			closingValueStr = v
		default:
			closingValueStr = "0.0"
	}

	// Convert closing value to float64
	closingValue, err = strconv.ParseFloat(closingValueStr, 64)
	if err != nil {
		log.Fatal(err)
	}

	return closingValue
}
