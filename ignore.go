package the_platinum_searcher

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/monochromegane/go-gitignore"
	"github.com/monochromegane/go-home"
)

type ignoreMatchers []gitignore.IgnoreMatcher

func (im ignoreMatchers) Match(path string, isDir bool) bool {
	for _, ig := range im {
		if ig == nil {
			return false
		}
		if ig.Match(path, isDir) {
			return true
		}
	}
	return false
}

func newIgnoreMatchers(path string, ignores []string) ignoreMatchers {
	var matchers ignoreMatchers
	for _, i := range ignores {
		if matcher, err := gitignore.NewGitIgnore(filepath.Join(path, i), path); err == nil {
			matchers = append(matchers, matcher)
		}
	}
	return matchers
}

func globalGitIgnore(base string) gitignore.IgnoreMatcher {
	if homeDir := home.Dir(); homeDir != "" {
		globalIgnore := globalGitIgnoreName()
		if globalIgnore != "" {
			if matcher, err := gitignore.NewGitIgnore(filepath.Join(homeDir, globalIgnore), base); err == nil {
				return matcher
			}
		}
	}
	return nil
}

func globalGitIgnoreName() string {
	gitCmd, err := exec.LookPath("git")
	if err != nil {
		return ""
	}

	file, err := exec.Command(gitCmd, "config", "--get", "core.excludesfile").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(filepath.Base(string(file)))
}

func homePtIgnore(base string) gitignore.IgnoreMatcher {
	if homeDir := home.Dir(); homeDir != "" {
		if matcher, err := gitignore.NewGitIgnore(filepath.Join(homeDir, ".ptignore"), base); err == nil {
			return matcher
		}
	}
	return nil
}
