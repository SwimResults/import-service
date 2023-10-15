package importer

import (
	"bytes"
	"github.com/ledongthuc/pdf"
	importModel "github.com/swimresults/import-service/model"
)

// ReadPdf opens pdf under given path and reads plain text to a buffer and returns content as string
func ReadPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	buf.ReadFrom(b)
	return buf.String(), nil
}

// ImportPdfStartListFile takes the path to a pdf file that contains a start
// list. All events, teams, athletes, heats and starts will be imported.
// If exclude is set, given event numbers will be excluded from import.
// If include is set, only given event numbers will be imported.
//
// For import process details see documentation on GitHub.
func ImportPdfStartListFile(file string, meeting string, exclude []int, include []int) (*importModel.ImportFileStats, error) {
	return nil, nil
}
