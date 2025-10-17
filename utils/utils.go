package utils

import (
	"fmt"
	"strconv"
	"strings"

	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/logging"
)

func IncrementArtifactVersion(artifactVersion string) (string, error) {
	parts := strings.Split(artifactVersion, ".")
	if len(parts) != 3 {
		logging.ERROR("invalid artifactVersion format: " + artifactVersion)
		return "", fmt.Errorf("invalid artifactVersion format: %s", artifactVersion)
	}

	major, err1 := strconv.Atoi(parts[0])
	minor, err2 := strconv.Atoi(parts[1])
	patch, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		logging.ERROR(fmt.Sprintf("failed to parse artifactVersion numbers: %v, %v, %v", err1, err2, err3))
		return "", fmt.Errorf("failed to parse artifactVersion numbers: %v, %v, %v", err1, err2, err3)
	}

	// Increment logic
	if artifactVersion == "0.0.0" {
		// Modify to 1.0.0
		major = 1
		minor = 0
		patch = 0
	} else {
		// Change from 1.x.0 to 1.(x + 1).0
		minor += 1
	}

	// Construct new version string
	newArtifactVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	return newArtifactVersion, nil
}
