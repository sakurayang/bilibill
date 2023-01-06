package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

type GiftListResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		List         []Gift `json:"list"`
		HasMore      int    `json:"has_more"`
		TotalHamster int    `json:"total_hamster"`
	} `json:"data"`
}

func (g *GiftListResponse) GetLastId() int {
	length := len(g.Data.List)
	return g.Data.List[length-1].Id
}

type Gift struct {
	Uid              int    `json:"uid"`
	Uname            string `json:"uname"`
	Time             string `json:"time"`
	GiftId           int    `json:"gift_id"`
	GiftName         string `json:"gift_name"`
	GiftImg          string `json:"gift_img"`
	GiftNum          int    `json:"gift_num"`
	Hamster          int    `json:"hamster"`
	Gold             int    `json:"gold"`
	Silver           int    `json:"silver"`
	IosHamster       int    `json:"ios_hamster"`
	NormalHamster    int    `json:"normal_hamster"`
	IosGold          int    `json:"ios_gold"`
	NormalGold       int    `json:"normal_gold"`
	IsHybrid         bool   `json:"is_hybrid"`
	Id               int    `json:"id"`
	IsOpenPlatfrom   int    `json:"is_open_platfrom"`
	OpenPlatfromRate int    `json:"open_platfrom_rate"`
	ReceiveTitle     string `json:"receive_title"`
	RoomId           int    `json:"room_id"`
}

type AllList struct {
	List         *[]Gift `json:"list"`
	TotalHamster int     `json:"total_hamster"`
}

type GiftTypesResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Message string      `json:"message"`
	Data    []GiftTypes `json:"data,omitempty"`
}
type GiftTypes struct {
	GiftId   int    `json:"gift_id"`
	GiftName string `json:"gift_name"`
	Price    int    `json:"price,omitempty"`
}

var client = resty.New().
	SetBaseURL("https://api.live.bilibili.com").
	SetDebug(c.AppConfig.Debug).
	SetHeader("Cookie", c.BiliConfig.Cookie)

func GetGiftTypes() (*[]GiftTypes, error) {
	res, err := client.R().
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

func getGiftList(lastId int, day int) (*GiftListResponse, error) {
	res, err := client.R().
		SetResult(GiftListResponse{}).
		SetQueryParams(map[string]string{
			"limit":      "100",
			"coin_type":  "0",
			"begin_time": fmt.Sprintf("%04d-%02d-%02d", time.Now().Year(), time.Now().Month(), day),
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

func GetDailyGiftList(day int) (*AllList, error) {
	var list []Gift
	var allList AllList

	flag := true
	lastId := 0

	for flag {
		result, err := getGiftList(lastId, day)
		if err != nil {
			return &AllList{}, nil
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

func GetMonthGiftList() (*AllList, error) {
	var list []Gift
	var allList AllList

	day := time.Now().Day()
	for i := 1; i < day; i++ {
		l, err := GetDailyGiftList(i)
		if err != nil {
			return &allList, err
		}
		allList.TotalHamster += l.TotalHamster
		list = append(list, *l.List...)
	}
	allList.List = &list
	return &allList, nil
}
