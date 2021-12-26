package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	readability "github.com/go-shiori/go-readability"
	gofeed "github.com/mmcdole/gofeed"
	mail "github.com/xhit/go-simple-mail/v2"
)

type SECREAT_DATA struct {
	SERVER_USERNAME string
	SERVER_PWD      string
	SERVER_HOST     string
	EMAIL_ADDTO     string
	EMAIL_ADDCC     string
	EMAIL_KINDLE    string
	PATH_1          string
	PATH_2          string
}
type any = interface{}

func sendEmail(paths []string) {
	server := mail.NewSMTPClient()

	// SMTP Server
	server.Host = _data.SERVER_HOST
	server.Port = 587
	server.Username = _data.SERVER_USERNAME
	server.Password = _data.SERVER_PWD
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// SMTP client
	smtpClient, err := server.Connect()

	if err != nil {
		log.Fatal(err)
	}

	// New email simple html with inline and CC
	email := mail.NewMSG()
	email.SetFrom("convert " + _data.SERVER_USERNAME).
		AddTo(_data.EMAIL_ADDTO).
		AddCc(_data.EMAIL_ADDCC).
		SetSubject("convert")

	for _, path := range paths {
		email.Attach(&mail.File{FilePath: path, Name: path, Inline: true})
	}

	// always check error after send
	if email.Error != nil {
		log.Fatal(email.Error)
	} else {
		log.Println("OOOOK")
	}

	// Call Send and pass the client
	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email Sent ")
	}
}

func takeHtmlElement(urls []string, titles []string) {
	e := os.Mkdir("tmp", 0700)
	if e != nil {
		log.Fatalln(e)
	}
	for i, url := range urls {
		article, err := readability.FromURL(url, 30*time.Second)
		if err != nil {
			log.Fatalf("failed to parse %s, %v\n", url, err)
		}

		dstHTMLFile, _ := os.Create(fmt.Sprintf("./tmp/html-%02d.html", i+1))
		defer dstHTMLFile.Close()

		dstHTMLFile.WriteString(article.Content)

		htmlBytes, e := ioutil.ReadFile(fmt.Sprintf("./tmp/html-%02d.html", i+1))
		if e != nil {
			log.Fatalln(e)
		}
		createPDFFromHtml(string(htmlBytes), i, titles[i])
	}
}

func createPDFFromHtml(_html string, i int, title string) {

	// For use wkhtml, install first -> https://wkhtmltopdf.org/downloads.html
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		log.Fatalln(err)
	}

	pdfg.AddPage(wkhtml.NewPageReader(strings.NewReader(_html)))

	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	dN := "./tmp/" + title
	err = pdfg.WriteFile(dN)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
	fmt.Println("URL CREATE PDF: ", title)
	fmt.Println("")
}

func formatterTitle(title string) string {
	title = strings.ReplaceAll(title, "ì", "i")
	title = strings.ReplaceAll(title, "è", "e")
	title = strings.ReplaceAll(title, "ò", "o")
	title = strings.ReplaceAll(title, "ù", "u")
	title = strings.ReplaceAll(title, ":", "")
	title = strings.ReplaceAll(title, ",", "")
	title = strings.ReplaceAll(title, ".", "")
	title = strings.ReplaceAll(title, "-", "")
	title = strings.ReplaceAll(title, " ", "_")
	title = strings.ReplaceAll(title, "__", "_")
	return title
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

func orchestrator(url string) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)
	items := feed.Items

	cryptoNews := []string{}
	cryptoTitles := []string{}

	// feeds := takeFeeds()

	// for _, feed := range feeds {
	for index, item := range items {
		if index <= 5 {
			cryptoNews = append(cryptoNews, item.Link)
			cryptoTitles = append(cryptoTitles, formatterTitle(item.Title)+".pdf")
		}
	}
	// }

	takeHtmlElement(cryptoNews, cryptoTitles) // -> load application

	if cryptoTitles != nil {
		// sendEmail(cryptoTitles)
		removeContents("tmp")
	}
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	os.Remove(dir)
	return nil
}

var _data SECREAT_DATA

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

func main() {
	readData("keys")
	orchestrator(_data.PATH_1)
}