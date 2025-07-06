package dto

type ImportCertificatesRequestDto struct {
	Directory string `json:"directory"`
	Meeting   string `json:"meeting"`
}
