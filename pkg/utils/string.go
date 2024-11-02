package utils

import (
	"strings"
)

func RegexPrepareName(name string) string {
	var builder strings.Builder

	charClasses := map[rune]string{
		'a': "[aáàãâä]",
		'e': "[eéèẽêë]",
		'i': "[iíìĩîï]",
		'o': "[oóòõôö]",
		'u': "[uúùũûü]",
		'n': "[nñ]",
		'c': "[cç]",
		's': "[sśš]",
		'z': "[zźżž]",
	}

	for _, char := range strings.ToLower(name) {
		if class, exists := charClasses[char]; exists {
			builder.WriteString(class)
		} else {
			builder.WriteRune(char)
		}
	}

	return builder.String()
}
