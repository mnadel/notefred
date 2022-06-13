package ext

import "strings"

const (
	// X_NOTEFRED is the notefred extension prefix.
	X_NOTEFRED = "x-nf"
)

// CreateKeyValue creates an extension-specific key and value string.
func CreateKeyValue(key, value string) string {
	b := strings.Builder{}

	WriteKeyValue(&b, key, value)

	return b.String()
}

// KeyValue write an extension-specific key and value to the given Builder.
func WriteKeyValue(b *strings.Builder, key string, value string) {
	b.WriteString(X_NOTEFRED)
	b.WriteString("-")
	b.WriteString(key)
	b.WriteString(":")
	b.WriteString(value)
}
