package util

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func FindExecPath() (string, error) {
	// parameter 1 because its called from main and not inside main... hacky as f...
	_, exePath, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("failed to determine caller information")
	}
	dir := filepath.Dir(exePath)
	return dir, nil
}

func ExistsPath(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func FileNameWithoutExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func ZipFolder(source, target string) error {
	// Thanks to ChatGTP for this
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Get the base directory name to include in the zip
	baseDir := filepath.Base(source)

	// Walk through the source directory
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a header based on the file info
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Ensure the header's name includes the base directory
		relativePath := strings.TrimPrefix(path, source)
		if relativePath == "" {
			return nil // Skip the source directory itself
		}
		header.Name = filepath.Join(baseDir, relativePath)

		// Use the file's original mode
		header.Method = zip.Deflate

		if info.IsDir() {
			header.Name += "/"
		} else {
			// Write the file to the zip
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			return err
		}
		return nil
	})

	return err
}

func Unzip(src, dest string) error {
	// Thanks to ChatGTP for this
	// Open the ZIP file
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// Ensure the destination directory exists
	os.MkdirAll(dest, 0755)

	// Iterate through the files in the archive
	for _, f := range r.File {
		// Determine the path for the file
		filePath := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", filePath)
		}

		// Create directory tree
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		} else {
			if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}
		}

		// Open the file within the ZIP archive
		inFile, err := f.Open()
		if err != nil {
			return err
		}
		defer inFile.Close()

		// Create the file in the destination directory
		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		// Copy the file content to the destination file
		if _, err := io.Copy(outFile, inFile); err != nil {
			outFile.Close()
			return err
		}
		outFile.Close()
	}

	return nil
}

func LoadJSON[T any](path string) (T, error) {
	var data T

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func SaveJSON[T any](path string, data T) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
