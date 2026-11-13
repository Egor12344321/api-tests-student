package models

type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}

type Advertisement struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type AdvertisementCreateResponse struct {
	ID     string          `json:"id"`
	Photos []PhotoResponse `json:"photos"`
}

type AdvertisementDetails struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type PhotoResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
