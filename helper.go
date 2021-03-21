package gotemplater

import (
	"errors"
	"strings"
)

func GetAbsolutePath(fromPath []string, relativePath string) ([]string, error) {

	mergeSegments := append(fromPath)

	segments := strings.Split(relativePath, "/")
	for _, segment := range segments {

		if segment == "." {

		} else if segment == ".." {
			numSegment := len(mergeSegments)
			if numSegment > 0 {
				mergeSegments = mergeSegments[0:numSegment]
			} else {
				return nil, errors.New("path error")
			}
		} else if len(segment) > 0 {
			mergeSegments = append(mergeSegments, segment)
		}
	}

	return mergeSegments, nil
}
