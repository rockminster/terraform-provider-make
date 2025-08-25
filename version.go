package main

// Version information that is set via build flags
var (
	// Version is the main version number that is being run at the moment.
	Version = "dev"

	// VersionPrerelease is a pre-release marker for the version.
	VersionPrerelease = ""

	// GitCommit is the git commit that was compiled. This will be filled in by the compiler.
	GitCommit = ""
)
