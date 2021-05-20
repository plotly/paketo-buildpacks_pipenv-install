package pipenvinstall

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type PipfileParser struct{}

func NewPipfileParser() PipfileParser {
	return PipfileParser{}
}

func (p PipfileParser) ParseVersion(path string) (version string, err error) {

	fp, err := os.Open(filepath.Join(path, "Pipfile"))
	if err != nil {
		return "", err
	}

	var Pipfile struct {
		Requires struct {
			PythonVersion string `toml:"python_version"`
		} `toml:"requires"`
	}

	err = toml.NewDecoder(fp).Decode(&Pipfile)
	if err != nil {
		return "", err
	}

	return Pipfile.Requires.PythonVersion, nil
}
