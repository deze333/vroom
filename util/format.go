package util

import (
    "fmt"
    "bytes"
    "strings"
)

//------------------------------------------------------------
// Formatting utils
//------------------------------------------------------------

// Replaces newlines with <br>.
func Filter_NewlinesToHtml(s string) string {

	// Standard is CRLF
	s = strings.Replace(s, "\r\n", "<br>", -1)

	// If for some reason LFCR strange combo is present
	if strings.Index(s, "\n\r") != -1 {
		s = strings.Replace(s, "\n\r", "<br>", -1)
	}

	// If non standard
	if strings.Index(s, "\n") != -1 {
		s = strings.Replace(s, "\n", "<br>", -1)
	}

	return s
}

// Formats integer into human readable format.
// Example: 123456 --> 123,456
func Fmt_IntAsThousands(v interface{}) string {
    s := fmt.Sprint(v)
    var buf bytes.Buffer
    for i := len(s); i > 3; {
        i -= 3
        buf.WriteString(s[:i])
        buf.WriteByte(',')
        buf.WriteString(s[i:])
    }
    return buf.String()
}
