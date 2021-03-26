package main

import (
	"encoding/json"
	_ "fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

type Match struct {
    Time time.Time
    Opponent string
    Venue string
    Competition string
    StreamLink string
    StreamTitle string
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	files, err := ioutil.ReadDir("./data")
    if err != nil {
        log.Fatal(err)
    }

	var matches []Match
    for _, file := range files {
		dat, _ := ioutil.ReadFile("./data/" + file.Name())
		var fileMatches []Match
		if err := json.Unmarshal(dat, &fileMatches); err != nil {
			panic(err)
		}
		matches = append(matches, fileMatches...)
    }

	sort.Slice(matches[:], func(i, j int) bool {
		return matches[i].Time.Before(matches[j].Time)
	})

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		y, m, d := time.Now().Date()
		var today = time.Date(y, m, d, 0, 0, 0, 0, time.Now().Location())
		var index int
		var element Match
		for index, element = range matches {
			if element.Time.After(today) {
				break;
			}
		}
		maxLength:= len(matches)
		if maxLength > index + 5 {
			maxLength = index + 5
		}
		activeMatches := matches[index:maxLength]
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H {
			"Message": "",
			"Markup": "<b>Another Test</b>",
			"Matches": activeMatches,
			"Matchday": false,
		})
	})

	router.Run(":" + port)
}
