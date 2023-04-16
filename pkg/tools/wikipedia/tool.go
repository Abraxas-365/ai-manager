package wikipedia

import "github.com/Abraxas-365/ai-manager/pkg/tools"

type WikipediaTool struct {
	name        string
	description string
}

func NewWikipediaTool() tools.Tool {
	return &WikipediaTool{
		name: "Wikipedia",
		description: `
        "A wrapper around Wikipedia. "
        "Useful for when you need to answer general questions about "
        "people, places, companies, historical events, or other subjects. "
        "Input should be a search query."
		"Priority: 1"
		`,
	}

}

func (s *WikipediaTool) Name() string {
	return s.name
}

func (s *WikipediaTool) Description() string {
	return s.description
}

func (s *WikipediaTool) Run(query string) string {
	resutl, err := fetchWikipediaSummary(query)
	if err != nil {
		return "No good Wikipedia Search Result was found"
	}
	return resutl
}
