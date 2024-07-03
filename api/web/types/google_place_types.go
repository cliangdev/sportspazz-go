package types

type GooglePlaceResponse struct {
	Result Result `json:"result"`
	Status string `json:"status"`
}

type Result struct {
	FormattedAddress     string       `json:"formatted_address"`
	FormattedPhoneNumber string       `json:"formatted_phone_number"`
	Name                 string       `json:"name"`
	OpeningHours         OpeningHours `json:"opening_hours"`
	Photos               []Photo      `json:"photos"`
	Rating               int          `json:"rating"`
	Website              string       `json:"website"`
}

type OpeningHours struct {
	WeekdayText []string `json:"weekday_text"`
}

type Photo struct {
	Height           int      `json:"height"`
	HtmlAttributions []string `json:"html_attributions"`
	PhotoReference   string   `json:"photo_reference"`
	Width            int      `json:"width"`
}
