// Package buildinfo is a subpackage that holds information about the current
// binary. This is just used to not populate the zpages package namespace.
package buildinfo

// Build Info that can be set via linker flags:
// go build -ldflags "-X github.com/awbraunstein/zpages/buildinfo.BuildTime=$(date +"%Y.%m.%d.%H%M%S") -X github.com/awbraunstein/zpages/buildinfo.CommitHash=$(git --no-pager log --pretty=format:'%h' -n 1)"
//
// Additionally, these can be set manually in your code:
// import "github.com/awbraunstein/zpages/buildinfo"
//
// init() {
//   buildinfo.BuildTime = "today"
//   buildinfo.CommitHash = "xxxx"
// }
var (
	BuildTime  string
	CommitHash string
)
