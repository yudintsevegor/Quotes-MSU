package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// restrictions:
// wall.get — 5000 вызовов в сутки. -> ~1 req per 20seconds
// https://vk.com/dev/data_limits

const (
	method  = "wall.get"
	version = "5.102"
	vkUrl   = "https://vk.com/"
)

func getVkPosts(groupID string) <-chan string {
	/**/
	count := "5"
	offset := "0"

	u := url.Values{}
	u.Set("count", count)
	u.Set("offset", offset)

	u.Set("owner_id", groupID)
	u.Set("access_token", ServiceKey)

	u.Set("v", version)
	u.Set("extended", "1") // is it really important?

	url := fmt.Sprintf("https://api.vk.com/method/%v?", method)
	path := url + u.Encode()

	out := make(chan string)
	go func() {
		var isFirstReq = true
		var corner int
		var zeroLevel int
		for {
			body, err := getPosts(path)
			if err != nil {
				log.Println(err)
				continue
			}

			// only one groupID in request before
			if len(body.Groups) != 1 {
				log.Println(errors.New("empty info about group"))
				continue
			}

			corner = body.Count - zeroLevel
			if isFirstReq {
				isFirstReq = false
				zeroLevel = body.Count
				corner = len(body.Items)
			}

			for i := corner - 1; i >= 0; i-- {
				out <- makeMessage(body.Items[i], groupID)
			}
			time.Sleep(20 * time.Second)
		}
	}()

	return out
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

func makeMessage(v Data, groupID string) string {
	return v.Text + "\n\n" + makeLink(string(v.ID), groupID)
}

func makeLink(id, groupID string) string {
	return vkUrl + "wall" + groupID + "_" + id
}
