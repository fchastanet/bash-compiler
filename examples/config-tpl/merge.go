package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	sidekick "cuelang.org/go/cmd/cue/cmd"
	"github.com/goccy/go-yaml"
)

//go:embed binFile.cue
var binFileCueSchema string

func main() {
	// declare two map to hold the yaml content

	// compute current directory
	currentDir, err := os.Getwd()
	check(err)

	// read one yaml file
	currentMap := map[string]interface{}{}
	referenceDir := filepath.Join(currentDir, "examples/configReference")
	referenceDirs := []string{referenceDir}
	err = decodeFile(
		"examples/configReference/shellcheckLint.yaml",
		referenceDirs,
		&currentMap,
	)
	check(err)

	resultMap := map[string]interface{}{}
	if _, ok := currentMap["extends"]; ok {
		extends := currentMap["extends"].([]interface{})
		for _, file := range extends {
			fileAbs := filepath.Join(referenceDir, file.(string))
			if _, err := os.Stat(fileAbs); err != nil {
				panic(err)
			}
			extendsMap := map[string]interface{}{}
			err = decodeFile(fileAbs, referenceDirs, &extendsMap)
			if err != nil {
				panic(err)
			}
			resultMap = mergeMaps(resultMap, extendsMap)
		}
	}

	resultMap = mergeMaps(resultMap, currentMap)
	delete(resultMap, "extends")

	// write result to temp file
	tempYamlFile, err := os.CreateTemp("", "config*.yaml")
	check(err)
	defer os.RemoveAll(tempYamlFile.Name())
	yamlResult, _ := yaml.Marshal(resultMap)
	tempYamlFile.Write(yamlResult)
	log.Printf("Temp file containing resulting yaml file : %s\n", tempYamlFile.Name())

	// write cue file to temp file
	tempCueFile, err := os.CreateTemp("", "binFile*.cue")
	check(err)
	defer os.RemoveAll(tempCueFile.Name())
	tempCueFile.Write([]byte(binFileCueSchema))
	log.Printf("Temp file containing cue file : %s\n", tempCueFile.Name())

	// transform using cue
	cmd, err := sidekick.New([]string{
		"export",
		"-l", "input:", tempYamlFile.Name(),
		tempCueFile.Name(),
		"--out", "yaml", "-e", "output",
	})
	check(err)

	// outputs result
	var resultWriter bytes.Buffer
	cmd.SetOutput(&resultWriter)
	err = cmd.Run(cmd.Context())
	check(err)
	fmt.Printf("%s\n", resultWriter.String())
}

func check(e error) {
	if e != nil {
		// notice that we're using 1, so it will actually log where
		// the error happened, 0 = this function, we don't want that.
		_, filename, line, _ := runtime.Caller(1)
		log.Fatal(fmt.Sprintf("[error] %s:%d %v", filename, line, e))
	}
}

func decodeFile(
	file string, referenceDirs []string, myMap *map[string]interface{},
) error {
	fileReader, err := os.Open(file)
	if err != nil {
		return err
	}

	dec := yaml.NewDecoder(fileReader, yaml.ReferenceDirs(referenceDirs...))
	if err := dec.Decode(myMap); err != nil {
		return err
	}
	return nil
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
