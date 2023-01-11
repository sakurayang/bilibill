package json

import (
	"github.com/sakurayang/bilibill/biliapi"
	"github.com/sakurayang/bilibill/format"
)

func WriteJSON(mode format.BillMode, list *biliapi.AllList) error {
	formatList := (*list).JSON()
	return format.WriteFile(mode, format.JSON, formatList, (*list.List)[0].Time)
}
