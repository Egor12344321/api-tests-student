package models

type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}

type Advertisement struct {
	Title string `json:"title"`
}

type AdvertisementCreateResponse struct {
	ID string `json:"id"`
}