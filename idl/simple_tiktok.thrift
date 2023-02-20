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
  1: optional string latest_time  // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
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

struct GetFriendRequest {
  1: i64 user_id  // 用户id
  2: string token // 用户鉴权token
}

struct GetFriendResponse {
  1: i64 status_code              // 状态码，0: 成功，其他值: 失败
  2: optional string status_msg   // 返回状态描述
  3: list<FriendUser> user_list   // 用户列表
}

struct UploadVideoRequest {
  1: string token     // 用户鉴权token
  2: string title     // 视频标题
}

struct UploadVideoResponse {
  1: i64 status_code             // 状态码，0 - 成功，其他值 - 失败
  2: optional string status_msg  // 返回状态描述
}

struct GetPublishRequest {
  1: string token     // 用户鉴权token
  2: string user_id     // 用户id
}

struct GetPublishResponse {
  1: i64 status_code             // 状态码，0 - 成功，其他值 - 失败
  2: optional string status_msg  // 返回状态描述
  3: list<Video> video_list      // 视频列表
}

struct CommentRequest {
  1: string token     // 用户鉴权token
  2: string video_id     // 视频id
  3: string action_type // 发布或删除评论
  4: string comment_text //评论内容
  5: string comment_id //评论id
}

struct CommentResponse {
  1: i64 status_code             // 状态码，0 - 成功，其他值 - 失败
  2: optional string status_msg  // 返回状态描述
  3: Comment comment 
}

struct GetCommentRequest {
  1: string token     // 用户鉴权token
  2: string video_id     // 视频id
}

struct GetCommentResponse {
  1: i64 status_code             // 状态码，0 - 成功，其他值 - 失败
  2: optional string status_msg  // 返回状态描述
  3: list<Comment> comment_list //评论列表 
}

struct FavoriteActionRequest{
  1: string token                                 // 用户鉴权token
  2: i64 video_id                                 // 点赞视频的id
  3: i32 action_type (api.vd="$ == 1 || $ == 2")  // 1-点赞，2-取消点赞
}

struct FavoriteActionResponse {
  1: i64 status_code             // 状态码，0: 成功，其他值: 失败
  2: optional string status_msg  // 返回状态描述
}

struct GetFavoriteListRequest{
  1: i64 user_id                 // 用户id
  2: string token                // 用户鉴权token
}

struct GetFavoriteListResponse{
  1: i64 status_code             // 状态码，0: 成功，其他值: 失败
  2: optional string status_msg  // 返回状态描述
  3: list<Video> video_list      // 点赞过的视频列表
}

struct GetChatRequest {
  1: string token       // 用户鉴权token
  2: i64 to_user_id     // 对方用户id
  3: i64 pre_msg_time   //上次最新消息的时间（新增字段-apk更新中）
}

struct GetChatResponse {
  1: i64 status_code              // 状态码，0 - 成功，其他值 - 失败
  2: optional string status_msg   // 返回状态描述
  3: list<Message> message_list   // 消息列表
}

struct ChatMessageActionRequest {
  1: string token                       // 用户鉴权token
  2: i64 to_user_id                     // 对方用户id
  3: i32 action_type (api.vd="$ == 1")  // 1-发送消息
  4: string content                     // 消息内容
}

struct ChatMessageActionResponse {
  1: i64 status_code              // 状态码，0 - 成功，其他值 - 失败
  2: optional string status_msg   // 返回状态描述
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
  1: i64 id                             // 用户id
  2: string name                        // 用户名称
  3: optional i64 follow_count          // 关注总数
  4: optional i64 follower_count        // 粉丝总数
  5: bool is_follow                     // true - 已关注，false - 未关注
  6: optional string avatar             // 用户头像
  7: optional string background_image   // 用户个人页顶部大图
  8: optional string signature          // 个人简介
  9: optional i64 total_favorited       // 获赞数量
  10: optional i64 work_count           // 作品数量
  11: optional i64 favorite_count       // 点赞数量
}

struct FriendUser {
  1: i64 id                       // 用户id
  2: string name                  // 用户名称
  3: optional i64 follow_count    // 关注总数
  4: optional i64 follower_count  // 粉丝总数
  5: bool is_follow               // true - 已关注，false - 未关注
  6: string avatar                // 用户头像Url
  7: optional string message      // 和该好友的最新聊天消息
  8: i64 msgType                  // message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
}

struct Comment {
  1: i64 id //评论id
  2: User user //评论用户
  3: string content //评论内容
  4: string create_date  //评论创建日期
}

struct Message {
  1: i64 id                       // 消息id
  2: i64 to_user_id               // 该消息接收者的id
  3: i64 from_user_id             // 该消息发送者的id
  4: string content               // 消息内容
  5: optional i64 create_time     // 消息创建时间
}

struct DownloadRequest {
  1: string location
}

// 用户服务
service UserService {
  CreateUserResponse CreateUser(1: CreateUserRequest req) (api.post="/douyin/user/register/")
  CheckUserResponse CheckUser(1: CheckUserRequest req) (api.post="/douyin/user/login/")
  GetUserResponse GetUser(1: GetUserRequest req) (api.get="/douyin/user/")

  FollowUserResponse FollowUser(1: FollowUserRequest req) (api.post="/douyin/relation/action/")
  GetFollowResponse GetFollow(1: GetFollowRequest req) (api.get="/douyin/relation/follow/list/")
  GetFollowerResponse GetFollower(1: GetFollowerRequest req) (api.get="/douyin/relation/follower/list/")

  GetFriendResponse GetFriend(1: GetFriendRequest req) (api.get="/douyin/relation/friend/list/")
}

// 视频服务
service VideoService {
  FeedResponse Feed(1: FeedRequest req) (api.get="/douyin/feed/")
  # 登录用户选择视频上传。
  UploadVideoResponse UploadVideo(1: UploadVideoRequest req) (api.post="/douyin/publish/action/")
  GetPublishResponse GetPublishList(1: GetPublishRequest req) (api.get="/douyin/publish/list/")
}

service CommentService {
  CommentResponse UploadComment(1: CommentRequest req) (api.POST="/douyin/comment/action/")
  GetCommentResponse GetCommentList(1: GetCommentRequest req) (api.get="/douyin/comment/list/")
}

service FavoriteService {
  FavoriteActionResponse FavoriteAction(1: FavoriteActionRequest req)(api.post="/douyin/favorite/action/")
  GetFavoriteListResponse GetFavoriteList(1: GetFavoriteListRequest req)(api.get="/douyin/favorite/list/")
}

service MessageService {
  GetChatResponse GetChat(1: GetChatRequest req) (api.get="/douyin/message/chat/")
  ChatMessageActionResponse ChatMessageAction(1: ChatMessageActionRequest req) (api.post="/douyin/message/action/")
}