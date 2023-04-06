package locale

// entry is a type that represents an entry.
type entry map[Language]string

// update updates the localized string of the given language.
func (l entry) update(lang Language, value string) {
	l[lang] = value
}

// get returns the localized string of the given language.
func (l entry) get(lang Language) string {
	return l[lang]
}

// newLocale returns a new entry.
func newLocale() *entry {
	return &entry{}
}
