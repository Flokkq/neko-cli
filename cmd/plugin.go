package cmd

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nekoman-hq/neko-cli/pkg/dispatcher"
	"github.com/nekoman-hq/neko-cli/pkg/plugin"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	defaultPluginRegistry = "https://api.github.com/repos/nekoman-hq/neko-cli/releases"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage neko plugins",
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed plugins",
	RunE:  runPluginList,
}

var pluginAvailableCmd = &cobra.Command{
	Use:   "available",
	Short: "List available plugins from the registry",
	RunE:  runPluginAvailable,
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install [plugin-name]",
	Short: "Install a plugin from the registry",
	Args:  cobra.ExactArgs(1),
	RunE:  runPluginInstall,
}

var pluginUninstallCmd = &cobra.Command{
	Use:   "uninstall [plugin-name]",
	Short: "Uninstall a plugin",
	Args:  cobra.ExactArgs(1),
	RunE:  runPluginUninstall,
}

var (
	installVersion string
)

func init() {
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginAvailableCmd)
	pluginCmd.AddCommand(pluginInstallCmd)
	pluginCmd.AddCommand(pluginUninstallCmd)

	pluginInstallCmd.Flags().StringVar(&installVersion, "version", "latest", "Version to install")
}

func runPluginList(cmd *cobra.Command, args []string) error {
	d := dispatcher.NewDispatcher(pluginDir)

	manifests, err := d.ListPlugins()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No plugins installed.")
			fmt.Println("Use 'neko plugin available' to see available plugins.")
			return nil
		}
		return fmt.Errorf("failed to list plugins: %w", err)
	}

	if len(manifests) == 0 {
		fmt.Println("No plugins installed.")
		fmt.Println("Use 'neko plugin available' to see available plugins.")
		return nil
	}

	fmt.Printf("%-15s %-10s %-40s %s\n", "NAME", "VERSION", "DESCRIPTION", "AUTHOR")
	for _, m := range manifests {
		fmt.Printf("%-15s %-10s %-40s %s\n", m.Name, m.Version, truncate(m.Description, 40), m.Author)
	}

	return nil
}

func runPluginAvailable(cmd *cobra.Command, args []string) error {
	plugins, err := fetchAvailablePlugins()
	if err != nil {
		return fmt.Errorf("failed to fetch available plugins: %w", err)
	}

	if len(plugins) == 0 {
		fmt.Println("No plugins available.")
		return nil
	}

	fmt.Printf("%-15s %-15s %s\n", "NAME", "LATEST VERSION", "STATUS")

	// Get installed plugins for comparison
	d := dispatcher.NewDispatcher(pluginDir)
	installedManifests, _ := d.ListPlugins()
	installedMap := make(map[string]string)
	for _, m := range installedManifests {
		installedMap[m.Name] = m.Version
	}

	for _, p := range plugins {
		status := "not installed"
		if v, ok := installedMap[p.Name]; ok {
			if v == p.Version {
				status = "installed"
			} else {
				status = fmt.Sprintf("installed (%s)", v)
			}
		}
		fmt.Printf("%-15s %-15s %s\n", p.Name, p.Version, status)
	}

	return nil
}

func runPluginInstall(cmd *cobra.Command, args []string) error {
	pluginName := args[0]

	fmt.Printf("Installing plugin '%s'...\n", pluginName)

	// Determine version to install
	version := installVersion
	if version == "latest" {
		latestVersion, err := getLatestVersion()
		if err != nil {
			return fmt.Errorf("failed to get latest version: %w", err)
		}
		version = latestVersion
	}

	// Build download URL
	downloadURL, err := getPluginDownloadURL(pluginName, version)
	if err != nil {
		return fmt.Errorf("failed to get download URL: %w", err)
	}

	// Download and extract
	if err := downloadAndInstallPlugin(pluginName, downloadURL); err != nil {
		return fmt.Errorf("failed to install plugin: %w", err)
	}

	fmt.Printf("Plugin '%s' installed successfully!\n", pluginName)
	return nil
}

func runPluginUninstall(cmd *cobra.Command, args []string) error {
	pluginName := args[0]

	installPath := filepath.Join(pluginDir, pluginName)
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		return fmt.Errorf("plugin '%s' is not installed", pluginName)
	}

	if err := os.RemoveAll(installPath); err != nil {
		return fmt.Errorf("failed to uninstall plugin: %w", err)
	}

	fmt.Printf("Plugin '%s' uninstalled successfully!\n", pluginName)
	return nil
}

// AvailablePlugin represents a plugin available in the registry
type AvailablePlugin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func fetchAvailablePlugins() ([]AvailablePlugin, error) {
	latestVersion, err := getLatestVersion()
	if err != nil {
		return nil, err
	}

	// Get release assets
	url := fmt.Sprintf("%s/tags/%s", defaultPluginRegistry, latestVersion)

	resp, err := httpGetWithAuth(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch release: %s", resp.Status)
	}

	var release struct {
		Assets []struct {
			Name string `json:"name"`
		} `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	// Parse plugin names from assets
	pluginMap := make(map[string]bool)
	for _, asset := range release.Assets {
		// Plugin assets follow pattern: plugin-{name}_{OS}_{Arch}.tar.gz
		if strings.HasPrefix(asset.Name, "plugin-") {
			parts := strings.Split(asset.Name, "_")
			if len(parts) >= 1 {
				pluginName := strings.TrimPrefix(parts[0], "plugin-")
				pluginMap[pluginName] = true
			}
		}
	}

	var plugins []AvailablePlugin
	for name := range pluginMap {
		plugins = append(plugins, AvailablePlugin{
			Name:    name,
			Version: latestVersion,
		})
	}

	return plugins, nil
}

func getLatestVersion() (string, error) {
	url := fmt.Sprintf("%s/latest", defaultPluginRegistry)

	resp, err := httpGetWithAuth(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch latest release: %s", resp.Status)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func getPluginDownloadURL(pluginName, version string) (string, error) {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Map arch names to match goreleaser output
	archName := arch
	if arch == "amd64" {
		archName = "x86_64"
	}

	// Capitalize OS name
	caser := cases.Title(language.English)
	osName = caser.String(osName)

	assetName := fmt.Sprintf("plugin-%s_%s_%s.tar.gz", pluginName, osName, archName)

	url := fmt.Sprintf("%s/tags/%s", defaultPluginRegistry, version)
	resp, err := httpGetWithAuth(url)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var release struct {
		Assets []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
			ID                 int    `json:"id"`
		} `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	for _, asset := range release.Assets {
		if asset.Name == assetName {
			return asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("plugin '%s' not found for %s/%s in version %s", pluginName, osName, archName, version)
}

func downloadAndInstallPlugin(pluginName, downloadURL string) error {
	resp, err := httpGetWithAuth(downloadURL)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download plugin: %s", resp.Status)
	}

	// Remove existing plugin directory if it exists
	installPath := filepath.Join(pluginDir, pluginName)
	if err = os.RemoveAll(installPath); err != nil {
		return fmt.Errorf("failed to remove existing plugin: %w", err)
	}

	// Create plugin directory
	if err = os.MkdirAll(installPath, 0755); err != nil {
		return err
	}

	// Extract tar.gz
	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer func(gzr *gzip.Reader) {
		_ = gzr.Close()
	}(gzr)

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Skip empty names or current directory entries
		name := filepath.Clean(header.Name)
		if name == "" || name == "." {
			continue
		}

		// Get just the base name (in case archive has nested structure)
		// This flattens any directory structure in the archive
		baseName := filepath.Base(name)
		target := filepath.Join(installPath, baseName)

		switch header.Typeflag {
		case tar.TypeDir:
			// Skip directories - we already created installPath
			continue
		case tar.TypeReg:
			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", target, err)
			}
			if _, err = io.Copy(f, tr); err != nil {
				_ = f.Close()
				return err
			}
			if err = f.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func httpGetWithAuth(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add GitHub token if available (for private repos)
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	return http.DefaultClient.Do(req)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// GetInstalledPluginManifest returns the manifest for an installed plugin
func GetInstalledPluginManifest(pluginName string) (*plugin.Manifest, error) {
	manifestPath := filepath.Join(pluginDir, pluginName, "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	var manifest plugin.Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}
