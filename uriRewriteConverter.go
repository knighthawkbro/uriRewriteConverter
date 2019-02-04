package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/knighthawkbro/urlRewrite/lib"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	var usage = `uriRewriteConverter takes in a flag and a file and converts it from ht.acl 
rewrite format to Microsoft web.config XML or vice versa.

Usage: uriRewriteConverter [-v] {-a|-x} <FileName>
`
	verbose := flag.Bool("v", false, "Sends output to STDOUT instead of a file")
	htacl := flag.Bool("a", false, "Converts a HT.ACL file to Microsoft Web.Config XML")
	web := flag.Bool("x", false, "Converts a Web.Config XML back to a HT.ACL")

	flag.Usage = func() {
		fmt.Printf(usage)
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("ERROR: No file specified...")
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}

	x := &lib.Configuration{}
	a := &lib.HTACL{}

	file := flag.Arg(0)
	f, err := os.Open(file)
	lib.CheckErr("File not found", err)
	defer f.Close()

	if *htacl {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			s := strings.TrimSpace(scanner.Text())
			a.Unmarshal(strings.Split(s, " "))
		}
		x = a.ToWebConfig()
	} else if *web {
		input, err := ioutil.ReadAll(f)
		lib.CheckErr("Unable to read file", err)
		x = lib.Unmarshal(input)
		a = x.ToHTACL()
	} else {
		fmt.Println("ERROR: No file format parameter set...")
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}

	if *verbose {
		if *htacl {
			fmt.Println(x.Marshal())
		} else if *web {
			fmt.Println(a.Marshal())
		}
	} else {
		if *htacl {
			file, err := os.Create("./web.config")
			lib.CheckErr("Could not create file ./web.config", err)
			defer file.Close()
			bytes, err := file.WriteString(x.Marshal())
			lib.CheckErr("Error writing to file ./web.config", err)
			fmt.Printf("Wrote %d bytes to file ./web.config\n", bytes)
			file.Sync()
		} else if *web {
			file, err := os.Create("./ht.acl")
			lib.CheckErr("Could not create file ./ht.acl", err)
			defer file.Close()
			bytes, err := file.WriteString(a.Marshal())
			lib.CheckErr("Error writing to file ./ht.acl", err)
			fmt.Printf("Wrote %d bytes to file ./ht.acl\n", bytes)
			file.Sync()
		}
	}
}
