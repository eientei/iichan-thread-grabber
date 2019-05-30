package model

import (
	"github.com/eientei/iichan-thread-grabber/common"
)

var (
	DatabaseConnection = common.EnvResolveString("DATABASE_CONNECTION", "postgres://iigrabber:iigrabber@127.0.0.1/iigrabber?sslmode=disable")
)
