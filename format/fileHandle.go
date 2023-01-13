package format

import (
	"fmt"
	"github.com/sakurayang/bilibill/biliapi"
	"github.com/sakurayang/bilibill/config"
	"github.com/sakurayang/bilibill/util"
	"os"
	"path"
	"strings"
	"time"
)

var C = config.C

type BillMode int

const (
	MonthBill BillMode = iota
	DailyBill
)

type OFormat int

const (
	CSV = iota
	JSON
)

func GetFileName(mode BillMode, format OFormat, date string) string {
	var extra string
	if C.Debug {
		extra = "_debug"
	} else {
		extra = ""
	}

	ext := "csv"
	switch format {
	default:
	case CSV:
		ext = "csv"
		break
	case JSON:
		ext = "json"
		break
	}

	t, _ := util.TimeParse(date)
	y, m, d := t.Date()
	if mode == MonthBill {
		return fmt.Sprintf("%04d-%02d流水账单%s.%s", y, m, extra, ext)
	} else if mode == DailyBill {
		return fmt.Sprintf("%04d-%02d-%02d流水账单%s.%s", y, m, d, extra, ext)
	} else {
		return GetFileName(DailyBill, format, date)
	}
}

func getFullPath(name string) string {
	return path.Join(C.Output, name)
}

func isExist(name string) bool {
	_, err := os.Stat(getFullPath(name))
	return err == nil
}

func WriteFile(mode BillMode, format OFormat, list *biliapi.OutStrings, date string) error {
	now := time.Now().Unix()
	filename := GetFileName(mode, format, date)
	if isExist(filename) {
		if err := os.Rename(getFullPath(filename), strings.Replace(filename, ".", fmt.Sprintf("_%d.", now), -1)); err != nil {
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

	for _, value := range *list {
		b := []byte(value)
		_, err = file.Write(b)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}
