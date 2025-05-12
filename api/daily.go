package api

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/robfig/cron/v3"
)

func (server *Server) scheduleDailyAnalytics() {
	// Create a new cron scheduler (uses the local time zone by default)
	c := cron.New(cron.WithLocation(Loc))

	// Schedule the job to run every day at midnight.
	// "@daily" is equivalent to "0 0 0 * * *"
	var err error
	_, err = c.AddFunc("@daily", func() {
		// Use time.Now() to set the current day at midnight

		if _, err := server.store.CreateDailyAnalytics(context.Background(), pgtype.Date{Time: Date, Valid: true}); err != nil {
			log.Printf("Failed to create daily analytics: %v\n", err)
		} else {
			log.Println("Daily analytics created successfully at midnight.")
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling daily analytics: %v\n", err)
	}

	// Start the cron scheduler in its own goroutine
	c.Start()
}

func (server *Server) deactivateAds() {
	// Create a new cron scheduler (uses the local time zone by default)
	c := cron.New(cron.WithLocation(Loc))

	// Schedule the job to run every day at midnight.
	// "@daily" is equivalent to "0 0 0 * * *"
	var err error
	_, err = c.AddFunc("1 0 * * *", func() {
		// Get the current time
		now := time.Now().In(Loc)
		log.Print("deactivateAds job is running..")

		// Create a context with timeout for database operations
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// List all active ads
		log.Print("Listing active ads...")
		ads, err := server.store.ListAds(ctx, 1000)
		if err != nil {
			log.Printf("Failed to list ads: %v\n", err)
			return
		}

		log.Print("Deactivating expired ads...")
		for _, ad := range ads {
			if ad.EndDate.Valid {
				nowYear, nowMonth, nowDay := now.Date()
				endYear, endMonth, endDay := ad.EndDate.Time.In(Loc).Date() // convert ad end date to local time
				todayDate := time.Date(nowYear, nowMonth, nowDay, 0, 0, 0, 0, Loc)
				adEndDate := time.Date(endYear, endMonth, endDay, 0, 0, 0, 0, Loc)

				// If today's date is the same as or after the ad's end date, deactivate the ad
				if !todayDate.Before(adEndDate) {
					_, err := server.store.DeactivateAd(ctx, ad.ID)
					if err != nil {
						log.Printf("Failed to deactivate ad %v: %v\n", ad.ID, err)
					} else {
						log.Printf("Successfully deactivated expired ad %v\n", ad.ID)
					}
				}
			}
		}
	})

	if err != nil {
		log.Fatalf("Error setting up cron job for deactivating expired ads: %v\n", err)
	}

	// Start the cron scheduler in its own goroutine
	c.Start()
}
