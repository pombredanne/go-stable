package gopkg

import (
	"io"

	"github.com/mxk/go-flowrate/flowrate"
	"gopkg.in/src-d/go-git.v2"
	"gopkg.in/src-d/go-git.v2/clients/common"
)

type Fetcher struct {
	repository *Repository
	remote     *git.Remote
	auth       common.AuthMethod
}

func NewFetcher(r *Repository, auth common.AuthMethod) *Fetcher {
	return &Fetcher{repository: r, auth: auth}
}

func (f *Fetcher) Info() (*common.GitUploadPackInfo, error) {
	var err error
	f.remote, err = git.NewAuthenticatedRemote(f.repository.URL, f.auth)
	if err != nil {
		return nil, err
	}

	if err := f.remote.Connect(); err != nil {
		return nil, err
	}

	return f.remote.Info(), nil
}

func (f *Fetcher) Fetch(w io.Writer) (*flowrate.Status, error) {
	if _, err := f.Info(); err != nil {
		return nil, err
	}

	r, err := f.remote.FetchDefaultBranch()
	if err != nil {
		return nil, err
	}

	flow := flowrate.NewReader(r, -1)
	_, err = io.Copy(w, flow)
	flow.Close()

	status := flow.Status()
	return &status, err
}