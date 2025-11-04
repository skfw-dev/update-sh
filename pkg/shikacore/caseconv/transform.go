package caseconv

// ToTitleCase ...
func ToTitleCase(s string) string {
	return toMixedCase(s, " ", false)
}

// ToPascalCase ...
func ToPascalCase(s string) string {
	return toMixedCase(s, "", false)
}

// ToCamelCase ...
func ToCamelCase(s string) string {
	return toMixedCase(s, "", true)
}

// ToSnakeCase ...
func ToSnakeCase(s string) string {
	return toUniformCase(s, "_", false)
}

// ToUpperSnakeCase ...
func ToUpperSnakeCase(s string) string {
	return toUniformCase(s, "_", true)
}

// ToKebabCase ...
func ToKebabCase(s string) string {
	return toUniformCase(s, "-", false)
}

// ToUpperKebabCase ...
func ToUpperKebabCase(s string) string {
	return toUniformCase(s, "-", true)
}
