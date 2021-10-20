package lib

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type Fetcher interface {
	Fetch(packageName string, branchName string) (string, error)
}

func NewFetcher(tmpFolder string) Fetcher {
	return &fetcher{tmpFolder: tmpFolder}
}

type fetcher struct {
	tmpFolder string
}

func (f *fetcher) Fetch(name string, branch string) (string, error) {
	gitAddress := fmt.Sprintf("https://%s", name)
	folder := fmt.Sprintf("%s/%s", f.tmpFolder, name)
	fmt.Println("fetch to", folder, f.tmpFolder, name)
	_, err := git.PlainClone(folder, false, &git.CloneOptions{
		URL:           gitAddress,
		Depth:         1,
		SingleBranch:  true,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
	})

	if err != nil && err != git.ErrRepositoryAlreadyExists {
		go func() {
			if err := os.RemoveAll(folder); err != nil {
				log.WithField("folder", folder).Error(err)
			}
		}()

		return "", err
	}

	return folder, nil
}
