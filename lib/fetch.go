package lib

import (
	"fmt"
	"log"
	"os"

	git "gopkg.in/src-d/go-git.v4"
)

type Fetcher interface {
	Fetch(packageName string) error
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

func (f *fetcher) Fetch(name string) error {
	gitAddress := fmt.Sprintf("https://%s", name)
	folder := fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), name)

	_, err := git.PlainClone(folder, false, &git.CloneOptions{
		URL:          gitAddress,
		Depth:        1,
		SingleBranch: true,
	})

	if err != nil && err != git.ErrRepositoryAlreadyExists {
		go func() {
			if err := os.RemoveAll(folder); err != nil {
				log.Printf("error on remove: %s", err)
			}
		}()

		return err
	}

	return nil
}
