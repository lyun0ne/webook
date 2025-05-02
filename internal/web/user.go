package web

import (
	"fmt"
	"net/http"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

// const (
// 	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
// 	// 和上面比起来，用 ` 看起来就比较清爽
// 	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
// )

type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	const (
		emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		// 和上面比起来，用 ` 看起来就比较清爽
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)

	return &UserHandler{
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}
func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.GET("/signup", u.SignUp)
	ug.GET("/login", u.Login)
	ug.GET("/edit", u.Edit)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String((http.StatusOK, "两次输入的密码不一致"))
	}
	ctx.String(http.StatusOK, "注册成功")
	fmt.Printf("%v", req)

}

func (u *UserHandler) Login(ctx *gin.Context) {}

func (u *UserHandler) Edit(ctx *gin.Context) {}
