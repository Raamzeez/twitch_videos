package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type video struct {
	ID            string
	Stream_id     string
	User_id       string
	User_login    string
	User_name     string
	Title         string
	Description   string
	Created_at    time.Time
	Published_at  time.Time
	Url           string
	Thumbnail_url string
	Viewable      bool
	View_count    int
	Language      string
	Type          string
	Duration      string
}

type videoData struct {
	Data []video
}

type user struct {
	ID                string
	Login             string
	Display_name      string
	Type              string
	Broadcaster_type  string
	Description       string
	Profile_image_url string
	Offline_image_url string
	View_count        int
	Create_at         time.Time
}

type userData struct {
	Data []user
}

type Response struct {
	Status  int
	Message []string
	Error   []string
}

func sendResponse(c *gin.Context, response Response) {
	if len(response.Message) > 0 {
		c.JSON(response.Status, map[string]interface{}{"message": strings.Join(response.Message, "; ")})
	} else if len(response.Error) > 0 {
		c.JSON(response.Status, map[string]interface{}{"error": strings.Join(response.Error, "; ")})
	}
}

func fetchUserId(context *gin.Context, username string, baseUserURL string, bearer string, CLIENT_ID string) string {
	var user_data userData
	apiURL := baseUserURL + "?login=" + username
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Client-id", CLIENT_ID)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errorResponse := Response{Status: resp.StatusCode, Message: []string{"Unable to make request to get user id"}, Error: []string{err.Error()}}
		sendResponse(context, errorResponse)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorResponse := Response{Status: resp.StatusCode, Message: []string{"Unable to read data from API Response"}, Error: []string{err.Error()}}
		sendResponse(context, errorResponse)
	}
	fmt.Println(string(body))
	json.Unmarshal([]byte(body), &user_data)
	return user_data.Data[0].ID
}

func main() {

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	CLIENT_ID := os.Getenv("CLIENT_ID")
	baseVideoURL := os.Getenv("BASE_VIDEO_URL")
	baseUserURL := os.Getenv("BASE_USER_URL")
	bearer := os.Getenv("BEARER_HEADER")

	router := gin.Default()
	router.GET("/user/:username", func(context *gin.Context) {
		var json_data videoData
		twitch_user_name := context.Param("username")
		twitch_user_id := fetchUserId(context, twitch_user_name, baseUserURL, bearer, CLIENT_ID)
		apiURL := baseVideoURL + "?user_id=" + twitch_user_id
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			errorResponse := Response{Status: http.StatusUnauthorized, Message: []string{"Cannot get video data from Twitch API for videos"}, Error: []string{err.Error()}}
			sendResponse(context, errorResponse)
		}
		req.Header.Add("Authorization", bearer)
		req.Header.Add("Client-id", CLIENT_ID)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errorResponse := Response{Status: http.StatusUnauthorized, Message: []string{"Unable to receive data from Twitch API for videos"}, Error: []string{err.Error()}}
			sendResponse(context, errorResponse)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errorResponse := Response{Status: http.StatusUnauthorized, Message: []string{"Cannot read API Data for videos from Twitch"}, Error: []string{err.Error()}}
			sendResponse(context, errorResponse)
		}
		json.Unmarshal([]byte(body), &json_data)
		context.JSON(200, gin.H{
			"response": string(body),
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
