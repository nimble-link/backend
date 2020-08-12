package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func OAuth2Handler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	if err != nil {
		c.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	code := string(body[:])
	savedUser, err := exchangeCode(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	token := savedUser.GenerateAccessToken()

	c.JSON(http.StatusOK, token)
}

func exchangeCode(code string) (*models.User, error) {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
		RedirectURL:  "postmessage",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println("TOKEN", err.Error())
		return nil, err
	}
	client := conf.Client(oauth2.NoContext, token)

	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Println("CLIENT", err.Error())
		return nil, err
	}
	defer userInfo.Body.Close()

	data, _ := ioutil.ReadAll(userInfo.Body)
	user := new(models.User)
	json.Unmarshal(data, user)

	savedUser := models.FindUserByEmail(user.Email)
	if savedUser.ID == 0 {
		user.Save()
		savedUser = user
	}

	return savedUser, nil
}
