// Code generated by go-bindata.
// sources:
// conf/app.ini
// conf/locale/locale_en-US.ini
// conf/locale/locale_fr-FR.ini
// DO NOT EDIT!

package bindata

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

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _confAppIni = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x56\xdd\x72\xe2\x3a\x12\xbe\xd7\x53\xf4\xb0\x35\x7b\x92\x2d\x92\x90\xe4\x6c\x26\x0b\xc3\xd6\x71\xc0\x24\xae\x01\xcc\xda\x66\x72\xb2\x53\x53\x1e\x61\xb7\x41\x15\x59\x62\x24\x91\xc0\x5e\xec\xeb\xec\x83\xec\x8b\x9d\x92\x6c\x33\x90\xb9\x39\x49\x55\x08\xed\xee\xcf\xad\xef\xeb\x1f\x79\xb3\x59\x3a\xf5\x26\x3e\xf4\x61\xc9\xcc\xb9\xd9\x1a\x32\xf0\xa6\x69\xe4\xdf\x07\x71\xe2\x47\xd0\x07\xa3\x36\x48\xbc\x69\x38\x7d\x9a\x84\xf3\x38\x1d\x44\xbe\x97\xf8\x8d\x9d\x7c\xd1\xa8\x5e\x50\x7d\x25\xb3\x28\x4c\xc2\x41\x38\x86\x3e\xac\x8c\x59\x93\x61\x38\xf1\x82\x29\xf4\x81\xcb\x8c\xf2\x95\xd4\x86\x44\x61\x98\xa4\xf3\xc8\xba\xbc\x3f\x69\xfc\x4f\x75\xf7\xe2\xe2\xfd\x49\xe5\x7e\xaa\xbb\xef\x4f\x1e\x92\x64\x96\xce\xc2\x28\x39\xd5\x17\x64\x3e\x0d\x7e\x4f\xe3\x70\xf0\xc9\x4f\xd2\x99\x1f\x4d\x82\x38\x0e\x42\x0b\x7b\x73\x73\x43\x9c\xa7\x37\x1c\xda\x34\x3b\xe7\xee\x97\xec\xa3\xa1\x0f\xbf\x76\x3a\x1d\x32\x0c\x62\xef\x6e\xec\xa7\x51\x38\x4f\xfc\x28\x1d\x87\xf7\xd0\x87\x82\x72\x8d\xa4\x07\xf3\xf5\x1a\x15\x70\x7c\x41\x0e\xb2\x00\x83\xe5\x9a\x53\x83\x40\x45\x0e\xda\x50\xc3\x32\x28\x18\x47\x58\x53\xb3\x22\x3d\xc8\xb1\xa0\x1b\x6e\x80\x69\x30\xab\xca\x0a\xaf\x2b\x54\xd8\x90\x67\x9f\xe0\x16\xb3\x8d\xc1\x9c\xc4\x89\x97\x04\x83\xd4\x1d\x7b\xe6\x25\x0f\xd0\x27\xe4\x4b\x4e\x0d\x5d\x50\x8d\x5f\x49\x0f\x7c\x66\x56\xa8\xa0\x55\xee\xf4\x77\xde\x6a\x43\x6b\x2d\xb5\x59\x2a\xd4\x2d\x90\x0a\x5a\xfa\x3b\x67\x06\xaf\x5b\x6d\xd8\xc9\x0d\x64\x54\x40\x26\x85\xc0\xcc\x80\x91\x90\xb0\xe1\x1d\xbc\x32\xb3\x82\xc9\x2e\xfe\xd7\x18\xd6\x4a\x1a\x99\x49\x4e\x86\x77\x69\xf2\x34\xb3\x0a\xd5\xf1\xe4\x21\x8c\x2d\x1b\x97\x57\x1f\x1c\x47\x97\xdd\xeb\xeb\xce\x0d\x69\x54\x97\x4b\x4d\xe6\xb1\x93\x5a\x49\x69\xc8\xcc\x8b\xe3\xc7\x21\xf4\x49\x0f\x46\x36\x8b\x83\x9c\x04\xdf\xb5\x01\xeb\xa4\x73\xa6\xe9\x82\xa3\x4d\x5b\xe1\xf7\x0d\x53\x58\x65\xfd\x82\x8a\x15\xbb\xb3\x62\xc3\x79\x8b\xc4\xf1\x38\x9d\x84\x43\xfb\xa2\xda\xbf\x81\x6d\x0e\xe7\xa8\x6e\x19\x96\x2f\x5a\x6d\xd8\x68\x04\xba\xd0\x92\x6f\xcc\x0f\x76\x85\x3b\xbe\x36\x54\x19\xa0\x1a\x6c\xc1\xb1\x0c\x49\xc5\x68\x43\xfc\x79\xbe\x20\xe4\x8b\xc2\xb5\xd4\xcc\x48\xb5\xb3\xf4\x46\x52\x9a\x0a\xa5\x90\x0a\xb4\x91\x8a\x89\x25\xec\x7d\x18\xea\x5f\x34\x58\x3d\xda\x87\xc2\xb6\xfe\x7b\xf1\x71\xa3\x51\x09\x5a\xe2\x3f\x2f\x96\xcc\x6c\xcd\xd9\x61\x4c\xcb\xd5\xb1\xd3\x92\xcb\xe5\x57\x72\x24\x6f\x0f\x06\x54\xc0\x02\xa1\x95\x49\xa1\x25\xc7\xfa\x7c\xb6\x88\x5a\xc7\xaf\x69\x1c\x6c\x15\x6a\x84\x4c\x96\x25\xb5\xc2\x6a\x5c\x53\x65\x6b\xb0\xdc\x70\xc3\xd6\x1c\xa1\x94\x39\xea\x36\xe0\xf9\xf2\x7c\x1f\xd6\x76\x75\xd9\x22\x35\xb9\xb5\x95\xf4\xe0\x6e\x53\x14\xae\xa4\xc5\xd2\xac\x6c\x4d\x67\x2b\x2a\x04\xf2\x36\x3c\x23\xae\x81\x39\x0e\x99\x4b\x81\x15\x8e\xd9\x5c\x8a\x5f\x0c\x3c\x0b\xf9\x0a\xaf\x2b\x6a\xaa\x87\xe7\xe4\x6e\x3e\x1a\xd9\x7e\xf1\x6d\xb7\x5d\x76\x3a\x07\xf5\x9a\x28\x9a\x39\xe1\x03\x51\x48\xfb\xf9\x48\x95\xb0\x9f\xbe\x52\x52\xd9\x7f\x46\xd4\x50\xfe\xe6\xc0\x55\x14\x19\xfb\x9f\x7d\x3b\x02\xdc\x57\xd2\x94\xc3\x9e\x2e\x7b\x58\x57\x69\x8e\xde\xf3\xda\x6e\xe5\xe4\x48\x5f\x10\xb0\x5c\x9b\x9d\xe5\x89\x89\x15\x2a\x66\x1a\xbc\x3d\x92\xe3\xe5\x2d\x8c\x35\xfe\x09\x8c\x1e\x24\x2b\xdb\xc1\xc2\x56\xaa\x06\xba\x31\xb2\xa4\x06\x73\xe0\x72\x09\x4a\x1a\x2b\xcb\x89\x7e\x65\x26\x73\xd4\x16\x92\x73\xf9\x6a\x8b\x4a\xae\x0d\x93\x42\x9f\x92\x71\x78\x9f\x46\x61\x72\x30\x22\x7b\x10\xe3\xb2\x44\x61\x1c\x48\x4e\x19\xdf\x91\xa1\x17\x8c\x9f\x7e\xf2\x9b\xd0\x2d\x68\xf6\x1f\x04\xbd\x62\x85\xb1\x2f\xd0\x4c\x2c\x39\x3a\xa9\x8f\xb8\xbc\xba\x85\x12\xa9\xd0\x70\x09\x1f\x3f\xc2\xd5\x6d\x1b\xae\xfe\x7e\x33\xb9\x23\x13\xef\xf7\x34\x0e\xfe\xed\xa7\xf1\x43\x30\xb2\x3d\x7f\x75\x5b\xe3\x72\x26\x10\xc4\xa6\x5c\xa0\x7a\x03\xec\x82\xc6\xc1\xd4\x8f\x2b\x9d\xed\x8f\xd5\x7a\xbb\x66\x0a\x73\xc8\xe9\x4e\xdb\x08\x9b\xbc\x1b\x85\x27\x39\x72\xb4\x23\xb2\x30\xa8\xa0\xa4\x5b\xe7\x72\xea\x60\x86\xde\x93\x45\xf9\xb0\x17\x43\x73\x9a\x3d\xff\xa4\x86\xb3\xfe\x29\x39\x1e\x71\xb1\x92\xf2\x19\xe6\xd1\x98\xb8\xcd\x51\x75\xdd\xf9\x56\xaa\xd2\x4d\x50\xa7\x54\x95\x98\xd3\x87\x49\x41\xde\xf2\x1a\x55\xc2\xe1\x0b\xaa\x9d\x4d\xb6\x76\x48\x9d\x0a\x3f\xb9\x49\x91\xd5\x80\x4e\x0b\xdc\x66\xa8\x35\x6a\xd8\xc2\x01\xc1\xfb\x96\x98\xd0\x2d\x2b\x37\x65\x45\x93\x91\x55\x93\x71\xb9\x5c\xa2\x72\x20\xfa\x90\x97\x6b\xb7\x2f\xb5\x66\x52\x1c\x8d\x7f\x2c\xa5\xda\xd9\xb6\xa9\xc7\x84\x65\x4e\x61\xce\xf4\x9b\x0e\xaa\x1d\xed\xae\xfd\x1c\x0c\xdd\xc8\xae\x4c\xa4\x07\x33\x25\x5f\x58\x8e\xca\x8e\x82\x82\xed\x2b\x92\xf4\x6a\x97\x2e\x08\x69\x60\x65\xf9\xa6\x62\xd7\x78\xed\xd0\x90\x9e\x4b\xb4\x0b\x75\x66\x3f\x16\x5e\x3d\x72\xbe\xd9\x11\x79\x51\x3f\xd5\xdf\x48\x0f\x5c\x6e\x5d\x10\x68\x5e\xa5\x7a\xee\x9b\x6c\xdd\xa6\x79\xae\xfa\xdd\x9b\xeb\x0f\xff\x68\xaf\xa9\xd6\xaf\x52\xe5\xfd\x92\x66\x54\x49\xd1\xce\x17\xfd\x4e\x7b\x2d\x25\x4f\x2d\xa3\xfd\xcb\x4e\xa7\xcd\x72\x8e\xa9\x61\x25\xca\x8d\xe9\x5f\xde\x5a\x26\xdd\x12\xec\xc2\x52\x9e\xe9\xef\xfc\x2c\x57\xec\x05\xd5\x85\x33\x42\xae\x45\x93\xb0\x36\x76\x88\x37\x99\xd9\x85\xd5\x6d\xde\xf7\x5b\x93\x63\x6a\x6c\x51\x7c\xdb\xd3\x94\x0e\xc2\xe9\x28\xb0\x6b\xff\xe8\x24\xae\x33\xab\x23\x67\x52\x3e\x33\x04\x3b\xf2\xc9\x20\x0c\x3f\x05\x7e\x73\x29\x62\x29\x67\xcf\x98\x2e\xe5\x92\x19\x1b\x11\x54\x33\xd3\x2e\xaa\x86\x2f\x26\xdc\x8d\x47\xd7\xfb\xf1\x40\xb0\xea\x96\x51\x03\xc6\xfe\x60\x1e\xf9\x07\x77\x8f\xba\x76\x35\x9a\xfa\xfd\x47\xb1\xae\x28\xfd\xa9\xbb\xba\xc4\x7e\x92\x56\x28\x87\x43\xa5\x7a\xfb\xfd\x00\x2c\x8f\xc0\x84\x41\xf5\x42\xf9\x11\xc8\xf5\x4d\xa7\x43\xee\x07\x69\x30\x4d\xfc\xe8\xb3\x37\x4e\x93\xc0\x9d\xca\xd9\x7f\x60\x70\x56\xa0\x43\x39\x0a\xbe\xbd\xf9\xb5\xd3\x21\xb1\xef\x6e\x5c\xe9\x38\x18\xf9\x4d\x78\xf5\xa4\x07\x83\x1f\xb4\xb9\x15\x3b\x88\xa3\x11\xb1\x7f\xd2\x63\x12\xd3\x4c\xab\xc2\xd5\x7e\xb6\x51\xcc\xb8\xe5\xfc\xee\xdd\xe0\xc1\x9b\xde\xfb\x90\x3c\x04\x31\x24\x21\x7c\xf2\xfd\x19\x3c\x85\xf3\x08\xdc\x6d\x64\xe8\x25\x1e\xc4\xde\xc8\x7f\xf7\x8e\xc4\xfe\x20\xf2\x93\xf4\x93\x6f\xbb\xf5\xdd\x5f\x7e\x1b\x0d\xfd\xc7\xc8\x7f\x8c\xfe\xfa\xb7\x13\x42\xbe\xb0\xcb\x5b\xf1\x95\x8c\xbd\xe9\xbd\xed\x31\x14\x67\xf3\xb8\x5d\xa8\xb3\x51\xe4\xae\x38\xd6\xe6\x8b\x25\x67\x7a\xd5\x1e\x29\x2a\xfe\xff\x3f\xca\x34\xf9\x23\x00\x00\xff\xff\x68\x92\x51\xa3\xf7\x0a\x00\x00")

func confAppIniBytes() ([]byte, error) {
	return bindataRead(
		_confAppIni,
		"conf/app.ini",
	)
}

func confAppIni() (*asset, error) {
	bytes, err := confAppIniBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "conf/app.ini", size: 2807, mode: os.FileMode(420), modTime: time.Unix(1494021357, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _confLocaleLocale_enUsIni = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func confLocaleLocale_enUsIniBytes() ([]byte, error) {
	return bindataRead(
		_confLocaleLocale_enUsIni,
		"conf/locale/locale_en-US.ini",
	)
}

func confLocaleLocale_enUsIni() (*asset, error) {
	bytes, err := confLocaleLocale_enUsIniBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "conf/locale/locale_en-US.ini", size: 0, mode: os.FileMode(420), modTime: time.Unix(1493941326, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _confLocaleLocale_frFrIni = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func confLocaleLocale_frFrIniBytes() ([]byte, error) {
	return bindataRead(
		_confLocaleLocale_frFrIni,
		"conf/locale/locale_fr-FR.ini",
	)
}

func confLocaleLocale_frFrIni() (*asset, error) {
	bytes, err := confLocaleLocale_frFrIniBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "conf/locale/locale_fr-FR.ini", size: 0, mode: os.FileMode(420), modTime: time.Unix(1493941287, 0)}
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
	"conf/app.ini": confAppIni,
	"conf/locale/locale_en-US.ini": confLocaleLocale_enUsIni,
	"conf/locale/locale_fr-FR.ini": confLocaleLocale_frFrIni,
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
	"conf": &bintree{nil, map[string]*bintree{
		"app.ini": &bintree{confAppIni, map[string]*bintree{}},
		"locale": &bintree{nil, map[string]*bintree{
			"locale_en-US.ini": &bintree{confLocaleLocale_enUsIni, map[string]*bintree{}},
			"locale_fr-FR.ini": &bintree{confLocaleLocale_frFrIni, map[string]*bintree{}},
		}},
	}},
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

