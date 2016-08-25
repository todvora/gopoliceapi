package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"time"
)

const APP_BASE_URL = "http://aplikace.policie.cz/patrani-vozidla/"
const REQUEST_TIMEOUT = 1 * time.Second

type Client struct {
	timeout   time.Duration
	httplient *http.Client
}

func constructSearchUrl(vin string, regNo string) string {
	return APP_BASE_URL + "default.aspx?__EVENTTARGET=&__EVENTARGUMENT=&__VIEWSTATE=%2FwEPDwULLTEzNzIzMjY0MDMPZBYCZg9kFgICBw9kFgICAQ9kFgoCBA8PZBYCHgpvbmtleXByZXNzBTZyZXR1cm4gb25LZXlwcmVzcyhldmVudCwnY3RsMDBfQXBwbGljYXRpb25fY21kSGxlZGVqJylkAgYPD2QWAh8ABTZyZXR1cm4gb25LZXlwcmVzcyhldmVudCwnY3RsMDBfQXBwbGljYXRpb25fY21kSGxlZGVqJylkAhQPDxYEHgRUZXh0BRxDZWxrb3bDvSBwb8SNZXQgesOhem5hbcWvOiAxHgdWaXNpYmxlZ2RkAhUPDxYCHwJoZGQCGQ8PFgIfAQU%2FRGF0YWLDoXplIGJ5bGEgbmFwb3NsZWR5IGFrdHVhbGl6b3bDoW5hIDxiPjYuIHByb3NpbmNlIDIwMTI8L2I%2BZGRkE2qlXWNJcxoc8%2FLZOQEi5oKrGzs%3D&__EVENTVALIDATION=%2FwEWBQL80qOCBwLQsb3%2BBgK9peeFDwLv%2BPyjBAL4oIjjDVvi8FJppOBh8gjuF1u%2Ft7viEDtA&ctl00%24Application%24txtSPZ=" + regNo + "&ctl00%24Application%24txtVIN=" + vin + "&ctl00%24Application%24cmdHledej=Vyhledat&ctl00%24Application%24CurrentPage=1"
}

func getDocument(client *http.Client, url string) (*goquery.Document, error) {

	resp, clientErr := client.Get(url) // body closed by goquery

	if clientErr != nil {
		//
		//if e,ok := clientErr.(net.Error); ok && e.Timeout() {
		//	fmt.Println("Failed to download page. Timeout occured!")
		//
		//} else if clientErr != nil {
		//	fmt.Println("Failed to download page. Net error occured!")
		//}
		return nil, clientErr
	}

	doc, goqueryErr := goquery.NewDocumentFromResponse(resp)

	if goqueryErr != nil {
		return nil, goqueryErr
	}

	return doc, nil
}

func getDetailUrls(c *http.Client, vin string, regNo string) ([]string, error) {

	if len(vin) == 0 && len(regNo) == 0 {
		return nil, errors.New("No search query provided!")
	}

	url := constructSearchUrl(vin, regNo)
	doc, err := getDocument(c, url)
	if err != nil {
		return nil, err
	}
	return doc.Find("table#celacr tr a").Map(func(i int, s *goquery.Selection) string {
		str, _ := s.Attr("href") // TODO: a without href?
		return APP_BASE_URL + str
	}), nil
}

func parseDetails(client *http.Client, url string, ret chan VehicleDetails) {
	if doc, err := getDocument(client, url); err != nil {
		ret <- VehicleDetails{error: err}
		return

	} else {
		attributes := make(map[string]string)

		doc.Find("table#searchTableResults tr span").Each(func(_ int, s *goquery.Selection) {
			id, _ := s.Attr("id")
			fieldName := strings.Replace(id, "ctl00_Application_lbl", "", 1)
			fieldName = strings.ToLower(fieldName)
			value := s.Text()
			key := TranslateKey(fieldName)
			attributes[key] = value
		})
		attributes["url"] = url
		attributes["stolendate"] = StandardizeDate(attributes["stolendate"])
		ret <- VehicleDetails{details: attributes}
		return
	}
}

type VehicleDetails struct {
	details map[string]string
	error   error
}

type Results struct {
	Results          []map[string]string `json:"results"`
	Count            int                 `json:"count"`
	Time             time.Time           `json:"time"`
	RequestsDuration string              `json:"requestsDuration,omitempty"`
	Error            string              `json:"error,omitempty"`
}

func toSuccessResults(details []map[string]string, duration time.Duration) Results {
	return Results{Count: len(details), Results: details, Time: time.Now(), RequestsDuration: duration.String()}
}

func toErrorResults(err error) Results {
	return Results{Error: err.Error(), Time: time.Now(), Results: make([]map[string]string, 0)}
}

func (c *Client) Search(vin string, regNo string) Results {
	start := time.Now()
	urls, err := getDetailUrls(c.httplient, vin, regNo)

	if err != nil {
		return toErrorResults(err)
	}

	resultsChannel := make(chan VehicleDetails, len(urls))

	for _, url := range urls {
		go parseDetails(c.httplient, url, resultsChannel)
	}

	details := make([]map[string]string, len(urls))
	for i := 0; i < len(urls); i++ {
		res := <-resultsChannel
		if res.error != nil {
			return toErrorResults(err)
		}
		details[i] = res.details
	}
	return toSuccessResults(details, time.Since(start))
}

func defaultClient() *Client {
	client := &http.Client{Timeout: REQUEST_TIMEOUT}
	return &Client{httplient: client}
}

func newClient(client *http.Client) *Client {
	return &Client{httplient: client}
}
