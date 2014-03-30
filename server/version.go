package hustle

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	// VersionString contains the git description
	VersionString string
	// RevisionString contains the full git revision
	RevisionString string
	// BuildTags contains all build tags provided at compile time
	BuildTags string
	// VersionPlusJSON contains version, revision, and build tag metadata
	// as a json string
	VersionPlusJSON string
)

func init() {
	VersionPlusJSON = versionPlusJSON()
}

type versionPlus struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Revision  string `json:"revision"`
	BuildTags string `json:"buildTags"`
}

func versionPlusJSON() string {
	vp := &versionPlus{
		Name:      "hustle",
		Version:   VersionString,
		Revision:  RevisionString,
		BuildTags: BuildTags,
	}

	jsonBytes, err := json.MarshalIndent(vp, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "hustle:ERROR: %v\n", err)
		return ""
	}

	return string(jsonBytes)
}
