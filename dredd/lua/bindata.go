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

var _cutoffLua = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x57\x6d\x6f\xda\xc8\x13\x7f\xef\x4f\x31\x7f\xf7\xc5\x1f\xa7\x36\x25\xb4\xb4\xa7\xe8\x38\xa9\x07\x5c\x85\x9a\xa4\x51\xa0\x0f\x51\x84\x4e\x8b\x3d\xc0\x8a\xf5\xda\xda\x5d\x43\xe8\xe9\xee\xb3\x9f\x76\xd7\x4f\x60\x9c\x53\x4f\xc7\x0b\x88\x76\x7e\x33\xb3\xf3\x9b\xa7\x4d\x10\x04\x4e\x10\x04\xf0\x01\x39\x0a\xa2\x30\x82\xe5\x01\x26\x71\x7c\xb8\xce\x48\x67\xa3\x54\x2a\xaf\x5e\xbd\x5a\x53\xb5\xc9\x96\xdd\x30\x89\x5f\xe5\x22\xcf\x28\x8d\x04\x16\x2a\xe1\x92\x32\x86\x4a\x75\x8d\x60\x4c\x14\xce\x69\x8c\x57\xd0\xef\xf5\x7b\x41\xaf\x1f\xf4\x06\xf0\xd3\xd5\xe0\x1d\x90\x6e\x6c\x20\x8e\x23\x30\xa2\xb2\x2b\x30\x65\x34\x24\x0a\x7f\x0f\x93\x38\x26\x3c\x92\x1d\xcf\x31\x36\x3e\x2b\xca\xa8\xa2\x28\x1d\x87\x25\x21\x61\xb0\xca\x78\xa8\x68\xc2\x41\x32\x1a\x62\x47\x2d\x99\x0f\x2b\x2a\xa4\xf2\x81\x11\xfd\x2d\x15\xa6\x9e\x03\x00\x60\x15\x0c\x2e\x82\x21\xfc\xf1\xa7\x63\x8e\x57\x89\x00\x0a\x43\xab\x05\x89\x80\x4b\xab\xaa\xff\x7c\x61\xec\x69\x13\x46\x00\x51\x62\x54\xf4\xc7\xda\x79\x7c\x91\xdb\x7b\x09\x97\x0b\x18\x82\x5a\xb2\x47\xba\x30\x20\xe4\x91\x75\x20\x50\x65\x22\xbf\x60\xe4\x98\xe3\x93\xbb\x87\x49\x9c\x66\x0a\x47\x99\x10\xc8\xd5\x1d\x0a\x9a\x44\xa3\x84\xcb\x2c\x4e\xb5\xbc\xb3\xa7\x3c\x4a\xf6\x33\xfa\x1d\x3d\xa7\x16\xca\x16\x0f\x12\x86\x79\xe4\x1f\x27\x0f\x33\x1f\x06\x3e\x0c\xe0\x25\xd4\x15\x2a\xfc\x8e\xb0\x0c\xb5\x86\x65\x39\x24\x8c\x75\xdc\x9b\x0f\x93\xb9\xeb\x43\xc6\x53\x12\x6e\x3b\xda\xa4\x77\x44\x57\x16\xc3\x10\x7a\x15\x53\x3c\xc2\x27\xdf\x9a\x02\xca\x81\xa6\x84\x0a\xd9\xb1\xa6\xbd\x3a\x41\xb9\x4f\x4d\x4a\xc2\xb3\x78\x89\xc2\xa2\xbc\x12\x41\x57\xb0\x83\xbf\x86\xc0\x29\x03\xb5\x41\x5e\x0a\x0c\xbf\xc6\xb1\xfe\x7e\x09\xbb\x52\xa2\xd9\xab\xff\x16\xdc\x66\xb1\x25\x56\xd7\xc8\xc7\xc9\x03\x8c\x27\xbf\x4d\x6f\xa7\xf3\xe9\xa7\xdb\x59\x4e\x76\x26\x51\xd4\x99\xc5\x68\x9c\x84\x30\xd4\xe8\xd9\xe3\xe5\x22\x47\x2d\x19\x09\xb7\xd7\x54\xaa\x8f\x78\x28\x84\xfd\x86\xf0\x0b\x0a\x49\x13\x5e\xc3\xbc\x5e\xd4\xdc\xdc\x50\x9e\x29\x64\x87\x71\x12\xd6\xb2\x58\x43\xbf\x29\x2d\x66\xe2\xc8\xd5\x60\x61\x43\x98\x8d\xee\xa7\x77\x73\x98\xde\xde\x7d\x9e\x17\x01\x10\xa6\x7f\x95\xb9\x77\x16\x23\x57\xa3\x24\xe3\xaa\x4e\xef\xfb\xfb\x0f\x5f\x1e\x2f\x17\x5e\xae\x10\xda\x38\x0b\x78\x03\xd9\x2f\x91\xc8\xa3\x4f\xab\xaf\xa6\x66\x74\x83\xce\x14\x89\xd3\x06\xfc\x75\x09\x5f\xb2\x24\xdc\x8e\x33\x41\x4c\xf1\x9e\xe2\xde\x94\xb8\xaa\x0a\x2d\x23\xb2\x81\x1d\x54\x36\x35\x13\x37\x19\x53\x34\x65\x14\x45\x03\xf9\xf6\x18\x69\x2f\xdb\x66\xf6\xdd\xc2\x8c\x8b\x54\x50\xae\x3a\xee\xd9\xcc\x5f\xb9\xd0\xed\x9e\x2f\x0a\xaf\x52\x3d\xcf\xb9\xd5\x3d\x2f\xab\x29\x9f\xf0\x6f\xb5\x4e\x0e\x6b\x70\x93\x84\x31\x39\x94\x19\xb0\x0a\x8d\x63\x3b\x65\x83\x1f\xfe\x18\xb5\xf9\x06\x21\x46\xa2\xfe\xf7\x6f\x8d\xd0\x55\xb3\xae\x86\xd0\xab\xfa\x37\xef\xc8\x5e\xd5\x8f\xb3\x2d\x4d\x21\xe3\x8c\xc6\x54\xef\x05\x12\x86\x28\xa5\x36\xd4\x56\xd1\x47\xf6\x82\xa0\xa9\x5b\xf3\xe3\x96\x42\xb7\xf0\x68\x1d\x12\x26\x90\x44\x07\xdb\xb4\x8c\x4a\xad\xae\xa7\x9b\x76\x5c\x9f\x7f\x93\x6f\xd3\xd9\x7c\xe6\xfa\x47\xad\xef\xe9\x4b\x5c\x36\x82\x72\x2b\xa3\x6e\x15\x5f\xd1\xee\x39\x31\xa6\xd9\x61\x43\x78\xc4\x28\x5f\x9f\xba\x9b\xde\x8e\xee\x7f\x7d\xd0\x03\xf7\xb9\x41\xe1\x37\xea\x44\x5f\xe8\x94\xf8\xda\xf5\xea\xf1\xdc\x4d\xef\x27\xff\xec\xe0\x6d\x0f\x2e\x9a\x1d\xea\x55\x61\x8d\x36\x18\x6e\x8f\x62\x4a\x76\x28\xb4\x53\x08\xed\xa2\xca\xd5\xf3\xbe\x4c\xeb\x7d\x54\x1b\x3a\xad\xd1\x37\x1a\xaf\x19\xb5\x66\xaf\xc5\xee\x33\x74\x04\x01\x7c\x45\xd8\x90\x1d\x02\x01\x8e\xfb\xe2\xbe\x8d\xc5\xaa\x0b\xa2\xb6\xeb\x5a\x61\xc3\x1f\xd8\xd0\x25\x91\xda\x6e\x1b\x27\x6d\x9e\x9c\x46\x3a\x67\x76\x3b\xb7\xd0\xd5\x62\xc6\x6b\xa9\x8a\xf7\xcf\xd8\x3a\xb7\x03\xf2\x6a\x68\x4d\xc2\x2f\xd0\x69\xe9\xe1\x73\xa5\x75\x94\x20\xba\xb2\x95\x94\x5f\x01\xa2\xc2\x66\x68\xf4\xa9\x04\x46\xc4\xda\xd4\x1b\xe1\x06\x6a\xe6\x3e\x48\xfa\x1d\xf5\x5b\x4c\x9f\x6c\xf1\x00\x1b\x22\xcb\x56\xcf\x24\x46\x40\xd5\xff\x65\x8e\xdd\x63\xd5\xfd\x85\x42\xe1\xdf\x22\xa8\x34\x43\x68\x8f\x11\x24\x3c\x44\xf3\xc2\x21\x1c\x50\x3f\x54\xcc\x23\xd6\x46\x01\x7b\xaa\x36\xd4\x5e\x63\x4d\x77\xc8\x0b\x07\xb6\xfe\xf3\xb7\xcc\x7f\xc0\x11\x5c\x9c\x2e\x42\x4f\x07\x5b\xcf\xe4\x7c\x7e\xad\x87\x55\xfe\x70\xf0\xe0\x67\xe8\x74\xce\xac\xc4\x0b\x78\xdb\xf3\x20\x80\x66\x65\x1a\x11\x5c\x40\xdf\xf3\x8e\x5f\x5d\x27\x65\x37\xf9\x76\x32\x14\xfd\xe3\xc5\xef\x83\xab\x12\x45\x98\x4e\x1d\xe0\x53\x88\x18\x61\xe4\x7a\x67\xcd\xe9\xbe\xaf\x5b\xab\x1e\x50\x75\xbc\x9d\xb1\x7a\xb6\xd6\x9f\x78\x67\x92\xe5\x03\x89\x22\x93\xa8\xa7\x94\x0a\xca\xd7\x39\x42\x17\x04\x5d\x01\x4f\x54\x59\x13\xf8\x44\xa5\x92\x45\x8a\xce\x0f\xfe\x92\xcb\xe3\xcd\xd3\xca\x49\x8e\xf7\xcf\x3c\x45\x7c\x70\xcd\x61\xce\xc3\xe9\xff\x00\xb9\xd0\xb6\x55\xb5\x54\x18\xac\x93\x24\x72\x9d\xbf\x03\x00\x00\xff\xff\xa3\x97\xf0\x40\x78\x0d\x00\x00")

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

	info := bindataFileInfo{name: "cutoff.lua", size: 3448, mode: os.FileMode(420), modTime: time.Unix(1615980337, 0)}
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
