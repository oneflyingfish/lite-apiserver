package main

import (
	"LiteKube/cmd/lite-apiserver/app"
	"flag"
	"fmt"
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
	defer klog.Flush()
	klog.MaxSize = 10240

	if err := os.MkdirAll("litekube-logs/lite-apiserver", os.ModePerm); err != nil {
		panic(err)
	}

	flag.Set("logtostderr", "false")
	year, month, day := time.Now().Date()
	flag.Set("log_file", fmt.Sprintf("litekube-logs/lite-apiserver/log-%d-%d-%d_%d-%d.log", year, month, day, time.Now().Hour(), time.Now().Minute()))
	flag.Set("alsologtostderr", "true")

	// Init Cobra command
	cmd := app.NewServerCommand()
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// Run LiteKube
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
