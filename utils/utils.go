package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func CalculateCarCarbonFootprint(carBrand, carModel string, distanceKm float64) (float64, error) {
	// call the gcp cloud function to calculate the carbon impact for the car
	return 0, nil
}

// GeocodingResponse represents the response structure from Google Geocoding API
type GeocodingResponse struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

// Co2Response represents the response structure from Impact CO₂ API
type Co2Response struct {
	Data []struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Value float64 `json:"value"`
	} `json:"data"`
}

func GetCarbonImpactByMode(modeID int, distanceKm float64) (float64, error) {
	baseURL := "https://impactco2.fr/api/v1/transport"
	params := url.Values{}
	params.Add("km", strconv.FormatFloat(distanceKm, 'f', 2, 64))
	params.Add("displayAll", "0")
	params.Add("transports", strconv.Itoa(modeID))
	params.Add("ignoreRadiativeForcing", "0")
	params.Add("occupencyRate", "1")
	params.Add("includeConstruction", "0")
	params.Add("language", "fr")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return 0, fmt.Errorf("failed to call Impact CO₂ API: %w", err)
	}
	defer resp.Body.Close()

	var result Co2Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to parse CO₂ response: %w", err)
	}

	if len(result.Data) == 0 {
		return 0, fmt.Errorf("no CO₂ data returned for transport ID: %d", modeID)
	}

	return result.Data[0].Value, nil
}

func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const EarthRadius = 6371 // Earth radius in kilometers
	lat1Rad, lon1Rad := lat1*math.Pi/180, lon1*math.Pi/180
	lat2Rad, lon2Rad := lat2*math.Pi/180, lon2*math.Pi/180

	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadius * c
}

func GetCoordinates(address, apiKey string) (float64, float64, error) {
	baseURL := "https://maps.googleapis.com/maps/api/geocode/json"
	params := url.Values{}
	params.Add("address", address)
	params.Add("key", apiKey)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return 0, 0, fmt.Errorf("failed to call geocoding API: %w", err)
	}
	defer resp.Body.Close()

	var result GeocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, fmt.Errorf("failed to parse geocoding response: %w", err)
	}

	if len(result.Results) == 0 {
		return 0, 0, fmt.Errorf("no results found for address: %s", address)
	}

	location := result.Results[0].Geometry.Location
	return location.Lat, location.Lng, nil
}

func CalculateDistance(startAddress, endAddress string) (float64, error) {
	key := os.Getenv("GOOGLE_MAPS_API_KEY")
	startLat, startLng, err := GetCoordinates(startAddress, key)
	if err != nil {
		return 0, fmt.Errorf("failed to get start coordinates: %w", err)
	}

	endLat, endLng, err := GetCoordinates(endAddress, key)
	if err != nil {
		return 0, fmt.Errorf("failed to get end coordinates: %w", err)
	}

	return HaversineDistance(startLat, startLng, endLat, endLng), nil
}
