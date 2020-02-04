package version

import "fmt"

var (
	Version  = "unknown"
	Commit   = "unknown"
	Date     = "unknown"
	Template = fmt.Sprintf("Adhesive version %s, commit %.7s, built on %s\n",
		Version, Commit, Date)
)
