package rhel

type Contact struct {
	Name  string `json:"name"`
	Team  string `json:"team"`
	Docs  string `json:"docs"`
	Email string `json:"email"`
	Url   string `json:"url"`
}

type System struct {
	Architecture string `json:"architecture"`
	Provider     string `json:"provider"`
	OS           string `json:"os"`
}

type Run struct {
	URL string `json:"url"`
	Log string `json:"log"`
}

type Test struct {
	Category  string   `json:"category"`
	Namespace string   `json:"namespace"`
	TestType  string   `json:"type"`
	Result    string   `json:"result"`
	Runtime   string   `json:"runtime"`
	XunitUrls []string `json:"xunit_urls"`
}

type Error struct {
	Reason string `json:"reason"`
}
