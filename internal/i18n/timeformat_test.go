package i18n

import (
	"testing"
	"time"
)

func TestFormatDate_English(t *testing.T) {
	tr := &Translator{
		strings:  map[string]string{"format.date": "Jan 2, 2006"},
		fallback: map[string]string{},
		lang:     "en",
	}
	tm := time.Date(2026, 2, 25, 0, 0, 0, 0, time.UTC)
	got := tr.FormatDate(tm)
	if got != "Feb 25, 2026" {
		t.Errorf("FormatDate en = %q, want %q", got, "Feb 25, 2026")
	}
}

func TestFormatDate_German(t *testing.T) {
	tr := &Translator{
		strings:  map[string]string{"format.date": "2. Jan 2006"},
		fallback: map[string]string{},
		lang:     "de",
	}
	tm := time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC)
	got := tr.FormatDate(tm)
	// monday translates "Apr" to "Apr" in German (same abbreviation)
	if got != "15. Apr 2026" {
		t.Errorf("FormatDate de = %q, want %q", got, "15. Apr 2026")
	}
}

func TestFormatDateTime_English(t *testing.T) {
	tr := &Translator{
		strings:  map[string]string{"format.datetime": "Jan 2, 2006 3:04 PM"},
		fallback: map[string]string{},
		lang:     "en",
	}
	tm := time.Date(2026, 2, 25, 14, 30, 0, 0, time.UTC)
	got := tr.FormatDateTime(tm)
	if got != "Feb 25, 2026 2:30 PM" {
		t.Errorf("FormatDateTime en = %q, want %q", got, "Feb 25, 2026 2:30 PM")
	}
}

func TestFormatDateTime_German(t *testing.T) {
	tr := &Translator{
		strings:  map[string]string{"format.datetime": "2. Jan 2006 15:04"},
		fallback: map[string]string{},
		lang:     "de",
	}
	tm := time.Date(2026, 2, 25, 14, 30, 0, 0, time.UTC)
	got := tr.FormatDateTime(tm)
	if got != "25. Feb 2026 14:30" {
		t.Errorf("FormatDateTime de = %q, want %q", got, "25. Feb 2026 14:30")
	}
}

func TestFormatTime_English(t *testing.T) {
	tr := &Translator{
		strings:  map[string]string{"format.time": "3:04 PM"},
		fallback: map[string]string{},
		lang:     "en",
	}
	tm := time.Date(2026, 1, 1, 9, 5, 0, 0, time.UTC)
	got := tr.FormatTime(tm)
	if got != "9:05 AM" {
		t.Errorf("FormatTime en = %q, want %q", got, "9:05 AM")
	}
}

func TestFormatTime_German(t *testing.T) {
	tr := &Translator{
		strings:  map[string]string{"format.time": "15:04"},
		fallback: map[string]string{},
		lang:     "de",
	}
	tm := time.Date(2026, 1, 1, 9, 5, 0, 0, time.UTC)
	got := tr.FormatTime(tm)
	if got != "09:05" {
		t.Errorf("FormatTime de = %q, want %q", got, "09:05")
	}
}

func TestFormatFallback_NoTranslation(t *testing.T) {
	tr := &Translator{
		strings:  map[string]string{},
		fallback: map[string]string{},
		lang:     "xx",
	}
	tm := time.Date(2026, 2, 25, 14, 30, 0, 0, time.UTC)

	// Should fall back to ISO formats
	if got := tr.FormatDate(tm); got != "2026-02-25" {
		t.Errorf("FormatDate fallback = %q, want %q", got, "2026-02-25")
	}
	if got := tr.FormatDateTime(tm); got != "2026-02-25 14:30" {
		t.Errorf("FormatDateTime fallback = %q, want %q", got, "2026-02-25 14:30")
	}
	if got := tr.FormatTime(tm); got != "14:30" {
		t.Errorf("FormatTime fallback = %q, want %q", got, "14:30")
	}
}

func TestLocaleFor_Known(t *testing.T) {
	tests := map[string]bool{"en": true, "de": true, "fr": true, "es": true}
	for lang := range tests {
		loc := localeFor(lang)
		if loc == "" {
			t.Errorf("localeFor(%q) returned empty", lang)
		}
	}
}

func TestLocaleFor_Unknown(t *testing.T) {
	got := localeFor("xx")
	want := localeFor("en")
	if got != want {
		t.Errorf("localeFor(xx) = %v, want en_US fallback %v", got, want)
	}
}
