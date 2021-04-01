package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

type Match struct {
	Time        time.Time
	Opponent    string
	Venue       string
	Competition string
	StreamLink  string
	StreamTitle string
}

type CalendarEvents struct {
	Kind             string        `json:"kind"`
	Etag             string        `json:"etag"`
	Summary          string        `json:"summary"`
	Updated          time.Time     `json:"updated"`
	Timezone         string        `json:"timeZone"`
	Accessrole       string        `json:"accessRole"`
	Defaultreminders []interface{} `json:"defaultReminders"`
	Items            []struct {
		Kind     string    `json:"kind"`
		Etag     string    `json:"etag"`
		ID       string    `json:"id"`
		Status   string    `json:"status"`
		Htmllink string    `json:"htmlLink"`
		Created  time.Time `json:"created"`
		Updated  time.Time `json:"updated"`
		Summary  string    `json:"summary"`
		Location string    `json:"location"`
		Creator  struct {
			Email string `json:"email"`
		} `json:"creator"`
		Organizer struct {
			Email       string `json:"email"`
			Displayname string `json:"displayName"`
			Self        bool   `json:"self"`
		} `json:"organizer"`
		Start struct {
			Datetime time.Time `json:"dateTime"`
		} `json:"start"`
		End struct {
			Datetime string `json:"dateTime"`
		} `json:"end"`
		Icaluid   string `json:"iCalUID"`
		Sequence  int    `json:"sequence"`
		Eventtype string `json:"eventType"`
	} `json:"items"`
}

func getCalendarEvents() CalendarEvents {
	fmt.Println("Calling API...")

	client := &http.Client{}
	apiKey := os.Getenv("API_KEY")
	fmt.Println(apiKey)
	requestUrl := fmt.Sprint("https://www.googleapis.com/calendar/v3/calendars/qnjbamj73cgtn2bcgjmuojejt0%40group.calendar.google.com/events?key=", apiKey, "&timeMin=2021-03-31T00:00:00-00:00&singleEvents=true&orderBy=startTime")
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject CalendarEvents
	json.Unmarshal(bodyBytes, &responseObject)
	fmt.Printf("API Response as struct %+v\n", responseObject)
	return responseObject
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	events := getCalendarEvents()
	competitionMatcher := regexp.MustCompile(`\((.*?)\)`)

	var matches []Match
	for _, event := range events.Items {
		if strings.Contains(event.Summary, "at") ||
		strings.Contains(event.Summary, "vs") ||
		strings.Contains(event.Summary, "VS") ||
		strings.Contains(event.Summary, "AT") {
			match := Match{}
			match.Venue = event.Location
			match.Time = event.Start.Datetime
			match.Competition = competitionMatcher.FindStringSubmatch(event.Summary)[1]
			fmt.Println(event.Summary)
			matches = append(matches, match)
		}
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
		matchday := false
		for index, element = range matches {
			if today.Before(element.Time) {
				if element.Time.Day() == d {
					matchday = true
				}
				fmt.Println("breaking!")
				break
			}
		}
		maxLength := len(matches)
		if maxLength > index+5 {
			maxLength = index + 5
		}
		activeMatches := matches[index:maxLength]
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"Matches":  activeMatches,
			"Matchday": matchday,
		})
	})

	router.Run(":" + port)
}
