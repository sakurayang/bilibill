package csv

import "C"
import (
	"fmt"
	"github.com/sakurayang/bilibill/biliapi"
	"github.com/sakurayang/bilibill/format"
)

func WriteCSV(mode format.BillMode, list *biliapi.AllList) error {
	formatList := *list.CSV()
	if !C.AppConfig.Debug {
		s := fmt.Sprintf("\"总计\",,,,,,=SUM(G2:G%d),\"=COUNTIF(H2:H%d,TRUE)\",=SUM(I2:I%d)\n",
			len(*list.List)+1, len(*list.List)+1, len(*list.List)+1)
		formatList = append(formatList, s)
	}
	return format.WriteFile(mode, format.CSV, &formatList, (*list.List)[0].Time)
}
