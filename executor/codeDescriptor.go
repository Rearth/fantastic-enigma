package executor

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// ProjConf contains all necessary information for one exercise
type ProjConf struct {
	Title        string   `xml:"Title"`
	File         string   `xml:"file,omitempty"`
	FilePath     string   `xml:"path,omitempty"`
	Description  string   `xml:"desc"`
	Number       string   `xml:"number"`
	TemplateName string   `xml:"template"`
	Inputs       []string `xml:"inputs"`
}

// GenerateTemplate generated a config with empty values, to manually edit later
func GenerateTemplate(filepath, filename string) {
	c := &ProjConf{
		Title:        "Insert Here",
		Description:  "Insert Here",
		Number:       "Insert Here",
		TemplateName: "template.tpl",
		Inputs:       []string{"5", "1", "3"},
	}

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("  ", "    ")
	if err := enc.Encode(c); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

// ReadConf reads a configuration file (xml file) and generates a config struct from that
func ReadConf(filepath, filename string) ProjConf {

	var c ProjConf

	if !strings.HasSuffix(filename, ".xml") {
		return c
	}
	dat, err := ioutil.ReadFile(filepath + filename)
	if err != nil {
		fmt.Println("error reading config: " + err.Error())
	}
	fmt.Print(string(dat))

	return c

}
