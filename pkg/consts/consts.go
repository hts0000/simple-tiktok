package consts

const (
	IdentityKeyID      = "id"
	IdentityKeyName    = "name"
	SecretKey          = "simple-tiktok"
	DBUser             = "root"
	DBPasswd           = "123456"
	DBHost             = "127.0.0.1"
	DBPort             = "13306"
	DBName             = "tiktok"
	MySQLDefaultDSN    = DBUser + ":" + DBPasswd + "@tcp(" + DBHost + ":" + DBPort + ")/" + DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	UserTableName      = "user"
	FollowTableName    = "follow"
	FavouriteTableName = "likes"
	MessageTableName   = "message"

	FollowUser      = 1
	UnFollowUser    = 2
	FavouriteAction = 1 //点赞状态
	LikeVideo       = 1
	UnLikeVideo     = 2
	SendMsg         = 1
)
