package factory

import (
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/app/models"
	"time"

	"github.com/jaswdr/faker/v2"
)

func FactoryAdvertistments(amount int, startAt *time.Time, endAt *time.Time) []models.Advertistment {
	tags := make([]models.Advertistment, amount)
	for i := 0; i < amount; i++ {
		tags = append(tags, generateAdvertistment(startAt, endAt))
	}
	return tags
}

func generateAdvertistment(startAt *time.Time, endAt *time.Time) models.Advertistment {
	var ad models.Advertistment
	faker := faker.New()

	title := faker.Beer().Name()
	ad.Title = &title

	if startAt == nil || endAt == nil {
		now := time.Now()
		nowadd := time.Now().Add(1 * time.Hour)
		ad.StartAt = &now
		ad.EndAt = &nowadd
	} else {
		ad.StartAt = startAt
		ad.EndAt = endAt
	}

	ageEnd := faker.IntBetween(1, 30)
	ad.AgeEnd = &ageEnd

	ageStart := faker.IntBetween(1, 30)
	ad.AgeStart = &ageStart

	gender := make(models.GenderModelArray, 0)
	gender = append(gender, models.MALE)
	ad.Gender = &gender

	country := make(models.CountryModelArray, 0)
	country = append(country, models.TW)
	ad.Country = &country

	platform := make(models.PlatformModelArray, 0)
	platform = append(platform, models.ANDROID)
	ad.Platform = &platform

	err := ad.Create()
	if err != nil {
		panic(err)
	}
	return ad
}
