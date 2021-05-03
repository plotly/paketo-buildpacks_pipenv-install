package pipenvinstall

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type PipfileLockParser struct{}

func NewPipfileLockParser() PipfileLockParser {
	return PipfileLockParser{}
}

func (p PipfileLockParser) ParseVersion(path string) (version string, err error) {

	pipFileLock, err := os.Open(filepath.Join(path, "Pipfile.lock"))
	if err != nil {
		return "", err
	}

	var PipfileLock struct {
		Meta struct {
			Requires struct {
				Version string `json:"python_version"`
			} `json:"requires"`
			Sources []struct {
				URL string
			}
		} `json:"_meta"`
		Default map[string]struct {
			Version string
		}
	}

	err = json.NewDecoder(pipFileLock).Decode(&PipfileLock)

	if err != nil {
		return "", err
	}

	return PipfileLock.Meta.Requires.Version, nil
}
