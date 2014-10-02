package the_platinum_searcher

import "testing"

type assert struct {
	depth    int
	patterns []string
	file     file
	expect   bool
}

type file struct {
	path  string
	isDir bool
	depth int
}

func TestGenericIgnoreMatch(t *testing.T) {
	asserts := []assert{
		assert{0, []string{"a.txt"}, file{"a.txt", false, 1}, true},
		assert{0, []string{"a.txt"}, file{"dir/a.txt", false, 1}, true},
		assert{0, []string{"dir"}, file{"dir", true, 1}, true},
		assert{0, []string{"dir"}, file{"dir/a.txt", false, 1}, false},
	}

	for _, assert := range asserts {
		gi := genericIgnore(assert.patterns)
		result := gi.Match(assert.file.path, assert.file.isDir, assert.file.depth)
		if result != assert.expect {
			t.Errorf("Match should return %t, got %t on %v", assert.expect, result, assert)
		}

	}
}

func TestGitIgnoreMatch(t *testing.T) {

	asserts := []assert{
		assert{1, []string{"a.txt"}, file{"a.txt", false, 1}, true},
		assert{1, []string{"dir/a.txt"}, file{"dir/a.txt", false, 2}, true},
		assert{1, []string{"dir/*.txt"}, file{"dir/a.txt", false, 2}, true},
		assert{1, []string{"dir2/a.txt"}, file{"dir1/dir2/a.txt", false, 3}, true},
		assert{1, []string{"dir3/a.txt"}, file{"dir1/dir2/dir3/a.txt", false, 4}, true},
		assert{1, []string{"a.txt"}, file{"dir/a.txt", false, 2}, true},
		assert{1, []string{"a.txt"}, file{"dir1/dir2/a.txt", false, 3}, true},
		assert{1, []string{"dir2/a.txt"}, file{"dir1/dir2/a.txt", false, 3}, true},
		assert{1, []string{"dir"}, file{"dir", true, 1}, true},
		assert{1, []string{"dir/"}, file{"dir", true, 1}, true},
		assert{1, []string{"dir/"}, file{"dir", false, 1}, false},
		assert{1, []string{"/a.txt"}, file{"a.txt", false, 1}, true},
		assert{1, []string{"/a.txt"}, file{"dir/a.txt", false, 2}, false},
		assert{2, []string{"/a.txt"}, file{"dir/a.txt", false, 2}, true},
		assert{1, []string{"a.txt", "b.txt"}, file{"dir/b.txt", false, 2}, true},
		assert{1, []string{"*.txt", "!b.txt"}, file{"dir/b.txt", false, 2}, false},
		assert{1, []string{"dir/*.txt", "!dir/b.txt"}, file{"dir/b.txt", false, 2}, false},
		assert{1, []string{"dir/*.txt", "!/b.txt"}, file{"dir/b.txt", false, 2}, true},
	}

	for _, assert := range asserts {
		gi := newGitIgnore(assert.depth, assert.patterns)
		result := gi.Match(assert.file.path, assert.file.isDir, assert.file.depth)
		if result != assert.expect {
			t.Errorf("Match should return %t, got %t on %v", assert.expect, result, assert)
		}
	}
}
