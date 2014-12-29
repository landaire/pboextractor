package main

import (
	"fmt"
	"os"
	"path"

	"io"

	"strings"

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

		file, err := os.Create(path.Join(outdir, entry.Name))
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
