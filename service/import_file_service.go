package service

import (
	"fmt"
	"github.com/swimresults/import-service/importer"
	"github.com/swimresults/import-service/model"
)

func ImportFile(r model.ImportFileRequest) error {
	printImportInfo(r)
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
	settings, err := GetImportSettingByMeeting(r.Meeting)
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

func printImportInfo(r model.ImportFileRequest) {
	fmt.Println()
	fmt.Println()
	fmt.Printf("\t+----======================================----+\n")
	fmt.Printf("\t|       \033[36mFILE IMPORT GO ROUTINE STARTED!\033[0m        |\n")
	fmt.Printf("\t+----======================================----+\n")
	fmt.Printf("\n")
	fmt.Printf("\t\033[37mFile: \033[36m%s\033[0m\n", r.Url)
	fmt.Printf("\t\033[37mExtension: \033[36m%s\033[0m\n", r.FileExtension)
	fmt.Printf("\t\033[37mType: \033[36m%s\033[0m\n", r.FileType)
	fmt.Printf("\t\033[37mInclude: \033[36m%d\033[0m\n", r.IncludeEvents)
	fmt.Printf("\t\033[37mExclude: \033[36m%d\033[0m\n", r.ExcludeEvents)
	fmt.Println()
	fmt.Println()
}
