package bitbucketcloud

import "github.com/git-town/git-town/v18/internal/git/giturl"

// Detect indicates whether the current repository is hosted on a Bitbucket server.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "bitbucket.org"
}
