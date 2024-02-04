package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/app/http/controllers"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/app/http/resources"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/app/models"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/config"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/route"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/test/factory"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/test/unit"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestGetTagListJSONResult struct {
	Code    int                    `json:"code" `
	Message string                 `json:"message"`
	Data    []models.Advertistment `json:"data"`
}

func TestGetAdvertistmentList(t *testing.T) {
	cleanup := unit.SetupIntegrationDB(t)
	router := route.Setup()

	total := 5
	_ = factory.FactoryAdvertistments(total)

	w := unit.PerformRequest(router, "GET", "/api/v1/ad", nil)
	assert.Equal(t, http.StatusOK, w.Code, "wrong status code")

	var response TestGetTagListJSONResult
	err := json.Unmarshal([]byte(w.Body.Bytes()), &response)
	assert.NoError(t, err)

	assert.Equal(t, total, len(response.Data), "they should be equal")

	t.Cleanup(cleanup)
}

func TestCreateAdvertistment(t *testing.T) {
	cleanup := unit.SetupIntegrationDB(t)
	router := route.Setup()

	startAt, _ := time.Parse(time.RFC3339, "2023-12-10T03:00:00.000Z")
	endAt, _ := time.Parse(time.RFC3339, "2023-12-31T16:00:00.000Z")
	ageStart := 20
	ageEnd := 30
	gender := []string{"M", "F"}
	country := []string{"TW", "JP"}
	platform := []string{"ANDROID", "IOS"}

	formAddAdvertistmentCondition := controllers.FormAddAdvertistmentCondition{
		AgeStart: &ageStart,
		AgeEnd:   &ageEnd,
		Gender:   &gender,
		Country:  &country,
		Platform: &platform,
	}

	formAddAdvertistment := &controllers.FormAddAdvertistment{
		Title:                         "連線送大vava",
		StartAt:                       startAt,
		EndAt:                         endAt,
		FormAddAdvertistmentCondition: formAddAdvertistmentCondition,
	}

	serialized, _ := json.Marshal(formAddAdvertistment)

	w := unit.PerformRequest(router, "POST", "/api/v1/ad", bytes.NewReader(serialized))
	assert.Equal(t, http.StatusOK, w.Code, "wrong status code")

	var response resources.JSONResult
	err := json.Unmarshal([]byte(w.Body.Bytes()), &response)
	assert.NoError(t, err)

	res, err := models.GetAdvertistmentByTitle("連線送大vava")
	assert.NoError(t, err)
	assert.NotNil(t, res)

	t.Cleanup(cleanup)
}

func TestMaxAdEveryday(t *testing.T) {
	cleanup := unit.SetupIntegrationDB(t)
	router := route.Setup()

	for i := 0; i < config.AppConfig.Detail.MaxAdEveryday; i++ {
		title := fmt.Sprintf("Test AD %d", i)
		startAt, _ := time.Parse(time.RFC3339, "2023-12-10T03:00:00.000Z")
		endAt, _ := time.Parse(time.RFC3339, "2023-12-31T16:00:00.000Z")
		ageStart := 20
		ageEnd := 30
		gender := models.GenderModelArray{models.MALE, models.FEMALE}
		country := models.CountryModelArray{models.TW, models.JP}
		platform := models.PlatformModelArray{models.ANDROID, models.IOS}

		testData := models.Advertistment{
			Title:    &title,
			StartAt:  &startAt,
			EndAt:    &endAt,
			AgeStart: &ageStart,
			AgeEnd:   &ageEnd,
			Gender:   &gender,
			Country:  &country,
			Platform: &platform,
		}

		testData.Create()
	}

	startAt2, _ := time.Parse(time.RFC3339, "2023-12-10T03:00:00.000Z")
	endAt2, _ := time.Parse(time.RFC3339, "2023-12-31T16:00:00.000Z")
	ageStart2 := 20
	ageEnd2 := 30
	gender2 := []string{"M", "F"}
	country2 := []string{"TW", "JP"}
	platform2 := []string{"ANDROID", "IOS"}

	formAddAdvertistmentCondition := controllers.FormAddAdvertistmentCondition{
		AgeStart: &ageStart2,
		AgeEnd:   &ageEnd2,
		Gender:   &gender2,
		Country:  &country2,
		Platform: &platform2,
	}

	formAddAdvertistment := &controllers.FormAddAdvertistment{
		Title:                         "連線送大vava",
		StartAt:                       startAt2,
		EndAt:                         endAt2,
		FormAddAdvertistmentCondition: formAddAdvertistmentCondition,
	}

	serialized, _ := json.Marshal(formAddAdvertistment)

	w := unit.PerformRequest(router, "POST", "/api/v1/ad", bytes.NewReader(serialized))
	assert.Equal(t, http.StatusBadRequest, w.Code, "wrong status code")

	var response resources.JSONResult
	err := json.Unmarshal([]byte(w.Body.Bytes()), &response)
	assert.NoError(t, err)
	assert.Equal(t, "No more ads can be added today :(", response.Message, "wrong message")

	t.Cleanup(cleanup)
}

func TestMaxActiveAd(t *testing.T) {
	cleanup := unit.SetupIntegrationDB(t)
	router := route.Setup()

	for i := 0; i < config.AppConfig.Detail.MaxActiveAd-1; i++ {
		title := fmt.Sprintf("Test AD %d", i)
		startAt := time.Now()
		endAt := startAt.Add(time.Hour * 100)
		ageStart := 20
		ageEnd := 30
		gender := models.GenderModelArray{models.MALE, models.FEMALE}
		country := models.CountryModelArray{models.TW, models.JP}
		platform := models.PlatformModelArray{models.ANDROID, models.IOS}

		testData := models.Advertistment{
			Title:    &title,
			StartAt:  &startAt,
			EndAt:    &endAt,
			AgeStart: &ageStart,
			AgeEnd:   &ageEnd,
			Gender:   &gender,
			Country:  &country,
			Platform: &platform,
		}

		testData.Create()
	}

	startAt2 := time.Now()
	endAt2 := startAt2.Add(time.Hour * 24)
	ageStart2 := 20
	ageEnd2 := 30
	gender2 := []string{"M", "F"}
	country2 := []string{"TW", "JP"}
	platform2 := []string{"ANDROID", "IOS"}

	formAddAdvertistmentCondition := controllers.FormAddAdvertistmentCondition{
		AgeStart: &ageStart2,
		AgeEnd:   &ageEnd2,
		Gender:   &gender2,
		Country:  &country2,
		Platform: &platform2,
	}

	formAddAdvertistment := &controllers.FormAddAdvertistment{
		Title:                         "連線送大vava",
		StartAt:                       startAt2,
		EndAt:                         endAt2,
		FormAddAdvertistmentCondition: formAddAdvertistmentCondition,
	}

	serialized, _ := json.Marshal(formAddAdvertistment)

	w := unit.PerformRequest(router, "POST", "/api/v1/ad", bytes.NewReader(serialized))
	assert.Equal(t, http.StatusBadRequest, w.Code, "wrong status code")

	var response resources.JSONResult
	err := json.Unmarshal([]byte(w.Body.Bytes()), &response)
	assert.NoError(t, err)
	assert.Equal(t, "No more ads can be added :(, choose other start_at and end_at", response.Message, "wrong message")

	t.Cleanup(cleanup)
}
