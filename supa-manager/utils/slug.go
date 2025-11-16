package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// GenerateProjectRef creates a URL-friendly project reference from a name
// Format: <slug>-<random6chars>
// Example: "my project" -> "my-project-a1b2c3"
func GenerateProjectRef(name string) string {
	// Create slug from name
	slug := strings.ToLower(name)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Limit length
	if len(slug) > 20 {
		slug = slug[:20]
	}

	// Generate random 6-character suffix
	randomBytes := make([]byte, 3)
	rand.Read(randomBytes)
	randomSuffix := hex.EncodeToString(randomBytes)

	return fmt.Sprintf("%s-%s", slug, randomSuffix)
}
