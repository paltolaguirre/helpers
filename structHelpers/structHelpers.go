package structhelpers

import (
	"github.com/jinzhu/gorm"
)

type helper struct {
	gorm.Model
	Nombre      string `json:"nombre"`
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
	Activo      int    `json:"activo"`
}

func (helper) TableName() string {
	return "provincia"
}
