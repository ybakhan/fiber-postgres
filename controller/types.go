package controller

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BookController interface {
	CreateBook(ctx *fiber.Ctx) error
	DeleteBook(ctx *fiber.Ctx) error
	GetBook(ctx *fiber.Ctx) error
	GetBooks(ctx *fiber.Ctx) error
}

type bookController struct {
	DB    *gorm.DB
	Cache *memcache.Client
}

type Book struct {
	ID        uint   `gorm:"primary key; autoIncrement" json:"id,omitempty"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}
