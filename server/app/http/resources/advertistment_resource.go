package resources

import (
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/app/models"
	"time"
)

type StructAdvertistmentCollection struct {
	Title *string `json:"title"`
	EndAt *string `json:"endAt"`
}

func AdvertistmentCollection(items []models.Advertistment) []StructAdvertistmentCollection {
	var data []StructAdvertistmentCollection
	for _, v := range items {
		var endAt *string = nil
		if v.EndAt != nil {
			res := v.EndAt.Format(time.RFC3339)
			endAt = &res
		}
		d := StructAdvertistmentCollection{
			Title: v.Title,
			EndAt: endAt,
		}
		data = append(data, d)
	}
	if len(data) == 0 {
		return make([]StructAdvertistmentCollection, 0)
	}
	return data
}
