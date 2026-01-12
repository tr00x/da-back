package auth

import (
	"time"
)

type Config struct {
	ACCESS_KEY   string
	ACCESS_TIME  time.Duration
	REFRESH_KEY  string
	REFRESH_TIME time.Duration
}

var ENV Config

func Init(access_key string, access_time time.Duration, refresh_key string, refresh_time time.Duration) {
	ENV.ACCESS_KEY = access_key
	ENV.ACCESS_TIME = access_time
	ENV.REFRESH_KEY = refresh_key
	ENV.REFRESH_TIME = refresh_time
}
