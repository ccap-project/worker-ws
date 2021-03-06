/*
 *
 * Copyright (c) 2016, 2017, 2018 Alexandre Biancalana <ale@biancalanas.net>.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the <organization> nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package git

import (
	"errors"
	"fmt"
	"time"

	git2go "github.com/libgit2/git2go"
)

func certificateCheckCallback(cert *git2go.Certificate, valid bool, hostname string) git2go.ErrorCode {
	return 0
}

func credentialsCallback(url string, username string, allowedTypes git2go.CredType) (git2go.ErrorCode, *git2go.Cred) {
	ret, cred := git2go.NewCredSshKey("git", "/Users/ale/.ssh/id_rsa.pub", "/Users/ale/.ssh/id_rsa", "")
	return git2go.ErrorCode(ret), &cred
}

func Clone(dir string, url string, chkcert bool) error {

	cloneOptions := &git2go.CloneOptions{}

	cloneOptions.FetchOptions = &git2go.FetchOptions{
		RemoteCallbacks: git2go.RemoteCallbacks{
			CredentialsCallback:      credentialsCallback,
			CertificateCheckCallback: certificateCheckCallback,
		},
	}
	_, err := git2go.Clone(url, dir, cloneOptions)
	if err != nil {
		return err
	}

	//if err != nil && err.Error() != "remote repository is empty" {
	//	return err
	//}

	return nil
}

func Checkout(dir string, chkcert bool, tag string) error {

	var commit *git2go.Commit
	branchName := "master"

	repo, err := git2go.OpenRepository(dir)

	if err != nil {
		return err
	}

	checkoutOpts := &git2go.CheckoutOpts{
		Strategy: git2go.CheckoutSafe | git2go.CheckoutRecreateMissing | git2go.CheckoutAllowConflicts | git2go.CheckoutUseTheirs,
	}

	/* if tag is defined checkout that tag */
	if len(tag) > 0 {

		// Getting the reference for the remote branch
		remote, err := repo.References.Lookup("refs/tags" + tag)
		if err != nil || remote == nil {
			return fmt.Errorf("Failed to find tag %s, %v", tag, err)
		}
		defer remote.Free()

		// Lookup for commit from remote branch
		commit, err = repo.LookupCommit(remote.Target())
		if err != nil {
			return fmt.Errorf("Failed to find tag %s commit(%s), %v", tag, remote.Target(), err)

		} else if commit == nil {
			return fmt.Errorf("Can't find tag %s", tag)
		}
		defer commit.Free()
	} else {

		// Getting the reference for the remote branch
		remote, err := repo.LookupBranch("origin/"+branchName, git2go.BranchRemote)
		if err != nil {
			return fmt.Errorf("Failed to find remote branch: %s", branchName)
		}
		defer remote.Free()

		// Lookup for commit from remote branch
		commit, err = repo.LookupCommit(remote.Target())
		if err != nil {
			return fmt.Errorf("Failed to find remote branch commit: %s", branchName)
		}
		defer commit.Free()

		localBranch, err := repo.LookupBranch(branchName, git2go.BranchLocal)
		// No local branch, lets create one
		if localBranch == nil || err != nil {
			// Creating local branch
			localBranch, err = repo.CreateBranch(branchName, commit, false)
			if err != nil {
				return fmt.Errorf("Failed to create local branch: %s", branchName)
			}

			// Setting upstream to origin branch
			err = localBranch.SetUpstream("origin/" + branchName)
			if err != nil {
				return fmt.Errorf("Failed to create upstream to origin/%s", branchName)
			}
		}

		if localBranch == nil {
			return errors.New("Error while locating/creating local branch")
		}
		defer localBranch.Free()

		// Getting the tree for the branch
		commit, err = repo.LookupCommit(localBranch.Target())
		if err != nil {
			return fmt.Errorf("Failed to lookup for commit in local branch %s", branchName)
		}
	}

	tree, err := commit.Tree()
	if err != nil {
		return fmt.Errorf("Failed to lookup for tree %s", branchName)
	}

	// Checkout the tree
	err = repo.CheckoutTree(tree, checkoutOpts)
	if err != nil {
		return fmt.Errorf("Failed to checkout tree " + branchName)
	}
	// Setting the Head to point to our branch
	repo.SetHead("refs/heads/" + branchName)
	return nil
}

func Commit(dir string, message string) (commit *git2go.Oid, err error) {

	branchName := "master"

	signature := &git2go.Signature{
		Name:  "Config Manager",
		Email: "c@m.c",
		When:  time.Now(),
	}

	repo, err := git2go.OpenRepository(dir)

	if err != nil {
		return nil, err
	}

	idx, err := repo.Index()
	if err != nil {
		return nil, err
	}

	err = idx.AddAll([]string{}, git2go.IndexAddDefault, nil)
	if err != nil {
		return nil, err
	}

	treeId, err := idx.WriteTree()
	if err != nil {
		return nil, err
	}

	err = idx.Write()
	if err != nil {
		return nil, err
	}

	tree, err := repo.LookupTree(treeId)
	if err != nil {
		return nil, err
	}

	localBranch, err := repo.LookupBranch(branchName, git2go.BranchLocal)
	// No local branch, lets create one
	if localBranch == nil || err != nil {

		commit, err = repo.CreateCommit("HEAD", signature, signature, message, tree)
		if err != nil {
			return nil, err
		}

	} else {

		commitTarget, err := repo.LookupCommit(localBranch.Target())
		if err != nil {
			return nil, err
		}

		commit, err = repo.CreateCommit("HEAD", signature, signature, message, tree, commitTarget)
		if err != nil {
			return nil, err
		}
	}

	return commit, nil
}

func Push(dir string, ref string) error {

	var refs []string
	branchName := "master"

	refs = append(refs, "refs/heads/"+branchName)

	if len(ref) > 0 {
		refs = append(refs, "refs/tags/"+ref)
	}

	pushOptions := &git2go.PushOptions{
		RemoteCallbacks: git2go.RemoteCallbacks{
			CredentialsCallback:      credentialsCallback,
			CertificateCheckCallback: certificateCheckCallback,
		},
	}

	repo, err := git2go.OpenRepository(dir)

	if err != nil {
		return err
	}

	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		remote, err = repo.Remotes.Create("origin", repo.Path())
		if err != nil {
			return err
		}
	}

	err = remote.Push(refs, pushOptions)
	if err != nil {
		return err
	}

	return nil
}

func Tag(dir string, commit *git2go.Oid, message string) error {

	repo, err := git2go.OpenRepository(dir)

	if err != nil {
		return err
	}

	commitTarget, err := repo.LookupCommit(commit)
	if err != nil {
		return fmt.Errorf("Creating tag(%s), %v", message, err)
	}

	_, err = repo.Tags.CreateLightweight(message, commitTarget, false)

	if err != nil {
		return fmt.Errorf("Creating tag(%s), %v", message, err)
	}
	return nil
}
