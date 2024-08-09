package model

import (
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
)

func loadModel(referenceDir string, modelFilePath string, resultMap *map[string]interface{}) (err error) {
	// read one yaml file
	referenceDirs := []string{referenceDir}
	model := map[string]interface{}{}
	err = decodeFile(modelFilePath, referenceDirs, &model)
	if err != nil {
		return err
	}

	if _, ok := model["extends"]; ok {
		extends := model["extends"].([]interface{})
		for _, file := range extends {
			fileAbs := os.ExpandEnv(file.(string))
			slog.Debug("Try expanding vars", "original", file.(string), "expanded", fileAbs)
			if _, err := os.Stat(fileAbs); err != nil {
				fileAbs = filepath.Join(referenceDir, file.(string))
				slog.Debug("Try finding file in referenceDir", "referenceDir", referenceDir, "expanded", fileAbs)
				if _, err := os.Stat(fileAbs); err != nil {
					return err
				}
			}

			extendsMap := map[string]interface{}{}
			err = decodeFile(fileAbs, referenceDirs, &extendsMap)
			if err != nil {
				return err
			}
			*resultMap = mergeMaps(resultMap, &extendsMap)
		}
	}

	*resultMap = mergeMaps(resultMap, &model)
	delete(*resultMap, "extends")

	return nil
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

func mergeMaps(map1 *map[string]interface{}, map2 *map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(*map1))
	// copy map1 to out
	for k, v := range *map1 {
		out[k] = v
	}
	for k, map2v := range *map2 {
		// if key does not exists in map1, just append
		map1v, ok := out[k]
		if !ok {
			out[k] = map2v

			continue
		}
		if v2, ok := map2v.(map[string]interface{}); ok { //nolint:gocritic // simpler to write it without switch
			// map2v is a map
			if v1, ok := map1v.(map[string]interface{}); ok {
				// if map1v is a map  too, we merge with map2v
				out[k] = mergeMaps(&v1, &v2)

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
