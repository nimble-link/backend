package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/controllers"
	"github.com/nimble-link/backend/middlewares"
	"github.com/nimble-link/backend/pkg/ginutils"
)

func registerWeb(r *ginutils.ApplicationRouter, mids ...gin.HandlerFunc) {
	r.Middlewares(mids...)

	r.Register("POST", "/auth/storeauthcode", controllers.OAuth2Handler)
	r.Register("POST", "/auth/logout", middlewares.Authenticated(), controllers.Logout)
	r.Register("GET", "/auth/userinfo", middlewares.Authenticated(), controllers.GetUserInfo)

	r.Register("POST", "/auth/refresh_token", controllers.RefreshToken)
}
