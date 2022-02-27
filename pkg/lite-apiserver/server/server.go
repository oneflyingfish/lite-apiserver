package server

import (
	"LiteKube/pkg/lite-apiserver/options/serverRunOptions"
	"LiteKube/pkg/lite-apiserver/server/runtimes"
	"context"
	"sync"

	"k8s.io/klog/v2"
)

type LiteServer struct {
	serverRuntime *runtimes.ServerRuntime
	ctx           context.Context
	stop          context.CancelFunc
	wg            *sync.WaitGroup
}

func CreateLiteServer(opt *serverRunOptions.ServerRunOption) (*LiteServer, error) {
	ctx, stop := context.WithCancel(context.Background())

	serverRuntime, err := runtimes.CreateServerRuntime(opt.ServerOption, ctx) // create serverRuntime instance
	return &LiteServer{
		serverRuntime: serverRuntime,
		ctx:           ctx,
		stop:          stop,
		wg:            &sync.WaitGroup{},
	}, err
}

func (s *LiteServer) Run() error {
	defer s.wg.Done()
	s.wg.Add(1)

	// run lite apiserver to listen with HTTP(S)
	if err := s.serverRuntime.RunServer(); err != nil {
		return err
	}

	return nil
}

// close all go routine for LiteServer and exist
func (s *LiteServer) Stop() {
	defer s.wg.Wait()
	s.stop() // give stop signal to servers

	s.serverRuntime.WaitUtilExit()

	klog.Info("All the component is closed now, Goodbye!")
}
