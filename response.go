package main

import (
	"encoding/json"
	"time"
)

type Response struct {
	Body Body `json:"response"`
}

type Body struct {
	Items  []Data  `json:"items"`
	Groups []Group `json:"groups`
}

type Data struct {
	ID int `json:"id"`

	Date Date   `json:"date`
	Text string `json::text`
}

type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
	var t int64
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}

	tm := time.Unix(t, 0)
	*d = Date(tm)

	return nil
}

type Group struct {
	ID         int    `json:"id"`
	ScreenName string `json:screen_name`
}