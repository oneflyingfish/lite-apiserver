package runtimes

import (
	"LiteKube/pkg/common"
	options "LiteKube/pkg/lite-apiserver/options/serverOptions"
	"LiteKube/pkg/lite-apiserver/server/runtimes/ServerHandlers"
	"LiteKube/pkg/util"
	"context"
	"fmt"
	"net/http"
	"sync"

	"k8s.io/klog/v2"
)

type ServerRuntime struct {
	*options.ServerOption
	// BackendTimeout int

	// runtime args
	serverMux   *http.ServeMux
	httpServer  *http.Server
	httpsServer *http.Server

	// routines safe args
	ctx  context.Context
	stop context.CancelFunc
	wg   *sync.WaitGroup
}

func CreateServerRuntime(serverOptions *options.ServerOption, ctx_parent context.Context) (*ServerRuntime, error) {
	ctx, stop := context.WithCancel(ctx_parent)
	return &ServerRuntime{
		serverOptions,
		http.NewServeMux(),
		nil,
		nil,
		ctx,
		stop,
		&sync.WaitGroup{},
	}, nil
}

// run HTTP server faild will not return error, only give some tips.
func (s *ServerRuntime) RunServer() error {
	klog.Info("try to start lite-apiserver")
	s.InitHandlers()

	// run HTTP Server
	if s.InsecurePort > 0 {
		klog.Warningf("you have enable HTTP server at port:%d, which is not secure and you can disable by --insecure-port=-1. Accessing to service by HTTPS is suggested", s.InsecurePort)
		err := s.RunHttpServer()
		if err != nil {
			klog.Errorf("fail to run HTTP server at port:%d", s.InsecurePort)

		}
	}

	// run HTTPS Server
	if s.Port > 0 {
		err := s.RunHttpsServer()
		if err != nil {
			klog.Errorf("fail to run HTTPS server at port:%d", s.Port)

		}
	} else {
		klog.Errorf("If you have specified a bad port=%d, the HTTPS server will refuse to start, please respecify it by --port=", s.Port)
	}

	klog.Info("--------------------------------------------------------------------------------------------")
	// give running tips
	if s.httpServer == nil && s.httpsServer == nil {
		klog.Info("| no server success to start, process terminate directly.")
		util.Exit(0)
	} else if s.httpServer != nil && s.httpsServer == nil {
		klog.Warningf("| ==> HTTP server listens at port:%d, but some errors occur when run the HTTPS server. You can still get your work, but it is not recommended", s.InsecurePort)
	} else {

		if s.httpServer != nil {
			klog.Infof("| ==> HTTP Server listens at port:%d.", s.InsecurePort)
		}
		klog.Infof("| ==> HTTPS Server listens at port:%d.", s.Port)

	}
	klog.Info("--------------------------------------------------------------------------------------------")

	return nil
}

func (s *ServerRuntime) RunHttpServer() error {
	defer s.wg.Done()
	s.wg.Add(1)

	if s.httpServer != nil {
		klog.Error("Start the HTTP Server repeatedly")
		return fmt.Errorf("try to start the HTTP server repeatedly")
	}

	s.httpServer = &http.Server{
		//Addr:    fmt.Sprintf("%s:%d", s.Hostname, s.InsecurePort),
		Addr:    fmt.Sprintf(":%d", s.InsecurePort),
		Handler: s.serverMux,
	}

	// run http server in new routine
	go func() {
		defer s.wg.Done()
		s.wg.Add(1)

		defer func() { s.httpServer = nil }()

		//klog.Fatal(s.httpServer.ListenAndServe())
		if err := s.httpServer.ListenAndServe(); err == http.ErrServerClosed || err == nil {
			klog.Info("HTTP server is closed now.")
		} else {
			klog.Error(err.Error())
		}
	}()

	// read close singnal and close HTTP Server
	go func() {
		defer s.wg.Done()
		s.wg.Add(1)

		<-s.ctx.Done()

		common.CloseServer(s.httpServer, s.SyncDuration, "HTTP server is ready to close...")
	}()

	return nil
}

func (s *ServerRuntime) RunHttpsServer() error {
	defer s.wg.Done()
	s.wg.Add(1)

	if s.httpsServer != nil {
		klog.Error("Start the HTTPS Server repeatedly")
		return fmt.Errorf("try to start the HTTPS server repeatedly")
	}

	// pool := x509.NewCertPool()

	// caCrt, err := ioutil.ReadFile(s.CATLSKeyPair.CACertPath)
	// if err != nil {
	// 	klog.Errorf("Read ca file err: %v", err)
	// 	return err
	// }
	// pool.AppendCertsFromPEM(caCrt)

	// s.httpsServer = &http.Server{
	// 	Addr:    fmt.Sprintf(":%d", s.Port),
	// 	Handler: s.serverMux,
	// 	TLSConfig: &tls.Config{
	// 		ClientCAs:  pool,
	// 		ClientAuth: tls.RequireAndVerifyClientCert,
	// 	},
	// }

	// // run https server in new routine
	// go func() {
	// 	defer s.wg.Done()
	// 	defer func() { s.httpsServer = nil }()
	// 	s.wg.Add(1)

	// 	if err := s.httpsServer.ListenAndServeTLS(); err == http.ErrServerClosed || err == nil {
	// 		klog.Info("HTTPS server is closed now.")
	// 	} else {
	// 		klog.Error(err.Error())
	// 	}
	// }()

	// // read close singnal and close HTTPS Server
	// go func() {
	// 	defer s.wg.Done()
	// 	s.wg.Add(1)

	// 	<-s.ctx.Done()

	// 	common.CloseServer(s.httpsServer, s.SyncDuration, "HTTP server is ready to close...")
	// }()
	return nil
}

func (s *ServerRuntime) InitHandlers() error {
	s.serverMux.HandleFunc("/debug/hello", ServerHandlers.Hello)
	s.serverMux.Handle("/about", http.RedirectHandler("https://github.com/kubesys/LiteKube", http.StatusTemporaryRedirect))
	return nil
}

func (s *ServerRuntime) WaitUtilExit() {
	s.wg.Wait()
}
