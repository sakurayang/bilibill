package main

import "fmt"

var c = GetConfig()

func main() {
	glist, err := GetMonthGiftList()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%#v\n", glist)
	err = WriteCSV(MonthBill, glist)
	if err != nil {
		fmt.Println(err.Error())
	}
}
