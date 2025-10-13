package middleware

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func Upload() fiber.Handler {
	return func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		fmt.Println("form", form.File)
		if err != nil || form == nil || form.File == nil || form.File["images"] == nil {
			c.Locals("filenames", []string{})
			return c.Next()
		}

		files := form.File["images"]
		uploadDir := "uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			if err := os.Mkdir(uploadDir, 0755); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "failed to create upload directory",
				})
			}
		}

		var filenames []string
		for _, file := range files {
			filename := utils.UUIDv4() + filepath.Ext(file.Filename)
			filePath := filepath.Join(uploadDir, filename)
			if err := c.SaveFile(file, filePath); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "failed to save file",
				})
			}
			filenames = append(filenames, filename)
		}

		c.Locals("filenames", filenames)
		return c.Next()
	}
}
