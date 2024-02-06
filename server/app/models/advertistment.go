package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/config"
	"strings"
	"time"

	"gorm.io/gorm"
)

type GenderModelArray []GenderModel

type GenderModel string

const (
	MALE   GenderModel = "M"
	FEMALE GenderModel = "F"
)

var validGenders = map[string]GenderModel{
	"M": MALE,
	"F": FEMALE,
}

type CountryModelArray []CountryModel

type CountryModel string

const (
	TW CountryModel = "TW"
	JP CountryModel = "JP"
)

var validCountries = map[string]CountryModel{
	"TW": TW,
	"JP": JP,
}

type PlatformModelArray []PlatformModel

type PlatformModel string

const (
	ANDROID PlatformModel = "ANDROID"
	IOS     PlatformModel = "IOS"
	WEB     PlatformModel = "WEB"
)

var validPlatforms = map[string]PlatformModel{
	"ANDROID": ANDROID,
	"IOS":     IOS,
	"WEB":     WEB,
}

type Advertistment struct {
	ID        int `gorm:"primary_key"`
	Title     *string
	StartAt   *time.Time
	EndAt     *time.Time
	AgeStart  *int
	AgeEnd    *int
	Gender    *GenderModelArray
	Country   *CountryModelArray
	Platform  *PlatformModelArray
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func StringArrayToGenderModelArray(genders []string) (GenderModelArray, error) {
	c := make([]GenderModel, len(genders))
	for i, v := range genders {
		res, err := StringToGenderModel(v)
		if err != nil {
			return nil, err
		}
		c[i] = res
	}
	return c, nil
}

func StringArrayToCountryModelArray(countries []string) (CountryModelArray, error) {
	c := make([]CountryModel, len(countries))
	for i, v := range countries {
		res, err := StringToCountry(v)
		if err != nil {
			return nil, err
		}
		c[i] = res
	}
	return c, nil
}

func StringArrayToPlatformModelArray(platforms []string) (PlatformModelArray, error) {
	c := make([]PlatformModel, len(platforms))
	for i, v := range platforms {
		res, err := StringToPlatformModel(v)
		if err != nil {
			return nil, err
		}
		c[i] = res
	}
	return c, nil
}

func StringToGenderModel(gender string) (GenderModel, error) {
	if genderModel, ok := validGenders[gender]; ok {
		return genderModel, nil
	}
	return "", errors.New("invalid gender")
}

func StringToCountry(country string) (CountryModel, error) {
	if countryModel, ok := validCountries[country]; ok {
		return countryModel, nil
	}
	return "", errors.New("invalid country")
}

func StringToPlatformModel(platform string) (PlatformModel, error) {
	if platformModel, ok := validPlatforms[platform]; ok {
		return platformModel, nil
	}
	return "", errors.New("invalid platform")
}

func (g *GenderModelArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to change value:", value))
	}
	text := string(bytes)
	string_slice := strings.Split(text, ",")
	res, err := StringArrayToGenderModelArray(string_slice)
	if err != nil {
		return err
	}
	*g = res
	return nil
}

func (e GenderModelArray) Value() (driver.Value, error) {
	c := make([]string, 0)
	for _, v := range e {
		text := string(v)
		c = append(c, text)
	}
	text := strings.Join(c[:], ",")
	return text, nil
}

func (g *CountryModelArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to change value:", value))
	}
	text := string(bytes)
	string_slice := strings.Split(text, ",")
	res, err := StringArrayToCountryModelArray(string_slice)
	if err != nil {
		return err
	}
	*g = res
	return nil
}

func (e CountryModelArray) Value() (driver.Value, error) {
	c := make([]string, 0)
	for _, v := range e {
		text := string(v)
		c = append(c, text)
	}
	text := strings.Join(c[:], ",")
	return text, nil
}

func (g *PlatformModelArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to change value:", value))
	}
	text := string(bytes)
	string_slice := strings.Split(text, ",")
	res, err := StringArrayToPlatformModelArray(string_slice)
	if err != nil {
		return err
	}
	*g = res
	return nil
}

func (e PlatformModelArray) Value() (driver.Value, error) {
	c := make([]string, 0)
	for _, v := range e {
		text := string(v)
		c = append(c, text)
	}
	text := strings.Join(c[:], ",")
	return text, nil
}

func (ad Advertistment) Create() error {
	err := db.Create(&ad).Error
	if err != nil {
		return err
	}
	return nil
}

func WhereAdIsActived(startAt time.Time, endAt time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("start_at < ? AND end_at > ?", startAt, endAt)
	}
}

func GetAdvertistmentByTitle(title string) (*Advertistment, error) {
	var advertistment Advertistment
	err := db.Model(&Advertistment{}).Where("title = ?", title).First(&advertistment).Error
	if err != nil {
		return nil, err
	}
	return &advertistment, nil
}

func GetAdvertistmentList(offset int, limit int, age *int, gender *GenderModel, country *CountryModel, platform *PlatformModel) ([]Advertistment, int64, error) {
	var advertistments []Advertistment
	var count int64
	tx := db.Model(&Advertistment{})
	if age != nil {
		tx = tx.Where("age_start < ? and age_end > ?", age, age)
	}
	if gender != nil {
		tx = tx.Where("FIND_IN_SET(?,gender)>0", gender)
	}
	if country != nil {
		tx = tx.Where("FIND_IN_SET(?,country)>0", country)
	}
	if platform != nil {
		tx = tx.Where("FIND_IN_SET(?,platform)>0", platform)
	}
	now := time.Now()

	err := tx.Scopes(Paginate(offset, limit), WhereAdIsActived(now, now)).Order("end_at asc").Find(&advertistments).Error
	if err != nil {
		return advertistments, 0, err
	}

	err = tx.Count(&count).Error
	if err != nil {
		return advertistments, 0, err
	}

	return advertistments, count, nil
}

func GetTodayCreateAdCount() (int64, error) {
	var count int64
	err := db.Model(&Advertistment{}).Where("date(DATE_ADD(created_at, INTERVAL ? HOUR)) = date(DATE_ADD(NOW(), INTERVAL ? HOUR))", config.AppConfig.TimeZone, config.AppConfig.TimeZone).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func CheckStartAtAndEndAt(startAt time.Time, endAt time.Time) (int64, error) {
	// 計算在 startAt 和 endAt 這段期間內活耀的廣告有多少
	var count int64
	err := db.Model(&Advertistment{}).Where("(? between start_at and end_at) and (? between start_at and end_at)", startAt, endAt).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
