package importer

import "testing"

func TestReadPdf(t *testing.T) {
	str, err := ReadPdf("../assets/ME_KKJS.pdf")
	if err != nil {
		panic(err)
	}
	println(str)
}
