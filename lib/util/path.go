package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PathObject struct {
	Abs  string
	Rel  string
	Dir  string
	Base string
	Name string
	Ext  string
}

func GetPathObj(f string) PathObject {
	absolutePath, err := filepath.Abs(f)
	if err != nil {
		fmt.Println("Error creating Path Object")
		return PathObject{}
	}
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "/"
	}
	relativePath, err := filepath.Rel(cwd, absolutePath)
	if err != nil {
		relativePath = f
	}
	name := ""
	base := filepath.Base(absolutePath)
	ext := filepath.Ext(absolutePath)
	if ext != "" {
		name = strings.Replace(base, ext, "", 0)
	} else {
		name = base
	}
	return PathObject{
		Abs:  absolutePath,
		Rel:  relativePath,
		Dir:  filepath.Dir(absolutePath),
		Base: base,
		Name: name,
		Ext:  ext,
	}
}

func (p PathObject) UpdateExt(ext string) PathObject {
	return GetPathObj(fmt.Sprintf("%s%c%s%s", p.Dir, os.PathSeparator, p.Name, ext))
}
