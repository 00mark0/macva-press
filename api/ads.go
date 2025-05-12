package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/00mark0/macva-press/components"
	"github.com/00mark0/macva-press/db/services"
	"github.com/00mark0/macva-press/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func (server *Server) listActiveAds(ctx echo.Context) error {
	var req ListAdsReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listActiveAds:", err)
		return err
	}

	nextLimit := req.Limit + 20

	activeAds, err := server.store.ListActiveAds(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing active ads in listActiveAds:", err)
		return err
	}

	url := "/api/admin/ads/active?limit="

	return Render(ctx, http.StatusOK, components.Ads(int(nextLimit), activeAds, url))
}

func (server *Server) listInactiveAds(ctx echo.Context) error {
	var req ListAdsReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listInactiveAds:", err)
		return err
	}

	nextLimit := req.Limit + 20

	inactiveAds, err := server.store.ListInactiveAds(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing inactive ads in listInactiveAds:", err)
		return err
	}

	url := "/api/admin/ads/inactive?limit="

	return Render(ctx, http.StatusOK, components.Ads(int(nextLimit), inactiveAds, url))
}

type CreateAdReq struct {
	Title       string `form:"title" validate:"required,min=3,max=50"`
	Description string `form:"description" validate:"required,min=3,max=100"`
	TargetUrl   string `form:"target_url" validate:"required"`
	Placement   string `form:"placement" validate:"required"`
	Status      string `form:"status" validate:"required"`
	StartDate   string `form:"start_date" validate:"required"`
	EndDate     string `form:"end_date" validate:"required"`
}

func (server *Server) createAd(ctx echo.Context) error {
	var req CreateAdReq
	var createAddErr components.CreateAdErr

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in createAd:", err)
		return err
	}

	if req.TargetUrl != "" {
		// Check if the URL already has a protocol
		if !strings.HasPrefix(strings.ToLower(req.TargetUrl), "http://") &&
			!strings.HasPrefix(strings.ToLower(req.TargetUrl), "https://") {
			// Add https:// if no protocol is present
			req.TargetUrl = "https://" + req.TargetUrl
		}
	}

	if err := ctx.Validate(req); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			switch fieldErr.Field() {
			case "Title":
				createAddErr = "Naslov oglasa mora biti između 3 i 50 karaktera."
			case "Description":
				createAddErr = "Opis oglasa mora biti između 3 i 100 karaktera."
			case "TargetUrl":
				createAddErr = "URL oglasa je obavezan."
			case "Placement":
				createAddErr = "Mesto oglasa je obavezno."
			case "Status":
				createAddErr = "Status oglasa je obavezan."
			case "StartDate":
				createAddErr = "Datum početka oglasa je obavezan."
			case "EndDate":
				createAddErr = "Datum završetka oglasa je obavezan."
			}
		}

		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	file, err := ctx.FormFile("image_url")
	if err != nil {
		createAddErr = "Slika oglasa je obavezna."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	uploadsDir := "static/ads"
	// Ensure directory exists
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			log.Println("Error creating directory in updatePfp:", err)
			return err
		}
	}

	filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
	filePath := fmt.Sprintf("%s/%s", uploadsDir, filename)

	src, err := file.Open()
	if err != nil {
		log.Println("Error opening uploaded file in createAd:", err)
		return err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating destination file in createAd:", err)
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		log.Println("Error copying file in createAd:", err)
		return err
	}

	convertedPath, err := ConvertToWebPWithResize(filePath, 800, 600, 80)
	if err != nil {
		log.Println("Error converting file to WebP in createAd:", err)
	} else {
		filePath = convertedPath
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		log.Println("Error parsing start date in createAd:", err)
		return err
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		log.Println("Error parsing end date in createAd:", err)
		return err
	}

	// 1. Start date must be before end date.
	if startDate.After(endDate) {
		createAddErr = "Datum početka oglasa mora biti pre datuma završetka."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	// 2. Start date must not be in the past (you might consider allowing today by comparing to midnight or adding a small margin).
	now := time.Now().In(Loc)
	midnightNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, Loc)
	if startDate.Before(midnightNow) {
		createAddErr = "Datum početka oglasa mora biti veći od trenutnog datuma."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	// 3. End date must not be in the past.
	if endDate.Before(midnightNow) {
		createAddErr = "Datum završetka oglasa mora biti veći od trenutnog datuma."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	// 4. Start date must not be more than one year in the future.
	maxFutureDate := midnightNow.AddDate(1, 0, 0)
	if startDate.After(maxFutureDate) {
		createAddErr = "Datum početka oglasa ne moze biti više od godinu dana unapred."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	// 5. The duration between start and end must be at least 3 days.
	minDuration := 3 * 24 * time.Hour
	if endDate.Sub(startDate) < minDuration {
		createAddErr = "Razmak između početka i kraja oglasa mora biti barem 3 dana."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	// 6. End date must not be more than 5 years in the future.
	maxEndDate := midnightNow.AddDate(5, 0, 0)
	if endDate.After(maxEndDate) {
		createAddErr = "Datum kraja oglasa ne može biti više od 5 godina u budućnosti."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	scheduledAds, err := server.store.ListScheduledAds(ctx.Request().Context(), 100)
	if err != nil {
		log.Println("Error listing scheduled ads in createAd:", err)
		return err
	}

	for _, ad := range scheduledAds {
		if req.Status == "active" && req.Placement == ad.Placement.String && startDate.After(midnightNow.Add(24*time.Hour)) {
			createAddErr = "Već postoji zakazani oglas za ovu poziciju."
			ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
			return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
		}
	}

	activeAds, err := server.store.ListActiveAds(ctx.Request().Context(), 100)
	if err != nil {
		log.Println("Error listing active ads in createAd:", err)
		return err
	}

	for _, ad := range activeAds {
		if req.Status == "active" && req.Placement == ad.Placement.String && startDate.After(midnightNow) && startDate.Before(midnightNow.Add(24*time.Hour)) {
			createAddErr = "Već postoji aktivan oglas za ovu poziciju."
			ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
			return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
		}

		// Check for overlaps between new scheduled ads and existing active ads
		if req.Status == "active" {
			// If same position and the existing ad's end date is after the new ad's start date
			if req.Placement == ad.Placement.String && ad.EndDate.Time.After(startDate) {
				// Calc duration before adjustment to make up for snipped time
				duration := endDate.Sub(startDate)

				// Adjust start date to midnight of the day after the existing ad ends
				adjustedStartDate := time.Date(
					ad.EndDate.Time.Year(),
					ad.EndDate.Time.Month(),
					ad.EndDate.Time.Day()+1,
					0, 0, 0, 0,
					Loc,
				)

				// Update the start date
				startDate = adjustedStartDate

				// Also update end date to maintain the original duration
				endDate = adjustedStartDate.Add(duration)

				// You might want to inform the user about this adjustment
				ctx.Response().Header().Add("HX-Trigger", "adDatesAdjusted")
				break // Stop after finding the first conflict
			}
		}
	}

	if len(activeAds) >= 4 && startDate.After(midnightNow) && startDate.Before(midnightNow.Add(24*time.Hour)) {
		createAddErr = "Maksimalan broj aktivnih oglasa je 4."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	if len(scheduledAds) >= 4 && startDate.After(midnightNow.Add(24*time.Hour)) {
		createAddErr = "Maksimalan broj zakazanih oglasa je 4."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	arg := db.CreateAdParams{
		Title: pgtype.Text{String: req.Title, Valid: true},
		Description: pgtype.Text{
			String: req.Description,
			Valid:  true,
		},
		ImageUrl: pgtype.Text{
			String: "/" + filePath,
			Valid:  true,
		},
		TargetUrl: pgtype.Text{String: req.TargetUrl, Valid: true},
		Placement: pgtype.Text{String: req.Placement, Valid: true},
		Status:    pgtype.Text{String: req.Status, Valid: true},
		StartDate: pgtype.Timestamptz{
			Time:  startDate,
			Valid: true,
		},
		EndDate: pgtype.Timestamptz{
			Time:  endDate,
			Valid: true,
		},
	}

	_, err = server.store.CreateAd(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error creating ad in createAd:", err)
		return err
	}

	if req.Status == "active" && startDate.After(midnightNow) && startDate.Before(midnightNow.Add(24*time.Hour)) {
		ctx.Response().Header().Add("HX-Trigger", "createAdSuccess")
		return server.activeAdsList(ctx)
	} else if req.Status == "inactive" {
		ctx.Response().Header().Add("HX-Trigger", "createAdSuccess")
		return server.inactiveAdsList(ctx)
	} else {
		ctx.Response().Header().Add("HX-Trigger", "createAdSuccess")
		return server.scheduledAdsList(ctx)
	}
}

// UpdateAdReq - request structure for updating an ad
type UpdateAdReq struct {
	ID          string `form:"id"`
	Title       string `form:"title" validate:"required,min=3,max=50"`
	Description string `form:"description" validate:"required,min=3,max=100"`
	TargetUrl   string `form:"target_url" validate:"required,url"`
	Placement   string `form:"placement" validate:"required,oneof=header sidebar footer article"`
	Status      string `form:"status" validate:"required,oneof=active inactive"`
	StartDate   string `form:"start_date" validate:"required"`
	EndDate     string `form:"end_date" validate:"required"`
}

func (server *Server) updateAd(ctx echo.Context) error {
	var req UpdateAdReq
	var updateAdErr components.UpdateAdErr

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in updateAd:", err)
		return err
	}

	adIDStr := ctx.Param("id")
	adID, err := utils.ParseUUID(adIDStr, "ad ID")
	if err != nil {
		log.Println("Invalid ad ID format in updateAd:", err)
		return err
	}

	existingAd, err := server.store.GetAd(ctx.Request().Context(), adID)
	if err != nil {
		log.Println("Error retrieving ad in updateAd:", err)
		return err
	}

	if req.TargetUrl != "" {
		// Check if the URL already has a protocol
		if !strings.HasPrefix(strings.ToLower(req.TargetUrl), "http://") &&
			!strings.HasPrefix(strings.ToLower(req.TargetUrl), "https://") {
			// Add https:// if no protocol is present
			req.TargetUrl = "https://" + req.TargetUrl
		}
	}

	if err := ctx.Validate(req); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			switch fieldErr.Field() {
			case "Title":
				updateAdErr = "Naslov oglasa mora biti između 3 i 50 karaktera."
			case "Description":
				updateAdErr = "Opis oglasa mora biti između 3 i 100 karaktera."
			case "TargetUrl":
				updateAdErr = "URL oglasa je obavezan."
			case "Placement":
				updateAdErr = "Mesto oglasa je obavezno."
			case "Status":
				updateAdErr = "Status oglasa je obavezan."
			case "StartDate":
				updateAdErr = "Datum početka oglasa je obavezan."
			case "EndDate":
				updateAdErr = "Datum završetka oglasa je obavezan."
			}
		}

		ctx.Response().Header().Set("HX-Retarget", "#update-ad-modal")
		return Render(ctx, http.StatusOK, components.UpdateAdModal(updateAdErr, existingAd))
	}

	// Handle image upload if a new file is provided
	var filePath string
	file, err := ctx.FormFile("image_url")
	if err != nil {
		// No new file uploaded, keep the existing image
		filePath = existingAd.ImageUrl.String
	} else {
		// New file uploaded, process it
		uploadsDir := "static/ads"
		// Ensure directory exists
		if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
			if err := os.MkdirAll(uploadsDir, 0755); err != nil {
				log.Println("Error creating directory in updateAd:", err)
				return err
			}
		}

		filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
		filePath = fmt.Sprintf("%s/%s", uploadsDir, filename)

		src, err := file.Open()
		if err != nil {
			log.Println("Error opening uploaded file in updateAd:", err)
			return err
		}
		defer src.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			log.Println("Error creating destination file in updateAd:", err)
			return err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			log.Println("Error copying file in updateAd:", err)
			return err
		}

		convertedPath, err := ConvertToWebPWithResize(filePath, 800, 600, 80)
		if err != nil {
			log.Println("Error converting file to WebP in createAd:", err)
		} else {
			filePath = convertedPath
		}

		// Delete old image file if it's not the default image
		if existingAd.ImageUrl.String != "" && !strings.Contains(existingAd.ImageUrl.String, "default") {
			if err := os.Remove(existingAd.ImageUrl.String); err != nil {
				log.Println("Warning: could not delete old image file:", err)
				// Continue anyway, this is not a fatal error
			}
		}
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		log.Println("Error parsing start date in updateAd:", err)
		return err
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		log.Println("Error parsing end date in updateAd:", err)
		return err
	}

	// 1. Start date must be before end date.
	if startDate.After(endDate) {
		updateAdErr = "Datum početka oglasa mora biti pre datuma završetka."
		ctx.Response().Header().Set("HX-Retarget", "#update-ad-modal")
		return Render(ctx, http.StatusOK, components.UpdateAdModal(updateAdErr, existingAd))
	}

	// 2. Start date must not be in the past (you might consider allowing today by comparing to midnight or adding a small margin).
	now := time.Now().In(Loc)
	midnightNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, Loc)

	// For updates, allow the start date to be in the past if it's the same as the existing start date
	if startDate.Before(midnightNow) && !startDate.Equal(existingAd.StartDate.Time) {
		updateAdErr = "Datum početka oglasa mora biti veći od trenutnog datuma."
		ctx.Response().Header().Set("HX-Retarget", "#update-ad-modal")
		return Render(ctx, http.StatusOK, components.UpdateAdModal(updateAdErr, existingAd))
	}

	// 3. End date must not be in the past.
	if endDate.Before(midnightNow) {
		updateAdErr = "Datum završetka oglasa mora biti veći od trenutnog datuma."
		ctx.Response().Header().Set("HX-Retarget", "#update-ad-modal")
		return Render(ctx, http.StatusOK, components.UpdateAdModal(updateAdErr, existingAd))
	}

	// 4. Start date must not be more than one year in the future.
	maxFutureDate := midnightNow.AddDate(1, 0, 0)
	if startDate.After(maxFutureDate) {
		updateAdErr = "Datum početka oglasa ne moze biti više od godinu dana unapred."
		ctx.Response().Header().Set("HX-Retarget", "#update-ad-modal")
		return Render(ctx, http.StatusOK, components.UpdateAdModal(updateAdErr, existingAd))
	}

	// 5. The duration between start and end must be at least 3 days.
	minDuration := 3 * 24 * time.Hour
	if endDate.Sub(startDate) < minDuration {
		updateAdErr = "Razmak između početka i kraja oglasa mora biti barem 3 dana."
		ctx.Response().Header().Set("HX-Retarget", "#update-ad-modal")
		return Render(ctx, http.StatusOK, components.UpdateAdModal(updateAdErr, existingAd))
	}

	// 6. End date must not be more than 5 years in the future.
	maxEndDate := midnightNow.AddDate(5, 0, 0)
	if endDate.After(maxEndDate) {
		updateAdErr = "Datum kraja oglasa ne može biti više od 5 godina u budućnosti."
		ctx.Response().Header().Set("HX-Retarget", "#update-ad-modal")
		return Render(ctx, http.StatusOK, components.UpdateAdModal(updateAdErr, existingAd))
	}

	scheduledAds, err := server.store.ListScheduledAds(ctx.Request().Context(), 100)
	if err != nil {
		log.Println("Error listing scheduled ads in updateAd:", err)
		return err
	}

	for _, ad := range scheduledAds {
		if req.Status == "active" && req.Placement == ad.Placement.String && startDate.After(midnightNow.Add(24*time.Hour)) && ad.ID != existingAd.ID {
			updateAdErr = "Već postoji zakazani oglas za ovu poziciju."
			ctx.Response().Header().Set("HX-Retarget", "#update-ad-modal")
			return Render(ctx, http.StatusOK, components.UpdateAdModal(updateAdErr, existingAd))
		}
	}

	activeAds, err := server.store.ListActiveAds(ctx.Request().Context(), 100)
	if err != nil {
		log.Println("Error listing active ads in updateAd:", err)
		return err
	}

	for _, ad := range activeAds {
		if req.Status == "active" && req.Placement == ad.Placement.String && startDate.After(midnightNow) && startDate.Before(midnightNow.Add(24*time.Hour)) && ad.ID != existingAd.ID {
			updateAdErr = "Već postoji aktivan oglas za ovu poziciju."
			ctx.Response().Header().Set("HX-Retarget", "#update-ad-modal")
			return Render(ctx, http.StatusOK, components.UpdateAdModal(updateAdErr, existingAd))
		}

		// Check for overlaps between updated ad and existing active ads
		if req.Status == "active" && ad.ID != existingAd.ID {
			// If same position and the existing ad's end date is after the new ad's start date
			if req.Placement == ad.Placement.String && ad.EndDate.Time.After(startDate) {
				// Calc duration before adjustment to make up for snipped time
				duration := endDate.Sub(startDate)

				// Adjust start date to midnight of the day after the existing ad ends
				adjustedStartDate := time.Date(
					ad.EndDate.Time.Year(),
					ad.EndDate.Time.Month(),
					ad.EndDate.Time.Day()+1,
					0, 0, 0, 0,
					Loc,
				)

				// Update the start date
				startDate = adjustedStartDate

				// Also update end date to maintain the original duration
				endDate = adjustedStartDate.Add(duration)

				// You might want to inform the user about this adjustment
				ctx.Response().Header().Add("HX-Trigger", "adDatesAdjusted")
				break // Stop after finding the first conflict
			}
		}
	}

	// Count ads without including the current one
	var activeCount, scheduledCount int
	for _, ad := range activeAds {
		if ad.ID != adID {
			activeCount++
		}
	}

	for _, ad := range scheduledAds {
		if ad.ID != adID {
			scheduledCount++
		}
	}

	if activeCount >= 4 && req.Status == "active" && startDate.After(midnightNow) && startDate.Before(midnightNow.Add(24*time.Hour)) {
		updateAdErr = "Maksimalan broj aktivnih oglasa je 4."
		ctx.Response().Header().Set("HX-Retarget", "#update-ad-modal")
		return Render(ctx, http.StatusOK, components.UpdateAdModal(updateAdErr, existingAd))
	}

	if scheduledCount >= 4 && req.Status == "active" && startDate.After(midnightNow.Add(24*time.Hour)) {
		updateAdErr = "Maksimalan broj zakazanih oglasa je 4."
		ctx.Response().Header().Set("HX-Retarget", "#update-ad-modal")
		return Render(ctx, http.StatusOK, components.UpdateAdModal(updateAdErr, existingAd))
	}

	// Ensure ImageUrl does not have a double leading slash
	imagePath := filePath
	if !strings.HasPrefix(filePath, "/") {
		imagePath = "/" + filePath
	}

	// Prepare the update parameters
	arg := db.UpdateAdParams{
		ID: adID,
		Title: pgtype.Text{
			String: req.Title,
			Valid:  true,
		},
		Description: pgtype.Text{
			String: req.Description,
			Valid:  true,
		},
		ImageUrl: pgtype.Text{
			String: imagePath,
			Valid:  true,
		},
		TargetUrl: pgtype.Text{
			String: req.TargetUrl,
			Valid:  true,
		},
		Placement: pgtype.Text{
			String: req.Placement,
			Valid:  true,
		},
		Status: pgtype.Text{
			String: req.Status,
			Valid:  true,
		},
		StartDate: pgtype.Timestamptz{
			Time:  startDate,
			Valid: true,
		},
		EndDate: pgtype.Timestamptz{
			Time:  endDate,
			Valid: true,
		},
	}

	_, err = server.store.UpdateAd(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error updating ad:", err)
		return err
	}

	if req.Status == "active" && startDate.After(midnightNow) && startDate.Before(midnightNow.Add(24*time.Hour)) || req.Status == "active" && startDate.Before(midnightNow) {
		ctx.Response().Header().Add("HX-Trigger", "updateAdSuccess")
		return server.activeAdsList(ctx)
	} else if req.Status == "inactive" {
		ctx.Response().Header().Add("HX-Trigger", "updateAdSuccess")
		return server.inactiveAdsList(ctx)
	} else {
		ctx.Response().Header().Add("HX-Trigger", "updateAdSuccess")
		return server.scheduledAdsList(ctx)
	}
}

func (server *Server) deleteAd(ctx echo.Context) error {
	adIDStr := ctx.Param("id")
	adID, err := utils.ParseUUID(adIDStr, "ad ID")
	if err != nil {
		log.Println("Invalid ad ID format in deleteAd:", err)
		return err
	}

	ad, err := server.store.GetAd(ctx.Request().Context(), adID)
	if err != nil {
		log.Println("Error getting ad in deleteAd:", err)
		return err
	}

	filePath := strings.TrimPrefix(ad.ImageUrl.String, "/")
	if err := os.Remove(filePath); err != nil {
		log.Printf("Error removing file from filesystem at %s: %v", filePath, err)
	}

	err = server.store.DeleteAd(ctx.Request().Context(), adID)
	if err != nil {
		log.Println("Error deleting ad in deleteAd:", err)
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (server *Server) deactivateAd(ctx echo.Context) error {
	adIDStr := ctx.Param("id")
	adID, err := utils.ParseUUID(adIDStr, "ad ID")
	if err != nil {
		log.Println("Invalid ad ID format in deactivateAd:", err)
		return err
	}

	_, err = server.store.DeactivateAd(ctx.Request().Context(), adID)
	if err != nil {
		log.Println("Error deactivating ad in deactivateAd:", err)
		return err
	}

	return ctx.NoContent(http.StatusOK)
}
