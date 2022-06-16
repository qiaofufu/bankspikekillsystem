package middleware

//func IpLimit() gin.HandlerFunc {
//	return iplimiter.NewRateLimiterMiddleware(redis.NewClient(
//		&redis.Options{
//			Addr:     viper.GetString("redis.addr"),
//			Password: viper.GetString("redis.password"),
//			DB:       1,
//		}), "general", 200, time.Minute)
//}
