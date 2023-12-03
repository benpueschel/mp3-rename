package main

import (
	"flag"
)

type Config struct {
	outputPath string
	stdout bool
	copy bool
	format string
	printHelp bool
}

func parseFlags(config *Config) {
	flag.BoolVar(&config.stdout, "stdout", false, "Print new file names to stdouta instead of renaming them")
	flag.BoolVar(&config.copy, "c", false, "Copy files instead of renaming them. Only works if -o is false")
	flag.StringVar(&config.outputPath, "o", "", "The output path for file(s)")
	flag.StringVar(&config.format, "f", "{artist} - {title}", "The naming format for file(s)")
	flag.BoolVar(&config.printHelp, "h", false, "Print help")
	flag.Parse()
}

func printFlagHelp() {
	flag.PrintDefaults()
}
