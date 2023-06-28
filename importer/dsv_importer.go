package importer

import (
	"bytes"
	"github.com/konrad2002/dsvparser/model"
	"github.com/konrad2002/dsvparser/parser"
	"os"
)

func VeranstaltungsortPlz() string {
	dat, err := os.ReadFile("assets/Ergebnisdatei.dsv6")
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(dat)
	r := parser.NewReader(buf)
	res, err := r.Read()
	if err != nil {
		panic(err)
	}
	def := res.(*model.Wettkampfdefinitionsliste)
	return def.Veranstaltungsort.PLZ
}
