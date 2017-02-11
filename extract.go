package main

import (
	"fmt"
	"os"
	"path"

	"io"

	"strings"

	"path/filepath"

	"github.com/codegangsta/cli"
)

func Extract(c *cli.Context) {
	if len(c.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "Extract error: No outdir argument provided")
		return
	}

	// TODO: Verbose option
	verbose := true

	outdir := c.Args().Get(1)
	fmt.Println(outdir)
	if pboFile.HeaderExtension != nil {
		prefix, ok := pboFile.HeaderExtension.ExtendedFields["prefix"]
		if ok {
			prefix = strings.Replace(prefix, "\\", "/", -1)
			fmt.Println(prefix)
			outdir = path.Join(outdir, prefix)
			fmt.Println(outdir)
		}
	}

	exists, err := exists(outdir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Extract error: unknown error occurred.", err)
		return
	}

	if !exists {
		err := mkdirIntermediate(outdir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Extract error: could not create output directory.", err)
			return
		}
	}

	// Start creating/copying files
	for _, entry := range pboFile.Entries {

		if verbose {
			fmt.Println("Extracting", entry.Name)
		}

		outfile := path.Join(outdir, filepath.ToSlash(entry.Name))

		if verbose {
			fmt.Printf("Extract file path: %s", outfile)
		}

		err := os.MkdirAll(path.Dir(outfile), 0777)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Extract error: could not create folder %s.\nError: %s\n", path.Dir(outfile), err)
		}

		file, err := os.Create(outfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Extract error: could not create file %s.\nError: %s\n", entry.Name, err)
		}

		entry.Seek(0, os.SEEK_SET)

		io.Copy(file, entry)

		file.Close()
	}

	if verbose {
		fmt.Println("Done")
	}
}
