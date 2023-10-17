package importer

import "testing"

func TestDoTheMagic(t *testing.T) {
	//err := DoTheMagic()
	//if err != nil {
	//	panic(err)
	//}
}

func TestDoTheResultsMagic(t *testing.T) {
	err := DoTheResultsMagic()
	if err != nil {
		panic(err)
	}
}

func TestImportDsvResultFile(t *testing.T) {
	stats, err := ImportDsvResultFile("../assets/Ergebnisdatei.dsv6", "IESC19", nil, nil)

	stats.PrintReport()

	if err != nil {
		panic(err)
	}
}

func TestImportDsvDefinitionFile(t *testing.T) {
	stats, err := ImportDsvDefinitionFile("../assets/2023-12-10-Marienbe-Wk.dsv7", "IESC23", nil, nil)

	stats.PrintReport()

	if err != nil {
		panic(err)
	}
}

func TestReadPdf(t *testing.T) {
	str, err := ReadPdf("../assets/ME_KKJS.pdf")
	if err != nil {
		panic(err)
	}
	println(str)
}
