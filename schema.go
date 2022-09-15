// This file needs to live here, as go embed doesn't allow embedding files in parent directories.

package enably

import _ "embed"

//go:embed schema.toml

// Schema contains information about product categories and their fieldsets.
// It contains the contents of the schema.toml file.
var Schema []byte
