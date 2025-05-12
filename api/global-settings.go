package api

import (
	"log"
	"net/http"

	"github.com/00mark0/macva-press/db/services"
	"github.com/labstack/echo/v4"
)

type GlobalSettingsReq struct {
	DisableComments *bool `form:"disable_comments"`
	DisableLikes    *bool `form:"disable_likes"`
	DisableDislikes *bool `form:"disable_dislikes"`
	DisableViews    *bool `form:"disable_views"`
	DisableAds      *bool `form:"disable_ads"`
}

func (server *Server) updateGlobalSettings(ctx echo.Context) error {
	var req GlobalSettingsReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in updateGlobalSettings:", err)
		return err
	}

	globalSettings, err := server.store.GetGlobalSettings(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting global settings in updateGlobalSettings:", err)
		return err
	}

	arg := db.UpdateGlobalSettingsParams{
		DisableComments: func() bool {
			if req.DisableComments != nil {
				return *req.DisableComments
			}
			return globalSettings[0].DisableComments
		}(),
		DisableLikes: func() bool {
			if req.DisableLikes != nil {
				return *req.DisableLikes
			}
			return globalSettings[0].DisableLikes
		}(),
		DisableDislikes: func() bool {
			if req.DisableDislikes != nil {
				return *req.DisableDislikes
			}
			return globalSettings[0].DisableDislikes
		}(),
		DisableViews: func() bool {
			if req.DisableViews != nil {
				return *req.DisableViews
			}
			return globalSettings[0].DisableViews
		}(),
		DisableAds: func() bool {
			if req.DisableAds != nil {
				return *req.DisableAds
			}
			return globalSettings[0].DisableAds
		}(),
	}

	err = server.store.UpdateGlobalSettings(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error updating global settings in updateGlobalSettings:", err)
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (server *Server) resetGlobalSettings(ctx echo.Context) error {
	err := server.store.ResetGlobalSettings(ctx.Request().Context())
	if err != nil {
		log.Println("Error resetting global settings in resetGlobalSettings:", err)
		return err
	}

	return server.adminSettings(ctx)
}
