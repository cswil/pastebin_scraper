package main

import (
	"fmt"
	"os"
	"time"

	"github.com/anaskhan96/soup"
)

// function we use to check if we've already visited the link, dont want dupes
func contains(arr []string, str string) bool {
	for _, a := range arr { // linear search used :/
		if a == str {
			fmt.Println("returning true, ", a, " + ", str)
			return true

		}
	}
	return false
}

func main() {

	for true {
		var mylinks []string                          // new links
		var visitedlinks []string                     // links weve visited
		resp, err := soup.Get("https://pastebin.com") // fetch stuff on the page
		if err != nil {
			fmt.Print(err)
		}

		bod := soup.HTMLParse(resp)
		links := bod.Find("div", "id", "menu_2").FindAll("a") // scrape the links on the page
		fmt.Println(links)

		for _, link := range links {
			fmt.Println(link.Attrs()["href"])                                                // extract only the links
			fullurl := fmt.Sprintf("%s%s", "https://pastebin.com/raw", link.Attrs()["href"]) // concat pastebin link with href we go to raw so we dont have to reparse
			mylinks = append(mylinks, fullurl)

		}

		for _, link := range mylinks {
			resplink, err := soup.Get(link)
			if err != nil {
				fmt.Println("Error stage 2 of program : ", err)
			}
			body := soup.HTMLParse(resplink)
			filecontents := body.FullText() // on raw so we just get all text on the page
			skip := contains(visitedlinks, link)
			fmt.Println(skip)
			fmt.Print(link)
			if skip == false {
				f, err := os.OpenFile("/tmp/pastescrapes", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // writes all things to this file. change as you please
				if err != nil {
					fmt.Println(err)
				}
				defer f.Close()
				if _, err := f.WriteString(filecontents); err != nil {
					fmt.Println(err)
				}
				fmt.Println("Written to file ")
				visitedlinks = append(visitedlinks, link)
				time.Sleep(1 * time.Second)
			}

		}
		mylinks = nil // clear our links
		fmt.Println("WAITING 5 SECONDS")
		time.Sleep(5 * time.Second)
	}

}
