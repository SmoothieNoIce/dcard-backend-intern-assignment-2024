package factory

import (
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/app/models"

	"github.com/jaswdr/faker/v2"
)

func FactoryAdvertistments(amount int) []models.Advertistment {
	tags := make([]models.Advertistment, amount)
	for i := 0; i < amount; i++ {
		tags = append(tags, generateAdvertistment())
	}
	return tags
}

func FactoryDiscordServertag() models.Advertistment {
	return generateAdvertistment()
}

func generateAdvertistment() models.Advertistment {
	var ad models.Advertistment
	faker := faker.New()

	title := faker.Beer().Name()
	ad.Title = &title

	ageStart := faker.IntBetween(1, 30)
	ad.AgeStart = &ageStart

	ageEnd := faker.IntBetween(1, 30)
	ad.AgeEnd = &ageEnd

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
