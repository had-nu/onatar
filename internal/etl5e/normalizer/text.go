package normalizer

import (
	"regexp"
	"strings"
)

var (
	tagRE          = regexp.MustCompile(`\{@(\w+)\s+([^}|]+)(?:\|([^}]+))?\}`)
	diceRE         = regexp.MustCompile(`\{@dice\s+([^}]+)\}`)
	damageRE       = regexp.MustCompile(`\{@damage\s+([^}]+)\}`)
	conditionRE    = regexp.MustCompile(`\{@condition\s+([^}]+)\}`)
	spellRE        = regexp.MustCompile(`\{@spell\s+([^}|]+)(?:\|([^}]+))?\}`)
	itemRE         = regexp.MustCompile(`\{@item\s+([^}|]+)(?:\|([^}]+))?\}`)
	creatureRE     = regexp.MustCompile(`\{@creature\s+([^}]+)\}`)
	tableRE        = regexp.MustCompile(`\{@table\s+([^}|]+)(?:\|([^}]+))?\}`)
	bookRE         = regexp.MustCompile(`\{@book\s+([^}]+)\}`)
	skillRE        = regexp.MustCompile(`\{@skill\s+([^}]+)\}`)
	abilityRE      = regexp.MustCompile(`\{@ability\s+([^}]+)\}`)
	featRE         = regexp.MustCompile(`\{@feat\s+([^}]+)\}`)
	classRE        = regexp.MustCompile(`\{@class\s+([^}]+)\}`)
	subclassRE     = regexp.MustCompile(`\{@subclass\s+([^}]+)\}`)
	raceRE         = regexp.MustCompile(`\{@race\s+([^}]+)\}`)
	backgroundRE   = regexp.MustCompile(`\{@background\s+([^}]+)\}`)
	featRE2        = regexp.MustCompile(`\{@feat\s+([^}]+)\}`)
	actionRE       = regexp.MustCompile(`\{@action\s+([^}]+)\}`)
	damageTypeRE   = regexp.MustCompile(`\{@damage\s+([^}]+)\}`)
)

// CleanTags converts 5etools {@tag ...} syntax to plain text.
func CleanTags(s string) string {
	// Order matters: more specific patterns first
	s = tagRE.ReplaceAllStringFunc(s, func(match string) string {
		// {@type content|display} -> display or content
		matches := tagRE.FindStringSubmatch(match)
		if len(matches) >= 4 && matches[3] != "" {
			return matches[3]
		}
		if len(matches) >= 3 {
			return matches[2]
		}
		return matches[1]
	})

	s = diceRE.ReplaceAllString(s, "$1")
	s = damageRE.ReplaceAllString(s, "$1")
	s = conditionRE.ReplaceAllString(s, "$1")
	s = spellRE.ReplaceAllStringFunc(s, func(match string) string {
		matches := spellRE.FindStringSubmatch(match)
		if len(matches) >= 3 && matches[2] != "" {
			return matches[2]
		}
		if len(matches) >= 2 {
			return matches[1]
		}
		return match
	})
	s = itemRE.ReplaceAllStringFunc(s, func(match string) string {
		matches := itemRE.FindStringSubmatch(match)
		if len(matches) >= 3 && matches[2] != "" {
			return matches[2]
		}
		if len(matches) >= 2 {
			return matches[1]
		}
		return match
	})
	s = creatureRE.ReplaceAllString(s, "$1")
	s = tableRE.ReplaceAllStringFunc(s, func(match string) string {
		matches := tableRE.FindStringSubmatch(match)
		if len(matches) >= 3 && matches[2] != "" {
			return matches[2]
		}
		if len(matches) >= 2 {
			return matches[1]
		}
		return match
	})
	s = bookRE.ReplaceAllString(s, "$1")
	s = skillRE.ReplaceAllString(s, "$1")
	s = abilityRE.ReplaceAllString(s, "$1")
	s = featRE.ReplaceAllString(s, "$1")
	s = classRE.ReplaceAllString(s, "$1")
	s = subclassRE.ReplaceAllString(s, "$1")
	s = raceRE.ReplaceAllString(s, "$1")
	s = backgroundRE.ReplaceAllString(s, "$1")
	s = featRE2.ReplaceAllString(s, "$1")
	s = actionRE.ReplaceAllString(s, "$1")
	s = damageTypeRE.ReplaceAllString(s, "$1")

	// Clean up remaining curly braces
	s = strings.ReplaceAll(s, "{", "")
	s = strings.ReplaceAll(s, "}", "")

	// Normalize whitespace
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// CleanEntries converts 5etools entries array to plain text.
func CleanEntries(entries []any) string {
	var sb strings.Builder
	for i, e := range entries {
		switch v := e.(type) {
		case string:
			sb.WriteString(CleanTags(v))
		case map[string]any:
			if name := getString(v, "name"); name != "" {
				if i > 0 {
					sb.WriteString("\n\n")
				}
				sb.WriteString(name + ": ")
				if entries := getStringSlice(v, "entries"); len(entries) > 0 {
					for _, e := range entries {
						if s, ok := e.(string); ok {
							sb.WriteString("- " + CleanTags(s))
						}
					}
				}
			} else if typ := getString(v, "type"); typ != "" {
				if entries := getStringSlice(v, "entries"); len(entries) > 0 {
					for _, e := range entries {
						if s, ok := e.(string); ok {
							sb.WriteString(CleanTags(s))
						}
					}
				}
			}
		}
		if i < len(entries)-1 {
			sb.WriteString("\n\n")
		}
	}
	return strings.TrimSpace(sb.String())
}

func getString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getStringSlice(m map[string]any, key string) []any {
	if v, ok := m[key]; ok {
		if arr, ok := v.([]any); ok {
			var out []any
			for _, item := range arr {
				if s, ok := item.(string); ok {
					out = append(out, s)
				}
			}
			return out
		}
	}
	return nil
}