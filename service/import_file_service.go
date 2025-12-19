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
	case "LEF":
		go LenexImport(r)
		return nil
	case "PDF":
		switch r.FileType {
		case "START_LIST":
			go PdfStartListImport(r)
		case "RESULT_LIST":
			go PdfResultListImport(r)
		default:
			return fmt.Errorf("unknown file_type for PDF (%s)", r.FileType)
		}
		return nil
	case "PDF_TXT":
		switch r.FileType {
		case "START_LIST":
			go PdfTxtStartListImport(r)
		case "RESULT_LIST":
			go PdfTxtResultListImport(r)
		default:
			return fmt.Errorf("unknown file_type for PDF_TXT (%s)", r.FileType)
		}
		return nil
	default:
		return fmt.Errorf("unknown file extension (%s)", r.FileExtension)
	}
}

func DsvDefinitionImport(r model.ImportFileRequest) {
	if r.SessionID != "" {
		SendLog(r.SessionID, "Starting DSV definition import...", "info")
		SendProgress(r.SessionID, 10, "Initializing import")
	}

	stats, err := importer.ImportDsvDefinitionFile(r.Url, r.Meeting, r.ExcludeEvents, r.IncludeEvents)
	if err != nil {
		println(err.Error())
		if r.SessionID != "" {
			SendError(r.SessionID, err)
		}
		return
	}

	stats.PrintReport()

	if r.SessionID != "" {
		SendProgress(r.SessionID, 100, "Import completed")
		SendLog(r.SessionID, "DSV definition import finished successfully", "success")
		SendComplete(r.SessionID)
	}
}

func DsvResultListImport(r model.ImportFileRequest) {
	if r.SessionID != "" {
		SendLog(r.SessionID, "Starting DSV result list import...", "info")
		SendProgress(r.SessionID, 10, "Initializing import")
	}

	stats, err := importer.ImportDsvResultFile(r.Url, r.Meeting, r.ExcludeEvents, r.IncludeEvents)
	if err != nil {
		println(err.Error())
		if r.SessionID != "" {
			SendError(r.SessionID, err)
		}
		return
	}

	stats.PrintReport()

	if r.SessionID != "" {
		SendProgress(r.SessionID, 100, "Import completed")
		SendLog(r.SessionID, "DSV result list import finished successfully", "success")
		SendComplete(r.SessionID)
	}
}

func LenexImport(r model.ImportFileRequest) {
	if r.SessionID != "" {
		SendLog(r.SessionID, "Starting Lenex import...", "info")
		SendProgress(r.SessionID, 5, "Fetching import settings")
	}

	settings, err := GetImportSettingByMeeting(r.Meeting)
	if err != nil {
		println(err.Error())
		if r.SessionID != "" {
			SendError(r.SessionID, err)
		}
		return
	}

	if r.SessionID != "" {
		SendProgress(r.SessionID, 15, "Processing Lenex file")
	}

	stats, err := importer.ImportLenexFile(r.Url, r.Meeting, r.ExcludeEvents, r.IncludeEvents, settings)
	if err != nil {
		println(err.Error())
		if r.SessionID != "" {
			SendError(r.SessionID, err)
		}
		return
	}

	stats.PrintReport()

	if r.SessionID != "" {
		SendProgress(r.SessionID, 100, "Import completed")
		SendLog(r.SessionID, "Lenex import finished successfully", "success")
		SendComplete(r.SessionID)
	}
}

func PdfStartListImport(r model.ImportFileRequest) {
	if r.SessionID != "" {
		SendLog(r.SessionID, "Starting PDF start list import...", "info")
		SendProgress(r.SessionID, 5, "Fetching import settings")
	}

	settings, err := GetImportSettingByMeeting(r.Meeting)
	if err != nil {
		println(err.Error())
		if r.SessionID != "" {
			SendError(r.SessionID, err)
		}
		return
	}

	if r.SessionID != "" {
		SendProgress(r.SessionID, 15, "Processing PDF file")
	}

	stats, err := importer.ImportPdfStartListFile(r.Url, r.Meeting, r.ExcludeEvents, r.IncludeEvents, settings.PdfStartListSettings)
	if err != nil {
		println(err.Error())
		if r.SessionID != "" {
			SendError(r.SessionID, err)
		}
		return
	}

	stats.PrintReport()

	if r.SessionID != "" {
		SendProgress(r.SessionID, 100, "Import completed")
		SendLog(r.SessionID, "PDF start list import finished successfully", "success")
		SendComplete(r.SessionID)
	}
}

func PdfResultListImport(r model.ImportFileRequest) {
	if r.SessionID != "" {
		SendLog(r.SessionID, "Starting PDF result list import...", "info")
		SendProgress(r.SessionID, 5, "Fetching import settings")
	}

	settings, err := GetImportSettingByMeeting(r.Meeting)
	if err != nil {
		println(err.Error())
		if r.SessionID != "" {
			SendError(r.SessionID, err)
		}
		return
	}

	if r.SessionID != "" {
		SendProgress(r.SessionID, 15, "Processing PDF file")
	}

	stats, err := importer.ImportPdfResultListFile(r.Url, r.Meeting, r.ExcludeEvents, r.IncludeEvents, settings.PdfResultListSettings)
	if err != nil {
		println(err.Error())
		if r.SessionID != "" {
			SendError(r.SessionID, err)
		}
		return
	}

	stats.PrintReport()

	if r.SessionID != "" {
		SendProgress(r.SessionID, 100, "Import completed")
		SendLog(r.SessionID, "PDF result list import finished successfully", "success")
		SendComplete(r.SessionID)
	}
}

func PdfTxtStartListImport(r model.ImportFileRequest) {
	if r.SessionID != "" {
		SendLog(r.SessionID, "Starting PDF text start list import...", "info")
		SendProgress(r.SessionID, 5, "Fetching import settings")
	}

	settings, err := GetImportSettingByMeeting(r.Meeting)
	if err != nil {
		println(err.Error())
		if r.SessionID != "" {
			SendError(r.SessionID, err)
		}
		return
	}

	if r.SessionID != "" {
		SendProgress(r.SessionID, 15, "Processing text data")
	}

	stats, err := importer.ImportPdfStartList(r.Text, r.Meeting, r.ExcludeEvents, r.IncludeEvents, settings.PdfStartListSettings)
	if err != nil {
		println(err.Error())
		if r.SessionID != "" {
			SendError(r.SessionID, err)
		}
		return
	}

	stats.PrintReport()

	if r.SessionID != "" {
		SendProgress(r.SessionID, 100, "Import completed")
		SendLog(r.SessionID, "PDF text start list import finished successfully", "success")
		SendComplete(r.SessionID)
	}
}

func PdfTxtResultListImport(r model.ImportFileRequest) {
	if r.SessionID != "" {
		SendLog(r.SessionID, "Starting PDF text result list import...", "info")
		SendProgress(r.SessionID, 5, "Fetching import settings")
	}

	settings, err := GetImportSettingByMeeting(r.Meeting)
	if err != nil {
		println(err.Error())
		if r.SessionID != "" {
			SendError(r.SessionID, err)
		}
		return
	}

	if r.SessionID != "" {
		SendProgress(r.SessionID, 15, "Processing text data")
	}

	stats, err := importer.ImportPdfResultList(r.Text, r.Meeting, r.ExcludeEvents, r.IncludeEvents, settings.PdfResultListSettings)
	if err != nil {
		println(err.Error())
		if r.SessionID != "" {
			SendError(r.SessionID, err)
		}
		return
	}

	stats.PrintReport()

	if r.SessionID != "" {
		SendProgress(r.SessionID, 100, "Import completed")
		SendLog(r.SessionID, "PDF text result list import finished successfully", "success")
		SendComplete(r.SessionID)
	}
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
