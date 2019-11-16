package utils

import (
	"archive/zip"
	"fmt"
	"github.com/mholt/archiver/v3"
	"path/filepath"
)

// ArchiveFile does just that
func ArchiveFile(filenameWithExtension string) error {
	files := []string{
		filenameWithExtension,
	}
	var extension = filepath.Ext(filenameWithExtension)
	var fileName = filenameWithExtension[0 : len(filenameWithExtension)-len(extension)]
	// archive format is determined by file extension
	err := archiver.Archive(files, fileName+".zip")
	if err != nil {
		return err
	}
	return nil
}

// UnarchiveIt does just that
func UnarchiveIt(filenameWithExtention string, targetFolder string) (string, error) {
	// The archive format is determined automatically.
	err := archiver.Unarchive(filenameWithExtention, targetFolder)
	if err != nil {
		return "", err
	}

	archivedFilename := ""
	err = archiver.Walk(filenameWithExtention, func(f archiver.File) error {
		zfh, ok := f.Header.(zip.FileHeader)
		if ok {
			fmt.Println("Filename:", zfh.Name)
			archivedFilename = zfh.Name
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return archivedFilename, nil
}
