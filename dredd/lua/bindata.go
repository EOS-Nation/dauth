// Code generated for package lua by go-bindata DO NOT EDIT. (@generated)
// sources:
// cutoff.lua
package lua

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _cutoffLua = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xdc\x56\xdd\x6e\xdb\x38\x13\xbd\xd7\x53\xcc\xa7\xde\x58\xad\x95\xca\x69\xd3\xaf\x08\xe0\x02\x6d\xec\x2d\x8c\x36\x69\x36\x76\xff\x10\x18\x0b\x5a\x1a\xdb\x44\x28\x52\xa0\xa8\xa4\xde\xc5\xee\xb3\x2f\x48\xca\x12\x65\x49\xde\x6d\xd1\xab\xf5\x85\x6d\x88\x67\x86\x33\xe7\xcc\x8f\xc2\x30\xf4\xc2\x30\x84\xb7\xc8\x51\x12\x85\x09\xac\x76\x30\x4d\xd3\xdd\xfb\x82\x0c\xb6\x4a\x65\xf9\xf9\xd3\xa7\x1b\xaa\xb6\xc5\xea\x24\x16\xe9\xd3\xf2\x28\x30\x46\x17\x12\xf7\x26\xf1\x8a\x32\x86\x4a\x9d\x98\x83\x09\x51\xb8\xa0\x29\x9e\xc3\x69\x74\x1a\x85\xd1\x69\x18\x9d\xc1\xcb\xf3\xb3\xff\x03\x39\x49\x0d\xc4\xf3\x24\x26\x34\x3f\x91\x98\x31\x1a\x13\x85\xbf\xc5\x22\x4d\x09\x4f\xf2\x41\xe0\x19\x1f\x1f\x15\x65\x54\x51\xcc\x3d\x8f\x89\x98\x30\x58\x17\x3c\x56\x54\x70\xc8\x19\x8d\x71\xa0\x56\x6c\x08\x6b\x2a\x73\x35\x04\x46\xf4\x77\xae\x30\x0b\x3c\x00\x00\x6b\x60\x70\x09\x8c\xe1\x8f\x3f\x3d\xf3\x78\x2d\x24\x50\x18\x5b\x2b\x10\x12\x46\xd6\x54\xff\x7d\x64\xfc\x69\x17\xe6\x00\x12\x61\x4c\xf4\xc7\xfa\xb9\x7d\x54\xfa\x7b\x02\xa3\x25\x8c\x41\xad\xd8\x2d\x5d\x1a\x10\xf2\xc4\x5e\x20\x51\x15\xb2\x0c\x30\xf1\xcc\xe3\x83\xd8\x63\x91\x66\x85\xc2\x8b\x42\x4a\xe4\xea\x1a\x25\x15\xc9\x85\xe0\x79\x91\x66\xfa\x7c\xf0\x40\x79\x22\x1e\xe6\xf4\x77\x0c\x3c\x27\x95\x3b\xdc\xe5\x30\x2e\x33\x7f\x37\xfd\x3a\x1f\xc2\xd9\x10\xce\xe0\x09\xb8\x06\x35\xfe\x9e\xb0\x02\xb5\x85\x65\x39\x26\x8c\x0d\xfc\xcb\xb7\xd3\x85\x3f\x84\x82\x67\x24\xbe\x1b\x68\x97\x41\x83\xae\x22\x85\x31\x44\x35\x53\x3c\xc1\x6f\x43\xeb\x0a\x28\x07\x9a\x11\x2a\xf3\x81\x75\x1d\xb8\x04\x95\x77\x6a\x52\x04\x2f\xd2\x15\x4a\x8b\x0a\x2a\x04\x5d\xc3\x3d\xfc\x35\x06\x4e\x19\xa8\x2d\xf2\xea\xc0\xf0\x6b\x2e\xd6\xdf\x4f\xe0\xbe\x3a\xd1\xec\xb9\xbf\x7b\x6e\x8b\xd4\x12\x6b\xaa\x56\x60\x0e\x6a\x2b\x45\xb1\xd9\x6a\xbf\x90\x52\x5e\x28\x64\x3b\x88\x0d\xa7\x98\x40\x22\xe2\x5c\x47\x2f\xf1\x1e\x65\x8e\xa0\x04\xac\xe9\xa6\x90\x08\xa2\x50\xb0\x15\x0f\xc0\x04\xdf\x18\xe3\x22\x47\x09\xf9\x56\x14\x2c\x81\x15\xc2\x8a\x89\xf8\x0e\x13\x2f\x0c\x6f\x6f\xbb\x55\x7c\xa3\x11\x93\x42\x12\x23\x5d\x22\xe2\x5f\x0b\xa1\xc8\x10\xbe\x43\xc4\x86\x82\x46\xd4\x70\xf4\x1f\x13\xb2\x34\xd3\x8f\x5f\x8d\x61\xcf\x52\xdb\xda\x11\xb9\xa6\x04\x42\x1b\x7e\x03\xb9\x2f\x89\x23\x65\x12\xe9\x22\x59\x2e\x6d\x99\xbc\x9b\x7e\x85\xc9\xf4\x97\xd9\xd5\x6c\x31\xfb\x70\x35\x2f\x7b\x52\xcb\xed\x36\x20\x26\x13\x11\xc3\x58\xa3\xe7\xb7\xa3\x65\x89\x5a\x31\x12\xdf\xbd\xa7\xb9\x7a\x87\xbb\xfd\xe1\x69\xeb\xf0\x13\xca\x9c\x0a\xee\x60\x9e\x2d\x9d\x6b\x26\x84\xb2\xdd\x44\xc4\x4e\xa7\x3b\xd0\xe7\x2e\xf4\xb2\xac\xe0\x5e\xf4\x59\x99\xd4\xfc\xe2\x66\x76\xbd\x80\xd9\xd5\xf5\xc7\xc5\x3e\x25\xc2\xf4\xaf\x32\x99\x14\x29\x72\x75\x21\x0a\xae\x5c\x39\x5f\xdf\xbc\xfd\x74\x3b\x5a\x06\xa5\x41\xd9\x26\x7b\x78\x0b\x79\x5a\x21\x91\x27\x1f\xd6\x9f\x8d\x2e\x7a\xb2\xcf\x15\x49\xb3\x16\xfc\x59\x05\x5f\xb9\xad\xd1\xc2\x3d\xaf\x70\x8e\xd2\x87\xa0\xb3\xa5\xd9\x04\x99\xa4\x5c\x0d\xfc\x4e\xb5\xce\x7d\x38\x39\xe9\x16\x32\xa8\x4d\xbb\x59\xb1\xb6\xdd\x67\x8e\xf1\x01\x43\xd6\xea\xe0\xa1\x03\x37\x34\x4d\xc8\xae\xe2\xc8\x1a\xb4\x1e\xdb\x05\x1a\x7e\xf7\xc7\x98\x2d\xf4\xa4\x43\xa2\xfe\xf7\xa3\x4e\xe8\xba\xad\xfc\x18\xa2\xba\x27\xdd\x2e\xf2\xca\x95\x4e\xab\xb1\x6a\x4a\x12\xb6\x84\x27\x8c\xf2\x8d\xf6\xe6\x4e\xa7\xd9\xd5\xc5\xcd\x9b\xaf\x7a\x3e\xf5\x56\xfe\xb0\xc5\xa0\xbe\xff\x30\x24\x27\x9a\xda\xfb\xf4\xcb\xf5\xec\x66\xfa\x0f\xde\x5f\xbe\x78\x1e\x45\x41\x1d\xfc\x65\x73\x2d\x7c\x47\xfc\xbd\xed\xf8\xb3\x52\x38\x72\xc1\x8b\x08\x1e\x37\x76\x49\x95\xcf\xfc\x8e\x66\x50\x70\x46\x53\xaa\xdf\xbf\x48\x1c\x63\x9e\xeb\x3c\xfa\x06\x40\x43\xdc\x30\x6c\xdb\x3a\xa2\xfb\xd5\xa1\x6f\x6f\x3c\xe0\x67\xfa\x65\x36\x5f\xcc\xfd\x61\x63\x34\x9a\xec\x47\xad\x02\xf2\x09\x93\x48\x92\x1d\xac\x98\x5f\x87\x7f\x83\x29\xa1\x9c\xf2\x8d\xde\x06\x96\xa8\x4a\x8c\x7a\xdb\x9a\x11\x42\xf9\xc6\x99\x22\x9d\x7b\xb7\x3b\xe7\x21\x8c\xa2\x40\x87\xde\x72\xf3\xca\xe5\xa2\x1c\x6e\xdf\x32\x2a\x71\xa6\xaf\x68\xc1\x1f\xc3\x8b\xc8\x6b\x69\x38\x9f\x2e\xa6\x5f\x0e\x38\x18\x56\x7e\x86\xe0\x2b\xa1\x08\xd3\xf9\x01\x7e\x8b\x11\x13\x4c\xfc\xa0\xe5\x45\x97\x9a\xeb\xa4\x5e\x23\x41\xe3\x85\xd2\x2f\xe9\xd3\x0b\xcd\x46\x9c\xb9\xf3\xce\x19\xdf\xbd\x85\xdc\x1a\x90\xed\x02\xd6\x6c\xf5\xf8\x3d\x52\xd9\x61\x08\x9f\x11\xb6\xe4\x1e\x81\x00\xc7\x87\x63\x0b\xef\x0e\x77\x0e\xeb\x71\xcf\x4b\x70\x2d\xf4\xbf\x7a\x4d\xb6\x41\x58\x97\x87\x14\x2c\x66\x97\x53\xbf\x09\x71\xb4\xee\xdc\x69\x61\xbd\x86\x94\xd9\x94\x1d\xd2\x1f\xe1\xb4\x27\x62\xeb\xa6\x4f\xb5\x3e\xab\x32\xf0\xde\xe1\xd1\x11\xc0\x3e\xbd\x76\xdc\xd6\xf2\xf5\x91\xe0\xbb\xf8\x08\xaa\x19\xd0\x1d\xfc\x2b\xe8\x69\xc0\xe6\xec\xea\x9e\x84\xdd\x5d\xd4\x78\x77\xf8\x19\xad\x64\x59\xec\x18\xf3\x3d\x06\x07\xb3\x2c\x0c\x7f\x68\x09\x35\xd2\xd8\x07\x71\xf0\x9a\xea\x57\x73\xb1\x1e\x98\x0c\x36\x42\x24\xbe\xf7\x77\x00\x00\x00\xff\xff\x84\xb9\x95\xb9\x9b\x0f\x00\x00")

func cutoffLuaBytes() ([]byte, error) {
	return bindataRead(
		_cutoffLua,
		"cutoff.lua",
	)
}

func cutoffLua() (*asset, error) {
	bytes, err := cutoffLuaBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cutoff.lua", size: 3995, mode: os.FileMode(420), modTime: time.Unix(1612876165, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"cutoff.lua": cutoffLua,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"cutoff.lua": &bintree{cutoffLua, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
