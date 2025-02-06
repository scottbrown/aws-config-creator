package setlist

import (
	"strings"
)

// ParseNicknameMapping parses a nickname mapping string into a map.
// The expected format is "accountID1=nickname1,accountID2=nickname2".
func ParseNicknameMapping(mapping string) map[string]string {
	nicknameMapping := make(map[string]string)

	if len(mapping) == 0 {
		return nicknameMapping
	}

	tokens := strings.Split(mapping, ",")
	for _, token := range tokens {
		parts := strings.Split(token, "=")

		nicknameMapping[parts[0]] = parts[1]
	}

	return nicknameMapping
}
