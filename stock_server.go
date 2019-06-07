// Package main - stock_server implements stock price lookup on most exchanges.
//
// Server is on localhost:3000
// Responds to endpoint "/stock/"
//
// Specify the stock:  symbol=stock, stock, ...
// Exchanges are optional: stock_exchange=exc, exc, ..
// If an exchange is not specified the default is AMEX.
//
// Example:   http://localhost:3000/stock/symbol=MSFT,AAPL,FAX&stock_exchange=NASDAQ,AMEX
//
// Result of example:
//	stock: AAPL	price  153.30	exchange  NASDAQ
//	stock: FAX	price  3.99	    exchange  AMEX
//	stock: MSFT	price  105.68	exchange  NASDAQ
//
// Error Messages:  If a viable stock name is not found:
//     "Error! The requested stock(s) could not be found."
//
// Unrecoverable Errors are handled by the "log" package.
//
// This program uses the minimum standard libraries to accomplish its goals!
// "fmt"  - For displaying text.
// "io/ioutil" - For reading the response body.
// "net/http" - For REST Operations.
// "strings" - For manipulating strings.
// "unicode" - For determining text classes (ie string, number..).
// "log" - For reporting Fatal Errors!
//
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"unicode"
)

// parseData - Cracks returned data into an array of strings.
// Will not break on "_" or ".".
//
// Returns []string (ie Slice of Strings).
func parseData(data string) []string {
	var enp rune = 0x005F
	var prd rune = 0x002E

	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != enp && c != prd
	}
	fld := strings.FieldsFunc(data, f)
	return fld
}

// parseList - Reads input slice of string and returns a list of stocks
// and a list of exchanges.
//
// Returns int, []string, []string
func parseList(fields []string) (numb int, slist, elist string) {
	stockList := make([]string, 1)
	exchangeList := make([]string, 1)
	s, ss, e := 0, 0, 0
	l := len(fields)
	numb = 0
	for i := 0; i < l; i++ {
		if fields[i] == "stock" {
			if fields[i+1] == "exchange" {
				e = i + 2
				ss = i
			}
			if fields[i+1] == "symbol" {
				s = i + 2
			}
		}
	}
	if e == 0 {
		exchangeList = append(exchangeList, "AMEX")
		ss = l
	}
	if s < ss {
		for i := s; i < ss; i++ {
			stockList = append(stockList, fields[i])
			numb++
		}
	}
	if e > 0 && e < l {
		for i := e; i < l; i++ {
			exchangeList = append(exchangeList, fields[i])
		}
	}

	slist = strings.Join(stockList, ",")
	elist = strings.Join(exchangeList, " ")
	return numb, slist, elist
}

// getStock - Returns data on the a particular stock or set stocks.
//
// Results: A []byte buffer containing the response body.
func getStock(stock string) []byte {
	res := make([]byte, 1)
	response, err := http.Get("https://www.worldtradingdata.com/api/v1/stock?symbol=" + stock + "&api_token=i4KpwCAiPeKfvg1Pqt0TAmZzfurU0OfsIgh6nxqDRnIObQH812BuDnaflzC6")
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		res, err = ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		return res
	}
	return res
}

// stock_handler - Handles the response to "/stock/".
//          - extracts the stock symbols and the exchanges using local functions.
//          - Formats Stock Symbols, Prices, and exchanges return to the user.
//
func stock_handler(w http.ResponseWriter, r *http.Request) {
	total := r.URL.Path
	var enp int32 = 0x002E
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != enp
	}
	fields := strings.FieldsFunc(total, f)

	lnum, slist, elist := parseList(fields)
	body := getStock(slist)
	if strings.Contains(string(body), "Error") {
		fmt.Fprintf(w, "Error! The requested stock(s) could not be found.")
	}

	fld := parseData(string(body))
	el := strings.SplitN(elist[1:], " ", lnum)
	stockEntry := make([]string, 1)
	priceEntry := make([]string, 1)
	exchange := make([]string, 1)
	valid := make([]bool, 1)

	for j := 0; j < len(fld); j++ {
		if fld[j] == "symbol" {
			stockEntry = append(stockEntry, fld[j+1])
			valid = append(valid, false)
		}
		if fld[j] == "price" {
			priceEntry = append(priceEntry, fld[j+1])
		}
		if fld[j] == "stock_exchange_short" {
			exchange = append(exchange, fld[j+1])
		}
	}

	for i := 0; i < len(exchange); i++ {
		for j := 0; j < len(el); j++ {
			if exchange[i] == string(el[j]) {
				valid[i] = true
			}
		}
	}

	for i := 0; i < len(stockEntry); i++ {
		if valid[i] {
			fmt.Fprintf(w, "stock: %s	price  %s	exchange  %s\n", stockEntry[i], priceEntry[i], exchange[i])
		}
	}
}

// main - Sets up handler on ":3000/stock..." and waits on http traffic.
//
func main() {
	http.HandleFunc("/stock/", stock_handler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
