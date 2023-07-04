package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func getPage(page string) string {
	return "https://www.majalah.com/?allclassifieds.page." + page
}

func main() {

	filename := "majalah.txt"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", filename, err)
		return
	}

	defer f.Close()

	c := colly.NewCollector(colly.AllowedDomains("www.majalah.com", "majalah.com"))

	// Create another collector to scrape iklan details
	iklanCollector := c.Clone()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	iklanCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	iklanCollector.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	iklanCollector.OnHTML("#contentLeft > div", func(h *colly.HTMLElement) {
		h.ForEach("div", func(i int, h *colly.HTMLElement) {

			if checkUsableText(h.Text) && i != 0 && i != 1 {
				commentContent := removeUsernameDateTime(h.Text)
				f.WriteString(commentContent)
				f.WriteString("\n@@@~~~@@@\n")
			}
		})
	})

	c.OnHTML("#contentLeft > div > b > a", func(h *colly.HTMLElement) {

		// title := h.Text
		url := h.Attr("href")

		iklanCollector.Visit(url)
		// time.Sleep(2 * time.Second)

	})

	for i := 1; i < 10000; i++ {
		c.Visit(getPage(strconv.Itoa(i)))
	}

	// iklanJson := json.NewEncoder(os.Stdout)
	// iklanJson.SetIndent("", "  ")
	// iklanJson.Encode(iklanList)
}

func checkUsableText(s string) bool {
	unwantedText := []string{"Your Comment: Max 1000 characters.  Login Email: Password:",
		"function RotateImages",
		"Your Comment: Max 1000 characters.",
		"Disclaimer. Messages posted to our forum are solely the opinion and responsibility of the person posting the message.",
		"1. NEVER give UPFRONT PAYMENT (deposit) to any Money Lenders. Upfront Payment is 100% scam!",
		"Login Email:",
		"Password:",
		"(Total ",
		"Smartphone Biasa Kini Boleh Digunakan untuk Membuat Video Berkualiti Media Sosial",
	}

	for _, v := range unwantedText {
		if strings.Contains(s, v) {
			return false
		}
	}

	return len(strings.TrimSpace(s)) >= 2
}

func removeUsernameDateTime(s string) string {
	pattern := `(?m)^.*(?:Mon|Tue|Wed|Thu|Fri|Sat|Sun)\s\d{1,2}\/(?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)\/\d{4},\s\d{1,2}:\d{2}(?:am|pm)`

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	if regex.MatchString(s) {
		return regex.ReplaceAllString(s, "")
	} else {
		return s
	}
}
