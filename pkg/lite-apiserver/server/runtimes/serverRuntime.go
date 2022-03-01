package runtimes

import (
	"LiteKube/pkg/common"
	"LiteKube/pkg/lite-apiserver/cert"
	options "LiteKube/pkg/lite-apiserver/options/serverOptions"
	"LiteKube/pkg/lite-apiserver/server/runtimes/ServerHandlers/debug"
	handleTLS "LiteKube/pkg/lite-apiserver/server/runtimes/ServerHandlers/tls"
	"LiteKube/pkg/util"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

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

	time.Sleep(3 * time.Second) // wait 3s, maybe server run error
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
			klog.Errorf("HTTP server may meet some errors while runnning, error tips: %s", err.Error())
		}
	}()

	// read close singnal and close HTTP Server
	go func() {
		defer s.wg.Done()
		s.wg.Add(1)

		<-s.ctx.Done()

		if s.httpServer != nil {
			common.CloseServer(s.httpServer, s.SyncDuration, "HTTP server is ready to close...")
		}
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

	caCertPath, _, caValid := s.CATLSKeyPair.GetTLSKeyPair()
	severCertPath, serverKeyPath, serverValid := s.ServerTLSKeyPair.GetTLSKeyPair()
	if !caValid || !serverValid {
		return fmt.Errorf("loss certificate")
	}

	pool := x509.NewCertPool()
	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		klog.Errorf("Read ca file err: %v", err)
		return err
	}
	pool.AppendCertsFromPEM(caCrt)

	s.httpsServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", s.Port),
		Handler:        s.serverMux,
		IdleTimeout:    90 * time.Second, // matches http.DefaultTransport keep-alive timeout
		ReadTimeout:    4 * 60 * time.Minute,
		WriteTimeout:   4 * 60 * time.Minute,
		MaxHeaderBytes: 1 << 20,
		TLSConfig: &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	// run https server in new routine
	go func() {
		defer s.wg.Done()
		defer func() { s.httpsServer = nil }()
		s.wg.Add(1)

		if err := s.httpsServer.ListenAndServeTLS(severCertPath, serverKeyPath); err == http.ErrServerClosed || err == nil {
			klog.Info("HTTPS server is closed now.")
		} else {
			klog.Errorf("HTTPS server may meet some errors while runnning, error tips: %s", err.Error())
		}
	}()

	// read close singnal and close HTTPS Server
	go func() {
		defer s.wg.Done()
		s.wg.Add(1)

		<-s.ctx.Done()

		if s.httpsServer != nil {
			common.CloseServer(s.httpsServer, s.SyncDuration, "HTTP server is ready to close...")
		}
	}()
	return nil
}

func (s *ServerRuntime) InitHandlers() error {
	//s.serverMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "Welcome, here is LiteKube!\n") })
	s.serverMux.HandleFunc("/tls", sendClientTLS(s.CATLSKeyPair, s.Port))
	s.serverMux.Handle("/debug/", debug.NewDebugHandle(s.CATLSKeyPair, s.Port))
	return nil
}

func (s *ServerRuntime) WaitUtilExit() {
	s.wg.Wait()
}

func sendClientTLS(caTLSKeyPair *cert.TLSKeyPair, port int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode, err := handleTLS.TLSHandleFunc(caTLSKeyPair, port, true)(w, r)
		if statusCode != http.StatusOK {
			w.WriteHeader(statusCode)
		}
		if err != nil {
			fmt.Fprint(w, common.ErrorString(err.Error(), r.URL.Query().Get("format") != "json"))
		}
	}
}
