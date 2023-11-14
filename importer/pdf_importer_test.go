package importer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPdfFileContent(t *testing.T) {
	text, _ := GetPdfFileContent("../assets/melde_wk047.pdf")
	assert.Equal(t, "erzeugt mit \"EasyWk vom 16.05.2022\" - www.easywk.de25. Internationaler Erzgebirgs-Schwimmcup 2022Seite 1www.sv-zwickau04.deWettkampf 47 - 200m Schmetterling weiblichLauf 1 (ca. 15:01 Uhr)BahnSchwimmerinJg.VereinMeldezeitBahn 2 Kreißig, Lilly  2006 SV 1990 Zschopau 03:04,40Bahn 3 Epperlein, Linda  2008 SV Zwickau von 1904 03:18,54Lauf 2 (ca. 15:04 Uhr)BahnSchwimmerinJg.VereinMeldezeitBahn 1 Maraskova, Linda  2009 Plavecký klub Litvínov 03:02,55Bahn 2 Kunz, Jenny  2007 SV Zwickau von 1904 02:50,99Bahn 3 Richter, Nele  2002 ST Erzgebirge 02:57,21", text)
}
