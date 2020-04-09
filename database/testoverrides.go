package database

var overrides map[string]ICollection

func tryGetOverrideFor(database string, collection string) (ICollection, bool) {
	key := database + "_" + collection

	val, ok := overrides[key]
	return val, ok
}

// SetOverride allows an override ICollection to be set for testing
func SetOverride(database string, collection string, coll ICollection) {
	key := database + "_" + collection

	overrides[key] = coll
}
