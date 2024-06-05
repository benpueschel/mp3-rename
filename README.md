# mp3-rename
A simple go app to fetch mp3 metadata and 

## Usage

```bash
find <directory> | parallel --pipe "./mp3-rename [opts]"
```

## Command-Line Options

- `-h`: Display help
- `-c`: Copy files instead of renaming them.
- `-o`: Target output directory for renamed/copied files. 
- `-f`: Target file format. Will replace templates with the mp3-tags. 
- `stdout`: Print file new names to stdout instead of copying or renaming them. 

## Format options

- `{artist}`
- `{title}`
- `{album}`
- `{album_artists}`
- `{genre}`
- `{composer}`
- `{year}`
- `{comment}`
