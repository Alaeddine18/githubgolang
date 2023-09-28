package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func downloadRepository(c *fiber.Ctx) error {
	zipFilePath := "./reposCloned.zip"

	_, err := os.Stat(zipFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return c.Status(fiber.StatusNotFound).SendString("can't find the zip file")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("api error")
	}

	c.Set("Content-Disposition", "attachment; filename=reposCloned.zip")
	c.Set("Content-Type", "application/zip")

	return c.SendFile(zipFilePath)
}
