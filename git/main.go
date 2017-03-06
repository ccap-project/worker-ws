package git

import (
	"errors"
	"fmt"
	"time"

	git2go "gopkg.in/libgit2/git2go.v24"
)

func certificateCheckCallback(cert *git2go.Certificate, valid bool, hostname string) git2go.ErrorCode {
	return 0
}

func credentialsCallback(url string, username string, allowedTypes git2go.CredType) (git2go.ErrorCode, *git2go.Cred) {
	ret, cred := git2go.NewCredSshKey("git", "/Users/ale/.ssh/id_rsa.pub", "/Users/ale/.ssh/id_rsa", "")
	return git2go.ErrorCode(ret), &cred
}

func CreateBranch(repoPath string, branchName string) error {

	repo, err := git2go.OpenRepository(repoPath)

	if err != nil {
		return err
	}

	localBranch, err := repo.LookupBranch(branchName, git2go.BranchLocal)
	// No local branch, lets create one
	if localBranch == nil || err != nil {
		// Creating local branch
		localBranch, err = repo.CreateBranch(branchName, nil, false)
		if err != nil {
			return fmt.Errorf("Failed to create local branch: %s", branchName)
		}

		// Setting upstream to origin branch
		err = localBranch.SetUpstream("origin/" + branchName)
		if err != nil {
			return fmt.Errorf("Failed to create upstream to origin/%s", branchName)
		}
	}

	return nil
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

func Checkout(dir string, chkcert bool) error {

	branchName := "master"

	repo, err := git2go.OpenRepository(dir)

	if err != nil {
		return err
	}

	checkoutOpts := &git2go.CheckoutOpts{
		Strategy: git2go.CheckoutSafe | git2go.CheckoutRecreateMissing | git2go.CheckoutAllowConflicts | git2go.CheckoutUseTheirs,
	}

	// Getting the reference for the remote branch
	remoteBranch, err := repo.LookupBranch("origin/"+branchName, git2go.BranchRemote)
	if err != nil {
		return fmt.Errorf("Failed to find remote branch: %s", branchName)
	}
	defer remoteBranch.Free()

	// Lookup for commit from remote branch
	commit, err := repo.LookupCommit(remoteBranch.Target())
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
	localCommit, err := repo.LookupCommit(localBranch.Target())
	if err != nil {
		return fmt.Errorf("Failed to lookup for commit in local branch %s", branchName)
	}
	defer localCommit.Free()

	tree, err := repo.LookupTree(localCommit.TreeId())
	if err != nil {
		return fmt.Errorf("Failed to lookup for tree %s", branchName)
	}
	defer tree.Free()

	// Checkout the tree
	err = repo.CheckoutTree(tree, checkoutOpts)
	if err != nil {
		return fmt.Errorf("Failed to checkout tree " + branchName)
	}
	// Setting the Head to point to our branch
	repo.SetHead("refs/heads/" + branchName)
	return nil
}

func CommitAndPush(dir string, message string) error {

	branchName := "master"

	signature := &git2go.Signature{
		Name:  "Config Manager",
		Email: "c@m.c",
		When:  time.Now(),
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

	idx, err := repo.Index()
	if err != nil {
		return err
	}

	err = idx.AddAll([]string{}, git2go.IndexAddDefault, nil)
	if err != nil {
		return err
	}

	treeId, err := idx.WriteTree()
	if err != nil {
		return err
	}

	err = idx.Write()
	if err != nil {
		return err
	}

	tree, err := repo.LookupTree(treeId)
	if err != nil {
		return err
	}

	localBranch, err := repo.LookupBranch(branchName, git2go.BranchLocal)
	// No local branch, lets create one
	if localBranch == nil || err != nil {

		_, err := repo.CreateCommit("HEAD", signature, signature, message, tree)
		if err != nil {
			return err
		}

		//commitTarget, err := repo.LookupCommit(commit)
		//if err != nil {
		//	return err
		//}

		// Creating local branch
		//localBranch, err = repo.CreateBranch(branchName, commitTarget, false)
		//if err != nil {
		//	return fmt.Errorf("Failed to create local branch: %s, %v", branchName, err)
		//}

		// Setting upstream to origin branch
		//err = localBranch.SetUpstream("origin/" + branchName)
		//if err != nil {
		//	return fmt.Errorf("Failed to create upstream to origin/%s", branchName)
		//}

	} else {

		commitTarget, err := repo.LookupCommit(localBranch.Target())
		if err != nil {
			return err
		}

		_, err = repo.CreateCommit("HEAD", signature, signature, message, tree, commitTarget)
		if err != nil {
			return err
		}
	}

	fmt.Print("Here Before Remote\n")

	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		remote, err = repo.Remotes.Create("origin", repo.Path())
		if err != nil {
			return err
		}
	}

	fmt.Printf("Here (%s)\n", branchName)
	err = remote.Push([]string{"refs/heads/" + branchName}, pushOptions)
	if err != nil {
		return err
	}

	return nil
}
