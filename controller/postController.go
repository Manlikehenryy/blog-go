package controller

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/manlikehenryy/blog-backend/database"
	"github.com/manlikehenryy/blog-backend/helpers"
	"github.com/manlikehenryy/blog-backend/models"
	"gorm.io/gorm"
)

func CreatePost(c *fiber.Ctx) error {
	var blogPost models.Blog

	if err := c.BodyParser(&blogPost); err != nil {

		fmt.Println("Unable to parse body:", err)

		return helpers.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return helpers.SendError(c, fiber.StatusUnauthorized, "User ID not found in context")
	}

	blogPost.UserId = userId

	if err := database.DB.Create(&blogPost).Error; err != nil {

		fmt.Println("Database error:", err)

		return helpers.SendError(c, fiber.StatusInternalServerError, "Failed to create post")
	}

	return helpers.SendJSON(c, fiber.StatusCreated, fiber.Map{
		"data":    blogPost,
		"message": "Post created successfully",
	})
}

func AllPost(c *fiber.Ctx) error {

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("perPage", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	if page == 0 {
		page = 1
		limit = math.MaxInt
	}

	offset := (page - 1) * limit

	var total int64
	var blogPosts []models.Blog

	if err := database.DB.Preload("User").Offset(offset).Limit(limit).Find(&blogPosts).Error; err != nil {
		return helpers.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve posts")
	}

	if err := database.DB.Model(&models.Blog{}).Count(&total).Error; err != nil {
		return helpers.SendError(c, fiber.StatusInternalServerError, "Failed to count posts")
	}

	pageCount := int(math.Ceil(float64(total) / float64(limit)))
	hasNextPage := page < pageCount
	hasPrevPage := page > 1

	if limit == math.MaxInt {
		limit = int(total)
	}

	return helpers.SendJSON(c, fiber.StatusOK, fiber.Map{
		"data":    blogPosts,
		"message": "Posts fetched successfully",
		"meta": fiber.Map{
			"page":      page,
			"perPage":   limit,
			"total":     total,
			"pageCount": pageCount,
			"nextPage": func() int {
				if hasNextPage {
					return page + 1
				}
				return 0
			}(),
			"hasNextPage": hasNextPage,
			"hasPrevPage": hasPrevPage,
		},
	})
}

func DetailPost(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.SendError(c, fiber.StatusBadRequest, "Invalid post ID")
	}

	var blogPost models.Blog
	if err := database.DB.Where("id = ?", id).Preload("User").First(&blogPost).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.SendError(c, fiber.StatusNotFound, "Post not found")
		}
		return helpers.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve post")
	}

	return helpers.SendJSON(c, fiber.StatusOK, fiber.Map{
		"data": blogPost,
	})
}

func UpdatePost(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helpers.SendError(c, fiber.StatusBadRequest, "Invalid post ID")
	}

	var blog models.Blog
	if err := database.DB.First(&blog, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.SendError(c, fiber.StatusNotFound, "Post not found")
		}
		return helpers.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve post")
	}

	if err := c.BodyParser(&blog); err != nil {
		fmt.Println("Unable to parse body:", err)
		return helpers.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return helpers.SendError(c, fiber.StatusUnauthorized, "User ID not found in context")
	}

	if blog.UserId != userId {
		return helpers.SendError(c, fiber.StatusForbidden, "Unauthorized to update this post")
	}

	if err := database.DB.Model(&blog).Updates(blog).Error; err != nil {
		return helpers.SendError(c, fiber.StatusInternalServerError, "Failed to update post")
	}

	return helpers.SendJSON(c, fiber.StatusOK, fiber.Map{
		"message": "Post updated successfully",
	})
}

func UsersPost(c *fiber.Ctx) error {

	userId, ok := c.Locals("userId").(string)
	if !ok {
		return helpers.SendError(c, fiber.StatusUnauthorized, "User ID not found in context")
	}

	var blogs []models.Blog
	if err := database.DB.Model(&models.Blog{}).
		Where("user_id = ?", userId).
		Preload("User").
		Find(&blogs).Error; err != nil {

		return helpers.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve posts")
	}

	return helpers.SendJSON(c, fiber.StatusOK, fiber.Map{
		"data": blogs,
	})
}

func DeletePost(c *fiber.Ctx) error {

	id, _ := strconv.Atoi(c.Params("id"))

	var blog models.Blog
	if err := database.DB.First(&blog, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.SendError(c, fiber.StatusNotFound, "Post not found")
		}

		return helpers.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve post")
	}

	userId, ok := c.Locals("userId").(string)
	if !ok || userId != blog.UserId {
		return helpers.SendError(c, fiber.StatusForbidden, "Unauthorized to delete this post")
	}

	if err := database.DB.Delete(&blog).Error; err != nil {
		return helpers.SendError(c, fiber.StatusInternalServerError, "Failed to delete post")
	}

	return helpers.SendJSON(c, fiber.StatusOK, fiber.Map{
		"message": "Post deleted successfully",
	})
}
