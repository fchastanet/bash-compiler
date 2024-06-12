package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
)

func main() {
	// declare two map to hold the yaml content
	frameworkConfig := map[string]interface{}{}
	base := map[string]interface{}{}
	currentMap := map[string]interface{}{}

	// read one yaml file
	data, err := os.ReadFile("examples/configReference/frameworkConfig.yaml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, &frameworkConfig); err != nil {
		panic(err)
	}

	// read one yaml file
	data, err = os.ReadFile("examples/configReference/defaultCommand.yaml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, &base); err != nil {
		panic(err)
	}

	// read another yaml file
	data1, _ := os.ReadFile("examples/configReference/shellcheckLint.yaml")
	if err := yaml.Unmarshal(data1, &currentMap); err != nil {
		panic(err)
	}

	// merge both yaml data recursively
	base = mergeMaps(frameworkConfig, base)
	base = mergeMaps(base, currentMap)

	// print merged map
	yamlResult, _ := yaml.Marshal(base)
	fmt.Printf("%s\n", string(yamlResult))
}

func compareObjects(o1 interface{}, o2 interface{}) int {
	v1, ok1 := o1.(string)
	v2, ok2 := o2.(string)
	if ok1 && ok2 {
		return strings.Compare(v1, v2)
	}
	return 1
}

func mergeMaps(map1, map2 map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(map1))
	// copy map1 to out
	for k, v := range map1 {
		out[k] = v
	}
	for k, map2v := range map2 {
		// if key does not exists in map1, just append
		map1v, ok := out[k]
		if !ok {
			out[k] = map2v
			continue
		}
		if v2, ok := map2v.(map[string]interface{}); ok {
			// map2v is a map
			if v1, ok := map1v.(map[string]interface{}); ok {
				// if map1v is a map  too, we merge with map2v
				out[k] = mergeMaps(v1, v2)
				continue
			}
		} else if v2, ok := map2v.([]interface{}); ok {
			// map2v is an array, concat
			if out[k] == nil {
				out[k] = []any{}
			}
			out1 := out[k].([]interface{})
			arr := append([]any{}, out1...)
			arr = append(arr, v2...)
			// remove duplicates
			slices.SortFunc(arr, compareObjects)
			arr = slices.CompactFunc(arr, func(o1 interface{}, o2 interface{}) bool {
				return compareObjects(o1, o2) == 0
			})
			out[k] = arr
			continue
		}
		out[k] = map2v
	}
	return out
}
