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
			go func() {
				stats, err := importer.ImportDsvDefinitionFile(r.Url, r.Meeting, r.ExcludeEvents, r.IncludeEvents)
				if err != nil {
					panic(err)
				}
				stats.PrintReport()
			}()
		case "RESULT_LIST":
			go func() {
				stats, err := importer.ImportDsvResultFile(r.Url, r.Meeting, r.ExcludeEvents, r.IncludeEvents)
				if err != nil {
					panic(err)
				}
				stats.PrintReport()
			}()
		default:
			return fmt.Errorf("unknown file_type for DSV (%s)", r.FileType)
		}
		return nil
	case "PDF":
		switch r.FileType {
		case "START_LIST":
			go func() {
				stats, err := importer.ImportPdfStartListFile(r.Url, r.Meeting, r.ExcludeEvents, r.IncludeEvents)
				if err != nil {
					panic(err)
				}
				stats.PrintReport()
			}()
		case "RESULT_LIST":
			go func() {
				stats, err := importer.ImportDsvResultFile(r.Url, r.Meeting, r.ExcludeEvents, r.IncludeEvents)
				if err != nil {
					panic(err)
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
