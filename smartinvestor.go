package main

import (
    "fmt"
	"strconv"
	"log"
	"bufio"
	"net/http"
	"os"
	"io/ioutil"
	"github.com/jeffail/gabs"
	"strings"
)

// Globals
var baserurl string
var param map[string]string

func main() {

	var closingValue1 float64
	var closingValue2 float64
	var closingValue3 float64
	var closingValue4 float64
	var closingValue5 float64
	var closingValue6 float64
	var closingValue7 float64
	var closingValue8 float64
	var closingValue9 float64
	var closingValue10 float64
	var closingValue11 float64
	var closingValue12 float64
	var closingValue13 float64
	var closingValue14 float64

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
	_, err = fmt.Fprintf(writer, "%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n", "SCRIPT", "2000-08-25", "2002-09-27", "2003-09-29", "2007-04-27", "2008-05-30", "2009-03-06", "2010-04-23", "2011-05-06", "2011-07-22", "2011-08-19", "2012-02-03", "2014-09-18", "2014-10-15", "2014-10-31")
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
		closingValue1 = getQuotes(scanner.Text(), "2000-08-25", "2000-08-25")
		closingValue2 = getQuotes(scanner.Text(), "2002-09-27", "2002-09-27")
		closingValue3 = getQuotes(scanner.Text(), "2003-09-29", "2003-09-29")
		closingValue4 = getQuotes(scanner.Text(), "2007-04-27", "2007-04-27")
		closingValue5 = getQuotes(scanner.Text(), "2008-05-30", "2008-05-30")
		closingValue6 = getQuotes(scanner.Text(), "2009-03-06", "2009-03-06")
		closingValue7 = getQuotes(scanner.Text(), "2010-04-23", "2010-04-23")
		closingValue8 = getQuotes(scanner.Text(), "2011-05-06", "2011-05-06")
		closingValue9 = getQuotes(scanner.Text(), "2011-07-22", "2011-07-22")
		closingValue10 = getQuotes(scanner.Text(), "2011-08-19", "2011-08-19")
		closingValue11 = getQuotes(scanner.Text(), "2012-02-03", "2012-02-03")
		closingValue12 = getQuotes(scanner.Text(), "2014-09-18", "2014-09-18")
		closingValue13 = getQuotes(scanner.Text(), "2014-10-15", "2014-10-15")
		closingValue14 = getQuotes(scanner.Text(), "2014-10-31", "2014-10-31")

		// Write the data to output file
		_, err = fmt.Fprintf(writer, "%s,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f\n", scanner.Text(), closingValue1,closingValue2,closingValue3,closingValue4,closingValue5,
			closingValue6,closingValue7,closingValue8,closingValue9,closingValue10,closingValue11,closingValue12,closingValue13,closingValue14)
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
