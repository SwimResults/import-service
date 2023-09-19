package service

import (
	"errors"
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
			return errors.New("unknown file_type for DSV")
		}
		return nil
	case "PDF":
		return nil
	default:
		return errors.New("unknown file extension")
	}
}
