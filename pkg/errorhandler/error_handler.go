package errorhandler

import "github.com/gofiber/fiber/v2"

func Handle(c *fiber.Ctx, err error) error {
	if fiberErr, ok := err.(*fiber.Error); ok {
		return c.Status(fiberErr.Code).JSON(fiber.Map{
			"status":  fiberErr.Code,
			"message": fiberErr.Message,
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status":  fiber.StatusInternalServerError,
		"message": "Internal Server Error",
	})
}