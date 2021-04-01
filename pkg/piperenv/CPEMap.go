package piperenv

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// CPEMap represents the common pipelineEnvironment Map
type CPEMap map[string]interface{}

// Flatten replaces keys which contains the '/' character and creates the corresponding maps instead an nests
// them properly
func (c *CPEMap) Flatten() {
	*c = flattenMap(*c, nil)
}

// LoadFromDisk reads the given path from disk and populates it to the CPEMap.
func (c *CPEMap) LoadFromDisk(path string) error {
	resMap, err := dirToMap(path)
	if err != nil {
		return err
	}
	*c = resMap
	return nil
}

// WriteToDisk writes the CPEMap to a disk and uses rootDirectory as the starting point
func (c CPEMap) WriteToDisk(rootDirectory string) error {
	return writeMapToDisk(rootDirectory, c)
}

func writeMapToDisk(rootPath string, m map[string]interface{}) error {
	err := os.MkdirAll(rootPath, 0755)
	if err != nil {
		return err
	}

	for k, v := range m {
		// if v is a map create sub directory
		if vMap, ok := v.(map[string]interface{}); ok {
			err := writeMapToDisk(path.Join(rootPath, k), vMap)
			if err != nil {
				return err
			}
			continue
		}
		err := ioutil.WriteFile(path.Join(rootPath, k), []byte(v.(string)), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func flattenMap(m map[string]interface{}, dest map[string]interface{}) (resMap map[string]interface{}) {
	resMap = dest
	if dest == nil {
		resMap = make(map[string]interface{}, len(m))
	}
	for k, v := range m {
		if vMap, ok := v.(map[string]interface{}); ok {
			if previous, ok := resMap[k].(map[string]interface{}); ok {
				v = flattenMap(vMap, previous)
			} else {
				v = flattenMap(vMap, nil)
			}
		}
		subKeys := strings.Split(k, "/")
		if len(subKeys) == 1 {
			resMap[subKeys[0]] = v
			continue
		}
		cursor := resMap
		for index, subKey := range subKeys {
			if index == len(subKeys)-1 {
				cursor[subKey] = v
				break
			}
			subMap := cursor[subKey]
			if subMap == nil {
				subMap = map[string]interface{}{}
			}
			cursor[subKey] = subMap
			cursor = subMap.(map[string]interface{})
		}
	}
	return
}

func dirToMap(dirPath string) (map[string]interface{}, error) {
	if stat, err := os.Stat(dirPath); err != nil || !stat.IsDir() {
		return nil, fmt.Errorf("stat on '%s' failed. Not a dir?: %w", dirPath, err)
	}

	items, err := ioutil.ReadDir(dirPath)
	dirMap := make(map[string]interface{}, len(items))
	if err != nil {
		return nil, err
	}

	for _, dirItem := range items {
		if dirItem.IsDir() {
			// create a map for the sub directory
			toMap, err := dirToMap(path.Join(dirPath, dirItem.Name()))
			if err != nil {
				return nil, err
			}
			dirMap[dirItem.Name()] = toMap
			continue
		}
		// load file content and store a string
		content, err := ioutil.ReadFile(path.Join(dirPath, dirItem.Name()))
		if err != nil {
			return nil, err
		}
		dirMap[dirItem.Name()] = string(content)
	}
	return dirMap, nil

}
