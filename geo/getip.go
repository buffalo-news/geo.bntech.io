package geo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gapi.bntech.io/packages/lib"
	geoip2 "github.com/oschwald/geoip2-golang"
)

// GetIPData gets data for ip from the microservice
func GetIPData(hosturl string, ip string, attempts int) City {

	// Return empty json if all attemps are exhausted
	if attempts < 0 {
		return City{}
	}

	req, err := http.NewRequest("POST", hosturl+"/ip", bytes.NewBuffer([]byte("")))
	req.Header.Set("X-IP", ip)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		lib.LogNl(err.Error())
		attempts--
		return GetIPData(ip, attempts)
	}
	defer resp.Body.Close()

	// Read the body
	body, _ := ioutil.ReadAll(resp.Body)
	jString, _ := json.Marshal(lib.JsonFromBytes(body)["Body"])

	// Create the city to return
	var ipData geoip2.City
	err = json.Unmarshal([]byte(jString), &ipData)
	if err != nil {
		// @TODO: error handle, perhaps try the function again
	}

	if ipData.City.GeoNameID == 0 {
		return City{}
	}

	var cleanCity City
	cleanCity.City = ipData.City.Names["en"]
	cleanCity.Postal = ipData.Postal.Code
	cleanCity.Continent = ipData.Continent.Code
	cleanCity.Country = ipData.Country.IsoCode
	cleanCity.Location.AccuracyRadius = ipData.Location.AccuracyRadius
	cleanCity.Location.Latitude = ipData.Location.Latitude
	cleanCity.Location.Longitude = ipData.Location.Longitude
	cleanCity.Location.TimeZone = ipData.Location.TimeZone
	cleanCity.Traits.IsProxy = ipData.Traits.IsAnonymousProxy
	cleanCity.Traits.IsEU = ipData.Country.IsInEuropeanUnion
	cleanCity.Subdivision = ipData.Subdivisions[0].Names["en"]

	return cleanCity
}
