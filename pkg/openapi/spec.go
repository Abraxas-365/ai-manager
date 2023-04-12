package openapi

import (
	"errors"
	"strings"
)

type Any = interface{}
type Dict = map[string]Any
type List = []Any

type ReducedOpenAPISpec struct {
	Servers     []Dict
	Description string
	Endpoints   []Endpoint
}

type Endpoint struct {
	Name        string
	Description string
	Docs        Dict
}

func DereferenceRefs(specObj Dict, fullSpec Dict) (Any, error) {
	var retrieveRefPath = func(path string, fullSpec Dict) (Dict, error) {
		components := strings.Split(path, "/")
		if components[0] != "#" {
			return nil, errors.New("All $refs I've seen so far are uri fragments (start with hash)")
		}
		out := fullSpec
		for _, component := range components[1:] {
			out = out[component].(Dict)
		}
		return out, nil
	}

	var dereferenceRefs func(obj Any) (Any, error)
	dereferenceRefs = func(obj Any) (Any, error) {
		switch v := obj.(type) {
		case Dict:
			objOut := make(Dict)
			for k, val := range v {
				if k == "$ref" {
					refPath, err := retrieveRefPath(val.(string), fullSpec)
					if err != nil {
						return nil, err
					}
					return dereferenceRefs(refPath)
				} else if _, ok := val.([]Any); ok {
					list := val.([]Any)
					newList := make([]Any, len(list))
					for i, el := range list {
						newEl, err := dereferenceRefs(el)
						if err != nil {
							return nil, err
						}
						newList[i] = newEl
					}
					objOut[k] = newList
				} else if _, ok := val.(Dict); ok {
					dict := val.(Dict)
					newDict, err := dereferenceRefs(dict)
					if err != nil {
						return nil, err
					}
					objOut[k] = newDict
				} else {
					objOut[k] = val
				}
			}
			return objOut, nil
		case []Any:
			list := v
			newList := make([]Any, len(list))
			for i, el := range list {
				newEl, err := dereferenceRefs(el)
				if err != nil {
					return nil, err
				}
				newList[i] = newEl
			}
			return newList, nil
		default:
			return obj, nil
		}
	}

	return dereferenceRefs(specObj)
}

func ReduceEndpointDocs(docs Dict) Dict {
	out := make(Dict)
	if description, ok := docs["description"]; ok {
		out["description"] = description
	}
	if params, ok := docs["parameters"]; ok {
		out["parameters"] = []Dict{}
		for _, param := range params.([]Dict) {
			if required, ok := param["required"]; ok && required.(bool) {
				out["parameters"] = append(out["parameters"].([]Dict), param)
			}
		}
	}
	if responses, ok := docs["responses"]; ok {
		if response, ok := responses.(Dict)["200"]; ok {
			out["responses"] = response
		}
	}
	return out
}

func ReduceOpenAPISpec(spec Dict, dereference bool) (*ReducedOpenAPISpec, error) {
	endpoints := []Endpoint{}
	for route, operation := range spec["paths"].(Dict) {
		for operationName, docs := range operation.(Dict) {
			if operationName == "get" || operationName == "post" {
				name := strings.ToUpper(operationName) + " " + route
				description := docs.(Dict)["description"].(string)
				endpoint := Endpoint{
					Name:        name,
					Description: description,
					Docs:        docs.(Dict),
				}
				if dereference {
					dereferencedDocs, err := DereferenceRefs(docs.(Dict), spec)
					if err != nil {
						return nil, err
					}
					endpoint.Docs = dereferencedDocs.(Dict)
				}
				endpoint.Docs = ReduceEndpointDocs(endpoint.Docs)
				endpoints = append(endpoints, endpoint)
			}
		}
	}

	servers := []Dict{}
	for _, server := range spec["servers"].([]Any) {
		servers = append(servers, server.(Dict))
	}

	return &ReducedOpenAPISpec{
		Servers:     servers,
		Description: spec["info"].(Dict)["description"].(string),
		Endpoints:   endpoints,
	}, nil
}
