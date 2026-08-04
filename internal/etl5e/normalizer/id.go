package normalizer

import (
	"strings"
	"unicode"
)

// MakeID converts a name to a URL-safe slug.
func MakeID(name string) string {
	var sb strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			sb.WriteRune(unicode.ToLower(r))
		} else if r == ' ' || r == '-' || r == '_' {
			sb.WriteRune('-')
		}
	}
	// Collapse multiple dashes
	result := strings.ReplaceAll(sb.String(), "--", "-")
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}
	return strings.Trim(result, "-")
}

// isCanonicalSource returns true for sources that should use base IDs without suffix.
// These are the canonical 2024 (5.5e) sources that should use simple slugs.
func isCanonicalSource(source string) bool {
	switch source {
	case "PHB", "XPHB", "XMM", "XDMG", "TDCSR", "BMT", "PHB24", "XGTE", "TCE", "MTF", "VRGR", "IDROTF", "EGW", "DSODQ", "ERLW", "MPMM", "FTD", "GGR", "AAG", "SCC", "SATO", "SATO2", "SATO3", "SATO4", "SATO5", "SATO6", "SATO7", "SATO8", "SATO9", "SATO10":
		return true
	}
	return false
}

func MakeClassID(name, source string) string {
	base := MakeID(name)
	if isCanonicalSource(source) {
		return base
	}
	return base + "-" + strings.ToLower(source)
}

func MakeSubclassID(name, source string) string {
	base := MakeID(name)
	if isCanonicalSource(source) {
		return base
	}
	return base + "-" + strings.ToLower(source)
}

func MakeSpeciesID(name, source string) string {
	base := MakeID(name)
	if isCanonicalSource(source) {
		return base
	}
	return base + "-" + strings.ToLower(source)
}

func MakeBackgroundID(name, source string) string {
	base := MakeID(name)
	if isCanonicalSource(source) {
		return base
	}
	return base + "-" + strings.ToLower(source)
}

func MakeFeatureID(prefix, name string) string {
	return MakeID(prefix + "-" + name)
}