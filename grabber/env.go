package grabber

import (
	"github.com/eientei/iichan-thread-grabber/common"
	"time"
)

var ThreadCache = common.EnvResolveDuration("THREAD_CACHE", time.Minute*5)
var ImageCache = common.EnvResolveDuration("IMAGE_CACHE", time.Hour*24)
var UserAgent = common.EnvResolveString("USER_AGENT", "iichan-grabber")
