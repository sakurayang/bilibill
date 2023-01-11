package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func TimeParse(date string) (time.Time, error) {
	arr := strings.Split(date, " ")
	tn := time.Now()
	var d, t string
	if len(arr) == 0 {
		d, t = "", ""
	} else if len(arr) == 1 {
		d, t = arr[0], ""
	} else {
		d, t = arr[0], arr[1]
	}

	arr = strings.Split(d, "-")
	if arr[0] == "" {
		arr = []string{}
	}
	var year, month, day string
	if len(arr) == 0 {
		year, month, day = strconv.Itoa(tn.Year()), strconv.Itoa(int(tn.Month())), strconv.Itoa(tn.Day())
	} else if len(arr) == 1 {
		year, month, day = arr[0], strconv.Itoa(int(tn.Month())), strconv.Itoa(tn.Day())
	} else if len(arr) == 2 {
		year, month, day = arr[0], arr[1], strconv.Itoa(tn.Day())
	} else {
		year, month, day = arr[0], arr[1], arr[2]
	}

	arr = strings.Split(t, ":")
	if arr[0] == "" {
		arr = []string{}
	}
	var hour, minute, second string
	if len(arr) == 0 {
		hour, minute, second = strconv.Itoa(tn.Hour()), strconv.Itoa(tn.Minute()), strconv.Itoa(tn.Second())
	} else if len(arr) == 1 {
		hour, minute, second = arr[0], strconv.Itoa(tn.Minute()), strconv.Itoa(tn.Second())
	} else if len(arr) == 2 {
		hour, minute, second = arr[0], arr[1], strconv.Itoa(tn.Second())
	} else {
		hour, minute, second = arr[0], arr[1], arr[2]
	}

	p, err := time.Parse(time.RFC3339, fmt.Sprintf("%04s-%02s-%02sT%02s:%02s:%02s+08:00", year, month, day, hour, minute, second))
	return p, err
}
