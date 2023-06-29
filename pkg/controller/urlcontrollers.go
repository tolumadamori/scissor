package controller

import (
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tolumadamori/scissor/pkg/config"
	"github.com/tolumadamori/scissor/pkg/helper"
	"github.com/tolumadamori/scissor/pkg/models"
)

func ShortenURL(context *gin.Context) {
	var input models.ShortenRequest
	var id string

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := helper.CurrentUser(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Original URL validation
	if !govalidator.IsURL(input.OriginalURL) {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !helper.RemoveDomainError(input.OriginalURL) {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.OriginalURL = helper.EnforceHTTP(input.OriginalURL)

	//shortening logic
	if input.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = input.CustomShort
	}

	var shortenedURL models.URL
	shortenedURL.UserID = user.ID
	shortenedURL.OriginalURL = input.OriginalURL
	domain := os.Getenv("DOMAIN")
	shortenedURL.CustomShort = domain + "/" + id
	//we shorten the URL here.

	savedURL, err := shortenedURL.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"data": savedURL})
}

// Resolve URL
func ResolveURL(context *gin.Context) {
	url := context.Request.Host + context.Request.URL.Path
	var shortenedURL models.URL
	err := config.Database.Where("custom_short = ?", url).Find(&shortenedURL).Error
	value := shortenedURL.OriginalURL
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err, "url": url, "value": value})
		return
	}

	context.Redirect(301, value)

}

// Get All URLs for a user
func GetAllURLs(context *gin.Context) {
	user, err := helper.CurrentUser(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": user.URLs})
}

func Healthchecks(context *gin.Context) {
	context.JSON(http.StatusOK, 200)
}
