# Neko CLI - Plugin-Basierte Architektur

## Überblick

Diese Architektur trennt **Business Logic** (in Plugins) von **Präsentation** (in Renderern) durch einen standardisierten JSON-basierten Kommunikationskanal.

```
┌─────────────┐     ┌──────────────┐     ┌────────────┐     ┌──────────────┐
│   CLI Tool  │ ──> │   Dispatcher │ ──> │   Plugin   │ ──> │   Renderer   │
│   (Cobra)   │     │              │     │  (Binary)  │     │  (JSON→UI)   │
└─────────────┘     └──────────────┘     └────────────┘     └──────────────┘
                           │                    │                    │
                           │                    │                    │
                    Plugin-Discovery      JSON Output         Pretty Output
```

## Kernkonzepte

### 1. Plugin-System
- **Plugins als eigenständige Binaries**: Jedes Plugin ist ein kompiliertes Go-Binary
- **Standard Input/Output**: Plugins kommunizieren über STDIN/STDOUT mit JSON
- **Versionierung**: Jedes Plugin hat eine Manifest-Datei mit Metadaten
- **Discovery**: Automatisches Erkennen installierter Plugins

### 2. JSON als Kommunikationsprotokoll
- Strukturierte, typisierte Daten
- Einfach zu parsen und zu validieren
- Versionierbar und erweiterbar
- Language-agnostic für zukünftige Erweiterungen

### 3. Renderer-System
- Trennung von Daten und Darstellung
- Multiple Renderer pro Plugin möglich (Table, JSON, TUI, Web)
- Konfigurierbar über Flags

## Architektur-Details

### Plugin-Struktur

Jedes Plugin wird als eigenständiges Binary kompiliert und in einem Plugin-Verzeichnis installiert:

```
~/.neko/plugins/
├── release/
│   ├── manifest.json
│   └── plugin-release          # Ausführbare Binary
├── deploy/
│   ├── manifest.json
│   └── plugin-deploy
└── core/
    ├── manifest.json
    └── plugin-core
```

### Manifest-Datei

Jedes Plugin hat eine `manifest.json`:

```json
{
  "name": "release",
  "version": "1.0.0",
  "description": "Release management plugin",
  "author": "nekoman-hq",
  "commands": [
    {
      "name": "init",
      "description": "Initialize release system",
      "outputs": ["table", "json"]
    },
    {
      "name": "create",
      "description": "Create new release",
      "outputs": ["table", "json", "tui"]
    }
  ],
  "renderer_types": ["table", "json", "tui"]
}
```

### Plugin-Kommunikation

#### Request (STDIN)
```json
{
  "command": "create",
  "args": ["--type", "patch"],
  "flags": {
    "dry-run": true,
    "verbose": false
  },
  "context": {
    "working_dir": "/path/to/project",
    "user": "bsenekowitsch"
  }
}
```

#### Response (STDOUT)
```json
{
  "status": "success",
  "metadata": {
    "plugin": "release",
    "version": "1.0.0",
    "command": "create",
    "timestamp": "2026-02-02T10:30:00Z"
  },
  "data": {
    "type": "release_info",
    "release": {
      "version": "1.2.3",
      "previous_version": "1.2.2",
      "changes": [
        {
          "type": "feature",
          "message": "Add new plugin system",
          "commit": "abc123"
        }
      ],
      "files_modified": [
        "go.mod",
        "pkg/version/version.go"
      ]
    }
  },
  "renderer_hint": "table"
}
```

#### Error Response
```json
{
  "status": "error",
  "metadata": {
    "plugin": "release",
    "version": "1.0.0",
    "command": "create"
  },
  "error": {
    "code": "INVALID_VERSION",
    "message": "Version 1.2.3 already exists",
    "details": {
      "existing_tag": "v1.2.3",
      "created_at": "2025-01-15"
    }
  }
}
```

## Implementierung

### 1. Plugin Interface

```go
// pkg/plugin/types.go
package plugin

import "time"

// Request wird an das Plugin gesendet
type Request struct {
    Command string            `json:"command"`
    Args    []string          `json:"args"`
    Flags   map[string]any    `json:"flags"`
    Context Context           `json:"context"`
}

type Context struct {
    WorkingDir string `json:"working_dir"`
    User       string `json:"user"`
    Verbose    bool   `json:"verbose"`
}

// Response kommt vom Plugin zurück
type Response struct {
    Status       string                 `json:"status"` // "success" | "error"
    Metadata     ResponseMetadata       `json:"metadata"`
    Data         map[string]any         `json:"data,omitempty"`
    Error        *ResponseError         `json:"error,omitempty"`
    RendererHint string                 `json:"renderer_hint,omitempty"`
}

type ResponseMetadata struct {
    Plugin    string    `json:"plugin"`
    Version   string    `json:"version"`
    Command   string    `json:"command"`
    Timestamp time.Time `json:"timestamp"`
}

type ResponseError struct {
    Code    string         `json:"code"`
    Message string         `json:"message"`
    Details map[string]any `json:"details,omitempty"`
}

// Plugin Interface für externe Plugins
type Plugin interface {
    // Execute führt den Plugin-Befehl aus
    Execute(req Request) (*Response, error)
    
    // Manifest gibt Plugin-Metadaten zurück
    Manifest() Manifest
}

type Manifest struct {
    Name          string        `json:"name"`
    Version       string        `json:"version"`
    Description   string        `json:"description"`
    Author        string        `json:"author"`
    Commands      []Command     `json:"commands"`
    RendererTypes []string      `json:"renderer_types"`
}

type Command struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Outputs     []string `json:"outputs"`
}
```

### 2. Plugin Dispatcher

```go
// pkg/dispatcher/dispatcher.go
package dispatcher

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    
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

// Dispatch führt ein Plugin aus
func (d *Dispatcher) Dispatch(ctx context.Context, pluginName string, req plugin.Request) (*plugin.Response, error) {
    // 1. Plugin-Binary finden
    pluginPath, err := d.findPlugin(pluginName)
    if err != nil {
        return nil, fmt.Errorf("plugin not found: %w", err)
    }
    
    // 2. Request als JSON serialisieren
    reqJSON, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    // 3. Plugin-Prozess starten
    cmd := exec.CommandContext(ctx, pluginPath)
    cmd.Stdin = bytes.NewReader(reqJSON)
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    // 4. Plugin ausführen
    if err := cmd.Run(); err != nil {
        return nil, fmt.Errorf("plugin execution failed: %w\nStderr: %s", err, stderr.String())
    }
    
    // 5. Response parsen
    var resp plugin.Response
    if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
        return nil, fmt.Errorf("failed to parse plugin response: %w\nOutput: %s", err, stdout.String())
    }
    
    return &resp, nil
}

func (d *Dispatcher) findPlugin(name string) (string, error) {
    pluginPath := filepath.Join(d.pluginDir, name, fmt.Sprintf("plugin-%s", name))
    
    if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
        return "", fmt.Errorf("plugin '%s' not found at %s", name, pluginPath)
    }
    
    return pluginPath, nil
}

// ListPlugins listet alle installierten Plugins
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
```

### 3. Renderer System

```go
// pkg/renderer/renderer.go
package renderer

import (
    "encoding/json"
    "fmt"
    "io"
    
    "github.com/nekoman-hq/neko-cli/pkg/plugin"
)

type Renderer interface {
    Render(resp *plugin.Response, w io.Writer) error
}

// RendererFactory erstellt Renderer basierend auf Type
type RendererFactory struct {
    renderers map[string]Renderer
}

func NewRendererFactory() *RendererFactory {
    rf := &RendererFactory{
        renderers: make(map[string]Renderer),
    }
    
    // Standard-Renderer registrieren
    rf.Register("json", &JSONRenderer{})
    rf.Register("table", &TableRenderer{})
    rf.Register("text", &TextRenderer{})
    
    return rf
}

func (rf *RendererFactory) Register(name string, renderer Renderer) {
    rf.renderers[name] = renderer
}

func (rf *RendererFactory) Get(name string) (Renderer, error) {
    r, ok := rf.renderers[name]
    if !ok {
        return nil, fmt.Errorf("renderer '%s' not found", name)
    }
    return r, nil
}

// JSONRenderer gibt die Response als JSON aus
type JSONRenderer struct{}

func (r *JSONRenderer) Render(resp *plugin.Response, w io.Writer) error {
    encoder := json.NewEncoder(w)
    encoder.SetIndent("", "  ")
    return encoder.Encode(resp)
}

// TableRenderer formatiert die Response als Tabelle
type TableRenderer struct{}

func (r *TableRenderer) Render(resp *plugin.Response, w io.Writer) error {
    if resp.Status == "error" {
        return r.renderError(resp, w)
    }
    
    // Generische Tabellen-Darstellung basierend auf Data
    // Dies kann pro Plugin-Type spezialisiert werden
    return r.renderData(resp.Data, w)
}

func (r *TableRenderer) renderError(resp *plugin.Response, w io.Writer) error {
    fmt.Fprintf(w, "Error: %s\n", resp.Error.Message)
    fmt.Fprintf(w, "Code: %s\n", resp.Error.Code)
    if len(resp.Error.Details) > 0 {
        fmt.Fprintf(w, "\nDetails:\n")
        for k, v := range resp.Error.Details {
            fmt.Fprintf(w, "  %s: %v\n", k, v)
        }
    }
    return nil
}

func (r *TableRenderer) renderData(data map[string]any, w io.Writer) error {
    // Einfache Key-Value Darstellung
    for k, v := range data {
        fmt.Fprintf(w, "%s: %v\n", k, v)
    }
    return nil
}

// TextRenderer gibt eine einfache Text-Repräsentation aus
type TextRenderer struct{}

func (r *TextRenderer) Render(resp *plugin.Response, w io.Writer) error {
    if resp.Status == "error" {
        fmt.Fprintf(w, "Error: %s\n", resp.Error.Message)
        return nil
    }
    
    fmt.Fprintf(w, "Success\n")
    return nil
}
```

### 4. Cobra Integration

```go
// cmd/plugin.go
package cmd

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    
    "github.com/nekoman-hq/neko-cli/pkg/dispatcher"
    "github.com/nekoman-hq/neko-cli/pkg/plugin"
    "github.com/nekoman-hq/neko-cli/pkg/renderer"
    "github.com/spf13/cobra"
)

var (
    outputFormat string
    pluginDir    string
)

func init() {
    rootCmd.PersistentFlags().StringVar(&outputFormat, "output", "table", "Output format (json, table, text)")
    
    // Plugin-Verzeichnis aus Home-Dir oder Env-Var
    home, _ := os.UserHomeDir()
    defaultPluginDir := filepath.Join(home, ".neko", "plugins")
    pluginDir = os.Getenv("NEKO_PLUGIN_DIR")
    if pluginDir == "" {
        pluginDir = defaultPluginDir
    }
}

// CreatePluginCommand erstellt einen dynamischen Command für ein Plugin
func CreatePluginCommand(manifest plugin.Manifest) *cobra.Command {
    cmd := &cobra.Command{
        Use:   manifest.Name,
        Short: manifest.Description,
        RunE: func(cmd *cobra.Command, args []string) error {
            return executePlugin(manifest.Name, cmd, args)
        },
    }
    
    // Subcommands für jedes Plugin-Command
    for _, pluginCmd := range manifest.Commands {
        subCmd := &cobra.Command{
            Use:   pluginCmd.Name,
            Short: pluginCmd.Description,
            RunE: func(cmd *cobra.Command, args []string) error {
                return executePlugin(manifest.Name, cmd, args)
            },
        }
        cmd.AddCommand(subCmd)
    }
    
    return cmd
}

func executePlugin(pluginName string, cmd *cobra.Command, args []string) error {
    // 1. Dispatcher erstellen
    d := dispatcher.NewDispatcher(pluginDir)
    
    // 2. Request vorbereiten
    req := plugin.Request{
        Command: cmd.Name(),
        Args:    args,
        Flags:   extractFlags(cmd),
        Context: plugin.Context{
            WorkingDir: mustGetwd(),
            User:       os.Getenv("USER"),
            Verbose:    verbose,
        },
    }
    
    // 3. Plugin ausführen
    ctx := context.Background()
    resp, err := d.Dispatch(ctx, pluginName, req)
    if err != nil {
        return fmt.Errorf("failed to execute plugin: %w", err)
    }
    
    // 4. Response rendern
    rf := renderer.NewRendererFactory()
    
    // Renderer-Hint vom Plugin oder Flag verwenden
    rendererType := outputFormat
    if resp.RendererHint != "" && outputFormat == "table" {
        rendererType = resp.RendererHint
    }
    
    r, err := rf.Get(rendererType)
    if err != nil {
        return err
    }
    
    return r.Render(resp, os.Stdout)
}

func extractFlags(cmd *cobra.Command) map[string]any {
    flags := make(map[string]any)
    
    cmd.Flags().VisitAll(func(flag *pflag.Flag) {
        if flag.Changed {
            flags[flag.Name] = flag.Value.String()
        }
    })
    
    return flags
}

func mustGetwd() string {
    wd, _ := os.Getwd()
    return wd
}

// InitializePlugins lädt alle Plugins und registriert ihre Commands
func InitializePlugins() error {
    d := dispatcher.NewDispatcher(pluginDir)
    
    manifests, err := d.ListPlugins()
    if err != nil {
        return fmt.Errorf("failed to list plugins: %w", err)
    }
    
    for _, manifest := range manifests {
        cmd := CreatePluginCommand(manifest)
        rootCmd.AddCommand(cmd)
    }
    
    return nil
}
```

### 5. Plugin-Implementierung (Beispiel Release)

```go
// plugin/release/main.go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
    
    "github.com/nekoman-hq/neko-cli/pkg/plugin"
)

func main() {
    // Request von STDIN lesen
    var req plugin.Request
    if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
        writeError("PARSE_ERROR", fmt.Sprintf("failed to parse request: %v", err))
        os.Exit(1)
    }
    
    // Command ausführen
    var resp *plugin.Response
    var err error
    
    switch req.Command {
    case "init":
        resp, err = handleInit(req)
    case "create":
        resp, err = handleCreate(req)
    case "history":
        resp, err = handleHistory(req)
    default:
        resp, err = nil, fmt.Errorf("unknown command: %s", req.Command)
    }
    
    if err != nil {
        writeError("EXECUTION_ERROR", err.Error())
        os.Exit(1)
    }
    
    // Response als JSON ausgeben
    if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
        writeError("RESPONSE_ERROR", fmt.Sprintf("failed to encode response: %v", err))
        os.Exit(1)
    }
}

func handleInit(req plugin.Request) (*plugin.Response, error) {
    // Release-System initialisieren
    
    return &plugin.Response{
        Status: "success",
        Metadata: plugin.ResponseMetadata{
            Plugin:    "release",
            Version:   "1.0.0",
            Command:   "init",
            Timestamp: time.Now(),
        },
        Data: map[string]any{
            "message": "Release system initialized",
            "config_file": ".neko-release.json",
        },
        RendererHint: "text",
    }, nil
}

func handleCreate(req plugin.Request) (*plugin.Response, error) {
    // Neue Release erstellen
    
    return &plugin.Response{
        Status: "success",
        Metadata: plugin.ResponseMetadata{
            Plugin:    "release",
            Version:   "1.0.0",
            Command:   "create",
            Timestamp: time.Now(),
        },
        Data: map[string]any{
            "type": "release_info",
            "version": "1.2.3",
            "previous_version": "1.2.2",
            "changes": []map[string]string{
                {
                    "type": "feature",
                    "message": "Add plugin system",
                },
            },
        },
        RendererHint: "table",
    }, nil
}

func handleHistory(req plugin.Request) (*plugin.Response, error) {
    // Release-Historie anzeigen
    
    return &plugin.Response{
        Status: "success",
        Metadata: plugin.ResponseMetadata{
            Plugin:    "release",
            Version:   "1.0.0",
            Command:   "history",
            Timestamp: time.Now(),
        },
        Data: map[string]any{
            "releases": []map[string]any{
                {
                    "version": "1.2.2",
                    "date": "2025-01-15",
                    "author": "bsenekowitsch",
                },
                {
                    "version": "1.2.1",
                    "date": "2025-01-10",
                    "author": "bsenekowitsch",
                },
            },
        },
        RendererHint: "table",
    }, nil
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
    json.NewEncoder(os.Stdout).Encode(resp)
}
```

## Verzeichnisstruktur

```
neko-cli/
├── cmd/
│   ├── root.go              # Root Command
│   └── plugin.go            # Plugin Command Factory
├── pkg/
│   ├── plugin/
│   │   └── types.go         # Plugin Interfaces & Types
│   ├── dispatcher/
│   │   └── dispatcher.go    # Plugin Dispatcher
│   └── renderer/
│       ├── renderer.go      # Renderer Interface
│       ├── json.go          # JSON Renderer
│       ├── table.go         # Table Renderer
│       └── text.go          # Text Renderer
├── plugin/
│   ├── release/
│   │   ├── main.go          # Plugin Binary
│   │   ├── manifest.json    # Plugin Manifest
│   │   └── Makefile         # Build Script
│   ├── deploy/
│   │   ├── main.go
│   │   ├── manifest.json
│   │   └── Makefile
│   └── core/
│       ├── main.go
│       ├── manifest.json
│       └── Makefile
├── main.go
├── go.mod
└── Makefile
```

## Build & Installation

### Makefile für Plugins

```makefile
# plugin/release/Makefile
PLUGIN_NAME := release
BINARY := plugin-$(PLUGIN_NAME)
INSTALL_DIR := $(HOME)/.neko/plugins/$(PLUGIN_NAME)

.PHONY: build install clean

build:
	go build -o $(BINARY) main.go

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY) $(INSTALL_DIR)/
	cp manifest.json $(INSTALL_DIR)/

clean:
	rm -f $(BINARY)
	rm -rf $(INSTALL_DIR)
```

### Haupt-Makefile

```makefile
# Makefile (Root)
.PHONY: build install-plugins clean test

build:
	go build -o neko main.go

install-plugins:
	cd plugin/release && $(MAKE) install
	cd plugin/deploy && $(MAKE) install
	cd plugin/core && $(MAKE) install

clean:
	rm -f neko
	cd plugin/release && $(MAKE) clean
	cd plugin/deploy && $(MAKE) clean
	cd plugin/core && $(MAKE) clean

test:
	go test ./...
```

## Vorteile dieser Architektur

1. **Trennung von Logik und Präsentation**: Business Logic in Plugins, UI in Renderern
2. **Erweiterbarkeit**: Neue Plugins ohne Core-Änderungen
3. **Testbarkeit**: Plugins können unabhängig getestet werden
4. **Versionierung**: Jedes Plugin hat seine eigene Version
5. **Polyglot-Fähigkeit**: Plugins können in verschiedenen Sprachen geschrieben werden (solange sie JSON über STDIN/STDOUT unterstützen)
6. **Lazy Loading**: Nur installierte Plugins werden geladen
7. **Isolation**: Plugin-Crashes beeinflussen nicht das CLI-Tool
8. **Standard-Kommunikation**: JSON als lingua franca

## Nächste Schritte

1. ✅ Architektur definieren
2. 🔄 Core Interfaces implementieren (`pkg/plugin`, `pkg/dispatcher`)
3. 🔄 Renderer-System aufbauen
4. 🔄 Beispiel-Plugin (release) portieren
5. 🔄 Plugin-Discovery und Installation
6. 🔄 Testing-Framework
7. 🔄 Documentation & Best Practices