package service

import (
	"fmt"
	"github.com/swimresults/import-service/importer"
	"github.com/swimresults/import-service/model"
)

func ImportFile(r model.ImportFileRequest) error {
	switch r.FileExtension {
	case "DSV":
		switch r.FileType {
		case "DEFINITION":
			go DsvDefinitionImport(r)
		case "RESULT_LIST":
			go DsvResultListImport(r)
		default:
			return fmt.Errorf("unknown file_type for DSV (%s)", r.FileType)
		}
		return nil
	case "PDF":
		switch r.FileType {
		case "START_LIST":
			go PdfStartListImport(r)
		case "RESULT_LIST":
			go func() {
				stats, err := importer.ImportDsvResultFile(r.Url, r.Meeting, r.ExcludeEvents, r.IncludeEvents)
				if err != nil {
					println(err.Error())
				}
				stats.PrintReport()
			}()
		default:
			return fmt.Errorf("unknown file_type for PDF (%s)", r.FileType)
		}
		return nil
	default:
		return fmt.Errorf("unknown file extension (%s)", r.FileExtension)
	}
}

func DsvDefinitionImport(r model.ImportFileRequest) {
	stats, err := importer.ImportDsvDefinitionFile(r.Url, r.Meeting, r.ExcludeEvents, r.IncludeEvents)
	if err != nil {
		println(err.Error())
	}
	stats.PrintReport()
}

func DsvResultListImport(r model.ImportFileRequest) {
	stats, err := importer.ImportDsvResultFile(r.Url, r.Meeting, r.ExcludeEvents, r.IncludeEvents)
	if err != nil {
		println(err.Error())
	}
	stats.PrintReport()
}

func PdfStartListImport(r model.ImportFileRequest) {
	println("kappa")
	settings, err := GetImportSettingByMeeting(r.Meeting)
	println("kappa")
	if err != nil {
		println(err.Error())
		return
	}
	stats, err := importer.ImportPdfStartList(r.Url, r.Meeting, r.ExcludeEvents, r.IncludeEvents, settings.PdfStartListSettings)
	if err != nil {
		println(err.Error())
	}
	stats.PrintReport()
}
