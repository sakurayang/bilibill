package biliapi

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sakurayang/bilibill/config"
	"github.com/sakurayang/bilibill/util"
	"strconv"
	"time"
)

func getClient() *resty.Client {
	return resty.New().
		SetBaseURL("https://api.live.bilibili.com").
		SetDebug(config.C.Debug).
		SetHeader("Cookie", config.C.Cookie)
}

func GetGiftTypes() (*[]GiftTypes, error) {
	res, err := getClient().R().
		SetResult(GiftTypesResponse{}).
		Get("/gift/v1/master/getGiftTypes")
	if err != nil {
		return &[]GiftTypes{}, err
	}

	if res.StatusCode() != 200 {
		return &[]GiftTypes{}, fmt.Errorf(res.Status())
	}

	result := res.Result().(*GiftTypesResponse)

	if result.Code != 0 {
		return &result.Data, fmt.Errorf(result.Message)
	}

	return &result.Data, nil
}

func getGiftList(lastId int, date string) (*GiftListResponse, error) {
	t, err := util.TimeParse(date)
	if err != nil {
		return &GiftListResponse{}, err
	}
	res, err := getClient().R().
		SetResult(GiftListResponse{}).
		SetQueryParams(map[string]string{
			"limit":      "100",
			"coin_type":  "0",
			"begin_time": fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day()),
			"last_id":    strconv.Itoa(lastId),
		}).
		Get("/xlive/revenue/v1/giftStream/getReceivedGiftStreamNextList")

	if err != nil {
		return &GiftListResponse{}, err
	}

	if res.StatusCode() != 200 {
		return &GiftListResponse{}, fmt.Errorf(res.Status())
	}

	result := res.Result().(*GiftListResponse)

	if result.Code != 0 {
		return result, fmt.Errorf(result.Message)
	}

	return result, nil
}

func GetDailyGiftList(date string) (*AllList, error) {
	var list []Gift
	var allList AllList

	flag := true
	lastId := 0

	for flag {
		result, err := getGiftList(lastId, date)
		if err != nil {
			return &AllList{}, err
		}

		list = append(list, result.Data.List...)

		if result.Data.HasMore == 0 {
			flag = false
			allList.TotalHamster = result.Data.TotalHamster
			break
		}
		lastId = result.GetLastId()
	}

	allList.List = &list
	return &allList, nil
}

func GetMonthGiftList(date string) (*AllList, error) {
	var list []Gift
	var allList AllList

	t, err := util.TimeParse(date)
	if err != nil {
		return &AllList{}, err
	}
	day := t.Day()
	for i := 1; i <= day; i++ {
		l, err := GetDailyGiftList(t.Format("2006-01-02 15:04:05"))
		if err != nil {
			return &allList, err
		}
		allList.TotalHamster += l.TotalHamster
		list = append(list, *l.List...)
		t.Add(time.Hour * 24)
	}
	allList.List = &list
	return &allList, nil
}
