package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func getFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func findDuplicateFiles(directory string) ([]string, error) {
	filesDict := make(map[string]string)
	duplicates := []string{}

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileChecksum, err := getFileChecksum(path)
		if err != nil {
			return err
		}

		if duplicatePath, ok := filesDict[fileChecksum]; ok {
			// Eine Datei mit demselben Hash wurde bereits gefunden,
			// daher handelt es sich um eine Duplikatdatei.
			duplicates = append(duplicates, path, duplicatePath)
		} else {
			// Füge den Hash der Datei zum Dictionary hinzu.
			filesDict[fileChecksum] = path
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return duplicates, nil
}

func deleteDuplicateFiles(directory string) error {
	duplicates, err := findDuplicateFiles(directory)
	if err != nil {
		return err
	}

	if len(duplicates) == 0 {
		fmt.Println("Es wurden keine doppelten Dateien gefunden.")
		return nil
	}

	fmt.Printf("Es wurden %d doppelte Dateien gefunden:\n", len(duplicates)/2)
	for i := 0; i < len(duplicates); i += 2 {
		duplicateFile := duplicates[i]
		err := os.Remove(duplicateFile)
		if err != nil {
			fmt.Printf("Fehler beim Löschen von %s: %s\n", duplicateFile, err.Error())
		} else {
			fmt.Printf("%s wurde gelöscht.\n", duplicateFile)
		}
	}

	return nil
}

func main() {
	directoryPath := "C:\\Users\\rabraha\\Dropbox\\PC\\Desktop\\holograms_from_pdfs_1"

	err := deleteDuplicateFiles(directoryPath)
	if err != nil {
		fmt.Printf("Fehler beim Löschen der doppelten Dateien: %s\n", err.Error())

	}
}
