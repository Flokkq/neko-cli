// Package errors includes helper functions to display cli errors or warnings
package errors

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since     17.12.2025
*/

import (
	"encoding/json"
	"os"
	"time"

	"github.com/nekoman-hq/neko-cli/pkg/plugin"
)

// PluginName and PluginVersion can be set by plugins before using error functions
var (
	PluginName    = "cli"
	PluginVersion = "1.0.0"
)

// WriteError writes an error response to stdout and exits
func WriteError(code, message string) {
	resp := plugin.Response{
		Status: "error",
		Metadata: plugin.ResponseMetadata{
			Plugin:    PluginName,
			Version:   PluginVersion,
			Timestamp: time.Now(),
		},
		Error: &plugin.ResponseError{
			Code:    code,
			Message: message,
		},
	}
	_ = json.NewEncoder(os.Stdout).Encode(resp)
	os.Exit(1)
}

// WriteErrorWithDetails writes an error response with additional details to stdout and exits
func WriteErrorWithDetails(code, message string, details map[string]any) {
	resp := plugin.Response{
		Status: "error",
		Metadata: plugin.ResponseMetadata{
			Plugin:    PluginName,
			Version:   PluginVersion,
			Timestamp: time.Now(),
		},
		Error: &plugin.ResponseError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	_ = json.NewEncoder(os.Stdout).Encode(resp)
	os.Exit(1)
}

// WriteWarning writes a warning response to stdout (does not exit)
func WriteWarning(code, message string) *plugin.Response {
	return &plugin.Response{
		Status: "warning",
		Metadata: plugin.ResponseMetadata{
			Plugin:    PluginName,
			Version:   PluginVersion,
			Timestamp: time.Now(),
		},
		Error: &plugin.ResponseError{
			Code:    code,
			Message: message,
		},
	}
}

// NewErrorResponse creates an error response without writing/exiting (for use within handlers)
func NewErrorResponse(code, message string) *plugin.Response {
	return &plugin.Response{
		Status: "error",
		Metadata: plugin.ResponseMetadata{
			Plugin:    PluginName,
			Version:   PluginVersion,
			Timestamp: time.Now(),
		},
		Error: &plugin.ResponseError{
			Code:    code,
			Message: message,
		},
	}
}

// NewErrorResponseWithDetails creates an error response with details without writing/exiting
func NewErrorResponseWithDetails(code, message string, details map[string]any) *plugin.Response {
	return &plugin.Response{
		Status: "error",
		Metadata: plugin.ResponseMetadata{
			Plugin:    PluginName,
			Version:   PluginVersion,
			Timestamp: time.Now(),
		},
		Error: &plugin.ResponseError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}
