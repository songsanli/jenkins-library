package piperenv

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func Test_flattenMap(t *testing.T) {
	t.Parallel()
	testMap := CPEMap{
		"A/B": "Hallo",
		"sub": map[string]interface{}{
			"A/B": "Test",
		},
		"A": map[string]interface{}{
			"C": "World",
		},
	}

	testMap.Flatten()
	assert.Equal(t, "Hallo", testMap["A"].(map[string]interface{})["B"])
	assert.IsType(t, map[string]interface{}{}, testMap["sub"])
	subMap := testMap["sub"].(map[string]interface{})
	assert.Equal(t, "Test", subMap["A"].(map[string]interface{})["B"])
	assert.Equal(t, "World", testMap["A"].(map[string]interface{})["C"])

}

func Test_writeMapToDisk(t *testing.T) {
	t.Parallel()
	testMap := CPEMap{
		"A/B": "Hallo",
		"A": map[string]interface{}{
			"C": "CValue",
		},
		"sub": map[string]interface{}{
			"A/B": "Test",
		},
	}

	tmpDir, err := os.MkdirTemp(os.TempDir(), "test-data-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	testMap.Flatten()
	err = testMap.WriteToDisk(tmpDir)
	assert.NoError(t, err)

	testData := []struct {
		Path          string
		ExpectedValue string
	}{
		{
			Path:          "A/B",
			ExpectedValue: "Hallo",
		},
		{
			Path:          "sub/A/B",
			ExpectedValue: "Test",
		},
		{
			Path:          "A/C",
			ExpectedValue: "CValue",
		},
	}

	for _, testCase := range testData {
		t.Run(fmt.Sprintf("check path %s", testCase.Path), func(t *testing.T) {
			tPath := path.Join(tmpDir, testCase.Path)
			bytes, err := ioutil.ReadFile(tPath)
			assert.NoError(t, err)
			assert.Equal(t, testCase.ExpectedValue, string(bytes))
		})
	}
}

func TestCPEMap_LoadFromDisk(t *testing.T) {
	t.Parallel()
	tmpDir, err := os.MkdirTemp(os.TempDir(), "test-data-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	err = ioutil.WriteFile(path.Join(tmpDir, "Foo"), []byte("Bar"), 0644)
	assert.NoError(t, err)
	err = ioutil.WriteFile(path.Join(tmpDir, "Hello"), []byte("World"), 0644)
	assert.NoError(t, err)
	subPath := path.Join(tmpDir, "Batman")
	err = os.Mkdir(subPath, 0744)
	assert.NoError(t, err)
	err = ioutil.WriteFile(path.Join(subPath, "Bruce"), []byte("Wayne"), 0644)
	assert.NoError(t, err)

	cpe := CPEMap{}
	err = cpe.LoadFromDisk(tmpDir)
	assert.NoError(t, err)

	assert.Equal(t, "Bar", cpe["Foo"])
	assert.Equal(t, "World", cpe["Hello"])
	subMap, ok := cpe["Batman"].(map[string]interface{})
	assert.True(t, ok, "map[string]interface{} is expected")
	assert.Equal(t, "Wayne", subMap["Bruce"])
}
