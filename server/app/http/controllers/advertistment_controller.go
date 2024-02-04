package controllers

import (
	"encoding/json"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/app/http/resources"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/app/models"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/config"
	"github/SmoothieNoIce/dcard-backend-intern-assignment-2024/database/cache"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/redis/go-redis/v9"
)

type FormAddAdvertistment struct {
	Title                         string                        `json:"title" binding:"required" example:"酷酷廣告1"`
	StartAt                       time.Time                     `json:"startAt" binding:"required,ltfield=EndAt"  example:"2023-12-10T03:00:00.000Z"`
	EndAt                         time.Time                     `json:"endAt" binding:"required,gtfield=StartAt"  example:"2023-12-31T16:00:00.000Z"`
	FormAddAdvertistmentCondition FormAddAdvertistmentCondition `json:"conditions" binding:"required"`
}

type FormAddAdvertistmentCondition struct {
	AgeStart *int      `json:"ageStart" example:"20" binding:"omitempty,omitnil,gte=1,lte=100,ltfield=AgeEnd"`
	AgeEnd   *int      `json:"ageEnd" example:"30" binding:"omitempty,omitnil,gte=1,lte=100,gtfield=AgeStart"`
	Gender   *[]string `json:"gender" example:"M,F" binding:"omitempty,omitnil,dive,oneof='M' 'F'"`
	Country  *[]string `json:"country"  example:"TW,JP" binding:"omitempty,omitnil,dive,oneof='TW' 'JP'"`
	Platform *[]string `json:"platform" example:"ANDROID,IOS" binding:"omitempty,omitnil,dive,oneof='ANDROID' 'IOS' 'WEB'"`
}

// @Summary 廣告列表
// @Produce  json
// @tags ad
// @Param offset query int false "Offset"
// @Param limit query int false  "Limit"
// @Param age query int false "Age"
// @Param gender query string false "Gender" Enums(M, F)
// @Param country query string false "Country" Enums(TW, JP)
// @Param platform query string false "Platform" Enums(ANDROID, IOS, WEB)
// @success 200 {object} resources.JSONResult{data=interface{}}
// @Failure 400 {object} resources.JSONResult{data=interface{}}
// @Router /ad [get]
func GetAdvertistmentList(c *gin.Context) {
	var offset int
	if c.Query("offset") == "" {
		offset = 1
	} else {
		offsetInt, err := strconv.Atoi(c.Query("offset"))
		if err != nil {
			c.JSON(http.StatusBadRequest, resources.JSONResult{
				Code:    http.StatusBadRequest,
				Message: "offset is not validated",
				Data:    gin.H{},
			})
			return
		}
		offset = offsetInt
	}

	var limit int
	if c.Query("limit") == "" {
		limit = 5
	} else {
		limitInt, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.JSON(http.StatusBadRequest, resources.JSONResult{
				Code:    http.StatusBadRequest,
				Message: "limit is not validated",
				Data:    gin.H{},
			})
			return
		}
		if limit <= 1 && limit >= 100 {
			c.JSON(http.StatusBadRequest, resources.JSONResult{
				Code:    http.StatusBadRequest,
				Message: "limit needs to be between 1 and 100",
				Data:    gin.H{},
			})
			return
		}
		limit = limitInt
	}

	var age *int = nil
	if c.Query("age") != "" {
		ageInt, err := strconv.Atoi(c.Query("age"))
		if err != nil {
			c.JSON(http.StatusBadRequest, resources.JSONResult{
				Code:    http.StatusBadRequest,
				Message: "age is not validated",
				Data:    gin.H{},
			})
			return
		}
		if ageInt <= 1 && ageInt >= 100 {
			c.JSON(http.StatusBadRequest, resources.JSONResult{
				Code:    http.StatusBadRequest,
				Message: "age should be between 1 and 100",
				Data:    gin.H{},
			})
			return
		}
		age = &ageInt
	}

	var genderModel *models.GenderModel = nil
	if c.Query("gender") != "" {
		gender, err := models.StringToGenderModel(c.Query("gender"))
		if err != nil {
			c.JSON(http.StatusBadRequest, resources.JSONResult{
				Code:    http.StatusBadRequest,
				Message: "gender is not validated",
				Data:    gin.H{},
			})
			return
		}
		genderModel = &gender
	}

	var countryModel *models.CountryModel = nil
	if c.Query("country") != "" {
		country, err := models.StringToCountry(c.Query("country"))
		if err != nil {
			c.JSON(http.StatusBadRequest, resources.JSONResult{
				Code:    http.StatusBadRequest,
				Message: "country is not validated",
				Data:    gin.H{},
			})
			return
		}
		countryModel = &country
	}

	var platformModel *models.PlatformModel = nil
	if c.Query("platform") != "" {
		platform, err := models.StringToPlatformModel(c.Query("platform"))
		if err != nil {
			c.JSON(http.StatusBadRequest, resources.JSONResult{
				Code:    http.StatusBadRequest,
				Message: "platform is not validated",
				Data:    gin.H{},
			})
			return
		}
		platformModel = &platform
	}

	keyArr := make([]string, 0)
	keyArr = append(keyArr, c.Query("offset"), c.Query("limit"), c.Query("age"), c.Query("gender"), c.Query("country"), c.Query("platform"))
	key := strings.Join(keyArr[:], ",")

	val, err := cache.Rdb0.Get(cache.Ctx, key).Result()
	if err == redis.Nil {
		result, _, err := models.GetAdvertistmentList(offset, limit, age, genderModel, countryModel, platformModel)
		if err != nil {
			c.JSON(http.StatusBadRequest, resources.JSONResult{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Data:    gin.H{},
			})
			return
		}
		serialized, err := json.Marshal(result)
		err = cache.Rdb0.Set(cache.Ctx, key, serialized, 60*time.Minute).Err()
		if err != nil {
			c.JSON(http.StatusBadRequest, resources.JSONResult{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Data:    gin.H{},
			})
			return
		}
		c.JSON(http.StatusOK, resources.JSONResult{
			Code:    http.StatusOK,
			Message: "ok",
			Data:    resources.AdvertistmentCollection(result),
		})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, resources.JSONResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    gin.H{},
		})
		return
	}

	var deserialized []models.Advertistment
	err = json.Unmarshal([]byte(val), &deserialized)
	c.JSON(http.StatusOK, resources.JSONResult{
		Code:    http.StatusOK,
		Message: "ok",
		Data:    resources.AdvertistmentCollection(deserialized),
	})
}

// @Summary 新增廣告
// @Produce  json
// @tags ad
// @Param data body FormAddAdvertistment true "廣告資料"
// @success 200 {object} resources.JSONResult{data=interface{}}
// @Failure 400 {object} resources.JSONResult{data=interface{}}
// @Router /ad [post]
func AddAdvertistment(c *gin.Context) {
	formAddBot := FormAddAdvertistment{}
	err := c.ShouldBindBodyWith(&formAddBot, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, resources.JSONResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    gin.H{},
		})
		return
	}

	genderModelArr, err := models.StringArrayToGenderModelArray(*formAddBot.FormAddAdvertistmentCondition.Gender)
	if err != nil {
		c.JSON(http.StatusBadRequest, resources.JSONResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    gin.H{},
		})
		return
	}

	countryModelArr, err := models.StringArrayToCountryModelArray(*formAddBot.FormAddAdvertistmentCondition.Country)
	if err != nil {
		c.JSON(http.StatusBadRequest, resources.JSONResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    gin.H{},
		})
		return
	}

	platformModelArr, err := models.StringArrayToPlatformModelArray(*formAddBot.FormAddAdvertistmentCondition.Platform)
	if err != nil {
		c.JSON(http.StatusBadRequest, resources.JSONResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    gin.H{},
		})
		return
	}

	count, err := models.GetTodayCreateAdCount()
	if err != nil {
		c.JSON(http.StatusBadRequest, resources.JSONResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    gin.H{},
		})
		return
	}
	if count >= int64(config.AppConfig.Detail.MaxAdEveryday) {
		c.JSON(http.StatusBadRequest, resources.JSONResult{
			Code:    http.StatusBadRequest,
			Message: "No more ads can be added today :(",
			Data:    gin.H{},
		})
		return
	}

	activedCount, err := models.CheckStartAtAndEndAt(formAddBot.StartAt, formAddBot.EndAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, resources.JSONResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    gin.H{},
		})
		return
	}
	if activedCount >= int64(config.AppConfig.Detail.MaxActiveAd)-1 {
		c.JSON(http.StatusBadRequest, resources.JSONResult{
			Code:    http.StatusBadRequest,
			Message: "No more ads can be added :(, choose other start_at and end_at",
			Data:    gin.H{},
		})
		return
	}

	newAd := models.Advertistment{
		Title:    &formAddBot.Title,
		StartAt:  &formAddBot.StartAt,
		EndAt:    &formAddBot.EndAt,
		AgeStart: formAddBot.FormAddAdvertistmentCondition.AgeStart,
		AgeEnd:   formAddBot.FormAddAdvertistmentCondition.AgeEnd,
		Gender:   &genderModelArr,
		Country:  &countryModelArr,
		Platform: &platformModelArr,
	}

	err = newAd.Create()
	cache.Rdb0.FlushAll(cache.Ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, resources.JSONResult{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, resources.JSONResult{
		Code:    http.StatusOK,
		Message: "ok",
		Data:    gin.H{},
	})
}
