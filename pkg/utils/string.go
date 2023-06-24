package utils

import "strings"

func GenCodeByName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.Join(strings.Split(strings.Join(strings.Fields(name), " "), " "), "-")

	return name
}
