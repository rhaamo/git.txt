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

var _confAppIni = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x57\xed\x72\xdb\xba\xd1\xfe\x8f\xab\xd8\xe8\x9d\xf3\x1e\xbb\x23\xdb\xb2\x7d\xea\xb8\x56\xd4\x09\x2d\x51\x36\x27\x92\xe8\x92\x54\x72\xdc\x4c\x86\x81\xc9\x95\x84\x31\x04\x28\x00\x68\x5b\xfd\xd1\xdb\xe9\x85\xf4\xc6\x3a\x0b\x92\xb6\xec\x74\xa6\x27\x99\x91\x12\x60\x77\xb1\xfb\x3c\xfb\xa5\xe0\xe6\x26\x9f\x05\xd3\x10\x06\xb0\x14\xee\xd0\x3d\x39\x36\x0c\x66\x79\x12\x5e\x45\x69\x16\x26\x30\x00\x67\x2a\x64\xc1\x2c\x9e\xdd\x4e\xe3\x79\x9a\x0f\x93\x30\xc8\xc2\xf6\x9c\x7d\xb5\x68\x1e\xd0\x7c\x63\x37\x49\x9c\xc5\xc3\x78\x02\x03\x58\x39\xb7\x61\xa3\x78\x1a\x44\x33\x18\x80\xd4\x05\x97\x2b\x6d\x1d\x4b\xe2\x38\xcb\xe7\x09\x89\xfc\xb2\xd7\xca\xef\xdb\x8b\xa3\xa3\x5f\xf6\x6a\xf1\x7d\x7b\xf1\xcb\xde\x75\x96\xdd\xe4\x37\x71\x92\xed\xdb\x23\x36\x9f\x45\xbf\xe7\x69\x3c\xfc\x14\x66\xf9\x4d\x98\x4c\xa3\x34\x8d\x62\x32\x7b\x76\x76\xc6\xbc\x64\x30\x1a\x91\x9b\xbd\x43\xff\x97\x3d\x6b\xc3\x00\x7e\xeb\xf5\x7a\x6c\x14\xa5\xc1\xe5\x24\xcc\x93\x78\x9e\x85\x49\x3e\x89\xaf\x60\x00\x0b\x2e\x2d\xb2\x3e\xcc\x37\x1b\x34\x20\xf1\x01\x25\xe8\x05\x38\x5c\x6f\x24\x77\x08\x5c\x95\x60\x1d\x77\xa2\x80\x85\x90\x08\x1b\xee\x56\xac\x0f\x25\x2e\x78\x25\x1d\x08\x0b\x6e\x55\x9f\xc2\xe3\x0a\x0d\xb6\xe0\xd1\x0d\x3e\x61\x51\x39\x2c\x59\x9a\x05\x59\x34\xcc\x7d\xd8\x37\x41\x76\x0d\x03\xc6\xbe\x96\xdc\xf1\x3b\x6e\xf1\x1b\xeb\x43\x28\xdc\x0a\x0d\x74\xd6\x5b\xfb\x43\x76\xba\xd0\xd9\x68\xeb\x96\x06\x6d\x07\xb4\x81\x8e\xfd\x21\x85\xc3\xd3\x4e\x17\xb6\xba\x82\x82\x2b\x28\xb4\x52\x58\x38\x70\x1a\x32\x31\xba\x84\x47\xe1\x56\x30\xdd\xa6\x7f\x9b\xc0\xc6\x68\xa7\x0b\x2d\xd9\xe8\x32\xcf\x6e\x6f\x88\xa1\x46\x9f\x5d\xc7\x29\xa1\x71\x7c\xf2\xde\x63\x74\x7c\x71\x7a\xda\x3b\x63\x2d\xeb\x7a\x69\xd9\x3c\xf5\x54\x1b\xad\x1d\xbb\x09\xd2\xf4\xcb\x08\x06\xac\x0f\x63\xf2\x62\xc7\x27\x25\xb7\x5d\xc0\xc6\xe9\x52\x58\x7e\x27\x91\xdc\x36\xf8\xa3\x12\x06\x6b\xaf\x1f\xd0\x88\xc5\xf6\x60\x51\x49\xd9\x61\x69\x3a\xc9\xa7\xf1\x88\x1e\x6a\xe4\x5b\xb3\x6d\x70\x1e\xea\x8e\x13\xe5\x5d\xa7\x0b\x95\x45\xe0\x77\x56\xcb\xca\xbd\xa0\xab\x7c\xf8\xd6\x71\xe3\x80\x5b\xa0\x84\x13\x05\xb2\x1a\xd1\x16\xf8\xc3\xf2\x8e\xb1\xaf\x06\x37\xda\x0a\xa7\xcd\x96\xe0\x4d\xb4\x76\xb5\x95\x85\x36\x60\x9d\x36\x42\x2d\xe1\x59\x46\xa0\xfd\xd5\x02\xf1\xd1\xdd\x25\xb6\xf3\xcf\xa3\x0f\x95\x45\xa3\xf8\x1a\xff\x7a\xb4\x14\xee\xc9\x1d\xec\xea\x74\x7c\x1e\x7b\x7c\xa2\x85\x77\xed\x91\x2b\x4f\x49\x13\x21\x14\x52\x2b\x7a\xa9\xb2\xf4\x49\x09\xf9\x9c\x84\x3e\x3b\xaf\xa2\x6c\x27\x05\x6f\x1b\x6e\xed\x06\x0b\xb1\xd8\x02\x01\x57\x7b\xed\x34\x45\x07\x62\x01\x4a\xab\x03\xeb\xb8\x2a\xb9\x29\x59\x1f\xae\x84\x83\x3b\xa1\xb8\xd9\x92\xc3\x44\x0b\x21\x57\xfa\x30\xe9\x05\xaf\x56\x6d\xa4\xe6\xe5\xc1\x86\x17\xf7\xb0\x57\xac\xb0\xb8\xd7\x95\x03\xfd\x80\xb5\xcc\x3e\xbb\x8a\xb2\xfc\x32\x9a\x05\xc9\x6d\x0d\x23\x63\x5f\xa5\x5e\x7e\x63\xaf\xf2\xb5\x0f\x43\xae\xe0\x0e\xa1\x53\x68\x65\xb5\xc4\x86\x30\xaa\x8a\xce\x6b\xdc\x5a\x01\x2a\x2b\x8b\x50\xe8\xf5\x9a\x53\x0c\x16\x37\xdc\x50\x51\xad\x2b\xe9\xc4\x46\x22\xac\x75\x89\xb6\x0b\x78\xb8\x3c\x7c\x56\xeb\xfa\x42\xeb\xb0\x26\x5b\x9a\x53\xd6\x87\xcb\x6a\xb1\xf0\x35\xaa\x96\x6e\x45\x45\x5a\xac\xb8\x52\x28\xbb\x70\x8f\xb8\x01\xe1\x93\x42\x78\x17\x44\xcd\x47\xa9\xd5\xaf\x0e\xee\x95\x7e\x84\xc7\x15\x77\xf5\xe5\x21\xbb\x9c\x8f\xc7\xd4\x00\x42\x6a\x1f\xc7\xbd\xde\x4e\x01\x66\x86\x17\x3e\x93\x23\xb5\xd0\xf4\xfd\x85\x1b\x45\xdf\xa1\x31\xda\xd0\x3f\xc6\xdc\x71\xf9\x26\xe0\x5a\x8b\x4d\xc2\xcf\x21\xf5\x34\xff\x5f\xd6\xe6\xf7\x33\x5c\x14\xac\xe7\xc8\xc3\x7b\xd8\x9c\x53\x7e\x4a\xe4\x0f\x08\xb8\xde\xb8\x2d\xe1\x24\xd4\x0a\x8d\x70\xad\xbd\x67\x4b\x1e\x97\xb7\x66\xe8\xf0\x0f\xd8\xe8\x43\xb6\xa2\x96\xa4\x28\x31\x2d\xf0\xca\xe9\x35\x77\x58\x82\xd4\x4b\x30\xda\x11\x2d\x7b\xf6\x51\xb8\xc2\x43\xbb\xd0\x52\xea\x47\xca\x5a\xbd\x71\x42\x2b\xbb\xcf\x26\xf1\x55\x9e\xc4\xd9\x4e\xcf\xef\x43\x8a\xcb\x35\x2a\xe7\x8d\x94\x5c\xc8\x2d\x1b\x05\xd1\xe4\xf6\x27\xb9\x29\x7f\x02\x2b\xfe\x81\x60\x57\x62\xe1\xe8\x01\xaa\x08\x89\x9e\xea\x57\x58\x9e\x9c\xc3\x1a\xb9\xb2\x70\x0c\x1f\x3e\xc0\xc9\x79\x17\x4e\xfe\x7c\x36\xbd\x64\xd3\xe0\xf7\x3c\x8d\xfe\x1e\xe6\xe9\x75\x34\xa6\xaa\x39\x39\x6f\xec\x4a\xa1\x10\x54\xb5\xbe\x43\xf3\xc6\xb0\x57\x9a\x44\xb3\x30\xad\x79\xa6\x3f\xc4\xf5\xd3\x46\x18\x2c\xa1\xe4\x5b\x4b\x1a\xe4\xbc\xef\xed\x7b\x25\x4a\xa4\x9e\xbf\x70\x68\x60\xcd\x9f\xbc\xc8\xbe\x37\x33\x0a\x6e\xc9\xca\xfb\x67\x32\xac\xe4\xc5\xfd\x4f\x6c\xf8\xd3\x3f\x44\xc7\x17\xbc\x5b\x69\x7d\x0f\xf3\x64\xc2\xfc\x28\xac\xab\xee\xf0\x49\x9b\xb5\x1f\x09\x9e\xa9\xda\x31\xcf\x8f\xd0\x8a\xbd\xc5\x35\xa9\x89\xc3\x07\x34\x5b\x72\xb6\x11\xc8\x3d\x0b\x3f\x89\x69\x55\x34\x06\x3d\x17\xf8\x54\xa0\xb5\x68\xe1\x09\x76\x00\x7e\x2e\x89\x29\x7f\x12\xeb\x6a\x5d\xc3\xe4\x74\x5d\x64\x52\x2f\x97\x68\xbc\x11\xbb\x8b\xcb\xa9\x5f\x00\xac\x15\x5a\xbd\x9a\x67\xb8\xd6\x66\x4b\x65\xd3\xb4\x09\x42\xce\x60\x29\xec\x9b\x0a\x6a\x04\x69\x79\xf8\x1c\x8d\xfc\x0c\xf2\x0c\xf6\xe1\xc6\xe8\x07\x51\xa2\xa1\x46\xb0\x10\xcf\xf9\xc8\xfa\x50\xeb\x5c\x80\xd2\x0e\x56\x84\x36\x57\xdb\x56\x6a\x8b\x8e\xf5\xbd\x89\x0b\x68\xfc\x7a\x99\xdf\x4d\xc3\xf9\x4e\x1d\xff\xa8\xb9\xb5\xdf\x59\x1f\xbc\x67\x17\xa0\xd0\x3d\x6a\x73\x3f\x70\xc5\xa6\xcb\xcb\xd2\x0c\x2e\xce\x4e\xdf\xff\xa5\xbb\xe1\xd6\x3e\x6a\x53\x0e\xd6\xbc\xe0\x46\xab\x6e\x79\x37\xe8\x75\x37\x5a\xcb\x9c\xf0\x1c\x1c\xf7\x7a\x5d\x51\x4a\xcc\x9d\x58\xa3\xae\xdc\xe0\xf8\x9c\x70\xf4\x33\xfd\x02\x96\xfa\xc0\xfe\x90\x07\xa5\x11\x0f\x68\x8e\xfc\x21\x94\x56\xb5\x0e\x5b\x47\x33\xa9\xf5\x8c\xe6\xef\x45\xfb\xde\xc7\xd6\xc7\xdc\x51\x4a\x7c\x7f\x06\x29\x1f\xc6\xb3\x71\x44\x5b\xcc\xab\x48\x7c\x5d\xd6\x21\x17\x5a\xdf\x0b\x04\x9a\x60\x6c\x18\xc7\x9f\xa2\xb0\xdd\xf1\x44\x2e\xc5\x3d\xe6\x1b\xad\x04\xda\x97\x09\x46\x73\xb7\xc5\x4b\x28\xbf\xc0\xd9\x66\xdc\xef\xd0\x55\x4f\xac\xc6\x60\x1a\x0e\xe7\x49\xb8\x33\xc7\x9a\xcc\xb5\xe8\x9a\xf7\x5f\xe9\xfa\x94\x0c\x67\x7e\x08\xa6\x61\x96\xd7\x56\x76\x5b\x4a\xfd\xfa\xd5\x10\x08\x47\x10\xca\xa1\x79\xe0\xf2\x95\x91\xd3\xb3\x5e\x8f\x5d\x0d\xf3\x68\x96\x85\xc9\xe7\x60\x92\x67\x91\x8f\xca\x9f\xbf\xd8\x90\x62\x81\xde\xca\x2b\xe5\xf3\xb3\xdf\x7a\x3d\x96\x86\x7e\x81\xcc\x27\xd1\x38\x6c\xd5\xeb\x9b\x3e\x0c\x5f\x60\xf3\xa3\x74\x98\x26\x63\x46\x1f\xf9\x6b\x10\xf3\xc2\x9a\x85\xcf\xfc\xa2\x32\xc2\xf9\x5d\xe3\xdd\xbb\xe1\x75\x30\xbb\x0a\x21\xbb\x8e\x52\xc8\x62\xf8\x14\x86\x37\x70\x1b\xcf\x13\xf0\xcb\xd5\x28\xc8\x02\x48\x83\x71\xf8\xee\x1d\x4b\xc3\x61\x12\x66\xf9\xa7\x90\x6a\xf5\xdd\xff\x7d\x1c\x8f\xc2\x2f\x49\xf8\x25\xf9\xff\x3f\xed\xb1\x68\x96\x66\xc1\x64\x92\x4f\xe2\xe1\xa7\x1d\x70\x83\xca\xe9\x03\xa9\x97\x42\x81\xc1\x35\xfa\xae\x47\x15\x4a\xed\x39\xa2\x75\x7d\x1a\x4e\x2f\xc3\xe4\xa5\x5f\x35\x1e\xd3\xdb\x2f\xeb\xfd\x93\xcb\x0b\x94\x68\x9d\xe0\xad\xc0\xb3\xe6\x2b\x29\x59\x29\xfe\xbf\x88\xf6\x93\xbd\xe5\xba\x6e\x76\xa5\x28\xa8\xe1\xd0\xf6\x04\xb5\xb3\xb4\x43\x57\xb6\x25\xbe\x76\x96\xd6\x62\xfa\x21\xd1\x66\x40\x6d\xf7\xbf\xdc\xb5\x3e\x79\x53\x79\x63\x8a\x7d\x15\xc7\xe7\xea\x1b\x9b\x04\xb3\x2b\x0a\x15\xd5\xc1\x3c\xed\x2e\xcc\xc1\x38\xf1\x2b\x2d\x9d\x85\x6a\x29\x85\x5d\x75\xc7\x86\xab\x7f\xff\x8b\x0b\xd2\x2a\x4c\xd3\xa3\x6a\xef\x4d\xa5\xfc\x8e\x46\xc7\xe0\xb8\xbd\xb7\xb0\x41\x23\x34\x85\x20\xe5\xf6\xb0\x71\x79\xb4\xd3\x4d\x2b\xb5\x2b\xed\x17\xd3\x2b\x42\xab\x5e\x4d\xed\x21\x4b\xe6\xb3\x3c\xc8\x28\x86\xe4\x65\xc1\xa3\xbc\x92\xc8\x55\xb5\x79\xd9\x3c\xb7\xc0\x4d\xb1\x12\x0f\x68\x6b\xbf\x0e\xe9\x26\x6f\xce\xf2\xa2\x16\xff\xf6\xd6\x9e\xf7\x23\x1d\x5e\x87\xa3\xf9\x84\x60\xf9\x58\xf7\xff\x93\xdf\xe8\x97\x49\x46\x45\x53\x56\xc6\x0f\x0c\x62\xc3\x6f\x7c\xb4\x11\x35\x66\xc1\xae\x74\x25\x4b\x5a\xe5\xfc\x03\x58\xb2\x78\x42\xdd\x24\xbb\x0e\x68\x21\x22\x33\xff\x09\x00\x00\xff\xff\x04\x33\xc3\xaa\x0d\x0e\x00\x00")

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

	info := bindataFileInfo{name: "conf/app.ini", size: 3597, mode: os.FileMode(420), modTime: time.Unix(1494487952, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _confLocaleLocale_enUsIni = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x54\xcb\x8a\xec\x36\x10\xdd\xfb\x2b\x0a\x6f\xb2\xb9\xf4\x07\x04\xb2\x08\x93\x09\x0c\x84\x10\x26\x24\x10\x86\x41\xa8\xad\x1a\xbb\xb8\xb2\x64\x4a\xe5\x7e\xfc\x7d\x28\xc9\x0f\xb9\x67\xee\xaa\xe5\xf3\x50\xab\x4e\x95\xf4\xe6\x63\x4f\xe1\xbd\x11\x12\x8f\xf0\x0b\xb4\x7f\x53\x1f\xe0\x25\xb4\x4d\xa2\x3e\x18\x0a\x07\x6c\x4e\xc8\xc1\x8e\x59\xf8\xcf\xb2\xde\x51\x33\x79\xdb\xe1\x10\xbd\x43\x56\xc5\x7f\x71\x66\x58\xc9\x53\xdb\xe0\x68\xc9\x2b\xf1\xac\x8b\xe5\xfb\x4b\x53\x66\x4e\x6d\x33\xd9\x94\xae\x91\x9d\x12\x7f\x2d\xeb\x1d\xfd\xd2\x1a\xa8\x43\x58\x15\xa7\xb6\x61\x9c\xd0\x8a\xf9\x91\xe7\xd7\xe0\x40\x06\x84\xa4\x55\xc5\x80\x30\x20\x63\xb6\x8d\x38\x9e\x91\x4d\x29\xf6\x75\xf9\x04\xad\xb7\x79\x63\xec\x29\x09\x72\x95\xdb\xeb\x02\xa9\xb5\xac\x8e\x68\x88\x62\xac\xf7\xf1\x8a\x6e\x27\xd8\x0a\xc5\x00\x21\x0a\x2c\x9c\xee\xfe\x11\x79\x7c\xdf\xab\x54\xe7\x68\xa5\x1b\xea\x14\x12\xb8\x88\x29\xfc\x24\x90\xa9\xba\x09\xab\x8f\x42\x17\x99\xb1\x13\xf5\xbd\x84\x8b\xf5\xe4\xb6\x76\x40\xe4\x2d\xa5\xca\x7c\x46\x0c\x46\xec\x77\x0c\x75\x8b\xc1\x7a\x46\xeb\xee\x90\x99\x6f\x90\x22\xf3\xbd\x72\x31\x26\xe4\x4b\x29\x6c\xf3\xac\xe0\x37\x98\x3c\xda\x84\xd0\x0d\x31\x26\x04\x1b\xa2\x0c\xb8\x0f\xc6\xe1\xe8\x22\xc8\xc1\x3c\x64\xf5\xe9\xec\x8b\xae\x6d\x92\xbd\x94\xa1\xd5\xdf\x6e\xb0\xa1\xc7\xd4\x36\xcf\x87\x39\xfb\x19\xd6\x49\x43\xe6\xc8\xf5\x86\x65\x20\x35\x6e\x2b\x1a\x7c\x42\x11\x0a\x7d\xaa\xda\xaa\xe5\xc0\x8a\xb7\xcd\x3c\x39\x2b\x68\x26\x8e\x1f\xe4\xd1\xa4\xb9\xeb\x30\xa5\xdc\x98\x02\x81\x1e\xc9\xc1\x42\x7c\xcc\xde\xdf\xeb\x9d\x57\x67\xf5\x0f\x9b\x71\xfb\x93\xe6\xad\x27\xb9\x89\x09\x78\xad\x74\x7f\xe2\x15\x7a\x92\x93\xdc\xa4\x6d\xd4\x91\x4c\x17\x83\x60\x10\xd3\xd9\xa0\x91\xe1\x38\xc9\x5d\xa5\x4f\x05\x87\xf8\x01\x59\x08\x45\x00\x67\x84\xac\x69\x9b\x1c\x85\x39\xec\x92\x73\x1c\xe2\xec\x9d\x4e\xd5\x2e\x5d\x4e\xf3\xde\x38\x4c\x1d\xd3\xa4\x13\xfb\x78\x8b\x7e\xdb\xa9\xd3\xe9\x54\x8e\x97\x77\xff\xea\x59\xf8\x7d\xc1\xe1\x4a\x32\x00\xde\x04\x43\xa2\x18\x16\xd7\x5a\xd3\x83\x69\x2d\x69\xb9\x9e\x9a\xb2\x99\x98\x2e\x56\xf6\x01\x58\xbe\x57\x76\x3e\x7b\xea\x76\x32\x7f\xb6\x8d\x75\x4e\x83\xcd\xa7\xcb\x0f\x80\x73\x10\xf0\x9a\x83\xda\xa3\xbf\xd0\x21\xfb\x7f\xa9\x0e\x7f\x15\x79\x4a\x52\x89\xfe\xa0\x94\x23\x5f\x74\x55\x23\xd1\x51\x2d\x7c\x76\x24\x9f\x77\xeb\x55\x53\xfa\x42\x81\xc4\x30\x4e\x31\x91\x44\xbe\x97\x91\x25\x81\x0a\xca\xc2\xb5\x8f\x8c\x63\xbc\xe0\x83\xe3\xa9\xf4\xbc\x70\x95\xb5\x6d\x3a\x6d\x72\xbe\x64\x71\x2a\x17\xfd\x49\x91\xfc\x0c\x65\xa4\x16\xeb\xda\xe8\x43\x32\x4f\x62\x22\xef\x43\xf6\xba\xff\xdb\x42\xa3\xd3\x57\xe5\x30\x61\x1d\xa3\xde\x97\xb3\x8f\xe7\x5c\x79\xbe\x81\x19\xa4\xd0\x83\xc2\xab\xb2\x47\x31\x14\x1c\xde\x76\x5d\x5f\x6e\x04\xbc\x28\xbc\xea\xb4\x7f\x18\xa4\xd4\x58\x74\xd6\x39\x95\x65\xb4\xdd\x22\x74\x78\x33\x57\x26\x41\x23\x8c\xb8\xab\xaf\x83\x76\x5e\x19\x35\x29\xb7\x7a\x7c\x8c\xdf\xe7\xe9\x41\x5e\xc0\x83\x4e\x8f\x3a\xa0\x75\x87\x93\x82\x22\x8f\x0a\xd3\xc5\x71\x24\xf9\x2c\x84\x42\x6c\x31\x6d\xb2\xf2\x46\xad\xec\xff\x01\x00\x00\xff\xff\xbc\x09\xff\xac\x9e\x07\x00\x00")

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

	info := bindataFileInfo{name: "conf/locale/locale_en-US.ini", size: 1950, mode: os.FileMode(420), modTime: time.Unix(1494395450, 0)}
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

