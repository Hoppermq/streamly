package common

func ContainsKeyword(text, keyword string) bool {
	return len(text) >= len(keyword) &&
		text[:len(keyword)] == keyword ||
		FindSubstring(text, keyword)
}

func FindSubstring(text, substring string) bool {
	for i := 0; i <= len(text)-len(substring); i++ {
		if text[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}

func ApplyStringUpdate(existing, update *string, fieldName string, fields []string) []string {
	if update != nil && *update != *existing {
		*existing = *update
		return append(fields, fieldName)
	}

	return fields
}
