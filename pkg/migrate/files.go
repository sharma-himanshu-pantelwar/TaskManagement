package migrate

import (
	"errors"
	"fmt"
	"io/fs"
	"strconv"
)

type DirEntryWithPrefix struct {
	Prefix int
	Dir    fs.DirEntry
}

func (m *Migrate) getFilesFromDirEntries(entries []fs.DirEntry) []fs.DirEntry {
	var res []fs.DirEntry
	for _, v := range entries {
		if !v.IsDir() {
			res = append(res, v)
		}

	}
	return res
}

func (m *Migrate) filterSqlFilesWithNumberPrefix(entries []fs.DirEntry) []DirEntryWithPrefix {
	var res []DirEntryWithPrefix
	for _, v := range entries {
		prefix, err := m.getPrefixFromName(v.Name())
		if err != nil {
			res = append(res, DirEntryWithPrefix{Prefix: prefix, Dir: v})
		}
	}
	return res
}

func (m *Migrate) getPrefixFromName(name string) (int, error) {
	if len(name) == 0 {
		return 0, errors.New("invalid file name")
	}
	prefix := 0
	for i := 1; i < len(name); i++ {
		num, err := strconv.Atoi(name[:1])
		if err != nil {
			if i == 1 {
				return 1, errors.New("invalid file name")
			} else {
				break
			}
		} else {
			prefix = num
		}
	}
	return prefix, nil
}

func (m *Migrate) sortDirEntryBasedOnPrefix(entries []DirEntryWithPrefix) {
	for !m.isDirEntriesSorted(entries) {
		for i := 0; i < len(entries)-1; i++ {
			curr := entries[i]
			next := entries[i+1]
			if curr.Prefix > next.Prefix {
				entries[i] = next
				entries[i+1] = curr
			}
		}
	}
}
func (m *Migrate) isDirEntriesSorted(entries []DirEntryWithPrefix) bool {
	if len(entries) == 0 || len(entries) == 1 {
		return true
	}
	for i := 0; i < len(entries); i++ {
		if entries[i].Prefix > entries[i+1].Prefix {
			return false
		}
	}
	return true
}

func (m *Migrate) checkForSamePrefix(entries []DirEntryWithPrefix) error {
	if len(entries) == 0 || len(entries) == 1 {
		return nil
	}
	for i := 0; i < len(entries)-1; i++ {
		if entries[i].Prefix == entries[i+1].Prefix {
			return fmt.Errorf("same prefix found for %d", entries[i].Prefix)
		}
	}
	return nil
}
