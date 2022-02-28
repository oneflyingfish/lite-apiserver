package version

// Base version information.
//
// This is the fallback data used when version information from git is not
// provided via go ldflags. It provides an approximation of the Kubernetes
// version for ad-hoc builds (e.g. `go build`) that cannot get the version
// information from git.
//
// If you are looking at these fields in the git tree, they look
// strange. They are modified on the fly by the build process. The
// in-tree values are dummy values used for "git archive", which also
// works for GitHub tar downloads.
var (
	litekube     string = "alpha 0.1"
	gitBranch    string = "default-main" // branch of git
	gitVersion          = "v2.25.1"
	gitCommit           = "$HEAD"                 // sha1 from git, output of $(git rev-parse HEAD)
	gitTreeState        = "clean"                 // state of git tree, either "clean" or "dirty"
	buildDate           = "2022-02-00 T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	kubelet             = "v1.23.1"
)
