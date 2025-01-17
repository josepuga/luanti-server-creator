package main

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func copy(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.Chmod(dst, sourceInfo.Mode()); err != nil {
		return err
	}

	return nil
}

// TODO: Implement concurrency...
func copyDir(src, dst string) error {

	//mkdir
	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return err
	}

	// Recorrer el contenido del directorio fuente
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Crear la ruta correspondiente en el destino
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, relPath)

		// Si es un directorio, crearlo en el destino
		if info.IsDir() {
			if err := os.MkdirAll(destPath, info.Mode()); err != nil {
				return err
			}
			return nil
		}
		// Si es un archivo, copiarlo
		return copy(path, destPath)
	})
}

func deleteDir(src string) error {
	// Delete files first
	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil // Directories will be delete after
		}

		if err := os.Remove(path); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Remove directories(*)
	// (*)I dont use only RemoveAll() because it ends at the 1st error.
	return os.RemoveAll(src)
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		// If an error occurs (e.g., path does not exist), it's not a directory
		return false
	}
	return info.IsDir()
}

func getDirectories(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var directories []string
	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, entry.Name())
		}
	}
	return directories, nil
}

func saveToFile(dst string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}
