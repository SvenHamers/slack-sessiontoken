package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"bufio"
	"github.com/PuerkitoBio/goquery"
	"fmt"
)

const (
	baseURL = "https://<slackname>.slack.com"
)

var (
	username = ""
	password = ""
)

type App struct {
	Client *http.Client
}

func (app *App) getCrumb() string {
	loginURL := baseURL 
	client := app.Client

	response, err := client.Get(loginURL)

	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	token, _ := document.Find("input[name='crumb']").Attr("value")

	

	return token
}

func (app *App) GetToken() string {
	
	
	client := app.Client
	token := app.getCrumb()
	

	loginURL := baseURL 

	data := url.Values{
		"email":     {username},
		"password":  {password},
		"crumb": {token},
		"signin": {"1"},
		"redir": {""},
		"has_remember": {"1"},
		"remember": {"on"},

	}

	response, err := client.PostForm(loginURL, data)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil{
		log.Fatalln(err)
	}

	pageContent := string(body)
    scanner := bufio.NewScanner(strings.NewReader(pageContent))
    for scanner.Scan() {
		if strings.Contains(scanner.Text(), "api_token:"){
			return strings.Split(scanner.Text(), "\"")[1]			
		}     
    }	
	return ""
}


func main() {
	jar, _ := cookiejar.New(nil)

	app := App{
		Client: &http.Client{Jar: jar},
	}

	token := app.GetToken()
	fmt.Println(token);

}