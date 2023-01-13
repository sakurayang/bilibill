package biliapi

import (
	"strings"
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

type OutStrings []string

func (o *OutStrings) String() string {
	return strings.Join(*o, "")
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
