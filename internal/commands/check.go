package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/plexsystems/sinker/internal/docker"

	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCheckCommand(ctx context.Context, logger *log.Logger) *cobra.Command {
	cmd := cobra.Command{
		Use:   "check",
		Short: "Check for newer images in the source registry",

		RunE: func(cmd *cobra.Command, args []string) error {
			if err := viper.BindPFlag("images", cmd.Flags().Lookup("images")); err != nil {
				return fmt.Errorf("bind images flag: %w", err)
			}

			manifestPath := viper.GetString("manifest")
			if err := runCheckCommand(ctx, logger, manifestPath); err != nil {
				return fmt.Errorf("check: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringSliceP("images", "i", []string{}, "The fully qualified images to check if newer versions exist (e.g. myhost.com/myrepo:v1.0.0)")

	return &cmd
}

func runCheckCommand(ctx context.Context, logger *log.Logger, manifestPath string) error {
	client, err := docker.NewClient(logger)
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}

	var imagesToCheck []string
	if len(viper.GetStringSlice("images")) > 0 {
		imagesToCheck = viper.GetStringSlice("images")
	} else {
		manifest, err := GetManifest(manifestPath)
		if err != nil {
			return fmt.Errorf("get manifest: %w", err)
		}

		for _, image := range manifest.Images {
			imagesToCheck = append(imagesToCheck, image.String())
		}
	}

	images := getPathsFromImages(imagesToCheck)
	for _, image := range images {
		if image.Tag() == "" {
			continue
		}

		imageVersion, err := version.NewVersion(image.Tag())
		if err != nil {
			client.Logger.Printf("[CHECK] Image %s version did not parse correctly. Skipping ...", image)
			continue
		}

		tags, err := client.GetTagsForRepo(ctx, image.Host(), image.Repository())
		if err != nil {
			return fmt.Errorf("get tags: %w", err)
		}

		tags = filterTags(tags)

		newerVersions, err := getNewerVersions(imageVersion, tags)
		if err != nil {
			return fmt.Errorf("getting newer version: %w", err)
		}

		if len(newerVersions) == 0 {
			client.Logger.Printf("[CHECK] Image %s is up to date!", image)
			continue
		}

		client.Logger.Printf("[CHECK] New versions for %v found: %v", image, newerVersions)
	}

	return nil
}

func getNewerVersions(currentVersion *version.Version, foundTags []string) ([]string, error) {
	var newerVersions []string
	for _, foundTag := range foundTags {
		tag, err := version.NewVersion(foundTag)
		if err != nil {
			continue
		}

		if currentVersion.LessThan(tag) {
			newerVersions = append(newerVersions, tag.Original())
		}
	}

	if len(newerVersions) > 5 {
		newerVersions = newerVersions[len(newerVersions)-5:]
	}

	return newerVersions, nil
}

func filterTags(tags []string) []string {
	var filteredTags []string
	for _, tag := range tags {
		if strings.Count(tag, ".") > 1 && !strings.Contains(tag, "-") {
			filteredTags = append(filteredTags, tag)
		}
	}

	return filteredTags
}

func getPathsFromImages(images []string) []docker.RegistryPath {
	var paths []docker.RegistryPath
	for _, image := range images {
		paths = append(paths, docker.RegistryPath(image))
	}

	return paths
}
