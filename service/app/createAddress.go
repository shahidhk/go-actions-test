package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/shahidhk/gql"
)

const gmapsKey = ""
var gmapsURL = "https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=" + gmapsKey

type gmapsAPIResponse struct {
	Results []results
	Status  string
}

type results struct {
	Geometry geometry
}
type geometry struct {
	Location location
}
type location struct {
	Lat float64
	Lng float64
}

var createAddressHandler = func(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST requests only", http.StatusMethodNotAllowed)
		return
	}
	var req requestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("bad request")
		return
	}

	fmt.Println(req)

	// contact gmaps with text
	resp, err := http.Get(fmt.Sprintf(gmapsURL, url.QueryEscape(req.Input["text"].(string))))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gmapsResponseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var apiResp gmapsAPIResponse
	err = json.Unmarshal(gmapsResponseBytes, &apiResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(apiResp)

	var res responseBody

	if len(apiResp.Results) == 0 {
		res.Errors = []string{"cannot find address"}
	} else {

		var hasuraResponse struct {
			Data struct {
				InsertAddress struct {
					Returning []struct {
						UserId string  `json:"user_id"`
						Text   string  `json:"text"`
						Long   float64 `json:"long"`
						Lat    float64 `json:"lat"`
						Id     string  `json:"id"`
					} `json:"returning"`
				} `json:"insert_addresses"`
			} `json:"data"`
		}
		client := gql.NewClient("http://graphql-engine:8080/v1/graphql", nil)

		err = client.Execute(gql.Request{Query: "mutation insertAddress($user_id: uuid!, $text: String!, $lat: numeric!, $long: numeric!) {insert_addresses(objects: {user_id: $user_id, text: $text, long: $long, lat:$lat}) {returning {id lat long text user_id}}}", Variables: map[string]interface{}{
			"user_id": "4e843d61-7df0-44a9-9119-7468d04ab771",
			"text":    req.Input["text"],
			"lat":     apiResp.Results[0].Geometry.Location.Lat,
			"long":    apiResp.Results[0].Geometry.Location.Lng,
		}}, &hasuraResponse)
		fmt.Println(hasuraResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Data = hasuraResponse.Data.InsertAddress.Returning[0]
	}

	fmt.Println(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
