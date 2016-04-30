package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type fbAPI struct {
	apiURL   string
	apiToken string
}

func newAPI(account string, token string) *fbAPI {
	url := fmt.Sprintf("https://%s.freshbooks.com/api/2.1/xml-in", account)
	fb := fbAPI{apiURL: url, apiToken: token}
	return &fb
}

func (a *fbAPI) makeRequest(request interface{}) (*[]byte, error) {
	xmlRequest, err := xml.MarshalIndent(request, "", "  ")
	if err != nil {
		return nil, err
	}

	if *trace {
		fmt.Printf("makeRequest: %v\n", string(xmlRequest))
	}

	if *dryRun {
		return nil, err
	}

	req, err := http.NewRequest("POST", a.apiURL, bytes.NewBuffer(xmlRequest))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(a.apiToken, "X")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type invoice struct {
	ClientID string    `xml:"client_id"`
	InvNum   string    `xml:"number"`
	Date     string    `xml:"date"`
	PoNum    string    `xml:"po_number,omitempty"`
	InvLines *invLines `xml:"lines"`
}

// this sub type is needed so we get the
// nested XML tags:
// <lines>
//    <line></line>
// </lines>
// if carried to invoice level it doesn't work
type invLines struct {
	Lines *[]invLine `xml:"line"`
}
type invLine struct {
	//XMLName xml.Name `xml:"line"` - moved to invLines
	Name string `xml:"name"`
	Desc string `xml:"description"`
	Cost string `xml:"unit_cost"`
	Qty  string `xml:"quantity"`
	Type string `xml:"type"`
}

func loadLines(details [][]string) *[]invLine {
	var lines []invLine
	for _, v := range details {
		l := invLine{
			Name: *fbTask,
			Desc: v[0],
			Cost: v[1],
			Qty:  v[2],
			Type: "Time",
		}
		lines = append(lines, l)
	}
	return &lines
}

func newInvoice(record []string, details [][]string) (int, error) {
	var err error

	lines := loadLines(details)
	i := &invoice{
		ClientID: *fbClientID,
		InvNum:   record[1],
		Date:     record[0],
		PoNum:    *fbPONum,
		InvLines: &invLines{lines},
	}

	req := &struct {
		XMLName xml.Name `xml:"request"`
		Method  string   `xml:"method,attr"`
		Invoice *invoice `xml:"invoice"`
	}{
		Method:  "invoice.create",
		Invoice: i,
	}

	fb := newAPI(*fbAccount, *fbToken)
	result, err := fb.makeRequest(req)
	if err != nil {
		return 0, err
	}

	if *dryRun {
		return 0, nil
	}

	parsedInto := struct {
		Status    string `xml:"status,attr"`
		Error     string `xml:"error"`
		InvoiceID int    `xml:"invoice_id"`
	}{}

	if err := xml.Unmarshal(*result, &parsedInto); err != nil {
		return 0, err
	}
	if parsedInto.Status == "ok" {
		return parsedInto.InvoiceID, nil
	}
	return 0, errors.New(parsedInto.Error)
}
