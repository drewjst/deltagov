package api

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(app *fiber.App) {
	// API v1 group
	v1 := app.Group("/api/v1")

	// Bills endpoints
	bills := v1.Group("/bills")
	bills.Get("/", ListBills)
	bills.Get("/:id", GetBill)
	bills.Get("/:id/versions", GetBillVersions)

	// Diff endpoints
	diff := v1.Group("/diff")
	diff.Get("/:versionA/:versionB", ComputeDiff)
}

// ListBills returns a list of tracked bills
func ListBills(c *fiber.Ctx) error {
	// TODO: Implement database query
	return c.JSON(fiber.Map{
		"bills": []fiber.Map{},
		"total": 0,
	})
}

// GetBill returns a single bill by ID
func GetBill(c *fiber.Ctx) error {
	id := c.Params("id")
	// TODO: Implement database query
	return c.JSON(fiber.Map{
		"id":     id,
		"title":  "",
		"status": "not_found",
	})
}

// GetBillVersions returns all versions of a bill
func GetBillVersions(c *fiber.Ctx) error {
	id := c.Params("id")
	// TODO: Implement database query
	return c.JSON(fiber.Map{
		"bill_id":  id,
		"versions": []fiber.Map{},
	})
}

// ComputeDiff computes the diff between two bill versions
func ComputeDiff(c *fiber.Ctx) error {
	versionA := c.Params("versionA")
	versionB := c.Params("versionB")
	// TODO: Implement diff computation using diff_engine
	return c.JSON(fiber.Map{
		"version_a": versionA,
		"version_b": versionB,
		"delta":     nil,
	})
}
