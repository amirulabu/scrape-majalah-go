package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func getPage(page string) string {
	return "https://www.majalah.com/?allclassifieds.page." + page
}

type Iklan struct {
	Title   string
	Url     string
	Content string
	Comment []Comment
}

type Comment struct {
	Content  string
	Username string
	DateTime string
}

func main() {

	c := colly.NewCollector(colly.AllowedDomains("www.majalah.com", "majalah.com"))

	// Create another collector to scrape iklan details
	iklanCollector := c.Clone()

	iklanList := make([]Iklan, 0)

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
				fmt.Println("====" + strconv.Itoa(i) + "====")
				fmt.Println(h.Text)
				fmt.Println("@~@~@~@~@~@~@~@~@~@~@~@~@~@~@~@~@~@~@~@~@~@~@")
			}
		})
		// iklanList[len(iklanList)-1].Content = h.Text
		// commentList := make([]Comment, 0)
		// comment := Comment{}
		// comment.Content =

		// #contentLeft > div:nth-child(1) > div:nth-child(4)
		// #contentLeft > div:nth-child(1) > div:nth-child(5)
	})

	c.OnHTML("#contentLeft > div > b > a", func(h *colly.HTMLElement) {
		iklan := Iklan{}

		iklan.Title = h.Text
		iklan.Url = h.Attr("href")

		iklanCollector.Visit(h.Attr("href"))
		time.Sleep(2 * time.Second)

		// fmt.Println(h.Text)
		// fmt.Println(h.Attr("href"))

		iklanList = append(iklanList, iklan)
	})

	c.Visit(getPage("3000"))

	// iklanJson := json.NewEncoder(os.Stdout)
	// iklanJson.SetIndent("", "  ")
	// iklanJson.Encode(iklanList)
}

func checkUsableText(s string) bool {
	if strings.Contains(s, "Your Comment: Max 1000 characters.  Login Email: Password:") {
		return false
	}

	if strings.Contains(s, "function RotateImages") {
		return false
	}

	if strings.Contains(s, "Your Comment: Max 1000 characters.") {
		return false
	}

	if strings.Contains(s, "Disclaimer. Messages posted to our forum are solely the opinion and responsibility of the person posting the message.") {
		return false
	}

	if strings.Contains(s, "1. NEVER give UPFRONT PAYMENT (deposit) to any Money Lenders. Upfront Payment is 100% scam!") {
		return false
	}

	if strings.Contains(s, "Login Email:") {
		return false
	}

	if strings.Contains(s, "Password:") {
		return false
	}

	if strings.Contains(s, "(Total ") {
		return false
	}

	if strings.Contains(s, "Smartphone Biasa Kini Boleh Digunakan untuk Membuat Video Berkualiti Media Sosial") {
		return false
	}

	if len(strings.TrimSpace(s)) < 2 {
		return false
	}

	return true
}
