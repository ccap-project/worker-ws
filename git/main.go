package git

import (
	"crypto/tls"
	"net/http"

	"srcd.works/go-git.v4"
	"srcd.works/go-git.v4/plumbing/transport/client"
	githttp "srcd.works/go-git.v4/plumbing/transport/http"
	//"gopkg.in/src-d/go-git.v4"
	//"gopkg.in/src-d/go-git.v4/plumbing/transport/client"
	//githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func Clone(dir string, url string, chkcert bool) error {

	if chkcert {
		customClient := &http.Client{
			// accept any certificate (might be useful for testing)
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: chkcert},
			},
			//Timeout: 15 * time.Second,
		}
		client.InstallProtocol("https", githttp.NewClient(customClient))
	}

	//r, err := git.NewFilesystemRepository(dir)
	//if err != nil {
	//	return err
	//}

	//r.setIsBare(true)

	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:   url,
		Depth: 1,
	})

	if err != nil {
		return err
	}

	return nil
}

func Pull(dir string, url string) error {

	r, err := git.PlainOpen(dir)

	if err != nil {
		return err
	}

	err = r.Pull(&git.PullOptions{
		RemoteName: "master",
	})

	if err != nil {
		return err
	}

	return nil
}
