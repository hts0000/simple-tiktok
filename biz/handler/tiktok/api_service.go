// Code generated by hertz generator.

package tiktok

import (
	"context"

	tiktok "simple-tiktok/biz/model/tiktok"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// CreateUser .
// @router /douyin/user/register/ [POST]
func CreateUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req tiktok.CreateUserRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(tiktok.CreateUserResponse)

	c.JSON(consts.StatusOK, resp)
}

// CheckUser .
// @router /douyin/user/login/ [POST]
func CheckUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req tiktok.CheckUserRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(tiktok.CheckUserResponse)

	c.JSON(consts.StatusOK, resp)
}

// Feed .
// @router /douyin/feed/ [GET]
func Feed(ctx context.Context, c *app.RequestContext) {
	var err error
	var req tiktok.FeedRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(tiktok.FeedResponse)

	c.JSON(consts.StatusOK, resp)
}
