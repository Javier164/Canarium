package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

func ParseObservationData(config Config) (*Data, error) {
	url := fmt.Sprintf("https://api.weather.com/v1/location/%s:4:US/observations/current.json?language=en-US&units=e&apiKey=%s", config.ZIP, config.Key)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data Data
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func ParseCommonData(config Config) ([]Common, error) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/search?city=%s&state%s&postalcode=%s&format=json", config.City, config.State, config.ZIP)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var arr []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&arr)
	if err != nil {
		return nil, err
	}

	url = fmt.Sprintf("https://api.weather.com/v3/aggcommon/v3-wx-forecast-daily-5day?geocodes=%s,%s&language=en-US&units=e&format=json&apiKey=%s", arr[0]["lat"].(string)[0:6], arr[0]["lon"].(string)[0:6], config.Key)
	resp, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var common []Common
	err = json.NewDecoder(resp.Body).Decode(&common)
	if err != nil {
		return nil, err
	}

	return common, nil
}

func ParseRSSFeed(config Config) (*RSS, error) {
	resp, err := http.Get(config.Feed)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rss RSS
	err = xml.NewDecoder(resp.Body).Decode(&rss)
	if err != nil {
		return nil, err
	}

	return &rss, nil
}
