namespace go tiktok

struct CreateUserRequest {
  1: string username (api.qury="username", api.vd="len($) > 0 && len($) < 33 && email($); msg:sprintf('Invalid user name: %v', $)")
  2: string password (api.qury="password", api.vd="len($) > 0 && len($) < 33")
}

struct CreateUserResponse {
  1: i64 status_code             // 状态码，0: 成功，其他值: 失败
  2: optional string status_msg  // 返回状态描述
  3: i64 user_id                 // 用户id
  4: string token                // 用户鉴权token
}

struct CheckUserRequest {
  1: string username (api.qury="username", api.vd="len($) > 0 && len($) < 33 && email($); msg:sprintf('Invalid user name: %v', $)")
  2: string password (api.qury="password", api.vd="len($) > 0 && len($) < 33")
}

struct CheckUserResponse {
  1: i64 status_code             // 状态码，0: 成功，其他值: 失败
  2: optional string status_msg  // 返回状态描述
  3: i64 user_id                 // 用户id
  4: string token                // 用户鉴权token
}

struct FeedRequest {
  1: optional i64 latest_time  // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
  2: optional string token     // 可选参数，登录用户设置
}

struct FeedResponse {
  1: i64 status_code             // 状态码，0 - 成功，其他值 - 失败
  2: optional string status_msg  // 返回状态描述
  3: list<Video> video_list      // 视频列表
  4: optional i64 next_time      // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
}

struct GetUserRequest {
  1: i64 user_id  // 用户id
  2: string token // 用户鉴权token
}

struct GetUserResponse {
  1: i64 status_code              // 状态码，0-成功，其他值-失败
  2: optional string status_msg   // 返回状态描述
  3: User user                    // 用户信息
}

struct FollowUserRequest {
  1: string token     // 用户鉴权token
  2: i64 to_user_id   // 对方用户id
  3: i32 action_type (api.vd="$ == 1 || $ == 2")  // 1-关注，2-取消关注
}

struct FollowUserResponse {
  1: i64 status_code             // 状态码，0: 成功，其他值: 失败
  2: optional string status_msg  // 返回状态描述
}

struct GetFollowRequest {
  1: i64 user_id  // 用户id
  2: string token // 用户鉴权token
}

struct GetFollowResponse {
  1: i64 status_code              // 状态码，0: 成功，其他值: 失败
  2: optional string status_msg   // 返回状态描述
  3: list<User> user_list         // 用户信息列表
}

struct GetFollowerRequest {
  1: i64 user_id  // 用户id
  2: string token // 用户鉴权token
}

struct GetFollowerResponse {
  1: i64 status_code              // 状态码，0: 成功，其他值: 失败
  2: optional string status_msg   // 返回状态描述
  3: list<User> user_list         // 用户列表
}

struct UploadVideoRequest {
  1: string token     // 用户鉴权token2: list<byte> data 视频数据
  2: string title     // 视频标题
}

struct UploadVideoResponse {
  1: i64 status_code             // 状态码，0 - 成功，其他值 - 失败
  2: optional string status_msg  // 返回状态描述
}

struct Video {
  1: i64 id              // 视频唯一标识
  2: User author         // 视频作者信息
  3: string play_url     // 视频播放地址
  4: string cover_url    // 视频封面地址
  5: i64 favorite_count  // 视频的点赞总数
  6: i64 comment_count   // 视频的评论总数
  7: bool is_favorite    // true - 已点赞，false - 未点赞
  8: string title        // 视频标题
}

struct User {
  1: i64 id                       // 用户id
  2: string name                  // 用户名称
  3: optional i64 follow_count    // 关注总数
  4: optional i64 follower_count  // 粉丝总数
  5: bool is_follow               // true - 已关注，false - 未关注
}

// 用户服务
service UserService {
  CreateUserResponse CreateUser(1: CreateUserRequest req) (api.post="/douyin/user/register/")
  CheckUserResponse CheckUser(1: CheckUserRequest req) (api.post="/douyin/user/login/")
  GetUserResponse GetUser(1: GetUserRequest req) (api.get="/douyin/user/")

  FollowUserResponse FollowUser(1: FollowUserRequest req) (api.post="/douyin/relation/action/")
  GetFollowResponse GetFollow(1: GetFollowRequest req) (api.get="/douyin/relation/follow/list/")
  GetFollowerResponse GetFollower(1: GetFollowerRequest req) (api.get="/douyin/relation/follower/list/")
}

// 视频服务
service VideoService {
  FeedResponse Feed(1: FeedRequest req) (api.get="/douyin/feed/")
  # 登录用户选择视频上传。
  UploadVideoResponse UploadVideo(1: UploadVideoRequest req) (api.post="/douyin/publish/action/")
}
