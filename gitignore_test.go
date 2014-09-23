package the_platinum_searcher

import "testing"

type assert struct {
	patterns []string
	file     file
	expect   bool
}

type file struct {
	path  string
	isDir bool
	depth int
}

func TestGitIgnoreMatch(t *testing.T) {

	asserts := []assert{
		assert{[]string{"a.txt"}, file{"a.txt", false, 1}, true},
		assert{[]string{"dir/a.txt"}, file{"dir/a.txt", false, 2}, true},
		assert{[]string{"dir/*.txt"}, file{"dir/a.txt", false, 2}, true},
		assert{[]string{"dir2/a.txt"}, file{"dir1/dir2/a.txt", false, 3}, true},
		assert{[]string{"dir3/a.txt"}, file{"dir1/dir2/dir3/a.txt", false, 4}, true},
		assert{[]string{"a.txt"}, file{"dir/a.txt", false, 2}, true},
		assert{[]string{"a.txt"}, file{"dir1/dir2/a.txt", false, 3}, true},
		assert{[]string{"dir2/a.txt"}, file{"dir1/dir2/a.txt", false, 3}, true},
		assert{[]string{"dir"}, file{"dir", true, 1}, true},
		assert{[]string{"dir/"}, file{"dir", true, 1}, true},
		assert{[]string{"dir/"}, file{"dir", false, 1}, false},
		assert{[]string{"/a.txt"}, file{"a.txt", false, 1}, true},
		assert{[]string{"/a.txt"}, file{"dir/a.txt", false, 2}, false},
		assert{[]string{"a.txt", "b.txt"}, file{"dir/b.txt", false, 2}, true},
		assert{[]string{"*.txt", "!b.txt"}, file{"dir/b.txt", false, 2}, false},
		assert{[]string{"dir/*.txt", "!dir/b.txt"}, file{"dir/b.txt", false, 2}, false},
		assert{[]string{"dir/*.txt", "!/b.txt"}, file{"dir/b.txt", false, 2}, true},
	}

	for _, assert := range asserts {
		gi := newGitIgnore(1, assert.patterns)
		result := gi.Match(assert.file.path, assert.file.isDir, assert.file.depth)
		if result != assert.expect {
			t.Errorf("Match should return %t, got %t on %v", assert.expect, result, assert)
		}

	}

}
