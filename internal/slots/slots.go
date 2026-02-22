package slots

type Slots struct {
	Language string

	Intent string
	Topic  string
	Stage  string

	Entities []Entity

	Audience string
	Depth    string
	Style    string
	Format   string
}

type Entity struct {
	Text string
	Role string
}
