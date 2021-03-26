package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

type Match struct {
    Time string
    Opponent string
    Venue string
    StreamLink string
    StreamTitle string
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		var matches []Match
		matches = append(matches, Match {Time: time.Now().Format("Mon Jan 2 15:04"), Opponent: "US Soccer", Venue: "The Galaxy", StreamLink: "https://youtube.com", StreamTitle: "Youtube"})
		fmt.Printf("%+v\n", matches)
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H {
			"Message": "Test",
			"Markup": "<b>Another Test</b>",
			"Matches": matches,
		})
	})

	router.Run(":" + port)
}
