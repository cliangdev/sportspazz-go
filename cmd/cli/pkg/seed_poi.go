package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

var apiKey string
var sportspazzHost string
var sportspazzToken string
var city string
var sport string
var pages int

var preSeedPoiCmd = &cobra.Command{
	Use:   "seed-poi --api-key <api_key> --city <city>  --sport <sport>  --pages <number_pages>  --sportspazz-host <sportspazz_host> --sportspazz-token <request_token>",
	Short: "Pre-seed point of interest data to sportspazz.com",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if pois, err := searchAllPlaces(city, sport, pages); err != nil {
			for _, poi := range pois {
				fmt.Printf("Creating POI: %s", poi.Name)

				if err := createPOI(poi); err != nil {
					fmt.Printf("%v", err)
					break
				}
			}
		}
	},
}

func init() {
	preSeedPoiCmd.Flags().StringVar(&apiKey, "api-key", "", "Google Places API Key")
	preSeedPoiCmd.MarkFlagRequired("api-key")

	preSeedPoiCmd.Flags().StringVar(&sportspazzHost, "sportspazz-host", "localhost:4001", "Sportspazz Host")
	preSeedPoiCmd.MarkFlagRequired("sportspazz-host")

	preSeedPoiCmd.Flags().StringVar(&sportspazzToken, "sportspazz-token", "", "Request Token")
	preSeedPoiCmd.MarkFlagRequired("sportspazz-token")

	preSeedPoiCmd.Flags().StringVar(&city, "city", "", "City")
	preSeedPoiCmd.MarkFlagRequired("city")

	preSeedPoiCmd.Flags().StringVar(&sport, "sport", "", "Sport")
	preSeedPoiCmd.MarkFlagRequired("sport")

	preSeedPoiCmd.Flags().IntVar(&pages, "pages", 10, "Page size")
}

func searchPlaces(city, sport, nextPageToken string) (*PlacesResponse, error) {
	fmt.Printf("Pre-seeding %s in %s data to %s\n", sport, city, sportspazzHost)

	var googlePlaceUrl string
	if nextPageToken != "" {
		googlePlaceUrl = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/textsearch/json?pagetoken=%s&key=%s", url.QueryEscape(nextPageToken), apiKey)
	} else {
		query := fmt.Sprintf("%s in %s", sport, city)
		googlePlaceUrl = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/textsearch/json?query=%s&key=%s", url.QueryEscape(query), apiKey)
	}

	resp, err := http.Get(googlePlaceUrl)
	if err != nil {
		fmt.Println("Error making the API request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return nil, err
	}

	var placesResponse PlacesResponse
	if err := json.Unmarshal(body, &placesResponse); err != nil {
		fmt.Println("Error unmarshaling the response:", err)
		return nil, err
	}

	if len(placesResponse.Results) == 0 {
		fmt.Println("No results found.")
		return nil, err
	}

	for _, place := range placesResponse.Results {
		printPlaceDetail(place)
	}

	return &placesResponse, nil
}

func printPlaceDetail(place Place) {
	photoURL := ""
	if len(place.Photos) > 0 {
		photo := place.Photos[0]
		photoURL = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&photoreference=%s&key=%s", photo.PhotoReference, apiKey)
	}
	fmt.Printf("Name: %s\nAddress: %s\nPhoto: %s\n", place.Name, place.Address, photoURL)
	placeDetails, err := getPlaceDetails(apiKey, place.PlaceID)
	if err == nil {
		fmt.Printf("Rating: %.1f\n", placeDetails.Rating)
		fmt.Printf("Website: %s\n", placeDetails.Website)
		fmt.Printf("Phone Number: %s\n", placeDetails.PhoneNumber)
		fmt.Printf("Business Status: %s\n", placeDetails.BusinessStatus)
		fmt.Printf("Place Id: %s\n\n", placeDetails.PlaceID)

		if len(placeDetails.Reviews) > 0 {
			review := placeDetails.Reviews[0]
			fmt.Printf("Review by %s: %s (Rating: %d)\n", review.AuthorName, review.Text, review.Rating)
		}
		fmt.Println("\n")
	}
}

func searchAllPlaces(city, sport string, numPages int) ([]POI, error) {
	var allPlaces []Place
	nextPageToken := ""

	for i := 0; i < numPages; i++ {
		placesResponse, err := searchPlaces(city, sport, nextPageToken)
		if err != nil {
			break
		}
		if placesResponse == nil || placesResponse.Results == nil {
			fmt.Println("no response")
			break
		}
		allPlaces = append(allPlaces, placesResponse.Results...)
		nextPageToken = placesResponse.NextPageToken

		if placesResponse.NextPageToken == "" {
			fmt.Println("no next page")
			break
		}
	}

	cityPlaceId, _ := getCityPlaceID(apiKey, city)

	var pois []POI
	for _, place := range allPlaces {
		thumbnailURL := ""
		if len(place.Photos) > 0 {
			photo := place.Photos[0]
			thumbnailURL = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&photoreference=%s&key=%s", photo.PhotoReference, apiKey)
		}
		description := ""
		placeDetails, err := getPlaceDetails(apiKey, place.PlaceID)
		if err != nil {
			if place.Website != "" {
				description += fmt.Sprintf("Website: %s", placeDetails.Website)
			}
			if place.PhoneNumber != "" {
				description += fmt.Sprintf("Phone Number: %s", placeDetails.PhoneNumber)
			}
			description += fmt.Sprintf("Rating: %.1f\n", placeDetails.Rating)
		}

		if err := createPOI(POI{
			Name:         place.Name,
			Address:      place.Address,
			CityID:       cityPlaceId,
			SportType:    sport,
			ThumbnailURL: thumbnailURL,
			Description:  description,
		}); err != nil {
			fmt.Println(err)
			break
		}
	}

	return pois, nil
}

func getCityPlaceID(apiKey, cityName string) (string, error) {
	query := url.QueryEscape(cityName)
	endpoint := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/textsearch/json?query=%s&key=%s", query, apiKey)
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var result struct {
		Results []Place `json:"results"`
		Status  string  `json:"status"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if result.Status != "OK" {
		return "", fmt.Errorf("API request failed with status: %s", result.Status)
	}

	if len(result.Results) == 0 {
		return "", fmt.Errorf("no results found for city: %s", cityName)
	}

	return result.Results[0].PlaceID, nil
}

func getPlaceDetails(apiKey, placeID string) (*PlaceDetails, error) {
	endpoint := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/details/json?place_id=%s&key=%s", url.QueryEscape(placeID), apiKey)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result struct {
		Result PlaceDetails `json:"result"`
		Status string       `json:"status"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("API request failed with status: %s", result.Status)
	}

	return &result.Result, nil
}

func createPOI(poi POI) error {
	sportspazzUrl := fmt.Sprintf("%s/api/v1/pois", sportspazzHost)

	payload, err := json.Marshal(poi)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", sportspazzUrl, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sportspazzToken)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %v", err)
		}
		return fmt.Errorf("failed to create POI: %v", string(body))
	}

	return nil
}
