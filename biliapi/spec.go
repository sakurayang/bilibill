package biliapi

import (
	"fmt"
	"github.com/sakurayang/bilibill/config"
	"sort"
	"strconv"
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

func (l *AllList) CSV() *OutStrings {
	sort.Slice(*l.List, func(i, j int) bool {
		return (*l.List)[i].Id < (*l.List)[j].Id
	})
	var stringList OutStrings
	if config.C.AppConfig.Debug {
		stringList = append(stringList,
			"\"编号\",\"b站昵称\",\"b站uid\",\"日期\",\"礼物名\",\"礼物id\",\"礼物数量\","+
				"\"金仓鼠数量\",\"金瓜子价值\",\"银瓜子价值\",\"iOS金仓鼠\",\"普通金仓鼠\","+
				"\"iOS金瓜子\",\"普通金瓜子\",\"流水id\",\"是否开放平台\",\"开放平台比率\"\n",
		)
		for i, value := range *l.List {
			isOpenPlatform, _ := strconv.ParseBool(strconv.Itoa(value.IsOpenPlatfrom))
			s := fmt.Sprintf(
				"%d,\"%s\",\"%d\",%s,\"%s\",\"%d\",%d,%d,%d,%d,%d,%d,%d,%d,\"%d\",%t,%d\n",
				i+1,
				value.Uname,
				value.Uid,
				value.Time,
				value.GiftName,
				value.GiftId,
				value.GiftNum,
				value.Hamster,
				value.Gold,
				value.Silver,
				value.IosHamster,
				value.NormalHamster,
				value.IosGold,
				value.NormalGold,
				value.Id,
				isOpenPlatform,
				value.OpenPlatfromRate,
			)
			stringList = append(stringList, s)
		}
	} else {
		stringList = append(stringList,
			"\"编号\",\"b站昵称\",\"b站uid\",\"日期\",\"礼物名\",\"礼物数量\",\"电池价值\",\"是否ios\",\"金仓鼠收益\"\n",
		)
		for i, value := range *l.List {
			isIOS := value.IosHamster != 0 || value.IosGold != 0
			s := fmt.Sprintf(
				"%d,\"%s\",\"%d\",%s,\"%s\",%d,%d,%t,%d\n",
				i+1,
				value.Uname,
				value.Uid,
				value.Time,
				value.GiftName,
				value.GiftNum,
				value.Gold/100,
				isIOS,
				value.Hamster,
			)
			stringList = append(stringList, s)
		}
	}
	return &stringList
}

func (l *AllList) JSON() *OutStrings {
	j := OutStrings{"{\"list\":["}
	tail := "]}"
	if config.C.AppConfig.Debug {
		for i, value := range *l.List {
			isOpenPlatform, _ := strconv.ParseBool(strconv.Itoa(value.IsOpenPlatfrom))
			s := fmt.Sprintf(
				"{\"uname\":\"%s\",\"uid\":%d,\"time\":\"%s\",\"gift_name\":\"%s\",\"gift_id\":%d,\"gift_num\":%d,"+
					"\"hamster\":%d,\"gold\":%d,\"sliver\":%d,\"ios_hamster\":%d,\"normal_hamster\":%d,\"ios_gold\":%d,"+
					"\"normal_gold\":%d,\"id\":%d,\"is_open_platform\":%t,\"open_platform_rate\":%d}",
				value.Uname,
				value.Uid,
				value.Time,
				value.GiftName,
				value.GiftId,
				value.GiftNum,
				value.Hamster,
				value.Gold,
				value.Silver,
				value.IosHamster,
				value.NormalHamster,
				value.IosGold,
				value.NormalGold,
				value.Id,
				isOpenPlatform,
				value.OpenPlatfromRate,
			)
			if i != len(*l.List)-1 {
				s += ","
			}
			j = append(j, s)
		}
	} else {
		for i, value := range *l.List {
			isIOS := value.IosHamster != 0 || value.IosGold != 0
			s := fmt.Sprintf(
				"{\"uname\":\"%s\",\"uid\":%d,\"time\":\"%s\",\"gift_name\":\"%s\",\"gift_num\":%d,\"battery\":%d,\"is_ios\":%t,\"hamster\":%d}",
				value.Uname,
				value.Uid,
				value.Time,
				value.GiftName,
				value.GiftNum,
				value.Gold/100,
				isIOS,
				value.Hamster,
			)
			if i != len(*l.List)-1 {
				s += ","
			}
			j = append(j, s)
		}
	}

	j = append(j, tail)
	return &j
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
