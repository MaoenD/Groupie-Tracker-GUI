package search

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
)

type GeoLocation struct {
    Lat string `json:"lat"`
    Lon string `json:"lon"`
}

func ConvertAddressToGeoCoord(address string) (GeoLocation, error) {
    var geoLocation GeoLocation
    baseURL := "https://nominatim.openstreetmap.org/search"
    query := fmt.Sprintf("?q=%s&format=json&limit=1", url.QueryEscape(address))

    // Set custom User-Agent
    client := &http.Client{}
    req, err := http.NewRequest("GET", baseURL+query, nil)
    if err != nil {
        return geoLocation, err
    }
    req.Header.Set("User-Agent", "GroupieTrackerApp")

    resp, err := client.Do(req)
    if err != nil {
        return geoLocation, err
    }
    defer resp.Body.Close()

    var result []GeoLocation
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        return geoLocation, err
    }

    if len(result) > 0 {
        geoLocation = result[0]
    } else {
        return geoLocation, fmt.Errorf("no results found for address: %s", address)
    }

    return geoLocation, nil
}
