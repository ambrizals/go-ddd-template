# Shared Utilities (pkg)

This directory contains shared utility functions and helpers that are agnostic to the domain and infrastructure.

## Responsibilities
- Provide helper functions (string manipulation, math, etc.).
- Shared constants.
- Utility libraries.

## Example: Adding a new Utility
1. Create a new file `pkg/utils/string_utils.go`.
2. Implement the utility function:

```go
package utils

import "strings"

func ToSlug(s string) string {
    return strings.ToLower(strings.ReplaceAll(s, " ", "-"))
}
```
