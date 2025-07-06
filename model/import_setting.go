package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ImportSetting struct {
	Identifier            primitive.ObjectID          `json:"_id,omitempty" bson:"_id,omitempty"`
	Meeting               string                      `json:"meeting,omitempty" bson:"meeting,omitempty"`
	PdfStartListSettings  ImportPdfStartListSettings  `json:"pdf_start_list_settings,omitempty" bson:"pdf_start_list_settings,omitempty"`
	PdfResultListSettings ImportPdfResultListSettings `json:"pdf_result_list_settings,omitempty" bson:"pdf_result_list_settings,omitempty"`
	CertificateSettings   ImportCertificateSettings   `json:"certificate_settings,omitempty" bson:"certificate_settings,omitempty"`
}

type ImportCertificateSettings struct {
	AiSystemPrompt string `json:"ai_system_prompt,omitempty" bson:"ai_system_prompt,omitempty"` //Aus Urkundentext eines Schwimmwettkampfs einen Dateinamen erzeugen: Strecke + Stil + relevante Zusätze (z.B. Vorlauf, Finale, Masters, Punktbeste Leistung), ohne Name, Verein, Veranstaltung, Platz, Jahrgang, Ort, Datum. Zusätze nur, wenn es klare infos dazu gibt! Ausgabe nur der Name!
}
