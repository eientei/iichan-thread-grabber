package model

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var __1_initialize_down_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x4a\x29\xca\x2f\x50\x28\x49\x4c\xca\x49\x55\x48\x2f\x4a\x4c\x4a\x4a\x2d\x8a\x2f\xc9\x28\x4a\x4d\x4c\x29\x8e\xcf\xcc\x4d\x4c\x4f\x2d\xb6\xe6\x02\x2b\x29\x4e\x2d\x2c\x4d\xcd\x4b\xc6\xa5\x2a\xbe\x20\xbf\x28\x25\xb5\x28\xbe\x38\xb5\x90\x48\x0d\x79\xe8\x1a\x50\x1d\x81\x62\x39\x56\xf7\x59\x03\x02\x00\x00\xff\xff\x02\xfd\xe0\x17\xbb\x00\x00\x00")

func _1_initialize_down_sql() ([]byte, error) {
	return bindata_read(
		__1_initialize_down_sql,
		"1_initialize.down.sql",
	)
}

var __1_initialize_up_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x52\xb1\x6e\xea\x40\x10\xec\xfd\x15\xdb\x9d\x91\x28\xde\x2b\xd2\x40\x85\x90\xdb\x48\x01\x52\x5b\x6b\x6e\x39\x9f\x72\x77\x3e\xd6\xeb\x04\xfe\x3e\x3a\x3b\x81\x04\x8c\x40\x4a\xbb\xbb\x33\xa3\x99\x9d\x2d\x13\x0a\x81\x60\xe5\x08\x0c\x63\x55\x11\x97\x52\x33\xa1\x6e\x21\xcf\x00\x3a\x76\x20\x74\x10\x88\x6c\x3d\xf2\x11\xde\xe8\x38\xcd\x00\x06\xa0\x06\xb1\x9e\x5a\x41\x1f\x21\x34\x02\xa1\x73\x0e\x34\xed\xb0\x73\x02\xa1\xf9\xc8\x27\xd9\x64\x9e\x65\xa3\x2a\xd6\xa3\xa1\x41\xc4\xeb\xa7\x51\x91\x9d\x75\x14\x51\xea\x61\x79\x25\xa0\xd4\x6c\x96\x36\xe9\x54\xea\xce\x57\x0f\xde\xc6\xa6\x15\xab\xc1\x06\x21\x43\x7c\x7d\xfb\xaf\x27\x44\xd3\x3e\xc0\xc5\x28\x36\x18\xd8\xd6\x38\x42\xa4\xf6\xaa\xd7\x43\xa6\x20\xe5\xc9\xe5\x4d\xc2\x3e\xac\xe5\xaa\x58\x6c\x0a\x58\x17\x2f\xaf\xc5\xf3\xb2\xb8\xfc\xca\x57\x6e\x65\x6c\x58\x13\x97\x2d\xed\x61\xbd\x59\xac\x36\xf0\x7f\xfe\x28\x34\x8c\x40\xc7\x7f\xf4\x1b\xd8\xff\x6a\x18\x95\xdf\xbd\x48\xfe\xfa\xed\xc9\xde\x90\x70\x52\xb8\x9d\x70\xa0\x83\xbc\xa3\xcb\xd5\x5d\x73\x6a\x92\xf8\xc2\x1f\xf9\xc2\x05\x5f\x34\xdc\x74\xf1\x4e\x03\x7e\x94\x31\x3f\xbb\x9e\x9e\xed\xa6\x72\x7f\x06\x00\x00\xff\xff\xb5\x32\x1c\x7f\x41\x03\x00\x00")

func _1_initialize_up_sql() ([]byte, error) {
	return bindata_read(
		__1_initialize_up_sql,
		"1_initialize.up.sql",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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
var _bindata = map[string]func() ([]byte, error){
	"1_initialize.down.sql": _1_initialize_down_sql,
	"1_initialize.up.sql":   _1_initialize_up_sql,
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
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func     func() ([]byte, error)
	Children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"1_initialize.down.sql": &_bintree_t{_1_initialize_down_sql, map[string]*_bintree_t{}},
	"1_initialize.up.sql":   &_bintree_t{_1_initialize_up_sql, map[string]*_bintree_t{}},
}}
