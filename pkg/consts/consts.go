package consts

const (
	IdentityKey     = "id"
	SecretKey       = "simple-tiktok"
	DBUser          = "root"
	DBPasswd        = "123456"
	DBHost          = "127.0.0.1"
	DBPort          = "13306"
	DBName          = "tiktok"
	MySQLDefaultDSN = DBUser + ":" + DBPasswd + "@tcp(" + DBHost + ":" + DBPort + ")/" + DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	UserTableName   = "user"
	FollowTableName = "follow"
)
