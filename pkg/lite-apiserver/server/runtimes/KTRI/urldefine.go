package KTRI

import (
	"path/filepath"
)

const (
	LOGS string = "/logs"

	METRICS          string = "/metrics"
	METRICS_CADVISOR string = "/metrics/cadvisor"
	METRICS_PROBES   string = "/metrics/probes"
	METRICS_RESOURCE string = "/metrics/resource"

	STATS           string = "/stats"
	STATS_CONTAINER string = "/stats/container"
	STATS_SUMMARY   string = "/stats/summary"

	DEBUG       string = "/debug/pprof"
	DEBUG_FLAGS string = "/debug/flags/v"

	PODLIST      string = "/pods"
	RUNNING_PODS string = "/runningpods"

	RUN            string = "/run"
	EXEC           string = "/exec"
	ATTACH         string = "/attach"
	PORT_FORWARD   string = "/portForward"
	CONTAINER_LOGS string = "/containerLogs"

	SPEC    string = "/spec"
	CONFIGZ string = "/configz"
	HEALTHZ string = "/healthz"

	STD_IN_OUT_TTY string = "input=1&output=1&tty=1"
	CRI            string = "/cri/exec"
)

var (
	Log_Path string = filepath.Join(LOGS)
)
