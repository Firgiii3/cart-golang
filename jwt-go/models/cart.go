package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Qty    int    `json:"qty"`
}
