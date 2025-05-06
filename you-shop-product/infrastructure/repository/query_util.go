package repository

import (
	"fmt"
	"strings"
)

// Conver query with question mark holder into query params.
//
// E.g:
//
// Convert: INSERT INTO some_table (id, value) VALUES (?, ?)
//
// To: INSERT INTO some_table (id, value) VALUES ($1, $2)
func ConvertTemplate(template *string) {
	if template == nil {
		return
	}
	parts := strings.Split(*template, "?")
	var result strings.Builder
	for i, part := range parts {
		result.WriteString(part)
		if i < len(parts) - 1 {
			result.WriteString(fmt.Sprintf("$%d", i+1))
		}
	}
	*template = result.String()
}
