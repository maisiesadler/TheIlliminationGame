package database

var createOverrideFn func(string) ICollection

func tryGetOverrideFor(database string, collection string) (ICollection, bool) {
	key := database + "_" + collection

	if createOverrideFn != nil {
		override := createOverrideFn(key)
		return override, true
	}

	return nil, false
}

// SetOverride allows an override ICollection to be set for testing
func SetOverride(overrideFn func(string) ICollection) {
	createOverrideFn = overrideFn
}
