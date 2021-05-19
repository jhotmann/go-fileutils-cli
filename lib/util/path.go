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
	relativePath, err := filepath.Rel(GetWorkingDir(), absolutePath)
	if err != nil {
		relativePath = f
	}
	name := ""
	base := filepath.Base(absolutePath)
	ext := filepath.Ext(absolutePath)
	if ext != "" {
		name = strings.Replace(base, ext, "", -1)
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

func (p PathObject) UpdateName(name string) PathObject {
	return GetPathObj(fmt.Sprintf("%s%c%s%s", p.Dir, os.PathSeparator, name, p.Ext))
}

func (p PathObject) UpdateExt(ext string) PathObject {
	return GetPathObj(fmt.Sprintf("%s%c%s%s", p.Dir, os.PathSeparator, p.Name, ext))
}

func (p PathObject) UpdateDir(dir string) PathObject {
	return GetPathObj(fmt.Sprintf("%s%c%s%s", dir, os.PathSeparator, p.Name, p.Ext))
}

func GetWorkingDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "/"
	}
	return cwd
}
