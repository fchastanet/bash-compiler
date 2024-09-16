package model

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
)

const (
	extendsKeyword = "extends"
)

type circularDependencyError struct {
	error
	dependency string
	oldParent  string
	newParent  string
}

func (e *circularDependencyError) Error() string {
	return fmt.Sprintf(
		"Circular dependency '%s' - loaded twice by '%s' and '%s'",
		e.dependency, e.oldParent, e.newParent,
	)
}

type invalidFormatError struct {
	error
	element string
	message string
}

func (e *invalidFormatError) Error() string {
	return fmt.Sprintf(
		"Element '%s' - invalid format '%s'",
		e.element, e.message,
	)
}

func loadModel(
	referenceDir string,
	modelFilePath string,
	resultMap *map[string]any,
	loadedFiles *map[string]string,
	parentFile string,
) (err error) {
	if previousParentFile, exists := (*loadedFiles)[modelFilePath]; exists {
		return &circularDependencyError{nil, modelFilePath, previousParentFile, parentFile}
	}
	// read one yaml file
	referenceDirs := []string{referenceDir}
	model := map[string]any{}
	err = decodeFile(modelFilePath, referenceDirs, &model)
	if err != nil {
		return err
	}
	(*loadedFiles)[modelFilePath] = parentFile

	if _, ok := model[extendsKeyword]; ok {
		extends, ok := model[extendsKeyword].([]any)
		if !ok {
			return &invalidFormatError{error: nil, element: extendsKeyword, message: "array expected"}
		}

		for _, file := range extends {
			fileAbs, err := getFileAbs(file, referenceDir)
			if err != nil {
				return err
			}
			extendsMap := map[string]any{}
			err = loadModel(referenceDir, fileAbs, &extendsMap, loadedFiles, modelFilePath)
			if err != nil {
				return err
			}
			*resultMap = mergeMaps(resultMap, &extendsMap)
		}
	}

	*resultMap = mergeMaps(resultMap, &model)
	delete(*resultMap, extendsKeyword)

	return nil
}

func getFileAbs(file any, referenceDir string) (string, error) {
	fileAbs := os.ExpandEnv(file.(string))
	slog.Debug("Try expanding vars", "original", file.(string), "expanded", fileAbs)
	if _, err := os.Stat(fileAbs); err != nil {
		fileAbs = filepath.Join(referenceDir, file.(string))
		slog.Debug("Try finding file in referenceDir", "referenceDir", referenceDir, "expanded", fileAbs)
		if _, err := os.Stat(fileAbs); err != nil {
			return "", err
		}
	}
	return fileAbs, nil
}

func decodeFile(
	file string, referenceDirs []string, myMap *map[string]any,
) error {
	fileReader, err := os.Open(file)
	if err != nil {
		return err
	}

	dec := yaml.NewDecoder(fileReader, yaml.ReferenceDirs(referenceDirs...))
	return dec.Decode(myMap)
}

func compareObjects(o1 any, o2 any) int {
	v1, ok1 := o1.(string)
	v2, ok2 := o2.(string)
	if ok1 && ok2 {
		return strings.Compare(v1, v2)
	}

	return 1
}

func mergeMaps(map1 *map[string]any, map2 *map[string]any) map[string]any {
	out := make(map[string]any, len(*map1))
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
		if v2, ok := map2v.(map[string]any); ok { //nolint:gocritic // simpler to write it without switch
			// map2v is a map
			if v1, ok := map1v.(map[string]any); ok {
				// if map1v is a map  too, we merge with map2v
				out[k] = mergeMaps(&v1, &v2)

				continue
			}
		} else if v2, ok := map2v.([]any); ok {
			// map2v is an array, concat
			if out[k] == nil {
				out[k] = []any{}
			}
			out1, ok := out[k].([]any)
			if !ok {
				slog.Debug("Fail to cast - This element should be an array", "elementKey", k)
				continue
			}
			arr := append([]any{}, out1...)
			arr = append(arr, v2...)
			// remove duplicates
			slices.SortFunc(arr, compareObjects)
			arr = slices.CompactFunc(arr, func(o1 any, o2 any) bool {
				return compareObjects(o1, o2) == 0
			})
			out[k] = arr

			continue
		}
		out[k] = map2v
	}

	return out
}
