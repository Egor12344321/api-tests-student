package models

import advModels "api-tests-template/internal/managers/advertisements/models"

type FavoriteListResponse struct {
	Items []FavoriteItem `json:"items"`
}

type FavoriteItem struct {
	Advertisement advModels.Advertisement `json:"advertisement"`
}
