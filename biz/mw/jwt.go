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

var QueryUserFunc = db.QueryUser

func InitJWT() {
	JwtMiddleware, _ = jwt.New(&jwt.HertzJWTMiddleware{
		Key: []byte(consts.SecretKey),
		// 从哪获得jwt token
		// 从 header 中，根据 Authorization 字段获取
		// 从 query(url参数) 中，根据 token 字段获取
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// 如果是从 header 中获取，根据 jwt 的使用约定，token 前缀是 Bearer
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		Timeout:       time.Hour * 24,
		MaxRefresh:    time.Hour,
		// 根据这个key，可以从 c 中拿出 IdentityHandler 返回的数据
		IdentityKey: consts.IdentityKeyID,
		// IdentityHandler 作用在登录成功后的每次请求中
		IdentityHandler: func(c context.Context, ctx *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(c, ctx)
			return &tiktok.User{
				// 从 token 中解出 id 和 name 信息
				ID:   int64(claims[consts.IdentityKeyID].(float64)),
				Name: claims[consts.IdentityKeyName].(string),
			}
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			// 生成的 token 将会携带用户的 id 和 name 的信息
			if user, ok := data.(*tiktok.User); ok {
				return jwt.MapClaims{
					consts.IdentityKeyID:   user.ID,
					consts.IdentityKeyName: user.Name,
				}
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
			users, err := QueryUserFunc(c, req.Username)
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
				ID:   int64(users[0].ID),
				Name: users[0].Username,
			}
			ctx.Set(consts.IdentityKeyID, user)
			return user, nil
		},
		LoginResponse: func(c context.Context, ctx *app.RequestContext, code int, token string, expire time.Time) {
			user := ctx.Value(consts.IdentityKeyID).(*tiktok.User)
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
