package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ImportSetting struct {
	Identifier           primitive.ObjectID         `json:"_id,omitempty" bson:"_id,omitempty"`
	Meeting              string                     `json:"meeting,omitempty" bson:"meeting,omitempty"`
	PdfStartListSettings ImportPdfStartListSettings `json:"pdf_start_list_settings,omitempty" bson:"pdf_start_list_settings,omitempty"`
}
