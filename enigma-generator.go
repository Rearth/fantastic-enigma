package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Rearth/fantastic-enigma/executor"
	"github.com/urfave/cli"
)

func main() {

	var language string

	app := cli.NewApp()
	app.Name = "Lab Report Generator"
	app.Usage = "Used to generate lab reports from a LaTeX template and some pieces of code"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "lang, l",
			Value:       "english",
			Usage:       "language for the greeting",
			Destination: &language,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"gen"},
			Usage:   "generate a full lab report",
			Action: func(c *cli.Context) error {
				fmt.Println("generating report...")
				fmt.Println(language)
				return nil
			},
		},
		{
			Name:  "run",
			Usage: "run the code in the codesrc folder",
			Action: func(c *cli.Context) error {
				runCode(app)
				return nil
			},
		},
		{
			Name:  "conf",
			Usage: "Used to read a configuration vile, and check it for errors",
			Action: func(c *cli.Context) error {
				runConf(app)
				return nil
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Printf("language: %s\n", language)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runConf(app *cli.App) {
	fmt.Println("reading confings!")
}

func runCode(app *cli.App) {
	fmt.Printf("running code, args: %s\n", app.Version)

	//working dir
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	//get files
	p := filepath.FromSlash(dir + "/codesrc/")
	files, err := ioutil.ReadDir(p)
	if err != nil {
		log.Fatal(err)
	}

	responses := make([]string, 5)
	for i := 3; i < 5; i++ {
		responses[i] = strconv.Itoa(i)
	}

	for _, f := range files {
		//executor.ReadConf(p, f.Name())
		runFile(p, f.Name(), responses)
	}

	executor.GenerateTemplate(p, "test.xml")
}

func runFile(path, filename string, responses []string) {

	data := executor.CodePacket{
		Title:    "Test",
		Path:     path,
		Filename: filename,
		Inputs:   responses,
		Timeout:  10,
		RunData:  &executor.RunInfo{},
	}
	//fmt.Printf("created data: %+v\n", data)

	data.Run()
	fmt.Printf("result: %+s\n", data.RunData.OutStd)

}
