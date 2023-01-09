package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"regexp"
)

var c = GetConfig()

func main() {
	billMode := DailyBill
	date := "1970-01-01 00:00:00"
	app := &cli.App{
		Name:  "biliBill",
		Usage: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "BillMode",
				Usage:   "指定导出模式，day 为日账单，month 为月账单，默认为日账单",
				Aliases: []string{"m"},
				Action: func(ctx *cli.Context, v string) error {
					if v != "day" && v != "month" {
						return fmt.Errorf("导出模式应当为 'day' 或 'month' 而不是 %v", v)
					}
					if v == "month" {
						billMode = MonthBill
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:     "day",
				Usage:    "指定账单导出日期，输入 yyyy-MM-dd 时为指定日， 输入 MM-dd 时默认为今年， 仅输入日时默认为本月",
				Required: true,
				Aliases:  []string{"d"},
				Action: func(ctx *cli.Context, v string) error {
					dateRegex := regexp.MustCompile(`[0-9]{4}-[0-9]{1,2}-[0-9]{1,2}`)
					if dateRegex.MatchString(v) {
						return fmt.Errorf("日期格式错误，不应为%v，应为 2000-01-01", v)
					}
					return nil
				},
				Destination: &date,
			},
			//&cli.StringFlag{
			//	Name:     "Format",
			//	Usage:    "",
			//	Required: false,
			//	Aliases:  nil,
			//	Action:   nil,
			//},
		},
		Action: func(context *cli.Context) error {
			var list *AllList
			var err error
			switch billMode {
			case MonthBill:
				list, err = GetMonthGiftList(date)
				return err
			default:
			case DailyBill:
				list, err = GetDailyGiftList(date)
				return err
			}

			return WriteCSV(billMode, list)
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
