package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jasonlvhit/gocron"
	gofeed "github.com/mmcdole/gofeed"
)

func main() {
	createScheduler("08:30")
}

func orchestrator(numberOfArticlesForPaths int) {
	cryptoNews := []string{}
	cryptoTitles := []string{}
	filters := _data.FILTER_ROLE

	feeds := takeFeeds()

	for _, feed := range feeds {
		for index, item := range feed.Items {
			if !FilterNews(filters, item.Title) && index <= numberOfArticlesForPaths {
				cryptoNews = append(cryptoNews, item.Link)
				cryptoTitles = append(cryptoTitles, formatterTitle(item.Title)+".pdf")
			}
		}
	}

	takeHtmlElement(cryptoNews, cryptoTitles) // -> load application

	if cryptoTitles != nil {
		sendEmail(cryptoTitles)
		removeContents("tmp")
	}
}

func createScheduler(time string){
	fmt.Println("*________________________*")
	fmt.Println("*                        *")
	fmt.Println("*        WELCOME         *")
	fmt.Println("*                        *")
	fmt.Println("*________________________*")
	log.Printf("Run Scheduler")
	s := gocron.NewScheduler()
	s.Every(1).Day().At(time).Do(task)
	<-s.Start()
}

func task(){
	log.Printf("Task created")
	readData("keys")
	orchestrator(5)
}

func takeFeeds() []*gofeed.Feed {
	arrayUrls := []string{}
	arrayFeeds := []*gofeed.Feed{}
	arrayUrls = append(arrayUrls, _data.PATH_1, _data.PATH_2)
	for _, url := range arrayUrls {
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(url)
		arrayFeeds = append(arrayFeeds, feed)
	}
	return arrayFeeds
}

func readData(path string) {
	// read file
	data, err := ioutil.ReadFile(path + ".json")
	if err != nil {
		fmt.Print(err)
	}

	// unmarshall it
	err = json.Unmarshal(data, &_data)
	if err != nil {
		fmt.Println("error:", err)
	}
}




