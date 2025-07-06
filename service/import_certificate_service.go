package service

import (
	"fmt"
	"github.com/swimresults/import-service/dto"
	"github.com/swimresults/import-service/importer"
)

func ImportCertificates(r dto.ImportCertificatesRequestDto) error {
	printCertImportInfo(r)
	go RunCertificatesImport(r)
	return nil
}

func RunCertificatesImport(r dto.ImportCertificatesRequestDto) {
	settings, err := GetImportSettingByMeeting(r.Meeting)
	if err != nil {
		println(err.Error())
		return
	}
	_, err = importer.ImportCertificates(r.Directory, r.Meeting, settings.CertificateSettings)
	if err != nil {
		println(err.Error())
	}
}

func printCertImportInfo(r dto.ImportCertificatesRequestDto) {
	fmt.Println()
	fmt.Println()
	fmt.Printf("\t+----======================================----+\n")
	fmt.Printf("\t|       \033[36mCERT IMPORT GO ROUTINE STARTED!\033[0m        |\n")
	fmt.Printf("\t+----======================================----+\n")
	fmt.Printf("\n")
	fmt.Printf("\t\033[37mDirectory: \033[36m%s\033[0m\n", r.Directory)
	fmt.Printf("\t\033[37mMeeting: \033[36m%s\033[0m\n", r.Meeting)
	fmt.Println()
	fmt.Println()
}
