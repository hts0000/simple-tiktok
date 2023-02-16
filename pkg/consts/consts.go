package consts

const (
	IdentityKeyID      = "id"
	IdentityKeyName    = "name"
	SecretKey          = "simple-tiktok"
	DBUser             = "root"
	DBPasswd           = "123456"
	DBHost             = "127.0.0.1"
	DBPort             = "3306"
	DBName             = "tiktok"
	MySQLDefaultDSN    = DBUser + ":" + DBPasswd + "@tcp(" + DBHost + ":" + DBPort + ")/" + DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	UserTableName      = "user"
	FollowTableName    = "follow"
	FavouriteTableName = "likes"

	FollowUser      = 1
	UnFollowUser    = 2
	FavouriteAction = 1 //点赞状态
	DisFavour       = 2 //取消点赞
	LikeAction      = 1 //点赞的行为
)
