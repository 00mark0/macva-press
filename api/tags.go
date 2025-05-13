package api

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	"github.com/00mark0/macva-press/components"
	"github.com/00mark0/macva-press/db/services"
	"github.com/00mark0/macva-press/utils"
)

type CreateTagReq struct {
	TagName string `form:"tag_name" validate:"required"`
}

func (server *Server) createTag(ctx echo.Context) error {
	var req CreateTagReq
	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in createTag:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		message := "Ime je obavezno pri pravljenju novog taga."
		// Set header to indicate error and where to show it
		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	_, err := server.store.CreateTag(ctx.Request().Context(), req.TagName)
	if err != nil {
		message := "Tag sa ovim imenom već postoji."
		// Set header to indicate error and where to show it
		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	// Success case - get tags and render
	tags, err := server.store.ListTags(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Error listing tags in createTag:", err)
		return err
	}

	contentIDCookie, err := ctx.Cookie("content_id")
	if err != nil {
		var contentTags []db.Tag

		// Set header to indicate success and where to show updated tags
		ctx.Response().Header().Set("HX-Retarget", "#admin-tags")
		return Render(ctx, http.StatusOK, components.AdminTags(tags, contentTags))
	}

	contentIDString := contentIDCookie.Value
	contentID, err := utils.ParseUUID(contentIDString, "content ID")
	if err != nil {
		log.Println("Invalid content ID in createTag:", err)
		return err
	}

	contentTags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting tags by content in createTag:", err)
		return err
	}

	// Set header to indicate success and where to show updated tags
	ctx.Response().Header().Set("HX-Retarget", "#admin-tags")
	return Render(ctx, http.StatusOK, components.AdminTagsUpdate(tags, contentTags, contentID.String()))
}

type AddTagReq struct {
	TagID string `form:"tag_id" validate:"required"`
}

func (server *Server) addTagToContent(ctx echo.Context) error {
	var req AddTagReq

	contentIDCookie, err := ctx.Cookie("content_id")
	if err != nil {
		message := "Sadržaj nije pronađen. Da bi dodali tagove pritisnite sačuvaj ili objavi."

		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	contentIDString := contentIDCookie.Value
	contentID, err := utils.ParseUUID(contentIDString, "content ID")
	if err != nil {
		log.Println("Invalid content ID in addTagToContent:", err)
		return err
	}

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in addTagToContent:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		message := "Izaberite postojeći tag."
		// Set header to indicate error and where to show it
		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	tagID, err := utils.ParseUUID(req.TagID, "tag ID")
	if err != nil {
		log.Println("Invalid tag ID in addTagToContent:", err)
		return err
	}

	arg := db.AddTagToContentParams{
		ContentID: contentID,
		TagID:     tagID,
	}

	err = server.store.AddTagToContent(ctx.Request().Context(), arg)
	if err != nil {
		message := "Greška prilikom dodavanja taga."

		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	// Success case - get tags and render
	tags, err := server.store.ListTags(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Error listing tags in addTagToContent:", err)
		return err
	}

	contentTags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting tags by content in addTagToContent:", err)
		return err
	}

	// Set header to indicate success and where to show updated tags
	ctx.Response().Header().Set("HX-Retarget", "#admin-tags")
	return Render(ctx, http.StatusOK, components.AdminTagsUpdate(tags, contentTags, contentID.String()))

}

func (server *Server) addTagToContentUpdate(ctx echo.Context) error {
	var req AddTagReq

	contentIDStr := ctx.Param("id")
	contentID, err := utils.ParseUUID(contentIDStr, "content ID")
	if err != nil {
		log.Println("Invalid content ID in addTagToContentUpdate:", err)
		return err
	}

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in addTagToContentUpdate:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		message := "Izaberite postojeći tag."
		// Set header to indicate error and where to show it
		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	tagID, err := utils.ParseUUID(req.TagID, "tag ID")
	if err != nil {
		log.Println("Invalid tag ID in addTagToContentUpdate:", err)
		return err
	}

	arg := db.AddTagToContentParams{
		ContentID: contentID,
		TagID:     tagID,
	}

	err = server.store.AddTagToContent(ctx.Request().Context(), arg)
	if err != nil {
		message := "Greška prilikom dodavanja taga."

		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	// Success case - get tags and render
	tags, err := server.store.ListTags(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Error listing tags in addTagToContentUpdate:", err)
		return err
	}

	contentTags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting tags by content in addTagToContentUpdate:", err)
		return err
	}

	// Set header to indicate success and where to show updated tags
	ctx.Response().Header().Set("HX-Retarget", "#admin-tags")
	return Render(ctx, http.StatusOK, components.AdminTagsUpdate(tags, contentTags, contentID.String()))

}

func (server *Server) removeTagFromContent(ctx echo.Context) error {
	contentIDCookie, err := ctx.Cookie("content_id")
	if err != nil {
		message := "Sadržaj nije pronađen. Dodajte tagove sa uredi stranice."

		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	contentIDString := contentIDCookie.Value
	contentID, err := utils.ParseUUID(contentIDString, "content ID")
	if err != nil {
		log.Println("Invalid content ID format in removeTagFromContent:", err)
		return err
	}

	tagIDString := ctx.Param("id")
	tagID, err := utils.ParseUUID(tagIDString, "tag ID")
	if err != nil {
		log.Println("Invalid tag ID format in removeTagFromContent:", err)
		return err
	}

	arg := db.RemoveTagFromContentParams{
		ContentID: contentID,
		TagID:     tagID,
	}

	err = server.store.RemoveTagFromContent(ctx.Request().Context(), arg)
	if err != nil {
		message := "Greška prilikom uklanjanja taga."

		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	// Success case - get tags and render
	tags, err := server.store.ListTags(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Error listing tags in removeTagFromContent:", err)
		return err
	}

	contentTags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting tags by content in removeTagFromContent:", err)
		return err
	}

	// Set header to indicate success and where to show updated tags
	ctx.Response().Header().Set("HX-Retarget", "#admin-tags")
	return Render(ctx, http.StatusOK, components.AdminTagsUpdate(tags, contentTags, contentID.String()))
}

func (server *Server) removeTagFromContentUpdate(ctx echo.Context) error {
	contentIDStr := ctx.Param("content_id")
	contentID, err := utils.ParseUUID(contentIDStr, "content ID")
	if err != nil {
		log.Println("Invalid content ID format in removeTagFromContent:", err)
		return err
	}

	tagIDString := ctx.Param("tag_id")
	tagID, err := utils.ParseUUID(tagIDString, "tag ID")
	if err != nil {
		log.Println("Invalid tag ID format in removeTagFromContent:", err)
		return err
	}

	arg := db.RemoveTagFromContentParams{
		ContentID: contentID,
		TagID:     tagID,
	}

	err = server.store.RemoveTagFromContent(ctx.Request().Context(), arg)
	if err != nil {
		message := "Greška prilikom uklanjanja taga."

		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	// Success case - get tags and render
	tags, err := server.store.ListTags(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Failed to get tags in removeTagFromContent:", err)
		return err
	}

	contentTags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Failed to get tags by content in removeTagFromContent:", err)
		return err
	}

	// Set header to indicate success and where to show updated tags
	ctx.Response().Header().Set("HX-Retarget", "#admin-tags")
	return Render(ctx, http.StatusOK, components.AdminTagsUpdate(tags, contentTags, contentID.String()))
}

func (server *Server) deleteTag(ctx echo.Context) error {
	tagIDStr := ctx.Param("id")
	tagID, err := utils.ParseUUID(tagIDStr, "tag ID")
	if err != nil {
		log.Println("Invalid tag ID format in deleteTag:", err)
		return err
	}

	err = server.store.DeleteTag(ctx.Request().Context(), tagID)
	if err != nil {
		log.Println("Error deleting tag in deleteTag:", err)
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type ListTagsReq struct {
	Limit int32 `query:"limit"`
}

func (server *Server) listTags(ctx echo.Context) error {
	var req ListTagsReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listTags:", err)
		return err
	}

	nextLimit := req.Limit + 20

	tags, err := server.store.ListTags(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing tags in listTags:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.TagsList(int(nextLimit), tags))
}

type SearchTagsReq struct {
	SearchTerm string `query:"search_term" validate:"required"`
}

func (server *Server) listSearchTags(ctx echo.Context) error {
	var req SearchTagsReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listSearchTags:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in listSearchTags:", err)
		return err
	}

	arg := db.SearchTagsParams{
		Search: req.SearchTerm,
		Limit:  20,
	}

	tags, err := server.store.SearchTags(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error searching tags in listSearchTags:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.TagsList(20, tags))
}

func (server *Server) listTagsByContent(ctx echo.Context) error {
	contentIDStr := ctx.Param("id")
	contentID, err := utils.ParseUUID(contentIDStr, "content ID")
	if err != nil {
		log.Println("Invalid content ID in listTagsByContent:", err)
		return err
	}

	tags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting tags by content in listTagsByContent:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.TagsInArticleDetes(tags))
}

// ContentByTagsHandler modified to return HTML instead of JSON
func (server *Server) listContentByTagsUnderCategory(ctx echo.Context) error {
	categoryIDStr := ctx.Param("id")
	categoryID, err := utils.ParseUUID(categoryIDStr, "category ID")
	if err != nil {
		log.Println("Invalid category ID in listContentByTagsUnderCategory:", err)
		return err
	}

	// Get category info for the header
	category, err := server.store.GetCategoryByID(ctx.Request().Context(), categoryID)
	if err != nil {
		log.Println("Error fetching category in listContentByTagsUnderCategory:", err)
		return err
	}

	// Step 1: Get unique tags for the category
	tags, err := server.store.GetUniqueTagsByCategoryID(ctx.Request().Context(), categoryID)
	if err != nil {
		log.Println("Error fetching unique tags for category:", categoryID, err)
		return err
	}

	const limit = 6 // Show fewer items initially per tag
	const offset = 0

	// Step 2: Loop through tags and fetch content
	var contentByTags components.ContentByTagsList

	for _, tag := range tags {
		content, err := server.store.ListContentByTag(ctx.Request().Context(), db.ListContentByTagParams{
			TagName: tag.TagName,
			Limit:   limit,
			Offset:  offset,
		})

		for i := range content {
			if content[i].Thumbnail.String == "" {
				content[i].Thumbnail = pgtype.Text{String: ThumbnailURL, Valid: true}
			}
		}

		if err != nil {
			log.Println("Error fetching content for tag:", tag, err)
			continue // Skip this tag but continue with others
		}

		// Only add tags that have content
		if len(content) > 0 {
			contentByTags = append(contentByTags, components.ContentByTag{
				TagID:   tag.TagID.String(),
				TagName: tag.TagName,
				Content: content,
			})
		}
	}

	// Get global settings for display options
	globalSettings, err := server.store.GetGlobalSettings(ctx.Request().Context())
	if err != nil {
		log.Println("Error fetching global settings:", err)
		// Use default settings if we can't fetch them
		globalSettings[0] = db.GlobalSetting{
			DisableComments: false,
			DisableLikes:    false,
			DisableViews:    false,
		}
	}

	// Render the templ component
	return Render(ctx, http.StatusOK, components.ContentByTagsSection(contentByTags, globalSettings[0], category.CategoryName))
}

func (server *Server) listAllContentByTag(ctx echo.Context) error {
	component, err := server.GenerateRecentTagContentComponent(ctx)
	if err != nil {
		log.Println("Error generating component in listAllContentByTag:", err)
		return err
	}

	return Render(ctx, http.StatusOK, component)
}

func (server *Server) GenerateRecentTagContentComponent(ctx echo.Context) (templ.Component, error) {
	var req ListTagsReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listAllContentByTag:", err)
		return nil, err
	}

	tagIDStr := ctx.Param("id")
	tagID, err := utils.ParseUUID(tagIDStr, "tag ID")
	if err != nil {
		log.Println("Invalid tag ID format in listAllContentByTag:", err)
		return nil, err
	}

	tag, err := server.store.GetTag(ctx.Request().Context(), tagID)
	if err != nil {
		log.Println("Error fetching tag in listAllContentByTag:", err)
		return nil, err
	}

	nextLimit := req.Limit + 9

	arg := db.ListContentByTagLimitParams{
		TagName: tag.TagName,
		Limit:   nextLimit,
	}

	content, err := server.store.ListContentByTagLimit(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error fetching content by tag in listAllContentByTag:", err)
		return nil, err
	}

	for i := range content {
		if content[i].Thumbnail.String == "" {
			content[i].Thumbnail = pgtype.Text{String: ThumbnailURL, Valid: true}
		}
	}

	globalSettings, err := server.store.GetGlobalSettings(ctx.Request().Context())
	if err != nil {
		log.Println("Error fetching global settings in listAllContentByTag:", err)
		// Use default settings if we can't fetch them
		globalSettings[0] = db.GlobalSetting{
			DisableComments: false,
			DisableLikes:    false,
			DisableDislikes: true,
			DisableViews:    false,
		}
	}

	return components.TagsGrid(tagIDStr, content, int(nextLimit), globalSettings[0]), nil

}
