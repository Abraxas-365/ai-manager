package linkedin

import (
	"github.com/Abraxas-365/ai-manager/pkg/tools"
	"github.com/Abraxas-365/ai-manager/pkg/utils/serpapi"
)

//Its not working good
type LinkednProfileInfo struct {
	name        string
	description string
	serpapi     *serpapi.SerpapiWrapper
}

func NewLinkednProfileInfo() (tools.Tool, error) {

	serpapi, err := serpapi.NewSerpapiWrapper()
	if err != nil {
		return nil, err
	}
	return &LinkednProfileInfo{
		name: "Linkedn Info",
		description: `"Usefull Get the position of a linkedn user".
		"Usefull to get the current job of someone"
		"Input should be search query example 'Luis Fernando Miranda Linkedn'."
		"Priority: 3"
		`,
		serpapi: serpapi,
	}, nil

}

func (s *LinkednProfileInfo) Name() string {
	return s.name
}

func (s *LinkednProfileInfo) Description() string {
	return s.description
}

func (s *LinkednProfileInfo) Run(query string) string {
	result, _ := s.serpapi.SearchTitle(query)
	if len(result) == 0 {
		return "No good Google Search Result was found"
	}
	return result
}
