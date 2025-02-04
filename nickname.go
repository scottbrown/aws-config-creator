package awsconfigcreator

import (
	"strings"
)

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
