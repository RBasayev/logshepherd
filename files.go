package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

func fileOpenOrCreate(filePathAndName string) (f *os.File, err error) {
	// try to open an existing file
	f, err = os.OpenFile(filePathAndName, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		// try to simply create the file
		f, err = os.Create(filePathAndName)
		if _, badPath := err.(*os.PathError); badPath {
			// if create failed and error is of type PathError, create path first
			err = os.MkdirAll(filepath.Dir(filePathAndName), 0777)
			isOK(err)
			// then retry creating the file
			f, err = os.Create(filePathAndName)
		}
		//err = os.Chmod(filePathAndName, 0666)
		//check(err)
	}
	return
}

func pipeOpenOrCreate(pipePathAndName string) (f *os.File, err error) {
	// check whether the pipe exists and create if it doesn't
	fCheck, err := os.Stat(pipePathAndName)
	if os.IsNotExist(err) {
		if _, err = os.Stat(filepath.Dir(pipePathAndName)); os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(pipePathAndName), 0777)
			isOK(err)
		}
		err = syscall.Mknod(pipePathAndName, 0666|syscall.S_IFIFO, 0)
		isOK(err)
		fCheck, err = os.Stat(pipePathAndName)
		isOK(err)
		err = os.Chmod(pipePathAndName, 0666)
	}
	isOK(err)
	if fCheck.Mode()&os.ModeNamedPipe == 0 {
		panic("Input is not a named pipe :(")
	}

	// open input and output and prepare for work
	f, err = os.OpenFile(pipePathAndName, os.O_RDWR, os.ModeNamedPipe)
	return
}

func considerRotating(fh *os.File, sizeMB int) (newFile *os.File, err error) {
	rotateAt := int64(sizeMB * 1048576)
	fileName := fh.Name()
	if rotateAt > 0 {
		if stat, _ := fh.Stat(); stat.Size() > rotateAt {
			fh.Close()
			timestamp := time.Now().Format("2006-01-02__15_04_05.9999")
			fmt.Println("Rotating: " + fileName + " ==> " + fileName + ".rotated." + timestamp)
			err := os.Rename(fileName, fileName+".rotated."+timestamp)
			isOK(err)
			newFile, err = fileOpenOrCreate(fileName)
			return newFile, err
		}
	}
	return fh, nil
}
