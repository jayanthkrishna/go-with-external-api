package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	PhotoApi = "https://api.pexels.xom/v1"
	VideoApi = "https://api.pexels.com/videos"
)

type Client struct {
	token          string
	hc             http.Client
	remainingTimes int32
}

func NewClient(token string) *Client {
	c := http.Client{}
	return &Client{token: token, hc: c}
}

type SearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Photos       []Photo `json:"photos"`
}

type Photo struct {
	Id              int32       `json:"id"`
	Width           int32       `json:"width"`
	Height          int32       `json:"height"`
	Url             string      `json:"url"`
	Photographer    string      `json:"photographer"`
	PhotoGrapherUrl string      `json:"photographer_url"`
	Src             PhotoSource `json:"src"`
}

type PhotoSource struct {
	Original  string `json:"original"`
	Large     string `json:"large"`
	Large2x   string `json:"large2x"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Portrait  string `json:"portrait"`
	Square    string `json:"square"`
	Landscape string `json:"landscape"`
	tiny      string `json:"tiny"`
}

type CuratedResult struct {
	Page     int32   `json:"page"`
	PerPage  int32   `json:"per_page"`
	NextPage string  `json:"next_page"`
	Photos   []Photo `json:"photos"`
}

func (c *Client) SearchPhotos(query string, perPage, page int) (*SearchResult, error) {
	url := fmt.Sprintf(PhotoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	resp, _ := c.requestDoWithAuth("GET", url)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	var result SearchResult

	err = json.Unmarshal(data, &result)

	return &result, err
}

func (c *Client) CuratedPhotos(perPage, page int) (*CuratedResult, error) {

	url := fmt.Sprintf(PhotoApi+"/curated?per_page=%d&page=%d", perPage, page)

	resp, err := c.requestDoWithAuth("GET", url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CuratedResult

	err = json.Unmarshal(data, &result)

	return &result, err
}

func (c *Client) GetPhoto(id int32) (*Photo, error) {
	url := fmt.Sprintf(PhotoApi+"/photos/%d", id)

	resp, err := c.requestDoWithAuth("GET", url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	var result Photo

	err = json.Unmarshal(data, &result)

	return &result, err

}

func (c *Client) GetRandomPhoto() (*Photo, error) {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(1001)
	res, err := c.CuratedPhotos(1, randNum)

	if err == nil && len(res.Photos) == 1 {
		return &res.Photos[0], nil
	}

	return nil, err

}

func (c *Client) SearchVideo(query string, perPage, page int) (*VideoSearchResult, error) {

}
func (c *Client) requestDoWithAuth(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.token)
	resp, err := c.hc.Do(req)

	if err != nil {
		return nil, err
	}

	times, err := strconv.Atoi(resp.Header.Get("x-Ratelimit-Remaining"))

	if err != nil {
		return resp, nil
	}

	c.remainingTimes = int32(times)

	return resp, err
}

func main() {
	os.Setenv("PexelsToken", "563492ad6f9170000100000143c2cb74dbfb49e7bb9d58bba926bb8c")
	TOKEN := os.Getenv("PexelsToken")
	c := NewClient(TOKEN)

	res, err := c.SearchPhotos("waves", 15, 5)

	if err != nil {
		fmt.Errorf("Search Error : %v", err)
	}
	if res.Page == 0 {
		fmt.Errorf("Search result wrong")
	}

	fmt.Println(res)

}
