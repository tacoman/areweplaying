package main

import (
	"encoding/json"
	"fmt"
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
        fmt.Println(file.Name())
		dat, _ := ioutil.ReadFile("./data/" + file.Name())
		fmt.Printf("%+v\n", dat)
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
		fmt.Printf("%+v\n", matches)
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H {
			"Message": "",
			"Markup": "<b>Another Test</b>",
			"Matches": matches,
		})
	})

	router.Run(":" + port)
}
