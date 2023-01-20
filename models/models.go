package models

import "gorm.io/gorm"

type LogicalExpression struct {
    gorm.Model
    Expression  string `json:"expression" gorm:"text;not null;default:null`
}