// Code generated by hertz generator.

package tiktok

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"

	"simple-tiktok/biz/dal/db"
	tiktok "simple-tiktok/biz/model/tiktok"
	"simple-tiktok/biz/mw"
	"simple-tiktok/pkg/consts"
	"simple-tiktok/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
)

// CreateUser .
// @router /douyin/user/register/ [POST]
func CreateUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req tiktok.CreateUserRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		log.Printf("参数BindAndValidate失败: %v\n", err.Error())
		c.JSON(http.StatusBadRequest, tiktok.CreateUserResponse{
			StatusCode: errno.ParamErr.ErrCode,
			StatusMsg:  &errno.ParamErr.ErrMsg,
		})
		return
	}

	users, err := db.QueryUser(ctx, req.Username)
	if err != nil {
		log.Printf("查询用户失败: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, tiktok.CreateUserResponse{
			StatusCode: errno.ServiceErr.ErrCode,
			StatusMsg:  &errno.ServiceErr.ErrMsg,
		})
		return
	}
	if len(users) != 0 {
		c.JSON(http.StatusBadRequest, tiktok.CreateUserResponse{
			StatusCode: errno.UserAlreadyExistErr.ErrCode,
			StatusMsg:  &errno.UserAlreadyExistErr.ErrMsg,
		})
		return
	}

	h := md5.New()
	if _, err = io.WriteString(h, req.Password); err != nil {
		log.Printf("md5加密错误: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, tiktok.CreateUserResponse{
			StatusCode: errno.ServiceErr.ErrCode,
			StatusMsg:  &errno.ServiceErr.ErrMsg,
		})
		return
	}

	password := fmt.Sprintf("%x", h.Sum(nil))
	_, err = db.CreateUser(ctx, req.Username, password)
	if err != nil {
		log.Printf("创建用户失败: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, tiktok.CreateUserResponse{
			StatusCode: errno.ServiceErr.ErrCode,
			StatusMsg:  &errno.ServiceErr.ErrMsg,
		})
		return
	}

	mw.JwtMiddleware.LoginHandler(ctx, c)
}

// CheckUser .
// @router /douyin/user/login/ [POST]
func CheckUser(ctx context.Context, c *app.RequestContext) {
	mw.JwtMiddleware.LoginHandler(ctx, c)
}

// GetUser .
// @router /douyin/user/ [GET]
func GetUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req tiktok.GetUserRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, tiktok.GetUserResponse{
			StatusCode: errno.ParamErr.ErrCode,
			StatusMsg:  &errno.ParamErr.ErrMsg,
		})
		return
	}

	followers, err := db.QueryFollower(ctx, uint(req.UserID))
	if err != nil {
		log.Printf("查询用户: %d粉丝失败: %v\n", req.UserID, err.Error())
		c.JSON(http.StatusInternalServerError, tiktok.GetUserResponse{
			StatusCode: errno.ServiceErr.ErrCode,
			StatusMsg:  &errno.ServiceErr.ErrMsg,
		})
		return
	}

	follows, err := db.QueryFollow(ctx, uint(req.UserID))
	if err != nil {
		log.Printf("查询用户: %d关注失败: %v\n", req.UserID, err.Error())
		c.JSON(http.StatusInternalServerError, tiktok.GetUserResponse{
			StatusCode: errno.ServiceErr.ErrCode,
			StatusMsg:  &errno.ServiceErr.ErrMsg,
		})
		return
	}

	user := c.Value(consts.IdentityKeyID).(*tiktok.User)
	followersCount := int64(len(followers))
	followsCount := int64(len(follows))
	c.JSON(http.StatusOK, tiktok.GetUserResponse{
		StatusCode: errno.Success.ErrCode,
		StatusMsg:  &errno.Success.ErrMsg,
		User: &tiktok.User{
			ID:            req.UserID,
			Name:          user.Name,
			FollowCount:   &followsCount,
			FollowerCount: &followersCount,
			IsFollow:      true,
		},
	})
}

// FollowUser .
// @router /douyin/relation/action/ [POST]
func FollowUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req tiktok.FollowUserRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, tiktok.FollowUserResponse{
			StatusCode: errno.ParamErr.ErrCode,
			StatusMsg:  &errno.ParamErr.ErrMsg,
		})
		return
	}

	user := c.Value(consts.IdentityKeyID).(*tiktok.User)
	users, err := db.MGetUsers(ctx, []int64{
		user.ID,
		req.ToUserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, tiktok.FollowUserResponse{
			StatusCode: errno.ServiceErr.ErrCode,
			StatusMsg:  &errno.ServiceErr.ErrMsg,
		})
		return
	}
	if len(users) != 2 {
		log.Printf("用户不存在, 查询用户: %v %v, 得到: ", user.ID, req.ToUserID)
		for _, user := range users {
			fmt.Printf("%d ", user.ID)
		}
		fmt.Println()
		c.JSON(http.StatusOK, tiktok.FollowUserResponse{
			StatusCode: errno.UserNotExistErr.ErrCode,
			StatusMsg:  &errno.UserNotExistErr.ErrMsg,
		})
		return
	}

	switch req.ActionType {
	case consts.FollowUser:
		err := db.FollowUser(ctx, uint(user.ID), uint(req.ToUserID))
		if err != nil {
			log.Printf("用户: %d 关注用户: %d失败: %v\n", user.ID, req.ToUserID, err.Error())
			c.JSON(http.StatusInternalServerError, tiktok.FollowUserResponse{
				StatusCode: errno.ServiceErr.ErrCode,
				StatusMsg:  &errno.ServiceErr.ErrMsg,
			})
			return
		}
	case consts.UnFollowUser:
		err := db.UnFollowUser(ctx, uint(user.ID), uint(req.ToUserID))
		if err != nil {
			log.Printf("用户: %d 取消关注用户: %d失败: %v\n", user.ID, req.ToUserID, err.Error())
			c.JSON(http.StatusInternalServerError, tiktok.FollowUserResponse{
				StatusCode: errno.ServiceErr.ErrCode,
				StatusMsg:  &errno.ServiceErr.ErrMsg,
			})
			return
		}
	}

	c.JSON(http.StatusOK, tiktok.FollowUserResponse{
		StatusCode: errno.Success.ErrCode,
		StatusMsg:  &errno.Success.ErrMsg,
	})
}

// GetFollow .
// @router /douyin/relation/follow/list/ [GET]
func GetFollow(ctx context.Context, c *app.RequestContext) {
	var err error
	var req tiktok.GetFollowRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	resp := new(tiktok.GetFollowResponse)

	c.JSON(http.StatusOK, resp)
}

// GetFollower .
// @router /douyin/relation/follower/list/ [GET]
func GetFollower(ctx context.Context, c *app.RequestContext) {
	var err error
	var req tiktok.GetFollowerRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	resp := new(tiktok.GetFollowerResponse)

	c.JSON(http.StatusOK, resp)
}