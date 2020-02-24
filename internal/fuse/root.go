// +build !netbsd
// +build !openbsd
// +build !solaris
// +build !windows

package fuse

import (
	"os"
	"time"

	"github.com/restic/restic/internal/debug"
	"github.com/restic/restic/internal/restic"

	"golang.org/x/net/context"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

// Config holds settings for the fuse mount.
type Config struct {
	OwnerIsRoot      bool
	Host             string
	Tags             []restic.TagList
	Paths            []string
	SnapshotTemplate string
}

// Root is the root node of the fuse mount of a repository.
type Root struct {
	repo          restic.Repository
	cfg           Config
	snapshots     restic.Snapshots
	blobSizeCache *BlobSizeCache

	lastCheck time.Time

	entries map[string]fs.Node

	uid, gid uint32
}

// ensure that *Root implements these interfaces
var _ = fs.HandleReadDirAller(&Root{})
var _ = fs.NodeStringLookuper(&Root{})

const rootInode = 1

// NewRoot initializes a new root node from a repository.
func NewRoot(ctx context.Context, repo restic.Repository, cfg Config) (*Root, error) {
	debug.Log("NewRoot(), config %v", cfg)

	root := &Root{
		repo:          repo,
		cfg:           cfg,
		blobSizeCache: NewBlobSizeCache(ctx, repo.Index()),
	}

	if !cfg.OwnerIsRoot {
		root.uid = uint32(os.Getuid())
		root.gid = uint32(os.Getgid())
	}

	root.entries = map[string]fs.Node{
		"snapshots": NewSnapshotsDir(root, fs.GenerateDynamicInode(rootInode, "snapshots"), "", ""),
		"tags":      NewTagsDir(root, fs.GenerateDynamicInode(rootInode, "tags")),
		"hosts":     NewHostsDir(root, fs.GenerateDynamicInode(rootInode, "hosts")),
		"ids":       NewSnapshotsIDSDir(root, fs.GenerateDynamicInode(rootInode, "ids")),
	}

	return root, nil
}

var _ = fs.HandleReadDirAller(&Root{})
var _ = fs.NodeStringLookuper(&Root{})

// Root is just there to satisfy fs.Root, it returns itself.
func (r *Root) Root() (fs.Node, error) {
	debug.Log("Root.Root()")
	return r, nil
}

// Attr returns the attributes for the root node.
func (d *Root) Attr(ctx context.Context, attr *fuse.Attr) error {
	attr.Inode = rootInode
	attr.Mode = os.ModeDir | 0555
	attr.Uid = d.uid
	attr.Gid = d.gid

	debug.Log("attr: %v", attr)
	return nil
}

// ReadDirAll returns all entries directly below the root node.
func (d *Root) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	debug.Log("Root.ReadDirAll()")
	items := []fuse.Dirent{
		{
			Inode: rootInode,
			Name:  ".",
			Type:  fuse.DT_Dir,
		},
		{
			Inode: rootInode,
			Name:  "..",
			Type:  fuse.DT_Dir,
		},
	}

	for name := range d.entries {
		items = append(items, fuse.Dirent{
			Inode: fs.GenerateDynamicInode(rootInode, name),
			Name:  name,
			Type:  fuse.DT_Dir,
		})
	}

	return items, nil
}

// Lookup returns a specific entry from the root node.
func (d *Root) Lookup(ctx context.Context, name string) (fs.Node, error) {
	debug.Log("Root.Lookup(%s)", name)

	if dir, ok := d.entries[name]; ok {
		return dir, nil
	}

	return nil, fuse.ENOENT
}
