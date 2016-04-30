# csv2i - CSV To Invoice

**CSV To Invoice** is a command line utility that imports CSV files to FreshBooks Invoices using [FreshBooks Invoice API](https://www.freshbooks.com/developers/docs/invoices)

## CSV Format

There are two types of files `csv2i` expects:

1. Invoices
2. Lines

### Invoices

Fields:

	DATE,INVOICE_NUM,AMOUNT

Example (saved in `LLS_invoices.csv`):

	2016-01-04,LLS-62,281.25
	2016-02-03,LLS-63,93.75
	2016-03-01,LLS-64,500.00
	2016-04-01,LLS-65,1406.25

### Lines

For each `INVOICE_NUM` in `LLS_invoices.csv` a **Lines** file should be created and named `INVOICE_NUM.csv`. With the following fields:

	"Description Can Have Spaces",UNIT_COST,QUANTITY

Example: **Lines** file for Invoice `LLS-65` (saved in `LLS-65.csv`):

	"[03/01/2016] 5057: COM error: password does not meet the password policy",225.00,1.00
	"[03/02/2016] 5060: Gmail Password Reset Web Service Call SSL Cert",225.00,1.00
	"[03/04/2016] 5060: ORA-29024: Certificate validation failure",225.00,1.00
	"[03/18/2016] 5065: COM GA pass sync: researched strange blank errors",225.00,.50
	"[03/27/2016] 5067: ORA-01033: ORACLE - idm.com:IDM - researched issue",225.00,2.00
	"[04/01/2016] 5069: Mar 2016 Monitoring Service (alert monitoring)",225.00,.75

## Command Line Options

	csv2i -h
	Usage of csv2i:
	  -account string
	    	FreshBooks Account Name
	  -csvDir string
	    	CSV File Directory (default "./testdata")
	  -csvFile string
	    	Invoice CSV File
	  -dryRun
	    	Don't Create FreshBooks Invoices just print XML Calls
	  -fbClientID string
	    	FreshBooks Client ID
	  -fbPONum string
	    	FreshBooks PO Number
	  -fbTask string
	    	FreshBooks Task Name
	  -fbToken string
	    	FreshBooks API Token
	  -trace
	    	Trace true|false

**CSV To Invoice** will prefix your FreshBooks `account` to FreshBooks API URL:

	https://%account%.freshbooks.com/api/2.1/xml-in

And make a basic authentication request using `fbToken`.

It'll then take `csvFile` (from `csvDir`) and use it as a driver parsing one invoice at a time and getting it's Lines file (INVOICE_NUM.csv) from `csvDir`.

Once the invoice and it's lines are assembled, this block of data is printed on the screen.

If `-trace=true` then we print an assembled FreshBooks XML call.

We then parse additional parameters as follows:

1. `fbClientID` as `invoice.client_id`
2. `fbPONum` as `invoice.po_number`
3. `fbTask` as `line.name`

If `-dryRun=true` then we move to the next invoice (no call to FreshBooks is made).

If `-dryRun=false` (default) then we call FreshBooks API `method=invoice.create` and parse it's response which could be either a newly created InvoiceID or an Error message.

## Examples

Does a dry-run (only prints XML no FreshBooks calls are made):

	csv2i -account=mycorp -csvFile=LLS_invoices.csv \
	   -fbClientID=499999 -fbPONum=PO-7654321 \
	   -fbTask=Hourly-225 \
	   -fbToken=********** \
	   -trace=true -dryRun=true

Creates invoices listed in `LLS_invoices.csv`

	csv2i -account=mycorp -csvFile=LLS_invoices.csv \
	   -fbClientID=499999 -fbPONum=PO-7654321 \
	   -fbTask=Hourly-225 \
	   -fbToken=**********

## License
MIT
