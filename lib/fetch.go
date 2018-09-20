package lib

import (
	"fmt"
	"os"
	"os/exec"
)

type Fetcher interface {
	Fetch(packageName string) (bool, error)
}

func NewFetcher() Fetcher {
	return new(fetcher)
}

type fetcher struct{}

func (fetcher) packageFound(name string) bool {
	dir := fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), name)
	if _, err := os.Stat(dir); err != nil {
		return false
	}
	return true
}

func (f *fetcher) Fetch(name string) (bool, error) {
	cmd := exec.Command("go", "get", "-d", "-insecure", name)
	cmd.Dir = os.Getenv("GOPATH")
	_, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return f.packageFound(name), err
}
