package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

const (
	Likes     Dimension = "likes"
	Comments            = "comments"
	Favorites           = "favorites"
	Retweets            = "retweets"
)

type Dimension string

type Media struct {
	Likes     int `json:"likes"`
	Comments  int `json:"comments"`
	Favorites int `json:"favorites"`
	Retweets  int `json:"retweets"`
	Timestamp int `json:"timestamp"`
}

type LikeResponse struct {
	TotalPosts    int `json:"total_posts"`
	MiniTimestamp int `json:"mini_timestamp"`
	MaxiTimestamp int `json:"maxi_timestamp"`
	AvgLikes      int `json:"avg_likes"`
}

type FavoriteResponse struct {
	TotalPosts    int `json:"total_posts"`
	MiniTimestamp int `json:"mini_timestamp"`
	MaxiTimestamp int `json:"maxi_timestamp"`
	AvgFavorites  int `json:"avg_favorites"`
}

type CommentResponse struct {
	TotalPosts    int `json:"total_posts"`
	MiniTimestamp int `json:"mini_timestamp"`
	MaxiTimestamp int `json:"maxi_timestamp"`
	AvgComments   int `json:"avg_comments"`
}

type RetweetResponse struct {
	TotalPosts    int `json:"total_posts"`
	MiniTimestamp int `json:"mini_timestamp"`
	MaxiTimestamp int `json:"maxi_timestamp"`
	AvgRetweet    int `json:"avg_retweet"`
}

func main() {
	r := gin.Default()
	r.GET("/analysis", analysisHandler)
	err := r.Run()

	if err != nil {
		panic(err)
	}
}

func analysisHandler(c *gin.Context) {

	duration := c.Query("duration")
	dimension := c.Query("dimension")

	durationParsed, dimensionParsed, err := parsingAnalysisParam(duration, dimension)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if durationParsed == nil || dimensionParsed == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to parse parameter",
		})
		return
	}

	tmp := *durationParsed

	ctx, cancel := context.WithTimeout(context.Background(), tmp)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://stream.upfluence.co/stream", nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	body, err := io.ReadAll(resp.Body)

	splitData := strings.Split(string(body), "data: ")[1:]

	var medias []Media

	// parse each media received
	for _, split := range splitData {
		println(split)
		var data map[string]Media
		err = json.Unmarshal([]byte(split), &data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		fmt.Printf("%+v\n", data)
		//We use get First Key because there is only one element in the map
		medias = append(medias, data[getFirstKey(data)])
	}

	analysis := analyzeMedias(medias, *dimensionParsed)

	c.Data(http.StatusOK, "application/json", analysis)
}

func analyzeMedias(medias []Media, dimension Dimension) []byte {
	var countDimension, lowerTimestamp, higherTimestamp int
	lowerTimestamp = math.MaxInt
	higherTimestamp = 0

	for _, media := range medias {
		switch dimension {
		case Likes:
			countDimension += media.Likes
		case Comments:
			countDimension += media.Comments
		case Favorites:
			countDimension += media.Favorites
		case Retweets:
			countDimension += media.Retweets
		}

		if media.Timestamp < lowerTimestamp {
			lowerTimestamp = media.Timestamp
		}
		if media.Timestamp > higherTimestamp {
			higherTimestamp = media.Timestamp
		}
	}
	var response []byte

	switch dimension {
	case Likes:
		tmp := LikeResponse{
			TotalPosts:    len(medias),
			MiniTimestamp: lowerTimestamp,
			MaxiTimestamp: higherTimestamp,
			AvgLikes:      countDimension / len(medias),
		}
		response, _ = json.Marshal(tmp)
	case Comments:
		tmp := CommentResponse{
			TotalPosts:    len(medias),
			MiniTimestamp: lowerTimestamp,
			MaxiTimestamp: higherTimestamp,
			AvgComments:   countDimension / len(medias),
		}
		response, _ = json.Marshal(tmp)
	case Favorites:
		tmp := FavoriteResponse{
			TotalPosts:    len(medias),
			MiniTimestamp: lowerTimestamp,
			MaxiTimestamp: higherTimestamp,
			AvgFavorites:  countDimension / len(medias),
		}
		response, _ = json.Marshal(tmp)
	case Retweets:
		tmp := RetweetResponse{
			TotalPosts:    len(medias),
			MiniTimestamp: lowerTimestamp,
			MaxiTimestamp: higherTimestamp,
			AvgRetweet:    countDimension / len(medias),
		}
		response, _ = json.Marshal(tmp)
	}
	return response
}

func parsingAnalysisParam(duration string, dimension string) (*time.Duration, *Dimension, error) {
	if duration == "" || dimension == "" {
		return nil, nil, errors.New("the query parameter duration or dimension is empty")
	}

	durationParsed, err := time.ParseDuration(duration)
	if err != nil {
		return nil, nil, errors.New("invalid duration parameter")
	}

	var dimensionParsed Dimension

	switch dimension {
	case "likes":
		dimensionParsed = Likes
	case "comments":
		dimensionParsed = Comments
	case "favorites":
		dimensionParsed = Favorites
	case "retweets":
		dimensionParsed = Retweets
	default:
		return nil, nil, errors.New("invalid dimension parameter")
	}

	return &durationParsed, &dimensionParsed, nil
}

func getFirstKey(data map[string]Media) string {
	for key, _ := range data {
		return key
	}
	return ""
}
