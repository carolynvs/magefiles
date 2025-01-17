package tools

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/carolynvs/magex/mgx"
	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/pkg/archive"
	"github.com/carolynvs/magex/pkg/downloads"
)

const (
	// Version of KIND to install if not already present
	DefaultKindVersion = "v0.10.0"
)

// Fail if the go version doesn't match the specified constraint
// Examples: >=1.16
func EnforceGoVersion(constraint string) {
	log.Printf("Checking go version against constraint %s...", constraint)

	value := strings.TrimPrefix(runtime.Version(), "go")
	version, err := semver.NewVersion(value)
	if err != nil {
		mgx.Must(fmt.Errorf("could not parse go version: '%s': %w", value, err))
	}
	versionCheck, err := semver.NewConstraint(constraint)
	if err != nil {
		mgx.Must(fmt.Errorf("invalid semver constraint: '%s': %w", constraint, err))
	}

	ok, _ := versionCheck.Validate(version)
	if !ok {
		mgx.Must(fmt.Errorf("your version of Go, %s, does not meet the requirement %s", version, versionCheck))
	}
}

// Install mage
func EnsureMage() error {
	return pkg.EnsureMage("")
}

// Install gh
func EnsureGitHubClient() {
	if ok, _ := pkg.IsCommandAvailable("gh", ""); ok {
		return
	}

	// gh cli unfortunately uses a different archive schema depending on the OS
	target := "gh_{{.VERSION}}_{{.GOOS}}_{{.GOARCH}}/bin/gh{{.EXT}}"
	if runtime.GOOS == "windows" {
		target = "bin/gh.exe"
	}

	opts := archive.DownloadArchiveOptions{
		DownloadOptions: downloads.DownloadOptions{
			UrlTemplate: "https://github.com/cli/cli/releases/download/v{{.VERSION}}/gh_{{.VERSION}}_{{.GOOS}}_{{.GOARCH}}{{.EXT}}",
			Name:        "gh",
			Version:     "1.8.1",
			OsReplacement: map[string]string{
				"darwin": "macOS",
			},
		},
		ArchiveExtensions: map[string]string{
			"linux":   ".tar.gz",
			"darwin":  ".tar.gz",
			"windows": ".zip",
		},
		TargetFileTemplate: target,
	}

	err := archive.DownloadToGopathBin(opts)
	mgx.Must(err)
}

// Install kind
func EnsureKind() {
	if ok, _ := pkg.IsCommandAvailable("kind", ""); ok {
		return
	}

	kindURL := "https://github.com/kubernetes-sigs/kind/releases/download/{{.VERSION}}/kind-{{.GOOS}}-{{.GOARCH}}"
	mgx.Must(pkg.DownloadToGopathBin(kindURL, "kind", getKindVersion()))
}

func getKindVersion() string {
	if version, ok := os.LookupEnv("KIND_VERSION"); ok {
		return version
	}
	return DefaultKindVersion
}

// Install the latest version of porter
func EnsurePorter() {
	err := pkg.DownloadToGopathBin("https://cdn.porter.sh/{{.VERSION}}/porter-{{.GOOS}}-{{.GOARCH}}{{.EXT}}", "porter", "latest")
	mgx.Must(err)
}
