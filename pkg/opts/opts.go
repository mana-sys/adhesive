package opts

import "strings"

func ParseKeyValueStringsMap(values []string) map[string]*string {
	args := make(map[string]*string, len(values))
	for _, value := range values {
		parts := strings.SplitN(value, "=", 2)
		if len(parts) == 2 {
			args[parts[0]] = &parts[1]
		} else {
			args[parts[0]] = nil
		}
	}

	return args
}
