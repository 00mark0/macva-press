package api

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/00mark0/macva-press/components"
	"github.com/00mark0/macva-press/db/services"
	"github.com/00mark0/macva-press/utils"
	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// ConvertToWebPWithResize converts an image file to WebP format, resizing it to fit within maxWidth and maxHeight.
// quality is a value from 0 to 100.
func ConvertToWebPWithResize(inputPath string, maxWidth, maxHeight int, quality float32) (string, error) {
	// Open the input file
	file, err := os.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("error opening input file: %v", err)
	}
	defer file.Close()

	// Decode the image based on its extension
	var img image.Image
	ext := strings.ToLower(filepath.Ext(inputPath))
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return inputPath, fmt.Errorf("unsupported image format: %s", ext)
	}
	if err != nil {
		return "", fmt.Errorf("error decoding image: %v", err)
	}

	// Resize the image if maxWidth or maxHeight is specified (> 0)
	if maxWidth > 0 || maxHeight > 0 {
		// imaging.Fit maintains aspect ratio and fits within the given bounds.
		img = imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)
	}

	// Generate WebP filename (replace original extension with .webp)
	webpPath := strings.TrimSuffix(inputPath, ext) + ".webp"

	// Create WebP output file
	output, err := os.Create(webpPath)
	if err != nil {
		return "", fmt.Errorf("error creating WebP output file: %v", err)
	}
	defer output.Close()

	// Encode to WebP (Lossy conversion with specified quality)
	if err := webp.Encode(output, img, &webp.Options{Lossless: false, Quality: quality}); err != nil {
		return "", fmt.Errorf("error encoding to WebP: %v", err)
	}

	// Optionally: Remove the original file if you don't need it
	if err := os.Remove(inputPath); err != nil {
		log.Printf("Warning: could not remove original file %s: %v", inputPath, err)
	}

	return webpPath, nil
}

// video
// Set file size limits based on media type
const maxImageSize = 20 * 1024 * 1024  // 20MB for images
const maxVideoSize = 100 * 1024 * 1024 // 100MB for videos
const maxAudioSize = 50 * 1024 * 1024  // 50MB for audio

// Add this to your server struct initialization or as a package-level variable
func (server *Server) getUploadSemaphore() chan struct{} {
	if server.uploadSemaphore == nil {
		server.uploadSemaphore = make(chan struct{}, 3) // Max 3 concurrent uploads
	}
	return server.uploadSemaphore
}

// OptimizeVideo transcodes videos to a more web-friendly format with reasonable size
func OptimizeVideo(inputPath string) (string, error) {
	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + "_optimized.mp4"

	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-vcodec", "libx264",
		"-crf", "18", // Lower CRF (higher quality), adjust if needed
		"-preset", "medium", // Use slower presets for better quality
		"-vf", "scale='if(gt(iw,1280),1280,-1)':'if(gt(ih,720),720,-1)'", // Scale down only if necessary, keep max 1280x720
		"-movflags", "+faststart", // For better streaming
		"-b:v", "2M", // Set a bitrate for more predictable quality, you can increase/decrease this based on your needs
		"-y", outputPath)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg error: %v - %s", err, stderr.String())
	}

	// Check if the output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return "", fmt.Errorf("output file not created: %v", err)
	}

	// Remove original file to save space
	os.Remove(inputPath)

	return outputPath, nil
}

func (server *Server) addMediaToNewContent(ctx echo.Context) error {
	contentIDCookie, err := ctx.Cookie("content_id")
	if err != nil {
		var emptyMedia []db.Medium
		return Render(ctx, http.StatusOK, components.InsertMedia(emptyMedia, ""))
	}
	contentIDStr := contentIDCookie.Value
	// Parse string UUID into proper UUID format
	contentID, err := utils.ParseUUID(contentIDStr, "content ID")
	if err != nil {
		log.Println("Invalid content ID in addMediaToNewContent:", err)
		return err
	}

	// Implement basic throttling to prevent concurrent heavy uploads
	uploadSemaphore := server.getUploadSemaphore()
	select {
	case uploadSemaphore <- struct{}{}:
		defer func() { <-uploadSemaphore }()
	default:
		return echo.NewHTTPError(http.StatusTooManyRequests, "Server is processing too many uploads")
	}

	// Get the file from the form
	file, err := ctx.FormFile("file_upload")
	if err != nil {
		log.Println("Error retrieving uploaded file in addMediaToNewContent:", err)
		return err
	}

	// Determine media type based on file extension
	mediaType := "image" // Default
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// Validate file type and size
	switch {
	case ext == ".mp4" || ext == ".mov" || ext == ".avi":
		mediaType = "video"
		// Validate video types
		allowedVideoTypes := map[string]bool{".mp4": true, ".mov": true, ".avi": true}
		if !allowedVideoTypes[ext] {
			return echo.NewHTTPError(http.StatusBadRequest, "Unsupported video format")
		}
		if file.Size > maxVideoSize {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Video file too large, maximum size is %d MB", maxVideoSize/1024/1024))
		}
	case ext == ".mp3" || ext == ".wav" || ext == ".ogg":
		mediaType = "audio"
		if file.Size > maxAudioSize {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Audio file too large, maximum size is %d MB", maxAudioSize/1024/1024))
		}
	default:
		// Image handling
		allowedImageTypes := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
		if !allowedImageTypes[ext] {
			return echo.NewHTTPError(http.StatusBadRequest, "Unsupported image format")
		}
		if file.Size > maxImageSize {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Image file too large, maximum size is %d MB", maxImageSize/1024/1024))
		}
	}

	uploadsDir := "static/uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Println("Error creating uploads directory in addMediaToNewContent:", err)
		return err
	}

	// Generate a unique filename to avoid collisions
	filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
	filePath := fmt.Sprintf("%s/%s", uploadsDir, filename)

	// Save the file to disk
	src, err := file.Open()
	if err != nil {
		log.Println("Error opening uploaded file in addMediaToNewContent:", err)
		return err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating destination file in addMediaToNewContent:", err)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		log.Println("Error copying file data in addMediaToNewContent:", err)
		return err
	}

	// Process files based on media type
	if mediaType == "image" {
		// Convert image to WebP with resize
		convertedPath, err := ConvertToWebPWithResize(filePath, 800, 600, 80)
		if err != nil {
			log.Println("Error converting image to WebP:", err)
			// Continue with original file
		} else {
			// Update filePath to the new WebP file
			filePath = convertedPath
		}
	} else if mediaType == "video" {
		// For videos, optimize using ffmpeg
		log.Println("Beginning video optimization for:", filePath)
		optimizedPath, err := OptimizeVideo(filePath)
		if err != nil {
			log.Println("Error optimizing video:", err)
			// Continue with original if optimization fails
		} else {
			filePath = optimizedPath
			log.Println("Video optimization complete, new path:", filePath)
		}
	}

	// Add timeout for database operations
	dbCtx, cancel := context.WithTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()

	// Get existing media to determine order
	existingMedia, err := server.store.ListMediaForContent(dbCtx, contentID)
	if err != nil {
		log.Println("Error listing existing media in addMediaToNewContent:", err)
		return err
	}

	nextOrder := int32(1)
	if len(existingMedia) > 0 {
		nextOrder = int32(len(existingMedia) + 1)
	}

	// Insert the media record into the database
	arg := db.InsertMediaParams{
		ContentID:    contentID,
		MediaType:    mediaType,
		MediaUrl:     "/" + filePath, // Store with leading slash for direct use in HTML
		MediaCaption: "",             // Empty caption by default
		MediaOrder:   nextOrder,
	}

	// Use the context with timeout
	_, err = server.store.InsertMedia(dbCtx, arg)
	if err != nil {
		log.Println("Error inserting media record in addMediaToNewContent:", err)
		return err
	}

	// Add first media as thumbnail if this is the first one
	if nextOrder == 1 {
		thumbnailArg := db.AddThumbnailParams{
			ContentID: contentID,
			Thumbnail: pgtype.Text{String: "/" + filePath, Valid: true},
		}
		_, err := server.store.AddThumbnail(dbCtx, thumbnailArg)
		if err != nil {
			log.Println("Error adding thumbnail in addMediaToNewContent:", err)
			// Continue despite error - not critical
		}
	}

	updatedMedia, err := server.store.ListMediaForContent(dbCtx, contentID)
	if err != nil {
		log.Println("Error listing updated media in addMediaToNewContent:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.InsertMedia(updatedMedia, contentID.String()))
}

func (server *Server) addMediaToUpdateContent(ctx echo.Context) error {
	contentIDStr := ctx.Param("id")
	// Parse string UUID into proper UUID format
	contentID, err := utils.ParseUUID(contentIDStr, "content ID")
	if err != nil {
		log.Println("Invalid content ID in addMediaToUpdateContent:", err)
		return err
	}

	// Implement basic throttling to prevent concurrent heavy uploads
	uploadSemaphore := server.getUploadSemaphore() // This would be a package-level semaphore
	select {
	case uploadSemaphore <- struct{}{}:
		defer func() { <-uploadSemaphore }()
	default:
		return echo.NewHTTPError(http.StatusTooManyRequests, "Server is processing too many uploads")
	}

	// Get the file from the form
	file, err := ctx.FormFile("file_upload")
	if err != nil {
		log.Println("Error retrieving uploaded file in addMediaToUpdateContent:", err)
		return err
	}

	// Determine media type based on file extension
	mediaType := "image" // Default
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// Validate file type and size
	switch {
	case ext == ".mp4" || ext == ".mov" || ext == ".avi":
		mediaType = "video"
		// Validate video types
		allowedVideoTypes := map[string]bool{".mp4": true, ".mov": true, ".avi": true}
		if !allowedVideoTypes[ext] {
			return echo.NewHTTPError(http.StatusBadRequest, "Unsupported video format")
		}
		if file.Size > maxVideoSize {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Video file too large, maximum size is %d MB", maxVideoSize/1024/1024))
		}
	case ext == ".mp3" || ext == ".wav" || ext == ".ogg":
		mediaType = "audio"
		if file.Size > maxAudioSize {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Audio file too large, maximum size is %d MB", maxAudioSize/1024/1024))
		}
	default:
		// Image handling
		allowedImageTypes := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
		if !allowedImageTypes[ext] {
			return echo.NewHTTPError(http.StatusBadRequest, "Unsupported image format")
		}
		if file.Size > maxImageSize {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Image file too large, maximum size is %d MB", maxImageSize/1024/1024))
		}
	}

	uploadsDir := "static/uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Println("Error creating uploads directory in addMediaToUpdateContent:", err)
		return err
	}

	// Generate a unique filename to avoid collisions
	filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
	filePath := fmt.Sprintf("%s/%s", uploadsDir, filename)

	// Save the file to disk
	src, err := file.Open()
	if err != nil {
		log.Println("Error opening uploaded file in addMediaToUpdateContent:", err)
		return err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating destination file in addMediaToUpdateContent:", err)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		log.Println("Error copying file data in addMediaToUpdateContent:", err)
		return err
	}

	// Process files based on media type
	if mediaType == "image" {
		// Convert image to WebP with resize
		convertedPath, err := ConvertToWebPWithResize(filePath, 800, 600, 80)
		if err != nil {
			log.Println("Error converting image to WebP:", err)
			// Continue with original file
		} else {
			// Update filePath to the new WebP file
			filePath = convertedPath
		}
	} else if mediaType == "video" {
		// For videos, optimize using ffmpeg
		log.Println("Beginning video optimization for:", filePath)
		optimizedPath, err := OptimizeVideo(filePath)
		if err != nil {
			log.Println("Error optimizing video:", err)
			// Continue with original if optimization fails
		} else {
			filePath = optimizedPath
			log.Println("Video optimization complete, new path:", filePath)
		}
	}

	// Add timeout for database operations
	dbCtx, cancel := context.WithTimeout(ctx.Request().Context(), 30*time.Second)
	defer cancel()

	// Get existing media to determine order
	existingMedia, err := server.store.ListMediaForContent(dbCtx, contentID)
	if err != nil {
		log.Println("Error listing existing media in addMediaToUpdateContent:", err)
		return err
	}

	nextOrder := int32(1)
	if len(existingMedia) > 0 {
		nextOrder = int32(len(existingMedia) + 1)
	}

	// Insert the media record into the database
	arg := db.InsertMediaParams{
		ContentID:    contentID,
		MediaType:    mediaType,
		MediaUrl:     "/" + filePath, // Store with leading slash for direct use in HTML
		MediaCaption: "",             // Empty caption by default
		MediaOrder:   nextOrder,
	}

	media, err := server.store.InsertMedia(dbCtx, arg)
	if err != nil {
		log.Println("Error inserting media record in addMediaToUpdateContent:", err)
		return err
	}

	// Set first uploaded media as thumbnail
	if media.MediaOrder == 1 {
		arg := db.AddThumbnailParams{
			ContentID: contentID,
			Thumbnail: pgtype.Text{String: "/" + filePath, Valid: true},
		}
		_, err := server.store.AddThumbnail(dbCtx, arg)
		if err != nil {
			log.Println("Error adding thumbnail in addMediaToUpdateContent:", err)
			return err
		}
	}

	updatedMedia, err := server.store.ListMediaForContent(dbCtx, contentID)
	if err != nil {
		log.Println("Error listing updated media in addMediaToUpdateContent:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.InsertMediaUpdate(updatedMedia, contentID.String()))
}

func (server *Server) listMediaForContent(ctx echo.Context) error {
	contentIDCookie, err := ctx.Cookie("content_id")
	if err != nil {
		var emptyMedia []db.Medium

		return Render(ctx, http.StatusOK, components.InsertMedia(emptyMedia, ""))
	}

	// Parse content ID from cookie
	contentIDString := contentIDCookie.Value
	contentID, err := utils.ParseUUID(contentIDString, "content ID")
	if err != nil {
		log.Println("Invalid content ID in listMediaForContent:", err)
		return err
	}

	media, err := server.store.ListMediaForContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error listing media for content in listMediaForContent:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.InsertMedia(media, contentIDCookie.Value))
}

func (server *Server) deleteMedia(ctx echo.Context) error {
	mediaIDStr := ctx.Param("id")
	mediaID, err := utils.ParseUUID(mediaIDStr, "media ID")
	if err != nil {
		log.Println("Invalid media ID format in deleteMedia:", err)
		return err
	}

	// Get the media record to find the file path before deleting
	media, err := server.store.GetMediaByID(ctx.Request().Context(), mediaID)
	if err != nil {
		log.Println("Error getting media record in deleteMedia:", err)
		return err
	}

	// Get the content ID to use for rendering updated media list
	contentID := media.ContentID
	contentIDStr := contentID.String()

	// Remove the file from filesystem
	// The filepath is stored with leading slash, so trim it for filesystem operations
	filePath := strings.TrimPrefix(media.MediaUrl, "/")
	if err := os.Remove(filePath); err != nil {
		log.Println("Error removing file from filesystem in deleteMedia:", err)
		// Continue with DB deletion even if file removal fails
	}

	// Delete the media record from the database
	if err := server.store.DeleteMedia(ctx.Request().Context(), mediaID); err != nil {
		log.Println("Error deleting media record in deleteMedia:", err)
		return err
	}

	// Get updated media list for rendering
	updatedMedia, err := server.store.ListMediaForContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error listing updated media in deleteMedia:", err)
		return err
	}

	// Render the updated media list component
	return Render(ctx, http.StatusOK, components.InsertMediaUpdate(updatedMedia, contentIDStr))
}

func (server *Server) listMediaForArticlePage(ctx echo.Context) error {
	articleIDStr := ctx.Param("id")
	articleID, err := utils.ParseUUID(articleIDStr, "article ID")
	if err != nil {
		log.Println("Invalid article ID in listMediaForArticlePage:", err)
		return err
	}

	media, err := server.store.ListMediaForContent(ctx.Request().Context(), articleID)
	if err != nil {
		log.Println("Error listing media for article in listMediaForArticlePage:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.ArticleMediaSlider(media))
}
