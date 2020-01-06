/*
Copyright 2019 BlackRock, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package context

import (
	"fmt"
	"io/ioutil"

	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	go_git_ssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"k8s.io/client-go/kubernetes"
)

const (
	// DefaultRemote is the default remote
	DefaultRemote = "origin"
	// DefaultBranch is the default branch
	DefaultBranch = "master"
)

var (
	fetchRefSpec = []config.RefSpec{
		"refs/*:refs/*",
		"HEAD:refs/heads/HEAD",
	}
)

// GitBuildContextReader is the context reader for pulling context from git
type GitBuildContextReader struct {
	k8sClient    kubernetes.Interface
	buildContext *v1alpha1.GitContext
}

func (contextReader *GitBuildContextReader) getRemote() string {
	if contextReader.buildContext.Remote != nil {
		return contextReader.buildContext.Remote.Name
	}
	return DefaultRemote
}

func getSSHKeyAuth(sshKeyFile string) (transport.AuthMethod, error) {
	var auth transport.AuthMethod
	sshKey, err := ioutil.ReadFile(sshKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read ssh key file. err: %+v", err)
	}
	signer, err := ssh.ParsePrivateKey([]byte(sshKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ssh key. err: %+v", err)
	}
	auth = &go_git_ssh.PublicKeys{User: "git", Signer: signer}
	return auth, nil
}

func (contextReader *GitBuildContextReader) getGitAuth() (transport.AuthMethod, error) {
	if contextReader.buildContext.SSHKeyPath != "" {
		return getSSHKeyAuth(contextReader.buildContext.SSHKeyPath)
	}
	username, err := common.ReadCredentials(contextReader.k8sClient, contextReader.buildContext.Username)
	if err != nil {
		return nil, err
	}
	password, err := common.ReadCredentials(contextReader.k8sClient, contextReader.buildContext.Password)
	if err != nil {
		return nil, err
	}
	if username == "" && password == "" {
		return nil, nil
	}
	return &http.BasicAuth{
		Username: username,
		Password: password,
	}, nil
}

func (contextReader *GitBuildContextReader) pullFromRepository(r *git.Repository) error {
	auth, err := contextReader.getGitAuth()
	if err != nil {
		return err
	}

	if contextReader.buildContext.Remote != nil {
		_, err := r.CreateRemote(&config.RemoteConfig{
			Name: contextReader.buildContext.Remote.Name,
			URLs: contextReader.buildContext.Remote.URLS,
		})
		if err != nil {
			return errors.Errorf("failed to create remote. err: %+v", err)
		}

		fetchOptions := &git.FetchOptions{
			RemoteName: contextReader.buildContext.Remote.Name,
			RefSpecs:   fetchRefSpec,
		}
		if auth != nil {
			fetchOptions.Auth = auth
		}

		if err := r.Fetch(fetchOptions); err != nil {
			return errors.Errorf("failed to fetch remote %s. err: %+v", contextReader.buildContext.Remote.Name, err)
		}
	}

	w, err := r.Worktree()
	if err != nil {
		return errors.Errorf("failed to get working tree. err: %+v", err)
	}

	fetchOptions := &git.FetchOptions{
		RemoteName: contextReader.getRemote(),
		RefSpecs:   fetchRefSpec,
	}
	if auth != nil {
		fetchOptions.Auth = auth
	}

	// In the case of a specific given ref, it isn't necessary to fetch anything
	// but the single ref
	if contextReader.buildContext.Ref != "" {
		fetchOptions.Depth = 1
		fetchOptions.RefSpecs = []config.RefSpec{config.RefSpec(contextReader.buildContext.Ref + ":" + contextReader.buildContext.Ref)}
	}

	if err := r.Fetch(fetchOptions); err != nil && err != git.NoErrAlreadyUpToDate {
		return errors.Errorf("failed to fetch. err: %v", err)
	}

	if err := w.Checkout(contextReader.getBranchOrTag()); err != nil {
		return errors.Errorf("failed to checkout. err: %+v", err)
	}

	// In the case of a specific given ref, it shouldn't be necessary to pull
	if contextReader.buildContext.Ref != "" {
		pullOpts := &git.PullOptions{
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			ReferenceName:     contextReader.getBranchOrTag().Branch,
			Force:             true,
		}
		if auth != nil {
			pullOpts.Auth = auth
		}

		if err := w.Pull(pullOpts); err != nil && err != git.NoErrAlreadyUpToDate {
			return errors.Errorf("failed to pull latest updates. err: %+v", err)
		}
	}
	return nil
}

func (contextReader *GitBuildContextReader) getBranchOrTag() *git.CheckoutOptions {
	opts := &git.CheckoutOptions{}
	opts.Branch = plumbing.NewBranchReferenceName(DefaultBranch)
	if contextReader.buildContext.Branch != "" {
		opts.Branch = plumbing.NewBranchReferenceName(contextReader.buildContext.Branch)
	}
	if contextReader.buildContext.Tag != "" {
		opts.Branch = plumbing.NewTagReferenceName(contextReader.buildContext.Tag)
	}
	if contextReader.buildContext.Ref != "" {
		opts.Branch = plumbing.ReferenceName(contextReader.buildContext.Ref)
	}
	return opts
}

func (contextReader *GitBuildContextReader) Read() (string, error) {
	r, err := git.PlainOpen(common.ContextDirectory)
	if err != nil {
		if err != git.ErrRepositoryNotExists {
			return "", errors.Errorf("failed to open repository. err: %+v", err)
		}

		cloneOpt := &git.CloneOptions{
			URL:               contextReader.buildContext.URL,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		}

		auth, err := contextReader.getGitAuth()
		if err != nil {
			return "", err
		}
		if auth != nil {
			cloneOpt.Auth = auth
		}

		// In the case of a specific given ref, it isn't necessary to have branch
		// histories
		if contextReader.buildContext.Ref != "" {
			cloneOpt.Depth = 1
		}

		r, err = git.PlainClone(common.ContextDirectory, false, cloneOpt)
		if err != nil {
			return "", errors.Errorf("failed to clone repository. err: %+v", err)
		}
	}
	if err := contextReader.pullFromRepository(r); err != nil {
		return "", errors.Errorf("failed to pull latest changes from the repository. err: %+v", err)
	}
	return common.ContextDirectory, nil
}

// NewGitBuildContextReader returns a build context stored on the git
func NewGitBuildContextReader(buildContext *v1alpha1.GitContext, k8sClient kubernetes.Interface) *GitBuildContextReader {
	return &GitBuildContextReader{
		k8sClient,
		buildContext,
	}
}
