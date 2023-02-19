// Code generated by hertz generator.

package tiktok

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
	"net/http"
	"simple-tiktok/biz/dal/db"
	tiktok "simple-tiktok/biz/model/tiktok"
	"simple-tiktok/pkg/consts"
	"simple-tiktok/pkg/errno"
)

// FavouriteAction .
// @router /douyin/favorite/action/ [POST]
func FavouriteAction(ctx context.Context, c *app.RequestContext) {
	var req tiktok.FavouriteActionRequest
	var user tiktok.User
	var err = c.BindAndValidate(&req)
	if err != nil {
		log.Printf("参数BindAndValidate失败: %v\n", err.Error())
		c.JSON(http.StatusBadRequest, tiktok.FavouriteActionResponse{
			StatusCode: errno.ParamErr.ErrCode,
			StatusMsg:  &errno.ParamErr.ErrMsg,
		})
		return
	}

	//没有找到对应的视频
	res, err1 := db.GetVideo(ctx, req.ToVideoID)
	if err1 != nil {
		log.Printf("点赞视频: %v不存在", res.Title)
		c.JSON(http.StatusOK, tiktok.FavouriteActionResponse{
			StatusCode: errno.VideoNotExistErr.ErrCode,
			StatusMsg:  &errno.VideoNotExistErr.ErrMsg,
		})
		return
	}

	//根据视频id判断是否已经点赞过这个视频
	ifLike, err2 := db.IsFavourite(ctx, req.ToVideoID, user.ID)
	if ifLike {
		log.Printf("点赞视频: %v失败: %v，已经点赞过该视频", req.ToVideoID, err2.Error())
		c.JSON(http.StatusInternalServerError, tiktok.FavouriteActionResponse{
			StatusCode: errno.ServiceErr.ErrCode,
			StatusMsg:  &errno.ServiceErr.ErrMsg,
		})
		return
	}

	//根据动作： 1点赞  2取消点赞
	switch req.ActionType {
	//case 1
	case consts.FavouriteAction:
		err3 := db.FavouriteAction(ctx, user.ID, req.ToVideoID)
		if err3 != nil {
			log.Printf("用户: %d 点赞视频: %d失败: %v\n", user.Name, req.ToVideoID, err3.Error())
			c.JSON(http.StatusInternalServerError, tiktok.FavouriteActionResponse{
				StatusCode: errno.ServiceErr.ErrCode,
				StatusMsg:  &errno.ServiceErr.ErrMsg,
			})
			return
		}
	//case 2
	case consts.DisFavour:
		err4 := db.FavouriteAction(ctx, user.ID, req.ToVideoID)
		if err4 != nil {
			log.Printf("用户: %d 取消点赞: %d失败: %v\n", user.Name, req.ToVideoID, err4.Error())
			c.JSON(http.StatusInternalServerError, tiktok.FavouriteActionResponse{
				StatusCode: errno.ServiceErr.ErrCode,
				StatusMsg:  &errno.ServiceErr.ErrMsg,
			})
			return
		}
	}

	c.JSON(http.StatusOK, tiktok.FavouriteActionResponse{
		StatusCode: errno.Success.ErrCode,
		StatusMsg:  &errno.Success.ErrMsg,
	})
}

// GetFavouriteList .
// @router /douyin/favorite/list/ [GET]
func GetFavouriteList(ctx context.Context, c *app.RequestContext) {
	var req tiktok.GetFavouriteListRequest
	var err = c.BindAndValidate(&req)
	//查询结果为：BindAndValidate参数失败
	if err != nil {
		log.Printf("参数BindAndValidate失败: %v\n", err.Error())
		c.JSON(http.StatusBadRequest, tiktok.GetFavouriteListResponse{
			StatusCode: errno.ParamErr.ErrCode,
			StatusMsg:  &errno.ParamErr.ErrMsg,
		})
		return
	}

	//返回userId的点赞列表
	videos, err1 := db.GetFavouriteList(ctx, req.UserID)
	//结果为：查询点赞列表失败
	if err1 != nil {
		log.Printf("查询用户: %v的点赞列表失败失败: %v\n", req.UserID, err1.Error())
		c.JSON(http.StatusInternalServerError, tiktok.GetFavouriteListResponse{
			StatusCode: errno.ServiceErr.ErrCode,
			StatusMsg:  &errno.ServiceErr.ErrMsg,
		})
		return
	}

	n := len(videos)
	videoLists := make([]*tiktok.Video, n)
	for i := 0; i < n; i++ {
		videoLists[i] = &tiktok.Video{
			Title:    videos[i].Title,
			ID:       int64(videos[i].ID),
			Author:   db.GetTiktokUser(ctx, videos[i].AuthorId),
			PlayURL:  videos[i].Play_url,
			CoverURL: videos[i].Cover_url,
		}
	}

	c.JSON(http.StatusOK, tiktok.GetFavouriteListResponse{
		StatusCode: errno.Success.ErrCode,
		StatusMsg:  &errno.Success.ErrMsg,
		VideoList:  videoLists,
	})
}
