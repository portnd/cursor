package helpers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
)

func Zip(outputFile string, files []string) error {
	// Create a buffer to store the zip file
	buf := new(bytes.Buffer)

	// Create a new zip writer
	zipWriter := zip.NewWriter(buf)

	// List of files to be added to the zip file
	// files := []string{"file1.txt", "file2.txt"}

	for _, filename := range files {
		if err := addFileToZip(zipWriter, filename); err != nil {
			fmt.Printf("Failed to add file '%s' to zip: %s\n", filename, err)
			return err
		}
	}

	// Close the zip writer
	if err := zipWriter.Close(); err != nil {
		fmt.Printf("Failed to close the zip writer: %s\n", err)
		return err
	}

	// Write the buffer to a zip file
	if err := writeBufferToFile(outputFile, buf.Bytes()); err != nil {
		fmt.Printf("Failed to write zip file: %s\n", err)
		return err
	}

	return nil
}

// addFileToZip reads the content of a file and adds it to the zip archive
func addFileToZip(zipWriter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get the file info
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Create a zip header
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("failed to create zip header: %w", err)
	}

	header.Method = zip.Deflate

	// Create a new entry in the zip archive
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("failed to create entry in zip archive: %w", err)
	}

	// Copy the file content to the zip entry
	_, err = io.Copy(writer, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content to zip entry: %w", err)
	}

	return nil
}

// writeBufferToFile writes the content of a bytes.Buffer to a file
func writeBufferToFile(filename string, content []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write content to file: %w", err)
	}

	return nil
}
