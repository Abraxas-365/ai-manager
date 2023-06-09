package googlesearch

import (
	"fmt"
	"strings"

	"github.com/Abraxas-365/ai-manager/pkg/tools"
	"github.com/Abraxas-365/ai-manager/pkg/utils/serpapi"
)

type SearchTool struct {
	name        string
	description string
	serpapi     *serpapi.SerpapiWrapper
}

func NewSearchTool() (tools.Tool, error) {
	serpapi, err := serpapi.NewSerpapiWrapper()
	if err != nil {
		return nil, err
	}
	return &SearchTool{
		name: "Google Search",
		description: `"A wrapper around Google Search. "
        "Useful for when you need to answer questions about current events. "
		"Always one of the first options when you need to find information on internet"
        "Input should be a search query."
		"Priority: 1"`,
		serpapi: serpapi,
	}, nil

}

func (s *SearchTool) Name() string {
	return s.name
}

func (s *SearchTool) Description() string {
	return s.description
}

func (s *SearchTool) Run(query string) string {
	result, _ := s.serpapi.Search(query)
	if len(result) == 0 {
		return "No good Google Search Result was found"
	}
	fmt.Println("QUERY", query)
	fmt.Println("RESULT", strings.Join(strings.Fields(result), " "))
	return strings.Join(strings.Fields(result), " ")

}
