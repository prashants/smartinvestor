package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	csvfile, err := os.Open("nifty.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvfile.Close()

	csvreader := csv.NewReader(csvfile)
	csvreader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := csvreader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	html := `
<!DOCTYPE HTML>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<title>Highstock Example</title>

		<script type="text/javascript" src="http://ajax.googleapis.com/ajax/libs/jquery/1.8.2/jquery.min.js"></script>
		<style type="text/css">
		</style>
		<script type="text/javascript">
			$(function () {
				// Create the chart
				$('#closing-container').highcharts('StockChart', {
					xAxis: {
						type: 'datetime',
						dateTimeLabelFormats: {
							month: '%e. %b',
							year: '%b'
						},
					},
					yAxis: {
						title: {
							text: 'Closing'
						},
					},
					series : [
						{
                					name : 'closing',
					                data : $closing,
						},
					]
				});
				// Create the chart
				$('#volume-container').highcharts('StockChart', {
					xAxis: {
						type: 'datetime',
						dateTimeLabelFormats: {
							month: '%e. %b',
							year: '%b'
						},
					},
					yAxis: {
						title: {
							text: 'Volume'
						},
					},
					series : [
						{
                					name : 'volume',
					                data : $volume,
						},
					]
				});
				// Create the chart
				$('#turnover-container').highcharts('StockChart', {
					xAxis: {
						type: 'datetime',
						dateTimeLabelFormats: {
							month: '%e. %b',
							year: '%b'
						},
					},
					yAxis: {
						title: {
							text: 'Turnover'
						},
					},
					series : [
						{
                					name : 'turnover',
					                data : $turnover,
						},
					]
				});
				// Create the chart
				$('#pe-container').highcharts('StockChart', {
					xAxis: {
						type: 'datetime',
						dateTimeLabelFormats: {
							month: '%e. %b',
							year: '%b'
						},
					},
					yAxis: {
						title: {
							text: 'P/E'
						},
					},
					series : [
						{
                					name : 'pe',
					                data : $pe,
						},
					]
				});
				// Create the chart
				$('#pb-container').highcharts('StockChart', {
					xAxis: {
						type: 'datetime',
						dateTimeLabelFormats: {
							month: '%e. %b',
							year: '%b'
						},
					},
					yAxis: {
						title: {
							text: 'P/B'
						},
					},
					series : [
						{
                					name : 'pb',
					                data : $pb,
						},
					]
				});
				// Create the chart
				$('#yield-container').highcharts('StockChart', {
					xAxis: {
						type: 'datetime',
						dateTimeLabelFormats: {
							month: '%e. %b',
							year: '%b'
						},
					},
					yAxis: {
						title: {
							text: 'Yield'
						},
					},
					series : [
						{
                					name : 'yield',
					                data : $yield,
						},
					]
				});
				// Create the chart
				$('#returns-container').highcharts('StockChart', {
					xAxis: {
						type: 'datetime',
						dateTimeLabelFormats: {
							month: '%e. %b',
							year: '%b'
						},
					},
					yAxis: {
						title: {
							text: 'Returns'
						},
					},
					series : [
						{
                					name : 'returns',
					                data : $returns,
						},
					]
				});
			});
		</script>
	</head>
	<body>
		<script src="./highstock/js/highstock.js"></script>
		<script src="./highstock/js/modules/exporting.js"></script>
		<div id="closing-container" style="height: 800px; min-width: 310px"></div>
		<div id="volume-container" style="height: 800px; min-width: 310px"></div>
		<div id="turnover-container" style="height: 800px; min-width: 310px"></div>
		<div id="pe-container" style="height: 800px; min-width: 310px"></div>
		<div id="pb-container" style="height: 800px; min-width: 310px"></div>
		<div id="yield-container" style="height: 800px; min-width: 310px"></div>
		<div id="returns-container" style="height: 800px; min-width: 310px"></div>
	</body>
</html>
`

	// Date
	var curts time.Time
	var ts []int64
	ts = make([]int64, len(rawCSVdata))
	for key, row := range rawCSVdata {
		if key == 0 {
			continue
		}
		curts, err = time.Parse("02-Jan-2006", row[0])
		if err != nil {
			log.Fatal(err)
		}
		ts[key] = curts.Unix() * 1000
	}

	// Closing
	var bufClosing string
	bufClosing = "["
	for key, row := range rawCSVdata {
		if key == 0 {
			continue
		}
		bufClosing += "[" + strconv.FormatInt(ts[key], 10) + "," + row[4] + "],"
	}
	bufClosing += "]"

	html = strings.Replace(html, "$closing", bufClosing, -1)

	// Volume
	var bufVolume string
	bufVolume = "["
	for key, row := range rawCSVdata {
		if key == 0 {
			continue
		}
		bufVolume += "[" + strconv.FormatInt(ts[key], 10) + "," + row[5] + "],"
	}
	bufVolume += "]"

	html = strings.Replace(html, "$volume", bufVolume, -1)

	// Turnover
	var bufTurnover string
	bufTurnover = "["
	for key, row := range rawCSVdata {
		if key == 0 {
			continue
		}
		bufTurnover += "[" + strconv.FormatInt(ts[key], 10) + "," + row[6] + "],"
	}
	bufTurnover += "]"

	html = strings.Replace(html, "$turnover", bufTurnover, -1)

	// P/E
	var bufPE string
	bufPE = "["
	for key, row := range rawCSVdata {
		if key == 0 {
			continue
		}
		bufPE += "[" + strconv.FormatInt(ts[key], 10) + "," + row[7] + "],"
	}
	bufPE += "]"

	html = strings.Replace(html, "$pe", bufPE, -1)

	// P/B
	var bufPB string
	bufPB = "["
	for key, row := range rawCSVdata {
		if key == 0 {
			continue
		}
		bufPB += "[" + strconv.FormatInt(ts[key], 10) + "," + row[8] + "],"
	}
	bufPB += "]"

	html = strings.Replace(html, "$pb", bufPB, -1)

	// Yield
	var bufYield string
	bufYield = "["
	for key, row := range rawCSVdata {
		if key == 0 {
			continue
		}
		bufYield += "[" + strconv.FormatInt(ts[key], 10) + "," + row[9] + "],"
	}
	bufYield += "]"

	html = strings.Replace(html, "$yield", bufYield, -1)

	// Returns
	var bufReturns string
	bufReturns = "["
	for key, row := range rawCSVdata {
		if key == 0 {
			continue
		}
		bufReturns += "[" + strconv.FormatInt(ts[key], 10) + "," + row[10] + "],"
	}
	bufReturns += "]"

	html = strings.Replace(html, "$returns", bufReturns, -1)

	// Open output file
	fo, err := os.Create("charts.html")
	if err != nil {
		panic(err)
	}
	// Close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	// Make a write buffer
	w := bufio.NewWriter(fo)
	// Write to output file
	w.Write([]byte(html))
	// Flush data
	if err = w.Flush(); err != nil {
		panic(err)
	}
}
