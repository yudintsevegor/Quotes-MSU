package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	groupID = "-65652356"
	method  = "wall.get"
	version = "5.102"
	vkUrl   = "https://vk.com/"
)

func main() {
	count := "3"
	offset := "1"

	u := url.Values{}
	u.Set("v", version)
	u.Set("count", count)
	u.Set("offset", offset)
	u.Set("owner_id", groupID)
	u.Set("access_token", ServiceKey)
	u.Set("extended", "1") // is it important?

	url := fmt.Sprintf("https://api.vk.com/method/%v?", method)
	path := url + u.Encode()

	body, err := getPosts(path)
	if err != nil {
		log.Fatal(err)
	}

	if len(body.Groups) != 1 {
		log.Fatal(errors.New("empty info about group"))
	}

	for _, v := range body.Items {
		log.Println(time.Time(v.Date))
		log.Println(v.Text)
		fmt.Println()
		id := strconv.Itoa(v.ID)
		link := makeLink(id)
		log.Println(link)
		fmt.Println()
	}

}

func makeLink(id string) string {
	return vkUrl + "wall" + groupID + "_" + id
}

func getPosts(path string) (Body, error) {
	response, err := http.Get(path)
	if err != nil {
		return Body{}, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Body{}, err
	}

	resp := new(Response)
	if err := json.Unmarshal(body, resp); err != nil {
		return Body{}, err
	}

	return resp.Body, nil
}
