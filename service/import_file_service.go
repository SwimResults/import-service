package service

import (
	"errors"
	"github.com/swimresults/import-service/importer"
	"github.com/swimresults/import-service/model"
)

func ImportFile(r model.ImportFileRequest) (*model.ImportFileStats, error) {
	switch r.FileExtension {
	case "DSV":
		return ImportDsvFile(r)
	case "PDF":
		return nil, nil
	}
	return nil, errors.New("unknown file extension")
}

func ImportDsvFile(r model.ImportFileRequest) (*model.ImportFileStats, error) {
	switch r.FileType {
	case "DEFINITION":
		return importer.ImportDsvDefinitionFile(r.Url, r.Meeting, r.ExcludeEvents, r.IncludeEvents)
	case "RESULT_LIST":
		return importer.ImportDsvResultFile(r.Url, r.Meeting, r.ExcludeEvents, r.IncludeEvents)
	}
	return nil, errors.New("unknown file_type for DSV")
}
