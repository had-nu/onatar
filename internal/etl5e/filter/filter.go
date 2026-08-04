package filter

// Filter decides whether a 5etools entity should be imported.
// Return true to include, false to skip.
type Filter func(name, source string) bool

// Filter2024Only includes only sources from D&D 2024 (5.5e).
// Sources: XPHB, XMM, XDMG, TDCSR, BMT, PHB24, etc.
func Filter2024Only(name, source string) bool {
	switch source {
	case "XPHB", "XMM", "XDMG", "TDCSR", "BMT", "PHB24", "XGTE", "TCE", "MTF", "VRGR", "IDROTF", "EGW", "DSODQ", "ERLW", "MPMM", "FTD", "GGR", "AAG", "SCC", "SATO", "SATO2", "SATO3", "SATO4", "SATO5", "SATO6", "SATO7", "SATO8", "SATO9", "SATO10":
		return true
	}
	return false
}

// FilterAll includes all sources (for debugging).
func FilterAll(name, source string) bool {
	return true
}