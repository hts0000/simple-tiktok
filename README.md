字节跳动后端青训营大作业

TODO:
1. 全局logger
2. 查询用户关注和粉丝的接口结果缓存到Redis
3. db查询时指定查询的列

部署时需要安装ffmpeg，linux系统直接运行install.sh即可，windows需要自行安装ffmpeg
另外，在自己电脑上部署的时候需要去[video.go](.\biz\dal\db\video.go)中修改一下p_url,和c_url的ip,把它改成自己的ip,否则看不了视频。