package locale

// Localizer is an interface that represents a struct that can provide localization strings.
type Localizer interface {
	Fetch(key string, lang Language) string
	Update(key string, lang Language, value string)
}

// localizer is a struct that implements Localizer.
type localizer struct {
	defaultLanguage Language
	locales         map[string]*entry
}

// Fetch returns the localized string of the given key.
func (g *localizer) Fetch(key string, lang Language) string {
	if locale, ok := g.locales[key]; ok {
		return locale.get(lang)
	}

	return ""
}

// Update updates the localized string of the given key.
func (g *localizer) Update(key string, lang Language, value string) {
	locale, ok := g.locales[key]

	if !ok {
		locale = newLocale()
		g.locales[key] = locale
	}

	locale.update(lang, value)
}

// NewLocalizer returns a new Localizer.
func NewLocalizer() Localizer {
	return &localizer{
		locales: make(map[string]*entry),
	}
}
