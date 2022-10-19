package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAnalyzer(t *testing.T) {
	t.Run("WithoutIgnoreList", func(t *testing.T) {
		packageName := "github.com/rodrigo-brito/gocity"
		branchName := "master"
		tmpFolder := t.TempDir()
		a := NewAnalyzer(packageName, branchName, tmpFolder)
		rawA, ok := a.(*analyzer)
		assert.True(t, ok)

		assert.Equal(t, rawA.PackageName, packageName)
		assert.Equal(t, rawA.BranchName, branchName)
		assert.Equal(t, rawA.tmpFolder, tmpFolder)
	})

	t.Run("WithIgnoreList", func(t *testing.T) {
		packageName := "github.com/rodrigo-brito/gocity"
		branchName := "master"
		tmpFolder := t.TempDir()
		ignoreList := []string{"/vendor/", "/third-party/", "/external/"}
		a := NewAnalyzer(packageName, branchName, tmpFolder, WithIgnoreList(ignoreList...))
		rawA, ok := a.(*analyzer)
		assert.True(t, ok)

		assert.Equal(t, rawA.PackageName, packageName)
		assert.Equal(t, rawA.BranchName, branchName)
		assert.Equal(t, rawA.tmpFolder, tmpFolder)
		assert.Equal(t, rawA.IgnoreNodes, ignoreList)
	})
}

func TestIsInvalidPath(t *testing.T) {
	a := &analyzer{
		IgnoreNodes: []string{"/vendor/"},
	}

	t.Run("ExactInvalidPath", func(t *testing.T) {
		assert.True(t, a.IsInvalidPath("/vendor/"))
	})

	t.Run("ChildOfInvalidPath", func(t *testing.T) {
		assert.True(t, a.IsInvalidPath("/vendor/github.com/foo/bar"))
	})

	t.Run("ContainsInvalidPath", func(t *testing.T) {
		assert.True(t, a.IsInvalidPath("/pkg/vendor/internal"))
	})

	t.Run("ValidPath", func(t *testing.T) {
		assert.False(t, a.IsInvalidPath("/pkg/lib"))
	})
}

func TestAnalyze(t *testing.T) {
	t.Run("EmptyDir", func(t *testing.T) {
		tmpFolder := t.TempDir()
		a := &analyzer{}
		result, err := a.Analyze(tmpFolder)
		assert.Nil(t, err)
		assert.Zero(t, len(result))
	})

	t.Run("NonexistentDir", func(t *testing.T) {
		tmpFolder := t.TempDir()
		a := &analyzer{}
		_, err := a.Analyze(filepath.Join(tmpFolder, "nonexistent"))
		assert.NotNil(t, err)
	})

	t.Run("NonGoFileOnly", func(t *testing.T) {
		a := &analyzer{}
		result, err := a.Analyze(filepath.Join(testDataDir(), "docs"))
		assert.Nil(t, err)
		assert.Zero(t, len(result))
	})

	t.Run("SingleGoFile", func(t *testing.T) {
		packagePath := filepath.Join(testDataDir(), "subpackage")
		a := NewAnalyzer("github.com/foo/bar", "main", packagePath)
		result, err := a.Analyze(packagePath)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(result))

		assert.Equal(t, expectedAnimal(), result[filepath.Join(packagePath, "example.go")])
		assert.Equal(t, expectedAnimalObj(), result[filepath.Join(packagePath, "example.go.(Animal)")])
	})

	t.Run("MultipleGoFiles", func(t *testing.T) {
		packagePath := filepath.Join(testDataDir())
		a := NewAnalyzer("github.com/foo/bar", "main", packagePath)
		result, err := a.Analyze(packagePath)
		assert.Nil(t, err)
		assert.Equal(t, 5, len(result))

		assert.Equal(t, expectedAnimal(), result[filepath.Join(packagePath, "subpackage", "example.go")])
		assert.Equal(t, expectedAnimalObj(), result[filepath.Join(packagePath, "subpackage", "example.go.(Animal)")])
		assert.Equal(t, expectedPerson(), result[filepath.Join(packagePath, "person.go")])
		assert.Equal(t, expectedPersonObj(), result[filepath.Join(packagePath, "person.go.(Person)")])
		assert.Equal(t, expectedEmployeeObj(), result[filepath.Join(packagePath, "person.go.(Employee)")])
	})
}

func testDataDir() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	project_root := filepath.Dir(filepath.Dir(wd))
	return filepath.Join(project_root, "testdata", "example")
}

func expectedAnimal() *NodeInfo {
	return &NodeInfo{
		NumberLines:      14,
		NumberMethods:    1,
		NumberAttributes: 3,
		Line:             10,
	}
}

func expectedAnimalObj() *NodeInfo {
	return &NodeInfo{
		ObjectName:       "Animal",
		NumberLines:      5,
		NumberMethods:    0,
		NumberAttributes: 3,
		Line:             12,
	}
}

func expectedPerson() *NodeInfo {
	return &NodeInfo{
		NumberLines:      10,
		NumberMethods:    1,
		NumberAttributes: 1,
		Line:             3,
	}
}

func expectedPersonObj() *NodeInfo {
	return &NodeInfo{
		ObjectName:       "Person",
		NumberLines:      4,
		NumberMethods:    0,
		NumberAttributes: 2,
		Line:             5,
	}
}

func expectedEmployeeObj() *NodeInfo {
	return &NodeInfo{
		ObjectName:       "Employee",
		NumberLines:      7,
		NumberMethods:    1,
		NumberAttributes: 2,
		Line:             10,
	}
}
