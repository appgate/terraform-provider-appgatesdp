package codegen

import (
	"regexp"
	"strings"
)

// The rules below mirror TerraformIdentifier in
// github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util,
// which converts OpenAPI property names into attribute names. The package is
// internal to the generator, so it is mirrored here; TestIdentifiersMatchIR
// fails if the two ever disagree.
var (
	lowerToUpper          = regexp.MustCompile(`([a-z])[A-Z]`)
	unsupportedCharacters = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	leadingNumbers        = regexp.MustCompile(`^(\d+)`)
)

// TerraformIdentifier converts an OpenAPI property name to the attribute name
// tfplugingen-openapi generates for it.
func TerraformIdentifier(s string) string {
	s = unsupportedCharacters.ReplaceAllString(s, "")
	s = leadingNumbers.ReplaceAllString(s, "")
	s = lowerToUpper.ReplaceAllStringFunc(s, func(m string) string {
		return m[:1] + "_" + strings.ToLower(m[1:])
	})
	return strings.ToLower(s)
}
