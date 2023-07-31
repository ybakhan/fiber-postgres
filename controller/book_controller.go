package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewBookController(DB *gorm.DB, cache *memcache.Client) BookController {
	return &bookController{DB, cache}
}

func (c *bookController) CreateBook(ctx *fiber.Ctx) error {
	fmt.Println("Create Book")

	book := &Book{}
	if err := ctx.BodyParser(book); err != nil {
		ctx.Status(http.StatusUnprocessableEntity).
			JSON(fiber.Map{"message": "create book failed"})
		return err
	}

	if err := c.DB.Create(book).Error; err != nil {
		ctx.Status(http.StatusInternalServerError).
			JSON(fiber.Map{"message": "create book failed"})
		return err
	}

	ctx.Status(http.StatusCreated).JSON(fiber.Map{"message": "created book"})

	go func() {
		bookStr, err := json.Marshal(book)
		if err != nil {
			fmt.Println(err)
			return
		}

		cacheItem := &memcache.Item{
			Key:        strconv.Itoa(int(book.ID)),
			Value:      bookStr,
			Expiration: 36000,
		}

		if err := c.Cache.Set(cacheItem); err != nil {
			fmt.Println(err)
		}
	}()

	return nil
}

func (c *bookController) DeleteBook(ctx *fiber.Ctx) error {
	bookModel := Book{}
	id := ctx.Params("id")
	if id == "" {
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	if err := c.DB.Delete(bookModel, id).Error; err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book delete successfully",
	})
	return nil
}

func (c *bookController) GetBook(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	// Retrieve the struct model from Memcached
	bookModel := &Book{}
	if cacheItem, err := c.Cache.Get(id); err == nil {
		if err = json.Unmarshal(cacheItem.Value, &bookModel); err == nil {
			ctx.Status(http.StatusOK).JSON(&fiber.Map{
				"message": "book id fetched successfully from cache",
				"data":    bookModel,
			})
			return nil
		}
	}

	if err := c.DB.Where("id = ?", id).First(bookModel).Error; err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not get the book"})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book id fetched successfully",
		"data":    bookModel,
	})
	return nil
}

func (c *bookController) GetBooks(ctx *fiber.Ctx) error {
	books := &[]Book{}
	if err := c.DB.Find(books).Error; err != nil {
		ctx.Status(http.StatusInternalServerError).
			JSON(fiber.Map{"message": "get books failed"})
		return err
	}

	ctx.Status(http.StatusOK).JSON(fiber.Map{"books": books})
	return nil
}
