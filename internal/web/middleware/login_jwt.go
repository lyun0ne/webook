package middleware

import (
	"net/http"
	"slices"
	"time"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lyun0ne/webook/internal/web"
)

// JWT登录校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}
func (l *LoginJWTMiddlewareBuilder) IngorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}
func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if slices.Contains(l.paths, ctx.Request.URL.Path) {
			return
		}

		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// seg := strings.SplitN(tokenHeader, " ", 2)
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			return []byte("rK5VZ3TsyVneRukCDYsPnBwTWzuSYyA7"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			//严重的安全问题
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("rK5VZ3TsyVneRukCDYsPnBwTWzuSYyA7"))
			if err != nil {
				// 记录日志
			}
			ctx.Header("x-jwt-token", tokenStr)
		}

		ctx.Set("userId", claims.Uid)

	}
}
