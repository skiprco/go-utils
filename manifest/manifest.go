package manifest

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/skiprco/go-utils/v2/errors"
)

// ManifestFileName is th file name of the manifest
const ManifestFileName = "manifest.json"

// Manifest contains metadata about the service
type Manifest struct {
	ServiceName string        `json:"service-name"`
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	Depends     []string      `json:"depends"`
	Contracts   []string      `json:"contracts"`
	Description string        `json:"description"`
	Data        []string      `json:"data"`
	Demo        []interface{} `json:"demo"`
	Deployable  string        `json:"deployable"`
}

// LoadManifest loads a manifest file from the current directory
func LoadManifest() (*Manifest, *errors.GenericError) {
	// Setup logging
	abs, _ := filepath.Abs(ManifestFileName)
	manifestLog := log.WithField("manifest_path", abs)

	// Load file
	file, err := ioutil.ReadFile(ManifestFileName)
	if err != nil {
		manifestLog.WithField("error", err).Error("Manifest file not found")
		return nil, errors.NewGenericError(404, "go_utils", "manifest", "manifest_file_not_found", nil)
	}

	// Parse manifest
	manifest := &Manifest{}
	err = json.Unmarshal([]byte(file), manifest)
	if err != nil {
		manifestLog.WithField("error", err).Error("Failed to unmarshal manifest file")
		return nil, errors.NewGenericError(500, "go_utils", "manifest", "unmarshal_manifest_failed", nil)
	}

	// Load manifest successful
	return manifest, nil
}
