package goplaces

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// AppID represents header X-Algolia-Application-Id
var AppID string

// AppKey represents header X-Algolia-API-Key
var AppKey string

// Parameters for query
type Parameters struct {
	Query     string `json:"query"`
	Countries string `json:"countries"`
	Type      string `json:"type"`
	AppID     string
	AppKey    string
}

// QueryResponse represents response from /places/query
type QueryResponse struct {
	Hits             []Hit  `json:"hits"`
	NbHits           int    `json:"nbHits"`
	ProcessingTimeMS int    `json:"processingTimeMS"`
	Query            string `json:"query"`
	Params           string `json:"params"`
}

// Hit is a part of QueryResponse, represents hits from response
type Hit struct {
	// Flags
	IsCountry bool `json:"is_country"`
	IsHighway bool `json:"is_highway"`
	IsCity    bool `json:"is_city"`
	IsSuburb  bool `json:"is_suburb"`
	IsPopular bool `json:"is_popular"`
	// Values
	Country        map[string]string   `json:"country"`
	City           map[string][]string `json:"city"`
	County         map[string]string   `json:"county"`
	LocaleNames    map[string][]string `json:"locale_names"`
	Administrative []string            `json:"administrative"`
	Suburb         []string            `json:"suburb"`
	Postcode       []string            `json:"postcode"`
	Tags           []string            `json:"tags"`
	Population     int                 `json:"population"`
	CountryCode    string              `json:"country_code"`
	Importance     int                 `json:"importance"`
	AdminLevel     int                 `json:"admin_level"`
	Geolocation    map[string]float64  `json:"_geoloc"`
}

// Address for decoding hits
type Address struct {
	Country  string
	Postcode string
	State    string
	City     string
	Street   string
}

// Query is a wrapper around /places/query request
func Query(p Parameters) (QueryResponse, error) {
	// Encode parameters
	pbytes, _ := json.Marshal(p)
	// Prepare request
	req, _ := http.NewRequest("POST", "https://places-dsn.algolia.net/1/places/query", bytes.NewBuffer(pbytes))
	if AppID != "" {
		req.Header.Add("X-Algolia-Application-Id", AppID)
	}
	if AppKey != "" {
		req.Header.Add("X-Algolia-API-Key", AppKey)
	}
	if p.AppID != "" {
		req.Header.Add("X-Algolia-Application-Id", p.AppID)
	}
	if p.AppKey != "" {
		req.Header.Add("X-Algolia-API-Key", p.AppKey)
	}
	// Get response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return QueryResponse{}, err
	}
	// Extract response object
	body, _ := ioutil.ReadAll(resp.Body)
	var qresp QueryResponse
	json.Unmarshal(body, &qresp)
	// Return
	return qresp, nil
}

// ExtractAddress is a function for extracting address from Hit record,
// hit records are not so obvious and hard to use
func ExtractAddress(hit Hit) Address {
	address := Address{}
	// Postcode
	if len(hit.Postcode) == 1 {
		address.Postcode = hit.Postcode[0]
	}
	// State
	if len(hit.Administrative) > 0 {
		address.State = hit.Administrative[0]
	}
	// City
	if hit.IsCity {
		address.City = hit.LocaleNames["default"][0]
	} else if len(hit.Suburb) > 0 {
		address.City = hit.Suburb[0]
	} else if len(hit.City) > 0 {
		address.City = hit.City["default"][0]
	}
	// Street
	if !hit.IsCity && !hit.IsSuburb && !hit.IsCountry && len(hit.LocaleNames) != 0 && len(hit.LocaleNames["default"]) > 0 {
		address.Street = hit.LocaleNames["default"][0]
	}
	return address
}

// ExtractAddresses same as ExtractAddress, but for slice
func ExtractAddresses(hits []Hit) []Address {
	addresses := []Address{}
	for _, hit := range hits {
		addresses = append(addresses, ExtractAddress(hit))
	}
	return addresses
}

// NewLabelFromAddress builds single string from address object
func NewLabelFromAddress(address Address) string {
	builder := []string{}
	if address.Street != "" {
		builder = append(builder, address.Street)
	}
	if address.City != "" {
		builder = append(builder, address.City)
	}
	if address.State != "" {
		builder = append(builder, address.State)
	}
	if address.Postcode != "" {
		builder = append(builder, address.Postcode)
	}
	if address.Country != "" {
		builder = append(builder, address.Country)
	}
	return strings.Join(builder, ", ")
}
