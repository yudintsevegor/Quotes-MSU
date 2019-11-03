package main

import (
	"encoding/json"
	"strconv"
	"time"
)

type Response struct {
	Body Body `json:"response"`
}

type Body struct {
	Count  int     `json:count`
	Items  []Data  `json:"items"`
	Groups []Group `json:"groups`
}

type Data struct {
	ID ID `json:"id"`

	Date Date   `json:"date`
	Text string `json::text`
}

type Group struct {
	ID         int    `json:"id"`
	ScreenName string `json:screen_name`
}

type ID string

func (id *ID) UnmarshalJSON(b []byte) error {
	var i int
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	*id = ID(strconv.Itoa(i))

	return nil
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
