package gitrans

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// -----------------------------------------------------------------------------

// FileMode represents the file mode in the Git repository.
type FileMode = filemode.FileMode

// opener defines an interface for opening a file.
type opener interface {
	Reader() (io.ReadCloser, error)
}

// File represents a file in the Git repository.
type File struct {
	// Name is the path of the file. It might be relative to a tree,
	// depending of the function that generates it.
	Name string
	// Mode is the file mode.
	Mode FileMode
	// Size of the (uncompressed) blob.
	Size int64
	// opener opens the file for reading.
	opener
}

// Content returns the content of the file as a byte slice.
func (p *File) Content__0() ([]byte, error) {
	switch o := p.opener.(type) {
	case *bytesOpener:
		return o.data, nil
	default:
		r, err := o.Reader()
		if err != nil {
			return nil, err
		}
		defer r.Close()
		return io.ReadAll(r)
	}
}

// Content sets the content of the file from a byte slice.
func (p *File) Content__1(b []byte) {
	p.Size = int64(len(b))
	p.opener = &bytesOpener{data: b}
}

// Content sets the content of the file from a string.
func (p *File) Content__2(s string) {
	p.Content__1(unsafe.Slice(unsafe.StringData(s), len(s)))
}

// Unchanged returns true if the file has NOT been modified.
func (p *File) Unchanged() bool {
	_, ok := p.opener.(*object.File)
	return ok
}

type bytesOpener struct {
	data []byte
}

func (p *bytesOpener) Reader() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(p.data)), nil
}

// -----------------------------------------------------------------------------

func writeFile(name string, r opener) (err error) {
	rc, err := r.Reader()
	if err != nil {
		return err
	}
	defer rc.Close()

	f, err := os.Create(name)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = io.Copy(f, rc)
	return
}

func chdirToGitRoot() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panicln("failed to get current directory:", err)
	}

	curr := cwd
	for {
		gitDir := filepath.Join(curr, ".git")
		info, err := os.Stat(gitDir)

		// Check if .git exists and is a directory
		if err == nil && info.IsDir() {
			if err := os.Chdir(curr); err != nil {
				log.Panicln("failed to change directory to git root:", err)
			}
			return
		}

		parent := filepath.Dir(curr)
		if parent == curr {
			log.Fatal("fatal: not a git repository (or any of the parent directories): .git")
		}
		curr = parent
	}
}

// -----------------------------------------------------------------------------

type applyer struct {
	upstream string
	handlers []handler
	noExec   bool
}

func newApplyer(p *App) *applyer {
	return &applyer{
		upstream: p.upstream,
		handlers: p.handlers,
		noExec:   p.noExec,
	}
}

func (p *applyer) applyFile(f *object.File) (err error) {
	name := f.Name
	file := &File{
		Name:   name,
		Mode:   f.Mode,
		Size:   f.Size,
		opener: f,
	}
	for _, h := range p.handlers {
		if matchPattern(h.pattern, name) {
			h.callback(file)
		}
	}
	if file.Unchanged() {
		return
	}
	if p.noExec {
		log.Println("edit", name)
		return
	}
	return writeFile(name, file.opener)
}

func (p *applyer) applyTrans(repo *git.Repository) {
	ref, err := repo.Reference(plumbing.NewBranchReferenceName(p.upstream), true)
	if err != nil {
		log.Panicln("failed to get upstream branch:", err)
	}
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		log.Panicln("failed to get commit object:", err)
	}
	tree, err := commit.Tree()
	if err != nil {
		log.Panicln("failed to get tree object:", err)
	}
	err = tree.Files().ForEach(p.applyFile)
	if err != nil {
		log.Panicln("failed to apply files:", err)
	}
}

func (p *applyer) run() {
	chdirToGitRoot()
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Panicln("failed to open git repository:", err)
	}
	p.applyTrans(repo)
}

// -----------------------------------------------------------------------------
