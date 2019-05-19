package server

import (
	"github.com/eientei/iichan-thread-grabber/common"
	"time"
)

var (
	ListenAddr       = common.EnvResolveString("LISTEN_ADDR", "127.0.0.1:8080")
	OutputDir        = common.EnvResolveString("OUTPUT_DIR", "/tmp/iich")
	PublicBase       = common.EnvResolveString("PUBLIC_BASE", "http://"+ListenAddr)
	CleanupOutputDir = common.EnvResolveDuration("CLEANUP_OUTPUT_DIR", time.Hour*24)
	CleanupFactor    = common.EnvResolveInt("CLEANUP_FACTOR", 10)
)

var (
	DownloaderWorkers     = common.EnvResolveInt("DOWNLOADER_WORKERS", 1)
	DownloaderBuffer      = common.EnvResolveInt("DOWNLOADER_BUFFER", 8*1024)
	DownloaderTimeout     = common.EnvResolveDuration("DOWNLOADER_TIMEOUT", time.Second*30)
	DownloaderNotifyRate  = common.EnvResolveDuration("DOWNLOADER_NOTIFY_RATE", time.Second)
	DownloaderExecuteRate = common.EnvResolveDuration("DOWNLOADER_EXECUTE_RATE", time.Second)
	GrabberWorkers        = common.EnvResolveInt("GRABBER_WORKERS", 1)
	GrabberBuffer         = common.EnvResolveInt("GRABBER_BUFFER", 8)
	GrabberTimeout        = common.EnvResolveDuration("GRABBER_TIMEOUT", time.Hour)
	GrabberNotifyRate     = common.EnvResolveDuration("GRABBER_NOTIFY_RATE", time.Second)
	GrabberExecuteRate    = common.EnvResolveDuration("GRABBER_EXECUTE_RATE", time.Second)
)
