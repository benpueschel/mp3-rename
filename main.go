package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/dhowden/tag"
)

func main() {
	config := new(Config)
	parseFlags(config)

	if config.printHelp {
		printFlagHelp()
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		handleError(err)

		line = strings.TrimRight(line, "\n")

		file, err := os.Open(line)
		handleError(err)
		defer file.Close()

		name, err := GetNewFileName(file, config.format)
		handleError(err)

		var newPath string
		if len(config.outputPath) > 0 {
			newPath = path.Join(config.outputPath, name)
		} else {
			newPath = path.Join(path.Dir(line), name)
		}

		// re-open file because we already read some data and may want to copy
		file.Close()
		file, err = os.Open(line)
		handleError(err)
		defer file.Close()

		if config.stdout {
			fmt.Printf("%s\n", newPath)
		} else {
			if config.copy {
				os.MkdirAll(path.Dir(newPath), os.ModePerm)

				newFile, err := os.Create(newPath)
				handleError(err)
				defer newFile.Close()

				_, err = io.Copy(newFile, file)
				handleError(err)
			} else {
				err = os.Rename(line, newPath)
				handleError(err)
			}
		}
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}

func GetNewFileName(file *os.File, format string) (string, error) {
	m, err := tag.ReadFrom(file)

	if err != nil {
		return "", err
	}

	name := format
	name = strings.ReplaceAll(name, "{artist}", m.Artist())
	name = strings.ReplaceAll(name, "{title}", m.Title())
	name = strings.ReplaceAll(name, "{album}", m.Album())
	name = strings.ReplaceAll(name, "{album_artist}", m.AlbumArtist())
	name = strings.ReplaceAll(name, "{genre}", m.Genre())
	name = strings.ReplaceAll(name, "{composer}", m.Composer())
	name = strings.ReplaceAll(name, "{year}", fmt.Sprint(m.Year()))
	name = strings.ReplaceAll(name, "{comment}", m.Comment())
	name += path.Ext(file.Name())

	return name, nil
}

