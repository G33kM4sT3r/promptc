package i18n

import (
	"time"

	"github.com/goodsign/monday"
)

// Predefined layout keys used in translation YAML files.
const (
	KeyFormatDate     = "format.date"
	KeyFormatDateTime = "format.datetime"
	KeyFormatTime     = "format.time"
)

// Fallback layouts (ISO 8601 style) used when no translation key is defined.
const (
	fallbackDate     = "2006-01-02"
	fallbackDateTime = "2006-01-02 15:04"
	fallbackTime     = "15:04"
)

// langToLocale maps ISO 639-1 language codes to monday locales.
var langToLocale = map[string]monday.Locale{
	"en": monday.LocaleEnUS,
	"de": monday.LocaleDeDE,
	"fr": monday.LocaleFrFR,
	"es": monday.LocaleEsES,
	"it": monday.LocaleItIT,
	"pt": monday.LocalePtBR,
	"nl": monday.LocaleNlBE,
	"ru": monday.LocaleRuRU,
	"pl": monday.LocalePlPL,
	"sv": monday.LocaleSvSE,
	"da": monday.LocaleDaDK,
	"fi": monday.LocaleFiFI,
	"nb": monday.LocaleNbNO,
	"nn": monday.LocaleNnNO,
	"tr": monday.LocaleTrTR,
	"hu": monday.LocaleHuHU,
	"ro": monday.LocaleRoRO,
	"bg": monday.LocaleBgBG,
	"uk": monday.LocaleUkUA,
	"cs": monday.LocaleCsCZ,
	"el": monday.LocaleElGR,
	"ja": monday.LocaleJaJP,
	"zh": monday.LocaleZhCN,
	"ko": monday.LocaleKoKR,
}

// localeFor returns the monday locale for a language code, defaulting to en_US.
func localeFor(lang string) monday.Locale {
	if loc, ok := langToLocale[lang]; ok {
		return loc
	}
	return monday.LocaleEnUS
}

// FormatDate formats a time using the translator's date layout and locale.
func (t *Translator) FormatDate(tm time.Time) string {
	layout := t.layoutOrFallback(KeyFormatDate, fallbackDate)
	return monday.Format(tm, layout, localeFor(t.lang))
}

// FormatDateTime formats a time using the translator's datetime layout and locale.
func (t *Translator) FormatDateTime(tm time.Time) string {
	layout := t.layoutOrFallback(KeyFormatDateTime, fallbackDateTime)
	return monday.Format(tm, layout, localeFor(t.lang))
}

// FormatTime formats a time using the translator's time layout and locale.
func (t *Translator) FormatTime(tm time.Time) string {
	layout := t.layoutOrFallback(KeyFormatTime, fallbackTime)
	return monday.Format(tm, layout, localeFor(t.lang))
}

// layoutOrFallback returns the translator value for key, or the fallback
// if the key resolves to itself (meaning it was not found in translations).
func (t *Translator) layoutOrFallback(key, fallback string) string {
	v := t.Get(key)
	if v == key {
		return fallback
	}
	return v
}
