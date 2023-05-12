package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func convertObsidianToMarkdown(inputPath, imagePath, outputPath string) error {
	files, err := os.ReadDir(inputPath)
	if err != nil {
		return err
	}

	outputMarkdownPath := filepath.Join(outputPath, "markdown")
	outputImagePath := filepath.Join(outputPath, "picture")

	if err := os.MkdirAll(outputMarkdownPath, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(outputImagePath, os.ModePerm); err != nil {
		return err
	}

	for _, file := range files {
		inputFilePath := filepath.Join(inputPath, file.Name())
		outputFilePath := filepath.Join(outputMarkdownPath, file.Name())

		inputFile, err := os.Open(inputFilePath)
		if err != nil {
			return err
		}
		defer inputFile.Close()

		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		re := regexp.MustCompile(`!\[\[(.*?)\]\]`)
		scanner := bufio.NewScanner(inputFile)
		writer := bufio.NewWriter(outputFile)

		for scanner.Scan() {
			line := scanner.Text()
			matches := re.FindAllStringSubmatch(line, -1)

			for _, match := range matches {
				imageName := match[1]
				imagePathWithExt := filepath.Join(imagePath, imageName)
				if _, err := os.Stat(imagePathWithExt); os.IsNotExist(err) {
					return fmt.Errorf("image not found: %s", imagePathWithExt)
				}

				// Copy image to outputImagePath
				number := regexp.MustCompile(`\d+`).FindString(imageName)
				newImageName := fmt.Sprintf("%s%s", number, filepath.Ext(imageName))
				newImagePath := filepath.Join(outputImagePath, newImageName)

				srcImage, err := os.Open(imagePathWithExt)
				if err != nil {
					return err
				}
				defer srcImage.Close()

				dstImage, err := os.Create(newImagePath)
				if err != nil {
					return err
				}
				defer dstImage.Close()

				if _, err := io.Copy(dstImage, srcImage); err != nil {
					return err
				}

				outputImagePathMD := filepath.Join("../picture", newImageName)

				line = strings.Replace(line, match[0], fmt.Sprintf("![%s](%s)", imageName, outputImagePathMD), 1)
			}

			writer.WriteString(line + "\n")
		}

		writer.Flush()
	}

	return nil
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: obsidian-md-converter <input-path> <image-path> <output-path>")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	imagePath := os.Args[2]
	outputPath := os.Args[3]

	if err := convertObsidianToMarkdown(inputPath, imagePath, outputPath); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
