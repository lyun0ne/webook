package middleware

import (
	"net/http"
	"time"

	"slices"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IngorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}
func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if slices.Contains(l.paths, ctx.Request.URL.Path) {
			return
		}

		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60 * 60,
		})
		now := time.Now().UnixMilli()
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}

		updateTimeval, _ := updateTime.(int64)

		if now-updateTimeval > 60*1000 {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
	}
}
