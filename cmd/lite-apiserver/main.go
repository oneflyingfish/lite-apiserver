package main

import (
	"LiteKube/cmd/lite-apiserver/app"
	"flag"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Init for global klog
	klog.InitFlags(nil)
	klog.Info("Welcome to LiteKube, a Pod deployment and monitoring system for edge weak configuration scenarios, which stay the same call-api with K8S.")
	defer klog.Flush()

	// Init Cobra command
	cmd := app.NewServerCommand()
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// Run LiteKube
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
