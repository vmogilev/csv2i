package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	fbAccount  = flag.String("account", "", "FreshBooks Account Name")
	fbToken    = flag.String("fbToken", "", "FreshBooks API Token")
	fbClientID = flag.String("fbClientID", "", "FreshBooks Client ID")
	fbPONum    = flag.String("fbPONum", "", "FreshBooks PO Number")
	fbTask     = flag.String("fbTask", "", "FreshBooks Task Name")
	csvDir     = flag.String("csvDir", "./testdata", "CSV File Directory")
	csvFile    = flag.String("csvFile", "", "Invoice CSV File")
	trace      = flag.Bool("trace", false, "Trace true|false")
	dryRun     = flag.Bool("dryRun", false, "Don't Create FreshBooks Invoices just print XML Calls")
)

func loadDetails(src string) ([][]string, error) {
	var s *os.File
	var err error

	if s, err = os.Open(src); err != nil {
		return nil, err
	}
	defer s.Close()

	r := csv.NewReader(s)

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func main() {
	var s *os.File
	var err error

	flag.Parse()

	if s, err = os.Open(filepath.Join(*csvDir, *csvFile)); err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	r := csv.NewReader(s)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print(err.Error())
		}

		fmt.Println(record)

		details, err := loadDetails(filepath.Join(*csvDir, record[1]+".csv"))
		if err != nil {
			log.Print(err)
		}
		for _, x := range details {
			fmt.Println(x)
		}

		inv, err := newInvoice(record, details)
		if err != nil {
			log.Print(err)
		}
		if inv > 0 {
			fmt.Printf("Created FreshBooks InvoiceID: %d\n", inv)
		}
	}
}
