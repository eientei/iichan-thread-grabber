package common

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func EnvResolveDuration(name string, def time.Duration) (res time.Duration) {
	defer func() {
		fmt.Printf("%30s: %v\n", name, res)
	}()
	if v, ok := os.LookupEnv(name); ok {
		if n, err := time.ParseDuration(v); err == nil {
			return n
		}
	}
	return def
}

func EnvResolveInt(name string, def int) (res int) {
	defer func() {
		fmt.Printf("%30s: %v\n", name, res)
	}()
	if v, ok := os.LookupEnv(name); ok {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return int(n)
		}
	}
	return int(def)
}

func EnvResolveBool(name string, def bool) (res bool) {
	defer func() {
		fmt.Printf("%30s: %v\n", name, res)
	}()
	if v, ok := os.LookupEnv(name); ok {
		if n, err := strconv.ParseBool(v); err == nil {
			return n
		}
	}
	return def
}

func EnvResolveString(name, def string) (res string) {
	defer func() {
		fmt.Printf("%30s: %v\n", name, res)
	}()
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return def
}
