package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const manifestFileName = ".images.yaml"

// Manifest is a collection of images to sync
type Manifest struct {
	Target Target        `yaml:"target"`
	Images []SourceImage `yaml:"sources,omitempty"`
}

// Target ...
type Target struct {
	Path Path `yaml:"path"`
	Auth Auth `yaml:"auth,omitempty"`
}

// Path is the registry and repository to sync images to
type Path string

func (p Path) String() string {
	return string(p)
}

// Host is the host of the registry
func (p Path) Host() string {
	if !strings.Contains(string(p), "/") {
		return ""
	}

	hostTokens := strings.Split(string(p), "/")

	return hostTokens[0]
}

// Repository is the repository of the registry
func (p Path) Repository() string {
	if p.Host() == "" {
		return string(p)
	}

	return strings.ReplaceAll(string(p), p.Host()+"/", "")
}

// SourceImage is a container image
type SourceImage struct {
	Path    Path   `yaml:"path,omitempty"`
	Target  Target `yaml:"target,omitempty"`
	Version string `yaml:"version"`
	Auth    Auth   `yaml:"auth,omitempty"`
}

// String returns the source image
func (c SourceImage) String() string {
	if c.Version == "" {
		return c.Path.String()
	}

	return c.Path.String() + ":" + c.Version
}

// TargetImage returns the target image
func (c SourceImage) TargetImage() string {
	var target string

	if c.Version != "" {
		target = ":" + c.Version
	} else {
		target = ":latest"
	}

	if c.Path.Repository() != "" {
		target = "/" + c.Path.Repository() + target
	}

	if c.Target.Path.Repository() != "" {
		target = "/" + c.Target.Path.Repository() + target
	}

	if c.Target.Path.Host() != "" {
		target = "/" + c.Target.Path.Host() + target
	}

	target = strings.TrimLeft(target, "/")

	return target
}

// Auth is a username and password to log into a registry
type Auth struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// NewManifest returns a new image manifest
func NewManifest(target string) Manifest {
	registry := Target{
		Path: Path(target),
	}

	manifest := Manifest{
		Target: registry,
	}

	return manifest
}

// NewAutodetectManifest returns a new image manifest with images found in the repository
func NewAutodetectManifest(target string, path string) (Manifest, error) {
	manifest := NewManifest(target)

	foundImages, err := getFromKubernetesManifests(path, manifest.Target)
	if err != nil {
		return Manifest{}, fmt.Errorf("get from kubernetes manifests: %w", err)
	}

	manifest.Images = foundImages

	return manifest, nil
}

// GetManifest returns the current manifest file in the working directory
func GetManifest() (Manifest, error) {
	manifestContents, err := ioutil.ReadFile(manifestFileName)
	if err != nil {
		return Manifest{}, fmt.Errorf("reading manifest: %w", err)
	}

	manifest, err := marshalManifest(manifestContents)
	if err != nil {
		return Manifest{}, fmt.Errorf("marshal manifest: %w", err)
	}

	return manifest, nil
}

func marshalManifest(manifestContents []byte) (Manifest, error) {
	var manifest Manifest
	if err := yaml.Unmarshal(manifestContents, &manifest); err != nil {
		return Manifest{}, fmt.Errorf("unmarshal current manifest: %w", err)
	}

	for i := range manifest.Images {
		if manifest.Images[i].Target.Path == "" {
			manifest.Images[i].Target = manifest.Target
		}
	}

	return manifest, nil
}

// WriteManifest writes the image manifest to disk
func WriteManifest(manifest Manifest) error {
	imageManifestContents, err := yaml.Marshal(&manifest)
	if err != nil {
		return fmt.Errorf("marshal image manifest: %w", err)
	}
	imageManifestContents = bytes.ReplaceAll(imageManifestContents, []byte(`"`), []byte(""))

	if err := ioutil.WriteFile(manifestFileName, imageManifestContents, os.ModePerm); err != nil {
		return fmt.Errorf("creating file: %w", err)
	}

	return nil
}
