# Migration Guide: Von Monolith zu Plugin-Architektur

## Übersicht

Dieser Guide zeigt, wie du deine bestehende Neko CLI von der aktuellen Struktur zur plugin-basierten Architektur migrierst.

## Phase 1: Grundlagen schaffen

### 1.1 Neue Package-Struktur erstellen

```bash
mkdir -p pkg/{plugin,dispatcher,renderer}
```

### 1.2 Plugin Types implementieren

Erstelle `pkg/plugin/types.go` mit allen grundlegenden Typen (siehe PLUGIN_ARCHITECTURE.md)

### 1.3 Dispatcher implementieren

Erstelle `pkg/dispatcher/dispatcher.go` - dieser ist das Herzstück der neuen Architektur

### 1.4 Basic Renderer implementieren

Starte mit JSON und Table Renderern in `pkg/renderer/`

## Phase 2: Erstes Plugin (Release) migrieren

### 2.1 Bestehenden Code analysieren

Dein aktuelles Release-Plugin ist in `plugin/release/pkg/` organisiert:

```
plugin/release/pkg/
├── cmd/          # Cobra Commands
├── config/       # Konfiguration
├── init/         # Initialisierung
├── release/      # Business Logic
└── history/      # History Management
```

### 2.2 Neue Plugin-Struktur

```
plugin/release/
├── main.go                    # Plugin Entry Point (neu)
├── manifest.json              # Plugin Manifest (neu)
├── internal/                  # Private Implementation
│   ├── handlers/             # Request Handlers
│   │   ├── init.go          # Aus pkg/init
│   │   ├── create.go        # Aus pkg/release
│   │   └── history.go       # Aus pkg/history
│   ├── config/              # Aus pkg/config (unverändert)
│   └── domain/              # Business Logic
│       ├── release.go
│       └── version.go
└── Makefile
```

### 2.3 Migration Steps

#### Schritt 1: main.go erstellen

```go
// plugin/release/main.go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    
    "github.com/nekoman-hq/neko-cli/pkg/plugin"
    "github.com/nekoman-hq/neko-cli/plugin/release/internal/handlers"
)

func main() {
    var req plugin.Request
    if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
        writeErrorResponse("PARSE_ERROR", err.Error())
        os.Exit(1)
    }
    
    resp, err := route(req)
    if err != nil {
        writeErrorResponse("EXECUTION_ERROR", err.Error())
        os.Exit(1)
    }
    
    if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
        writeErrorResponse("ENCODE_ERROR", err.Error())
        os.Exit(1)
    }
}

func route(req plugin.Request) (*plugin.Response, error) {
    switch req.Command {
    case "init":
        return handlers.HandleInit(req)
    case "create":
        return handlers.HandleCreate(req)
    case "validate":
        return handlers.HandleValidate(req)
    case "history":
        return handlers.HandleHistory(req)
    default:
        return nil, fmt.Errorf("unknown command: %s", req.Command)
    }
}

func writeErrorResponse(code, message string) {
    resp := plugin.Response{
        Status: "error",
        Error: &plugin.ResponseError{
            Code:    code,
            Message: message,
        },
    }
    json.NewEncoder(os.Stdout).Encode(resp)
}
```

#### Schritt 2: Handler erstellen

```go
// plugin/release/pkg/handlers/init.go
package handlers

import (
    "time"
    
    "github.com/nekoman-hq/neko-cli/pkg/plugin"
    "github.com/nekoman-hq/neko-cli/plugin/release/internal/domain"
)

func HandleInit(req plugin.Request) (*plugin.Response, error) {
    // Bestehende Logic aus pkg/cmd/init.go hierher verschieben
    // Aber statt Cobra-Output → JSON Response
    
    service := domain.NewReleaseService()
    
    result, err := service.Initialize(req.Context.WorkingDir)
    if err != nil {
        return nil, err
    }
    
    return &plugin.Response{
        Status: "success",
        Metadata: plugin.ResponseMetadata{
            Plugin:    "release",
            Version:   "1.0.0",
            Command:   "init",
            Timestamp: time.Now(),
        },
        Data: map[string]any{
            "config_file":    result.ConfigPath,
            "version_system": result.VersionSystem,
            "message":        "Release system initialized successfully",
        },
        RendererHint: "text",
    }, nil
}
```

#### Schritt 3: Bestehenden Code refactoren

**Vorher** (pkg/cmd/init.go):
```go
func InitCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use: "init",
        Run: func(cmd *cobra.Command, args []string) {
            // Direct Console Output
            fmt.Println("Initializing release system...")
            
            service := release.NewService()
            if err := service.Initialize(); err != nil {
                log.Error("Failed: %v", err)
                os.Exit(1)
            }
            
            fmt.Println("✓ Release system initialized")
        },
    }
    return cmd
}
```

**Nachher** (internal/handlers/init.go):
```go
func HandleInit(req plugin.Request) (*plugin.Response, error) {
    service := domain.NewReleaseService()
    
    result, err := service.Initialize(req.Context.WorkingDir)
    if err != nil {
        return nil, err
    }
    
    // Strukturierte Daten statt Console Output
    return &plugin.Response{
        Status:   "success",
        Data:     result.ToMap(),
        RendererHint: "text",
    }, nil
}
```

#### Schritt 4: Domain Logic separieren

```go
// plugin/release/pkg/domain/release.go
package domain

import (
    "fmt"
    "os"
    "path/filepath"
)

type ReleaseService struct {
    // Dependencies
}

type InitResult struct {
    ConfigPath    string
    VersionSystem string
    ProjectType   string
}

func (r *InitResult) ToMap() map[string]any {
    return map[string]any{
        "config_file":    r.ConfigPath,
        "version_system": r.VersionSystem,
        "project_type":   r.ProjectType,
    }
}

func NewReleaseService() *ReleaseService {
    return &ReleaseService{}
}

func (s *ReleaseService) Initialize(workingDir string) (*InitResult, error) {
    // Reine Business Logic - keine UI Concerns
    
    configPath := filepath.Join(workingDir, ".neko-release.json")
    
    // Check if already initialized
    if _, err := os.Stat(configPath); err == nil {
        return nil, fmt.Errorf("release system already initialized")
    }
    
    // Create config
    // ... (deine bestehende Logic)
    
    return &InitResult{
        ConfigPath:    configPath,
        VersionSystem: "semver",
        ProjectType:   "go",
    }, nil
}
```

### 2.4 Manifest erstellen

```json
{
  "name": "release",
  "version": "1.0.0",
  "description": "Release management plugin for Neko CLI",
  "author": "nekoman-hq",
  "commands": [
    {
      "name": "init",
      "description": "Initialize release management system",
      "outputs": ["text", "json"]
    },
    {
      "name": "create",
      "description": "Create a new release",
      "outputs": ["table", "json"]
    },
    {
      "name": "validate",
      "description": "Validate release configuration",
      "outputs": ["table", "json"]
    },
    {
      "name": "history",
      "description": "Show release history",
      "outputs": ["table", "json"]
    }
  ],
  "renderer_types": ["text", "table", "json"]
}
```

### 2.5 Makefile für Plugin

```makefile
PLUGIN_NAME := release
BINARY := plugin-$(PLUGIN_NAME)
INSTALL_DIR := $(HOME)/.neko/plugins/$(PLUGIN_NAME)

.PHONY: build install clean test

build:
	@echo "Building $(PLUGIN_NAME) plugin..."
	go build -o $(BINARY) main.go

install: build
	@echo "Installing $(PLUGIN_NAME) plugin to $(INSTALL_DIR)..."
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY) $(INSTALL_DIR)/
	cp manifest.json $(INSTALL_DIR)/
	@echo "✓ Plugin installed"

clean:
	rm -f $(BINARY)

test:
	go test ./internal/...

# Test plugin locally
test-run: build
	@echo '{"command":"init","args":[],"flags":{},"context":{"working_dir":"."}}' | ./$(BINARY)
```

## Phase 3: CLI Tool anpassen

### 3.1 main.go aktualisieren

```go
// main.go
package main

import (
    "log"
    
    "github.com/nekoman-hq/neko-cli/cmd"
)

func main() {
    // Plugins beim Start laden
    if err := cmd.InitializePlugins(); err != nil {
        log.Printf("Warning: Failed to load plugins: %v", err)
    }
    
    cmd.Execute()
}
```

### 3.2 cmd/root.go erweitern

```go
// cmd/root.go
package cmd

import (
    "os"
    "path/filepath"
    
    "github.com/nekoman-hq/neko-cli/internal/log"
    "github.com/nekoman-hq/neko-cli/pkg/dispatcher"
    "github.com/spf13/cobra"
)

var (
    outputFormat string
    pluginDir    string
)

var rootCmd = &cobra.Command{
    Use:   "neko-cli",
    Short: "Modular CLI tool for release and deployment management",
}

func init() {
    // Global Flags
    rootCmd.PersistentFlags().StringVar(&outputFormat, "output", "table", 
        "Output format (json, table, text)")
    rootCmd.PersistentFlags().BoolVarP(&log.Verbose, "verbose", "v", false, 
        "Enable verbose output")
    
    // Plugin-Verzeichnis
    home, _ := os.UserHomeDir()
    pluginDir = os.Getenv("NEKO_PLUGIN_DIR")
    if pluginDir == "" {
        pluginDir = filepath.Join(home, ".neko", "plugins")
    }
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

// InitializePlugins wird von main.go aufgerufen
func InitializePlugins() error {
    d := dispatcher.NewDispatcher(pluginDir)
    
    manifests, err := d.ListPlugins()
    if err != nil {
        return err
    }
    
    // Für jedes Plugin einen Command erstellen
    for _, manifest := range manifests {
        cmd := CreatePluginCommand(manifest)
        rootCmd.AddCommand(cmd)
    }
    
    return nil
}
```

## Phase 4: Testing

### 4.1 Unit Tests für Plugin

```go
// plugin/release/pkg/handlers/init_test.go
package handlers

import (
    "testing"
    
    "github.com/nekoman-hq/neko-cli/pkg/plugin"
)

func TestHandleInit(t *testing.T) {
    req := plugin.Request{
        Command: "init",
        Context: plugin.Context{
            WorkingDir: t.TempDir(),
        },
    }
    
    resp, err := HandleInit(req)
    if err != nil {
        t.Fatalf("HandleInit failed: %v", err)
    }
    
    if resp.Status != "success" {
        t.Errorf("Expected success, got %s", resp.Status)
    }
    
    if resp.Data["config_file"] == "" {
        t.Error("Expected config_file in response")
    }
}
```

### 4.2 Integration Tests

```go
// pkg/dispatcher/dispatcher_test.go
package dispatcher

import (
    "context"
    "testing"
    
    "github.com/nekoman-hq/neko-cli/pkg/plugin"
)

func TestDispatch(t *testing.T) {
    // Setup test plugin directory
    pluginDir := setupTestPlugins(t)
    
    d := NewDispatcher(pluginDir)
    
    req := plugin.Request{
        Command: "init",
        Context: plugin.Context{
            WorkingDir: t.TempDir(),
        },
    }
    
    resp, err := d.Dispatch(context.Background(), "release", req)
    if err != nil {
        t.Fatalf("Dispatch failed: %v", err)
    }
    
    if resp.Status != "success" {
        t.Errorf("Expected success, got %s", resp.Status)
    }
}
```

## Phase 5: Backward Compatibility (Optional)

Falls du eine Übergangsphase brauchst:

### 5.1 Hybrid Mode

```go
// cmd/release.go (temporary compatibility layer)
package cmd

import (
    "context"
    "os"
    
    "github.com/nekoman-hq/neko-cli/pkg/dispatcher"
    "github.com/nekoman-hq/neko-cli/pkg/plugin"
    "github.com/nekoman-hq/neko-cli/pkg/renderer"
    "github.com/spf13/cobra"
)

// Alte Command-Struktur behalten, aber intern Plugin verwenden
var releaseCmd = &cobra.Command{
    Use:   "release",
    Short: "Release management (plugin-based)",
}

var releaseInitCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize release system",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Dispatch to plugin
        d := dispatcher.NewDispatcher(pluginDir)
        
        req := plugin.Request{
            Command: "init",
            Args:    args,
            Context: plugin.Context{
                WorkingDir: mustGetwd(),
            },
        }
        
        resp, err := d.Dispatch(context.Background(), "release", req)
        if err != nil {
            return err
        }
        
        // Render response
        rf := renderer.NewRendererFactory()
        r, _ := rf.Get(outputFormat)
        return r.Render(resp, os.Stdout)
    },
}

func init() {
    releaseCmd.AddCommand(releaseInitCmd)
    rootCmd.AddCommand(releaseCmd)
}
```

## Phase 6: Rollout-Strategie

### 6.1 Step-by-Step

1. **Week 1**: Core Infrastructure
    - `pkg/plugin`, `pkg/dispatcher`, `pkg/renderer`
    - Basic testing framework

2. **Week 2**: Release Plugin Migration
    - Migriere `init` command
    - Teste End-to-End
    - Documentation

3. **Week 3**: Restliche Release Commands
    - `create`, `validate`, `history`
    - Integration tests

4. **Week 4**: Deploy Plugin Skeleton
    - Basic structure
    - One simple command (e.g., `status`)

5. **Week 5+**: Deploy Plugin Full Implementation
    - Iterative development

### 6.2 Testing Checklist

- [ ] Plugin binary wird korrekt gebaut
- [ ] Manifest wird korrekt gelesen
- [ ] Dispatcher findet und startet Plugin
- [ ] JSON Request/Response funktioniert
- [ ] Alle Renderer funktionieren (json, table, text)
- [ ] Error handling funktioniert
- [ ] Integration mit bestehenden Commands
- [ ] Backward compatibility (falls gewünscht)

## Phase 7: Documentation

### 7.1 Plugin Development Guide

```markdown
# Plugin Development Guide

## Creating a New Plugin

1. Create plugin directory: `plugin/myplugin/`
2. Implement `main.go` with request routing
3. Create `manifest.json`
4. Add Makefile for build & install
5. Implement handlers in `internal/handlers/`
6. Add tests

## Plugin Interface Contract

- Read JSON from STDIN
- Write JSON to STDOUT
- Exit with code 0 on success, 1 on error
- Follow response format specification
```

### 7.2 User Documentation

```markdown
# Using Plugins

## Installing a Plugin

```bash
cd plugin/release
make install
```

## Listing Available Plugins

```bash
neko plugin list
```

## Plugin Commands

Each plugin exposes its own commands:

```bash
neko release init
neko release create --type patch
neko deploy plan coredns
```
```

## Troubleshooting

### Problem: Plugin not found

```bash
# Check plugin directory
ls -la ~/.neko/plugins/

# Verify plugin binary is executable
chmod +x ~/.neko/plugins/release/plugin-release

# Test plugin directly
echo '{"command":"init","args":[],"flags":{},"context":{"working_dir":"."}}' | \
  ~/.neko/plugins/release/plugin-release
```

### Problem: JSON parsing error

```bash
# Enable verbose mode
neko --verbose release init

# Test plugin with manual JSON
echo '{"command":"init"}' | plugin-release
```

## Checkliste für vollständige Migration

- [ ] Package-Struktur erstellt (`pkg/plugin`, `pkg/dispatcher`, `pkg/renderer`)
- [ ] Basic types definiert (`plugin.Request`, `plugin.Response`)
- [ ] Dispatcher implementiert
- [ ] JSON Renderer implementiert
- [ ] Table Renderer implementiert
- [ ] Release plugin migriert
    - [ ] `init` command
    - [ ] `create` command
    - [ ] `validate` command
    - [ ] `history` command
- [ ] Plugin manifest erstellt
- [ ] Plugin Makefile erstellt
- [ ] CLI integration funktioniert
- [ ] Tests geschrieben
- [ ] Documentation aktualisiert
- [ ] Backward compatibility (optional)

## Zeitschätzung

- **Phase 1-2**: 2-3 Tage (Core + erstes Plugin)
- **Phase 3**: 1 Tag (CLI Integration)
- **Phase 4**: 1-2 Tage (Testing)
- **Phase 5**: Optional, 1 Tag
- **Phase 6-7**: Ongoing

**Total**: ~1 Woche für vollständige Migration des Release-Plugins