package pathx

import (
	pathmatcher "github.com/gohutool/boot4go-pathmatcher"
	"io/fs"
	"path/filepath"
)

func AntMatch(path string, patterns ...string) (bool, error) {
	for _, pattern := range patterns {
		match, err := pathmatcher.Match(pattern, path)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

func AntScan(root string, onlyDir bool, patterns ...string) (matchPaths map[string]string, err error) {
	matchPaths = map[string]string{}
	err = filepath.WalkDir(root, func(absPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && onlyDir {
			return nil
		}

		relPath, err := filepath.Rel(root, absPath)
		if err != nil {
			return err
		}

		match := false
		if len(patterns) == 0 {
			match = true
		} else {
			match, err = AntMatch(relPath, patterns...)
			if err != nil {
				return err
			}
		}
		if match {
			matchPaths[absPath] = relPath
		}
		return nil
	})
	return
}

func AntScanThenDo(root string, onlyDir bool, handler func(absPath, relPath string) error, patterns ...string) error {
	matchedPaths, err := AntScan(root, onlyDir, patterns...)
	if err != nil {
		return err
	}
	for absPath, relPath := range matchedPaths {
		err := handler(absPath, relPath)
		if err != nil {
			return err
		}
	}
	return nil
}
