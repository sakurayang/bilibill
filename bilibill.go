package main

import (
	"fmt"
	"github.com/sakurayang/bilibill/biliapi"
	"github.com/sakurayang/bilibill/config"
	"github.com/sakurayang/bilibill/format"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	billMode := format.DailyBill
	date := "1970-01-01 00:00:00"
	var outputPath, configPath cli.Path
	var cookie string
	var oFormat format.OFormat

	app := &cli.App{
		Name:    "biliBill",
		Usage:   "",
		Suggest: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "billmode",
				Usage:       "指定导出模式，day 为日账单，month 为月账单，默认为日账单",
				DefaultText: "day",
				Value:       "day",
				Aliases:     []string{"m"},
				Action: func(ctx *cli.Context, v string) error {
					if v != "day" && v != "month" {
						return fmt.Errorf("导出模式应当为 'day' 或 'month' 而不是 %v", v)
					}
					if v == "month" {
						billMode = format.MonthBill
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
					dateRegex := regexp.MustCompile(
						`([1-9][0-9]{3})-(((0[13578]|1[02])-(0[1-9]|[12][0-9]|3[01]))|((0[469]|11)-(0[1-9]|[12][0-9]|30))|(02-(0[1-9]|1[0-9]|2[0-9])))`,
					)
					if !dateRegex.MatchString(v) {
						return fmt.Errorf("日期格式错误，不应为%v，应为 2000-01-01", v)
					}
					return nil
				},
				Destination: &date,
			},
			&cli.PathFlag{
				Name:        "output",
				Usage:       "指定输出文件",
				Aliases:     []string{"o"},
				Destination: &outputPath,
			},
			&cli.PathFlag{
				Name:        "config",
				Usage:       "指定配置文件",
				Value:       "config.yaml",
				DefaultText: "config.yaml",
				HasBeenSet:  true,
				Destination: &configPath,
				Aliases:     []string{"C"},
				Action: func(context *cli.Context, path cli.Path) error {
					if !strings.HasSuffix(path, "yaml") {
						return fmt.Errorf("请输入 yaml 文件")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:        "cookie",
				Usage:       "指定需要使用的 cookie",
				Required:    true,
				Destination: &cookie,
				Action: func(context *cli.Context, s string) error {
					sl := strings.ToUpper(s)
					if !(strings.HasPrefix(sl, "SESSDATA") && strings.HasSuffix(s, ";")) {
						return fmt.Errorf("cookie 格式应当为 SESSDATA=xxxxxxxxx; 而不是 %s", s)
					}
					if len(strings.Split(s, "=")[1]) < 35 {
						return fmt.Errorf("cookie 长度错误")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:        "Format",
				Usage:       "指定输出的格式",
				DefaultText: "CSV",
				Value:       "CSV",
				Aliases:     []string{"f"},
				Action: func(context *cli.Context, s string) error {
					switch strings.ToUpper(s) {
					case "JSON":
						oFormat = format.JSON
						break
					case "CSV":
						oFormat = format.CSV
						break
					default:
						return fmt.Errorf("输出格式为 csv 或 json，%s 为未适配格式", s)
					}
					return nil
				},
			},
		},
		Action: func(context *cli.Context) error {
			config.C = config.GetConfig(configPath)
			var list *biliapi.AllList
			var err error
			switch billMode {
			case format.MonthBill:
				list, err = biliapi.GetMonthGiftList(date)
				if err != nil {
					return err
				}
				break
			default:
			case format.DailyBill:
				list, err = biliapi.GetDailyGiftList(date)
				if err != nil {
					return err
				}
				break
			}

			var outstr *biliapi.OutStrings
			switch oFormat {
			case format.JSON:
				outstr = list.JSON()
				break
			default:
			case format.CSV:
				outstr = list.CSV()
				break
			}
			if len(outputPath) == 0 {
				fmt.Println(outstr)
				return nil
			}
			return format.WriteFile(billMode, oFormat, outstr, (*list.List)[0].Time)
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
