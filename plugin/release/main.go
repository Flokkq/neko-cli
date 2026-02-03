// plugin/release/main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nekoman-hq/neko-cli/pkg/log"
	"github.com/nekoman-hq/neko-cli/pkg/plugin"
	"github.com/nekoman-hq/neko-cli/plugin/release/pkg/contributors"
	"github.com/nekoman-hq/neko-cli/plugin/release/pkg/history"
	initcmd "github.com/nekoman-hq/neko-cli/plugin/release/pkg/init"
	"github.com/nekoman-hq/neko-cli/plugin/release/pkg/release"
	"github.com/nekoman-hq/neko-cli/plugin/release/pkg/validate"

	// Register all release tools
	_ "github.com/nekoman-hq/neko-cli/plugin/release/pkg/release/tool"
)

func main() {
	// Read request from stdin
	var req plugin.Request
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		writeError("PARSE_ERROR", fmt.Sprintf("failed to parse request: %v", err))
		os.Exit(1)
	}

	// Set verbose mode from request context
	log.Verbose = req.Context.Verbose

	var resp *plugin.Response
	var err error

	switch req.Command {
	case "init":
		resp, err = initcmd.HandleInit(req)
	case "init-options":
		resp, err = initcmd.GetAvailableOptions()
	case "patch":
		resp, err = release.HandleRelease(req, release.Patch)
	case "minor":
		resp, err = release.HandleRelease(req, release.Minor)
	case "major":
		resp, err = release.HandleRelease(req, release.Major)
	case "history":
		resp, err = history.HandleHistory()
	case "contributors":
		resp, err = contributors.HandleContributors()
	case "validate":
		resp, err = validate.HandleValidate(req)
	default:
		resp, err = nil, fmt.Errorf("unknown command: %s", req.Command)
	}

	if err != nil {
		writeError("EXECUTION_ERROR", err.Error())
		os.Exit(1)
	}

	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		writeError("RESPONSE_ERROR", fmt.Sprintf("failed to encode response: %v", err))
		os.Exit(1)
	}
}

func writeError(code, message string) {
	resp := plugin.Response{
		Status: "error",
		Metadata: plugin.ResponseMetadata{
			Plugin:    "release",
			Version:   "1.0.0",
			Timestamp: time.Now(),
		},
		Error: &plugin.ResponseError{
			Code:    code,
			Message: message,
		},
	}
	_ = json.NewEncoder(os.Stdout).Encode(resp)
}
