// The MIT License (MIT)
//
// SmartInvestor - Smart Investor Tools
//
// Copyright (c) 2015 Prashant Shah <pshah.mumbai@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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

	var closingValue float64

	welcomeMsg()

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

	// Read dates
	var dates []string
	dateFile, err := os.Open("dates.txt")
	if err != nil {
        log.Fatal(err)
	}
	defer dateFile.Close()
	dateScanner := bufio.NewScanner(dateFile)
	// Setup slice containing all the dates from the file
	for dateScanner.Scan() {
		if len(dateScanner.Text()) < 1 {
			continue
		}
		dates = append(dates, dateScanner.Text())
	}
	if err := dateScanner.Err(); err != nil {
		log.Fatal(err)
	}
	if (len(dates) < 1) {
		log.Fatal("Please enter atleast one date in the format YYYY-MM-DD in the dates.txt file")
	}

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
	_, err = fmt.Fprintf(writer, "%s", "SCRIPT,")
	if err != nil {
		log.Fatal(err)
	}
	for _, eachdate := range dates {
		_, err = fmt.Fprintf(writer, "%s,", eachdate)
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = fmt.Fprintf(writer, "\n")
	if err != nil {
		log.Fatal(err)
	}
	// Flush data
	if err = writer.Flush(); err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		fmt.Printf("Processing : %s\n", scanner.Text())

		// Write script name
		_, err = fmt.Fprintf(writer, "%s,", scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		
		// Calculate and write closing value for each date to output file
		for _, eachdate := range dates {
			// Read the closing value
			closingValue = getQuotes(scanner.Text(), eachdate, eachdate)
			// Write closing value
			_, err = fmt.Fprintf(writer, "%f,", closingValue)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Write new line
		_, err = fmt.Fprintf(writer, "\n")
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

func welcomeMsg() {
	fmt.Println("")
	fmt.Println("*****************************************************************************")
	fmt.Println("")
	fmt.Println("                        WELCOME TO SMART INVESTOR !                        ")
	fmt.Println("")
	fmt.Println("- Create a symbols.txt file containing all the script symbols one per each line from Yahoo finance that you want to track.")
	fmt.Println("- Create a dates.txt file containing all the dates for which you want the closing price. Enter one date per line in the format YYYY-MM-DD.")
	fmt.Println("- Press Ctrl+C if you which to halt the program.")
	fmt.Println("*****************************************************************************")
	fmt.Println("")
	fmt.Println("")
}
