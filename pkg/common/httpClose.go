package common

import (
	"context"
	"net/http"
	"time"

	"k8s.io/klog/v2"
)

func CloseServer(s *http.Server, timesSecond int, tips string) {
	klog.Info(tips)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timesSecond)*time.Second) // As a child of root CTX, to ensure a normal exit
	defer cancel()
	s.Shutdown(ctx)
}
