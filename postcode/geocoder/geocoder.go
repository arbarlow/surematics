package geocoder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

type Postcode struct {
	Postcode  string  `json:"postcode"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type postcodeRequest struct {
	Postcodes []string `json:"postcodes"`
}

type postcodeResult struct {
	Status int `json:"status"`
	Result []struct {
		Query    string   `json:"query"`
		Postcode Postcode `json:"result"`
	} `json:"result"`
}

// DistanceBetweenCodes is a function that ships out a request to postcodes.io
// postcodes.io is free to use and open source. In a prod environment we'd probably
// dockerize postcodes.io and run it as an internal service to save latency,
// increase security etc.
//
// Notice the function could be variadiac and could accept up to 100 postcodes, Go will
// stream the encoding and decoding of large reponses, though response decoding
// is somewhat negated by dumping it into an in memory struct (still fun)
func DistanceBetweenCodes(p1, p2 string) (float64, error) {
	req := postcodeRequest{Postcodes: []string{p1, p2}}

	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(req)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(
		"https://api.postcodes.io/postcodes",
		"application/json", &b)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP status not 200, was %v", resp.Status)
	}
	defer resp.Body.Close()

	var res postcodeResult
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return 0, err
	}

	if len(res.Result) != 2 {
		return 0, fmt.Errorf("expected 2 results, got %v", len(res.Result))
	}

	return distance(res.Result[0].Postcode, res.Result[1].Postcode), nil
}

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance formula using haversine shamelesly stolen from the internet and
// made a bit more "structy"
func distance(p1, p2 Postcode) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = p1.Latitude * math.Pi / 180
	lo1 = p1.Longitude * math.Pi / 180
	la2 = p2.Latitude * math.Pi / 180
	lo2 = p2.Longitude * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}
