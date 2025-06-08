package models

type WeatherResponse struct {
	Country string `json:"country"`
	Date    string `json:"date"`
	Text    string `json:"text"`
}

type WeatherFindResult struct {
	DataResponse *WeatherResponse
	Err          error
}
