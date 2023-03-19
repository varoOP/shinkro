package mapping

import (
	"fmt"
	"regexp"
	"strconv"
)

func GetMovieMalID(guid string) (int, error) {
	r := regexp.MustCompile(`//(\d+ ?)`)

	if !r.MatchString(guid) {
		return -1, fmt.Errorf("unable to parse GUID: %v", guid)
	}

	m := r.FindStringSubmatch(guid)
	id, err := strconv.Atoi(m[1])
	if err != nil {
		return -1, err
	}

	return id, nil
}
