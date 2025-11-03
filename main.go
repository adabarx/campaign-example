package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"

	"campaign/templates"
)

type Donation struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Amount    float64   `json:"amount"`
	Message   string    `json:"message,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type AppState struct {
	donations []Donation
	nextID    int
	mu        sync.Mutex
}

var state = &AppState{
	donations: []Donation{},
	nextID:    1,
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--generate" {
		if err := generateStaticSite(); err != nil {
			fmt.Printf("Error generating static site: %v\n", err)
			os.Exit(1)
		}
		return
	}

	app := fiber.New()

	app.Static("/", "./public")

	app.Get("/api/stats", getStats)
	app.Get("/api/recent-donors", getRecentDonors)
	app.Post("/api/donations", createDonation)

	fmt.Println("🚀 Server running on http://localhost:3000")
	fmt.Println("   Static files: public/")
	fmt.Println("   API: GET /api/stats, GET /api/recent-donors, POST /api/donations")

	app.Listen(":3000")
}

func getStats(c *fiber.Ctx) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	total := 0.0
	for _, d := range state.donations {
		total += d.Amount
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return templates.DonationStats(total, len(state.donations)).Render(c.Context(), c.Response().BodyWriter())
}

func getRecentDonors(c *fiber.Ctx) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	donors := make([]templates.Donation, len(state.donations))
	for i, d := range state.donations {
		donors[i] = templates.Donation{
			ID:        d.ID,
			Name:      d.Name,
			Email:     d.Email,
			Amount:    d.Amount,
			Message:   d.Message,
			CreatedAt: d.CreatedAt.Format(time.RFC3339),
		}
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return templates.RecentDonors(donors).Render(c.Context(), c.Response().BodyWriter())
}

func createDonation(c *fiber.Ctx) error {
	donation := new(Donation)

	if err := c.BodyParser(donation); err != nil {
		return c.Status(400).SendString("Invalid form data")
	}

	if donation.Name == "" || donation.Email == "" || donation.Amount <= 0 {
		return c.Status(400).SendString("Name, email, and amount are required")
	}

	state.mu.Lock()
	donation.ID = state.nextID
	donation.CreatedAt = time.Now()
	state.nextID++
	state.donations = append(state.donations, *donation)
	state.mu.Unlock()

	c.Set("Content-Type", "text/html; charset=utf-8")
	c.Set("HX-Trigger", "donationComplete")
	return templates.DonationSuccess(donation.Name, donation.Amount).Render(c.Context(), c.Response().BodyWriter())
}
