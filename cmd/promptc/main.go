package main

var (
	version = "dev"
	commit  = "unknown"
	date    = ""
)

func main() {
	buildVersion = version
	buildCommit = commit
	buildDate = date
	Execute(version, commit, date)
}
