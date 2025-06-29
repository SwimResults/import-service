package importer

import (
	"github.com/stretchr/testify/assert"
	"github.com/swimresults/import-service/model"
	"os"
	"testing"
)

func TestGetPdfFileContent(t *testing.T) {
	text, _ := GetPdfFileContent("../assets/melde_wk047.pdf")
	assert.Equal(t, "erzeugt mit \"EasyWk vom 16.05.2022\" - www.easywk.de25. Internationaler Erzgebirgs-Schwimmcup 2022Seite 1www.sv-zwickau04.deWettkampf 47 - 200m Schmetterling weiblichLauf 1 (ca. 15:01 Uhr)BahnSchwimmerinJg.VereinMeldezeitBahn 2 Kreißig, Lilly  2006 SV 1990 Zschopau 03:04,40Bahn 3 Epperlein, Linda  2008 SV Zwickau von 1904 03:18,54Lauf 2 (ca. 15:04 Uhr)BahnSchwimmerinJg.VereinMeldezeitBahn 1 Maraskova, Linda  2009 Plavecký klub Litvínov 03:02,55Bahn 2 Kunz, Jenny  2007 SV Zwickau von 1904 02:50,99Bahn 3 Richter, Nele  2002 ST Erzgebirge 02:57,21", text)
}

func TestGetPdfFileContent_External(t *testing.T) {
	text, _ := GetPdfFileContent("https://bsv-sws.de/images/dateien/2025/BML/BM_lang_2025.pdf")
	println(text)
	assert.Contains(t, text, "Offene Bezirksjahrgangs- und Bezirksmeisterschaften 2025 - lange Strecken")
}

func TestImportPdfStartList_BMSWS24(t *testing.T) {
	settings := model.ImportPdfStartListSettings{
		OmitFirst:            nil,
		EventSeparator:       "Wettkampf ",
		EventSkipStrings:     nil,
		EventRequiredStrings: []string{"Meldezeit"},
		EventNumberSeparator: "-",
		DistanceSeparator:    "m",
		GenderMapping: [3][2]string{
			{"männlich", "MALE"},
			{"weiblich", "FEMALE"},
			{"mixed", "MIXED"},
		},
		StyleNameSkipStrings:   []string{"staffel", "Staffel", "beine", "Beine"},
		HeatSeparator:          "Lauf",
		HeatSkipStrings:        []string{"Finalabschnittes"},
		HeatRequiredStrings:    []string{"Uhr", "Meldezeit"},
		HeatNumberSeparator:    "/",
		HeatTimeLeftSeparator:  "ca.",
		HeatTimeRightSeparator: "Uhr",
		HeatTimeLayout:         "15:04",
		LaneSeparator:          "Bahn",
		LaneSkipStrings:        []string{"Uhr", "Meldezeit"},
		LaneNumberPattern:      "[0-9]+",
		YearPattern:            "[0-9]{4}",
		YearOpenString:         "Offen",
		SwimTimePattern:        "[0-9]{2}:[0-9]{2},[0-9]{2}",
	}

	err := os.Setenv("SR_NO_IMPORT", "true")
	assert.NoError(t, err)

	stats, err1 := ImportPdfStartListFile("../assets/ME_BMSWS24.pdf", "IESC23", nil, nil, settings)
	assert.NoError(t, err1)
	stats.PrintReport()

	assert.Equal(t, 73, stats.Found.Events)
	assert.Equal(t, 0, stats.Found.AgeGroups)
	assert.Equal(t, 1836, stats.Found.Teams)
	assert.Equal(t, 1830, stats.Found.Athletes)
	assert.Equal(t, 244, stats.Found.Heats)
	assert.Equal(t, 1836, stats.Found.Starts)
	assert.Equal(t, 1836, stats.Found.Results)
	assert.Equal(t, 0, stats.Found.Disqualifications)
}

func TestImportPdfStartList_IESC23(t *testing.T) {
	settings := model.ImportPdfStartListSettings{
		OmitFirst:            nil,
		EventSeparator:       "Wettkampf ",
		EventSkipStrings:     nil,
		EventRequiredStrings: []string{"Meldezeit"},
		EventNumberSeparator: "-",
		DistanceSeparator:    "m",
		GenderMapping: [3][2]string{
			{"männlich", "MALE"},
			{"weiblich", "FEMALE"},
			{"mixed", "MIXED"},
		},
		StyleNameSkipStrings:   []string{"staffel", "Staffel", "beine", "Beine"},
		HeatSeparator:          "Lauf",
		HeatSkipStrings:        []string{"Finalabschnittes"},
		HeatRequiredStrings:    []string{"Uhr", "Meldezeit"},
		HeatNumberSeparator:    "/",
		HeatTimeLeftSeparator:  "ca.",
		HeatTimeRightSeparator: "Uhr",
		HeatTimeLayout:         "15:04",
		LaneSeparator:          "Bahn",
		LaneSkipStrings:        []string{"Uhr", "Meldezeit"},
		LaneNumberPattern:      "[0-9]+",
		YearPattern:            "[0-9]{4}",
		YearOpenString:         "Offen",
		SwimTimePattern:        "[0-9]{2}:[0-9]{2},[0-9]{2}",
	}

	err := os.Setenv("SR_NO_IMPORT", "true")
	assert.NoError(t, err)

	stats, err1 := ImportPdfStartListFile("../assets/ME_26_IESC_2023.pdf", "IESC23", nil, nil, settings)
	assert.NoError(t, err1)
	stats.PrintReport()

	assert.Equal(t, 112, stats.Found.Events)
	assert.Equal(t, 0, stats.Found.AgeGroups)
	assert.Equal(t, 1813, stats.Found.Teams)
	assert.Equal(t, 1762, stats.Found.Athletes)
	assert.Equal(t, 479, stats.Found.Heats)
	assert.Equal(t, 1813, stats.Found.Starts)
	assert.Equal(t, 1813, stats.Found.Results)
	assert.Equal(t, 0, stats.Found.Disqualifications)
}

func TestSwimTimeToDuration(t *testing.T) {
	timeString := "01:02,64"
	println(SwimTimeToDuration(timeString))
}
