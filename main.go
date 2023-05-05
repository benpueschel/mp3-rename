package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/dhowden/tag"
)

func print_error(err error) {
	log.Fatalf("ERROR: %s\n", err.Error())
}

func main() {
	arguments := os.Args
	fmt.Println(arguments)
	
	config := new(Config)
	config.path = os.Args[len(os.Args)-1]
	parse_flags(config);

	stat, err := os.Stat(config.path)

	if err != nil {
		print_error(err)
		return
	}

	if stat.IsDir() {

		file_function := func(file_path string, file_info *fs.FileInfo) {
			file, err := os.Open(file_path)

			pattern := regexp.MustCompile(config.matchPattern)
			
			if pattern.MatchString(file.Name()) {

				if err != nil {
					print_error(err)
					return
				}
				
				name, err := get_new_file_name(file, config)

				if err != nil {
					print_error(err)
					return
				}


				output_path := path.Dir(file_path)
				if config.outputPath != "" {
					relative_path := file_path[len(config.path):]
					output_path = path.Join(config.outputPath, relative_path)
				}
				output_path = path.Join(output_path, name)

				if config.copy {
					copy_file(file_path, output_path)
				} else {
					os.Rename(file_path, output_path)
				}

			}

			file.Close()
		}
		walk_dir(config.path, config.recursive, file_function)

	} else {
		file, err := os.Open(config.path)

		if err != nil {
			print_error(err)
			return
		}
		
		name, err := get_new_file_name(file, config)

		if err != nil {
			print_error(err)
			return
		}

		output_path := config.outputPath
		if config.outputPath == "" {
			path.Join(path.Dir(config.path), name)
		}

		if config.copy {
			copy_file(config.path, output_path)
		} else {
			os.Rename(config.path, output_path)
		}
	}
}

func copy_file(src string, dest string) error {

	bytesRead, err := os.ReadFile(src)

    if err != nil {
		return err
    }

    err = os.WriteFile(dest, bytesRead, 0644)

    if err != nil {
		return err
    }
	return nil
}

type file_func func(path string, file_info *fs.FileInfo)

func walk_dir(file_path string, recursive bool, function file_func) error {
	files, err := os.ReadDir(file_path)
	if err != nil {
		return err
	}

	for _, de := range files {
		new_path := path.Join(file_path, de.Name())

		if(recursive && de.IsDir()) {
			walk_dir(new_path, recursive, function)
		}

		file_info, err := de.Info()
		if err != nil {
			return err
		}

		function(new_path, &file_info)
	}
	return nil
}

func get_new_file_name(file *os.File, config *Config) (string, error) {
	
	m, err := tag.ReadFrom(file)

	if err != nil {
		return "", err
	}

	dots := strings.Split(file.Name(), ".")
	ext := dots[len(dots)-1]

	name := config.format
	name = strings.ReplaceAll(name, "${artist}", m.Artist())
	name = strings.ReplaceAll(name, "${title}", m.Title())
	name = strings.ReplaceAll(name, "${album}", m.Album())
	name = strings.ReplaceAll(name, "${album_artist}", m.AlbumArtist())
	name = strings.ReplaceAll(name, "${genre}", m.Genre())
	name = strings.ReplaceAll(name, "${composer}", m.Composer())
	name = strings.ReplaceAll(name, "${year}", fmt.Sprint(m.Year()))
	name = strings.ReplaceAll(name, "${comment}", m.Comment())
	name += "." + ext

	return name, nil
}

type Config struct {
	path string
	outputPath string
	recursive bool
	copy bool
	format string
	matchPattern string
}

func parse_flags(config *Config) {

	flag.BoolVar(&config.recursive, "r", false, "Recursively iterate through directory")
	flag.BoolVar(&config.copy, "c", true, "Copy files instead of renaming them")
	flag.StringVar(&config.outputPath, "o", "", "The output path for file(s)")
	flag.StringVar(&config.format, "f", "${artist} - ${title}", "The naming format for file(s)")
	flag.StringVar(&config.matchPattern, "p", "", "The pattern to match files with")

}