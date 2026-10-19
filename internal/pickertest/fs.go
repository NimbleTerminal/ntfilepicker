package pickertest

import "nimbleterminal.dev/ntfilepicker/internal/memfs"

// FS is the in-memory filesystem a script's fs part builds; see memfs.
type FS = memfs.FS

// NewFS makes an FS of lines; see memfs.NewFS.
func NewFS(lines ...string) *FS { return memfs.NewFS(lines...) }

// Epoch is the modification time of an entry given none.
var Epoch = memfs.Epoch
