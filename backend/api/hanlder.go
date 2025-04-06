package api

import "gorm.io/gorm"

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	h := &Handler{
		db: db,
	}

	return h
}
