package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type BillMode int

const (
	MonthBill BillMode = iota
	DailyBill
)

func timeParse(date string) (time.Time, error) {
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

func getFileName(mode BillMode, date string) string {
	var extra string
	if c.AppConfig.Debug {
		extra = "_debug"
	} else {
		extra = ""
	}

	t, _ := timeParse(date)
	y, m, d := t.Date()
	if mode == MonthBill {
		return fmt.Sprintf("%04d-%02d流水账单%s.csv", y, m, extra)
	} else if mode == DailyBill {
		return fmt.Sprintf("%04d-%02d-%02d流水账单%s.csv", y, m, d, extra)
	} else {
		return getFileName(DailyBill, date)
	}
}

func getFullPath(name string) string {
	return path.Join(c.AppConfig.Output, name)
}

func isExist(name string) bool {
	_, err := os.Stat(getFullPath(name))
	return err == nil
}

func writeFile(mode BillMode, list *[]string, date string) error {
	now := time.Now().Unix()
	filename := getFileName(mode, date)
	if isExist(filename) {
		err := os.Rename(getFullPath(filename), strings.Replace(filename, ".csv", fmt.Sprintf("_%d.csv", now), -1))
		if err != nil {
			return err
		}
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(file)

	if c.AppConfig.Debug {
		_, err = file.WriteString(fmt.Sprint(
			"\"编号\",\"b站昵称\",\"b站uid\",\"日期\",\"礼物名\",\"礼物id\",\"礼物数量\"," +
				"\"金仓鼠数量\",\"金瓜子价值\",\"银瓜子价值\",\"iOS金仓鼠\",\"普通金仓鼠\"," +
				"\"iOS金瓜子\",\"普通金瓜子\",\"流水id\",\"是否开放平台\",\"开放平台比率\"\n",
		))
		if err != nil {
			return err
		}
	} else {
		_, err = file.WriteString(fmt.Sprint("\"编号\",\"b站昵称\",\"b站uid\",\"日期\",\"礼物名\",\"礼物数量\",\"电池价值\",\"是否ios\",\"金仓鼠收益\"\n"))
		if err != nil {
			return err
		}
	}

	for index, value := range *list {
		_, err = file.WriteString(fmt.Sprintf("\"%d\",%s", index+1, value))
		if err != nil {
			return err
		}
	}

	if !c.AppConfig.Debug {
		_, err = file.WriteString(
			fmt.Sprintf(
				"\"总计\",,,,,,=SUM(G2:G%d),\"=COUNTIF(H2:H%d,TRUE)\",=SUM(I2:I%d)\n",
				len(*list)+1, len(*list)+1, len(*list)+1,
			),
		)
	}

	if err != nil {
		return err
	}

	return nil
}

func GetCSVString(list *AllList) *[]string {
	var stringList []string
	if c.AppConfig.Debug {
		for _, value := range *list.List {
			isOpenPlatform, _ := strconv.ParseBool(strconv.Itoa(value.IsOpenPlatfrom))
			s := fmt.Sprintf(
				"\"%s\",\"%d\",%s,\"%s\",\"%d\",%d,%d,%d,%d,%d,%d,%d,%d,\"%d\",%t,%d\n",
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
		for _, value := range *list.List {
			isIOS := value.IosHamster != 0 || value.IosGold != 0
			s := fmt.Sprintf(
				"\"%s\",\"%d\",%s,\"%s\",%d,%d,%t,%d\n",
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

func WriteCSV(mode BillMode, list *AllList) error {
	formatList := GetCSVString(list)
	return writeFile(mode, formatList, (*list.List)[0].Time)
}
