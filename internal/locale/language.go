package locale

import (
	"errors"
	"strings"
)

// Language is an enum that represents a language.
type Language int

const (
	English Language = iota
	SimplifiedChinese
	Japanese
	Korean
)

// ErrInvalidLanguage is an error that represents an invalid language.
var ErrInvalidLanguage = errors.New("invalid language")

// ToLanguage converts a string to a Language.
func ToLanguage(lang string) (Language, error) {
	switch strings.ToLower(lang) {
	case "en", "enus", "en-us", "en_us":
		return English, nil
	case "zh", "cn", "zhcn", "zh-cn", "zh_cn":
		return SimplifiedChinese, nil
	case "ja", "jp", "jajp", "ja-jp", "ja_jp":
		return Japanese, nil
	case "ko", "kr", "kokr", "ko-kr", "ko_kr":
		return Korean, nil
	}

	return English, ErrInvalidLanguage
}
