package mw

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"simple-tiktok/biz/dal/db"
	"simple-tiktok/biz/model/tiktok"
	"simple-tiktok/pkg/consts"
	"simple-tiktok/pkg/errno"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
)

var JwtMiddleware *jwt.HertzJWTMiddleware

func InitJWT() {
	JwtMiddleware, _ = jwt.New(&jwt.HertzJWTMiddleware{
		Key:           []byte(consts.SecretKey),
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		Timeout:       time.Hour * 24,
		MaxRefresh:    time.Hour,
		// 根据这个key，可以从上下文中拿出存储的用户id，用于后续的查询
		IdentityKey: consts.IdentityKey,
		// IdentityHandler 作用在登录成功后的每次请求中
		IdentityHandler: func(c context.Context, ctx *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(c, ctx)
			user := claims[consts.IdentityKey].(*tiktok.User)
			return &tiktok.User{
				ID: user.ID,
			}
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			// 在上下文中存入用户id
			if user, ok := data.(*tiktok.User); ok {
				v := jwt.MapClaims{
					consts.IdentityKey: user,
				}
				return v
			}
			// 默认存储 token 的过期时间和创建时间
			return jwt.MapClaims{}
		},
		// 登录时触发，用于认证用户的登录信息
		// 这里返回的数据会存起来，因此出错时故意返回的字符串，让上面的断言失败
		Authenticator: func(c context.Context, ctx *app.RequestContext) (interface{}, error) {
			var req tiktok.CheckUserRequest
			if err := ctx.BindAndValidate(&req); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			if len(req.Username) == 0 || len(req.Password) == 0 {
				return "", jwt.ErrMissingLoginValues
			}

			// TODO: 微服务拆分，使用rpc请求auth服务，rpc.CheckUser
			users, err := db.QueryUser(c, req.Username)
			if err != nil {
				return "", errno.ServiceErr
			}
			if len(users) == 0 {
				return "", jwt.ErrFailedAuthentication
			}

			h := md5.New()
			if _, err = io.WriteString(h, req.Password); err != nil {
				log.Printf("md5加密错误: %v\n", err.Error())
				return "", errno.ServiceErr
			}

			password := fmt.Sprintf("%x", h.Sum(nil))
			if password != users[0].Password {
				return "", jwt.ErrFailedAuthentication
			}

			// 返回将要存储的数据
			user := &tiktok.User{
				ID: int64(users[0].ID),
			}
			ctx.Set(consts.IdentityKey, user)
			return user, nil
		},
		LoginResponse: func(c context.Context, ctx *app.RequestContext, code int, token string, expire time.Time) {
			user := ctx.Value(consts.IdentityKey).(*tiktok.User)
			ctx.JSON(http.StatusOK, tiktok.CheckUserResponse{
				StatusCode: errno.Success.ErrCode,
				StatusMsg:  &errno.Success.ErrMsg,
				UserID:     user.ID,
				Token:      token,
			})
		},
		Unauthorized: func(c context.Context, ctx *app.RequestContext, code int, message string) {
			ctx.JSON(http.StatusOK, tiktok.CheckUserResponse{
				StatusCode: errno.AuthorizationFailedErr.ErrCode,
				StatusMsg:  &message,
			})
		},
		HTTPStatusMessageFunc: func(e error, c context.Context, ctx *app.RequestContext) string {
			switch t := e.(type) {
			case errno.ErrNo:
				return t.ErrMsg
			default:
				return t.Error()
			}
		},
	})
}
