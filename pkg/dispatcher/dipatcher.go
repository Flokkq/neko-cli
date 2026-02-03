package dispatcher

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/nekoman-hq/neko-cli/pkg/plugin"
)

type Dispatcher struct {
	pluginDir string
}

func NewDispatcher(pluginDir string) *Dispatcher {
	return &Dispatcher{
		pluginDir: pluginDir,
	}
}

// Dispatch executes a plugin and captures both output and logs
func (d *Dispatcher) Dispatch(ctx context.Context, pluginName string, req plugin.Request) (*plugin.Response, error) {
	pluginPath, err := d.findPlugin(pluginName)
	if err != nil {
		return nil, fmt.Errorf("plugin not found: %w", err)
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	cmd := exec.CommandContext(ctx, pluginPath)
	cmd.Stdin = bytes.NewReader(reqJSON)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// Check if stdout contains a valid JSON response (error response from plugin)
		if stdout.Len() > 0 {
			var resp plugin.Response
			if jsonErr := json.Unmarshal(stdout.Bytes(), &resp); jsonErr == nil {
				// Valid response found, parse logs and return it
				resp.Logs = parseLogOutput(stderr.String())
				return &resp, nil
			}
		}
		return nil, fmt.Errorf("plugin execution failed: %w\nStderr: %s", err, stderr.String())
	}

	var resp plugin.Response
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		return nil, fmt.Errorf("failed to parse plugin response: %w\nOutput: %s", err, stdout.String())
	}

	// Parse stderr as structured logs
	resp.Logs = parseLogOutput(stderr.String())

	return &resp, nil
}

// parseLogOutput converts stderr lines into structured log entries
// Expected format: "15:04:05 [category] message" or plain text
func parseLogOutput(stderr string) []plugin.LogEntry {
	if stderr == "" {
		return nil
	}

	var logs []plugin.LogEntry
	scanner := bufio.NewScanner(strings.NewReader(stderr))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		entry := parseLogLine(line)
		logs = append(logs, entry)
	}

	return logs
}

// parseLogLine attempts to parse a structured log line
func parseLogLine(line string) plugin.LogEntry {
	// Try to parse: "15:04:05 [category] message"
	parts := strings.SplitN(line, " ", 3)

	if len(parts) >= 3 && strings.HasPrefix(parts[1], "[") && strings.HasSuffix(parts[1], "]") {
		category := strings.Trim(parts[1], "[]")
		level := inferLogLevel(parts[2])

		return plugin.LogEntry{
			Timestamp: parts[0],
			Level:     level,
			Category:  category,
			Message:   parts[2],
		}
	}

	// Fallback: plain text log
	return plugin.LogEntry{
		Timestamp: time.Now().Format("15:04:05"),
		Level:     "info",
		Category:  "plugin",
		Message:   line,
	}
}

// inferLogLevel determines log level based on message content
func inferLogLevel(msg string) string {
	msgLower := strings.ToLower(msg)
	if strings.Contains(msgLower, "error") || strings.Contains(msgLower, "failed") {
		return "error"
	}
	if strings.Contains(msgLower, "warn") {
		return "warn"
	}
	if strings.HasPrefix(msg, "V$") {
		return "verbose"
	}
	return "info"
}

func (d *Dispatcher) findPlugin(name string) (string, error) {
	pluginPath := filepath.Join(d.pluginDir, name, fmt.Sprintf("plugin-%s", name))
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return "", fmt.Errorf("plugin '%s' not found at %s", name, pluginPath)
	}
	return pluginPath, nil
}

func (d *Dispatcher) ListPlugins() ([]plugin.Manifest, error) {
	entries, err := os.ReadDir(d.pluginDir)
	if err != nil {
		return nil, err
	}

	var manifests []plugin.Manifest
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(d.pluginDir, entry.Name(), "manifest.json")
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			continue
		}

		var manifest plugin.Manifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			continue
		}

		manifests = append(manifests, manifest)
	}

	return manifests, nil
}
