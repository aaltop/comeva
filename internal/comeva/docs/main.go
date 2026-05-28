package docs

import (
	"comeva/internal/comeva/globals"
	"comeva/internal/utils/io"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

//go:embed docs/*
var docsFolder embed.FS

type baseDirectory interface {
	fmt.Stringer
	// GetBaseDirectory returns the directory specified by `path`. Returns non-nil error
	// if the directory does not exist or if `path` is of length zero.
	GetBaseDirectory(path ...string) (dir baseDirectory, e error)
	// DirectoriesString returns a string representation of the top-level directories
	// contained in this directory.
	DirectoriesString() string
}

type writeDirectory interface {
	baseDirectory
	// AddFile adds a file. `path` should be the location of the file.
	// Returns non-nil error when the file already exists or if `path` is
	// of length zero.
	AddFile(content string, path ...string) (e error)
	// GetWriteDirectory returns the directory specified by `path`. Returns non-nil error
	// if the directory does not exist or if `path` is of length zero.
	GetWriteDirectory(path ...string) (dir writeDirectory, e error)
}

type readDirectory interface {
	baseDirectory
	// GetFile returns the content of a file. Returns non-nil error
	// if the file does not exist or if `path` is of length zero.
	GetFile(path ...string) (content string, e error)
	// GetReadDirectory returns the directory specified by `path`. Returns non-nil error
	// if the directory does not exist or if `path` is of length zero.
	GetReadDirectory(path ...string) (dir readDirectory, e error)
}

type readWriteDirectory interface {
	writeDirectory
	readDirectory
	// GetReadWriteDirectory returns the directory specified by `path`. Returns non-nil error
	// if the directory does not exist or if `path` is of length zero.
	GetReadWriteDirectory(path ...string) (dir readWriteDirectory, e error)
}

type directory struct {
	Name string
	// Dirs represents the directories within the directory.
	Dirs map[string]*directory
	// Files represents the files within the directory.
	Files map[string]string
}

func newDefaultDirectory() (dir *directory) {
	return &directory{
		Name:  "/",
		Dirs:  make(map[string]*directory),
		Files: make(map[string]string),
	}
}

func newDirectory(name string) (dir *directory, e error) {
	dir = newDefaultDirectory()
	dir.Name = name
	return
}

func (dir *directory) AddFile(content string, path ...string) (e error) {

	switch {
	case len(path) == 0:
		return errors.New("`paths` should contain at least one item.")
	case len(path) == 1:

		var exists bool
		_, exists = dir.Files[path[0]]
		if exists {
			return errors.New("File already exists.")
		}
		dir.Files[path[0]] = content
		return
	default:
		var dirExists bool
		_, dirExists = dir.Dirs[path[0]]
		if !dirExists {
			dir.Dirs[path[0]], e = newDirectory(path[0])
			if e != nil {
				return
			}
		}

		return dir.Dirs[path[0]].AddFile(content, path[1:]...)
	}

}

func (dir *directory) GetFile(path ...string) (content string, e error) {
	var exists bool
	switch {
	case len(path) == 0:
		e = errors.New("`paths` should contain at least one item.")
		return
	case len(path) == 1:
		content, exists = dir.Files[path[0]]
		if !exists {
			e = errors.New("File does not exist.")
			return
		}
		return content, nil
	default:
		_, exists = dir.Dirs[path[0]]
		if !exists {
			e = fmt.Errorf("Directory `%v` does not exist", path[0])
			return
		}
		return dir.Dirs[path[0]].GetFile(path[1:]...)
	}
}

func (dir *directory) getDirectory(path ...string) (subDir *directory, e error) {
	switch {
	case len(path) == 0:
		e = errors.New("`paths` should contain at least one item")
	case len(path) > 0:
		var exists bool
		subDir, exists = dir.Dirs[path[0]]
		if !exists {
			e = errors.New("Directory does not exist.")
			return
		}
		if len(path) == 1 {
			return
		} else {
			return subDir.getDirectory(path[1:]...)
		}
	}
	return
}

func (dir *directory) GetReadDirectory(path ...string) (subDir readDirectory, e error) {
	return dir.getDirectory(path...)
}

func (dir *directory) GetBaseDirectory(path ...string) (subDir baseDirectory, e error) {
	return dir.getDirectory(path...)
}

func (dir *directory) GetWriteDirectory(path ...string) (subDir writeDirectory, e error) {
	return dir.getDirectory(path...)
}

func (dir *directory) GetReadWriteDirectory(path ...string) (subDir readWriteDirectory, e error) {
	return dir.getDirectory(path...)
}

func (dir *directory) String() string {
	var builder = &strings.Builder{}

	builder.WriteString("+ " + dir.Name + "\n")

	for filepath := range dir.Files {
		builder.WriteString(io.IndentLinesWithString("- "+filepath+"\n", 1, "| "))
	}

	for _, direc := range dir.Dirs {
		builder.WriteString(io.IndentLinesWithString(direc.String(), 1, "| "))
	}

	return builder.String()
}

func (dir *directory) DirectoriesString() string {
	var builder = &strings.Builder{}

	for _, direc := range dir.Dirs {
		builder.WriteString(direc.Name + "\n")
	}
	return builder.String()
}

var dir = newDefaultDirectory()

func init() {

	var subdirRegex = regexp.MustCompile(`docs/(?P<subdir>(?:.*?/)*.*\.md\z)`)
	// map the embedded file system to dir
	fs.WalkDir(docsFolder, "docs", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || !subdirRegex.MatchString(path) {
			return nil
		}
		var cont = subdirRegex.ReplaceAllString(path, "$subdir")
		var splitPath = strings.Split(cont, ".")
		var sepEntry = strings.Split(splitPath[0], "/")

		var content, e = docsFolder.ReadFile(path)
		if e != nil {
			return e
		}
		e = dir.AddFile(string(content), sepEntry...)

		return e
	})

}

// Get returns the documentation or structure for the given path if it exists,
// returning a non-nil error otherwise. . If path ends in an empty string,
// a directory is assumed. If it does not end in an empty string, a file is assumed.
// if path is empty or contains a single empty string, the root directory is assumed.
func Get(path ...string) (content string, e error) {

	var readDir readDirectory = dir

	switch {

	// root
	case len(path) == 0 || (len(path) == 1 && path[0] == ""):
		content = readDir.String()

	// double forward slash would work here, but it should still just error
	// in that case as it's not able to find anything
	// non-root directory
	case path[len(path)-1] == "":
		var subDir readDirectory
		subDir, e = dir.GetReadDirectory(path[:len(path)-1]...)
		if e == nil {
			content = subDir.String()
		} else {
			e = fmt.Errorf("Directory '%v' not found.", "/"+strings.Join(path, "/"))
		}

	default:
		content, e = readDir.GetFile(path...)
		if e != nil {
			e = fmt.Errorf("Documentation for '%v' does not exist.\n", "/"+strings.Join(path, "/"))
		}

	}
	return
}

var subdirRegex = regexp.MustCompile(`docs/(?P<subdir>.*)`)

// Create creates the documentation at the given path. The path should
// end in a forward slash.
func Create(path string) (e error) {
	if len(path) == 0 {
		return errors.New("Path was empty")
	}

	if path[len(path)-1] != '/' {
		return fmt.Errorf("Path should end in a forward slash, got '%v'", path)
	}

	e = os.MkdirAll(path, 0750)
	if e != nil {
		return
	}

	e = fs.WalkDir(docsFolder, "docs", func(_path string, d fs.DirEntry, err error) (e error) {

		if !subdirRegex.MatchString(_path) {
			return nil
		}

		var newPath = subdirRegex.ReplaceAllString(_path, fmt.Sprintf("%v$subdir", path))

		if d.IsDir() {
			globals.DebugLogger.Info().Printf("Creating directories '%v'\n", newPath)
			e = os.MkdirAll(newPath, 0750)
		} else {
			globals.DebugLogger.Info().Printf("Creating docs file at '%v'\n", newPath)
			var content []byte
			content, e = docsFolder.ReadFile(_path)
			if e != nil {
				return
			}
			e = os.WriteFile(newPath, content, 0660)
		}
		return
	})
	return
}
