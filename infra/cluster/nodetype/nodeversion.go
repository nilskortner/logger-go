package nodetype

import (
	"errors"
	"fmt"
	"loggergo/infra/lang"
	"regexp"
	"strconv"
	"strings"
)

var VERSION_PATTERN *regexp.Regexp = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(-SNAPSHOT)?$`)

const WRONG_VERSION_FORMAT_ERROR_MESSAGE = "The version string must follow the formats: \"major.minor.patch\" or \"major.minor.patch-SNAPSHOT\""
const (
	QUALIFIERRELEASE  int8 = 0
	QUALIFIERSNAPSHOT int8 = 1
)

type NodeVersion struct {
	major         int8
	minor         int8
	patch         int8
	qualififer    int8
	version       string
	versionNumber int
}

func NewNodeVersion(major, minor, patch, qualifier int8, version string) (*NodeVersion, error) {
	if major < 0 || minor < 0 || patch < 0 || qualifier < 0 {
		return nil, fmt.Errorf("wrong version format for node")
	}
	return &NodeVersion{
		major:         major,
		minor:         minor,
		patch:         patch,
		qualififer:    qualifier,
		version:       version,
		versionNumber: lang.FromInt8(major, minor, patch, qualifier),
	}, nil
}

func NewNodeVersionWithoutVersion(major, minor, patch, qualifier int8) (*NodeVersion, error) {
	version := ""
	if qualifier == 0 {
		version = strconv.Itoa(int(major)) + "." + strconv.Itoa(int(minor)) + "." + strconv.Itoa(int(patch))
	} else {
		version = strconv.Itoa(int(major)) + "." + strconv.Itoa(int(minor)) + "." + strconv.Itoa(int(patch)) + "-SNAPSHOT"
	}

	return NewNodeVersion(major, minor, patch, qualifier, version)
}

func Parse(version string) (*NodeVersion, error) {

	matches := VERSION_PATTERN.FindStringSubmatch(version)

	var major int8 = -1
	var minor int8 = -1
	var patch int8 = -1
	var qualifier int8 = -1
	var err error

	if matches != nil {
		str := matches[1]
		major = int8(str[0])
		str = matches[2]
		minor = int8(str[0])
		str = matches[3]
		patch = int8(str[0])
		if matches[4] != "" {
			str = matches[4]
			qualifier, err = parseQualifier(str)
			if err != nil {
				return nil, err
			}
		}
	} else {
		return nil, fmt.Errorf("version does not match the pattern")
	}

	return NewNodeVersionWithoutVersion(major, minor, patch, qualifier)
}

func parseQualifier(str string) (int8, error) {
	qualifiers := map[string]int8{
		"RELEASE":  QUALIFIERRELEASE,
		"SNAPSHOT": QUALIFIERSNAPSHOT,
	}

	for name, id := range qualifiers {
		if strings.HasSuffix(str, name) {
			return id, nil
		}
	}

	return 0, errors.New("invalid qualifier")
}
