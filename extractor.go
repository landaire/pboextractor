package main

import (
	"errors"
	"fmt"
	"os"

	"strings"

	"github.com/codegangsta/cli"
	"github.com/landaire/pbo"
)

const (
	version = "0.2.0"
)

var pboFile *pbo.Pbo

func main() {

	app := cli.NewApp()
	app.Name = "pboextractor"
	app.Usage = "Extract PBO archives used in games such as Arma 3"
	app.Author = "Lander Brandt"
	app.Email = "@landaire"
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "pbo",
			Usage: "PBO file to read",
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:      "extract",
			ShortName: "e",
			Usage:     "Extract the PBO to the given output directory",
			Before:    LoadPbo,
			Action:    Extract,
		},
		cli.Command{
			Name:   "header",
			Usage:  "Print header information to stdout",
			Before: LoadPbo,
			Action: PrintHeader,
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}

}

func LoadPbo(c *cli.Context) error {
	if c.GlobalString("pbo") == "" {
		return errors.New("No PBO provided")
	}

	if pboFile != nil {
		return nil
	}

	var err error
	pboFile, err = pbo.NewPbo(c.GlobalString("pbo"))

	if err != nil {
		return err
	}

	return nil
}

func PrintHeader(c *cli.Context) {
	// Print header extension info if it's present
	if pboFile.HeaderExtension != nil {
		fmt.Println("Header Extension:")
		lines := strings.Split(pboFile.HeaderExtension.String(), "\n")

		for _, line := range lines {
			fmt.Println("\t", line)
		}

		fmt.Println()

		fmt.Println("\tExtended Fields:")
		for key, val := range pboFile.HeaderExtension.ExtendedFields {
			fmt.Printf("\t\t%s: %s\n", key, val)
		}

		fmt.Println()
	}

	fmt.Println("Entries:")

	for _, entry := range pboFile.Entries {
		lines := strings.Split(entry.String(), "\n")

		for _, line := range lines {
			fmt.Println("\t", line)
		}

		fmt.Println()
	}
}
