// Copyright 2014-2017 Verizon Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Debug bool

var dbg Debug = false

func (d Debug) Printf(s string, a ...interface{}) {
	if d {
		fmt.Printf(s, a...)
	}
}

// xmlOutput reads the profile data from profile and generates an XML
// coverage report, writing it to outfile. If outfile is empty, it
// writes the report to a temporary file and opens it in a web browser.
func xmlOutput(profile, outfile string) error {
	profiles, err := ParseProfiles(profile)
	if err != nil {
		return err
	}

	currPkg := ""
	currFile := ""
	utc := time.Now().UTC().Unix()
	covEl := CoverageElement{Clover: "cover2clover.go", Generated: utc}
	covEl.Project.Timestamp = utc
	covEl.Project.Name = "" // TODO: add -p <project> flag
	var pkgEl *PackageElement
	var fileEl *FileElement

	var d templateData

	for _, profile := range profiles {
		// TODO: Can this change between profiles?
		if profile.Mode == "set" {
			d.Set = true
		}

		// pathname relative to $GOPATH
		path := profile.FileName
		// final element of the path is the file name
		fname := filepath.Base(path)
		// penultimate element of the path is the package name
		// TODO: Handle situation when package name != dir name
		pkg := filepath.Base(filepath.Dir(path))

		file, err := findFile(path)
		if err != nil {
			return err
		}
		if err != nil {
			return fmt.Errorf("can't read %q: %v", path, err)
		}
		if pkg != currPkg {
			// Aggregate metrics for the package element just completed
			if currPkg != "" {
				covEl.Project.Metrics.Complexity += pkgEl.Metrics.Complexity
				covEl.Project.Metrics.Statements += pkgEl.Metrics.Statements
				covEl.Project.Metrics.CoveredStatements += pkgEl.Metrics.CoveredStatements
				covEl.Project.Metrics.Elements += pkgEl.Metrics.Elements
				covEl.Project.Metrics.CoveredElements += pkgEl.Metrics.CoveredElements
				covEl.Project.Metrics.Classes += pkgEl.Metrics.Classes
				covEl.Project.Metrics.Loc += pkgEl.Metrics.Loc
				covEl.Project.Metrics.Ncloc += pkgEl.Metrics.Ncloc
			}
			currPkg = pkg
			pkgEl = &PackageElement{Name: pkg}
			covEl.Project.Package = append(covEl.Project.Package, pkgEl)
			covEl.Project.Metrics.Packages++
		}
		if file != currFile {
			currFile = file
			fileEl = &FileElement{Name: fname, Path: path}
			pkgEl.File = append(pkgEl.File, fileEl)
			pkgEl.Metrics.Files++
			covEl.Project.Metrics.Files++

			dbg.Printf("path=%s, pkg=%s, file=%s\n", path, pkg, fname)
		} else {
			dbg.Printf("*** Should never get here ***\n")
		}

		// For each file, process function declarations (as classes) and lines...
		err = xmlGen(file, profile, pkgEl, fileEl)
		if err != nil {
			return err
		}
		pkgEl.Metrics.Complexity += fileEl.Metrics.Complexity
		pkgEl.Metrics.Statements += fileEl.Metrics.Statements
		pkgEl.Metrics.CoveredStatements += fileEl.Metrics.CoveredStatements
		pkgEl.Metrics.Elements += fileEl.Metrics.Elements
		pkgEl.Metrics.CoveredElements += fileEl.Metrics.CoveredElements
		pkgEl.Metrics.Classes += fileEl.Metrics.Classes
		pkgEl.Metrics.Loc += fileEl.Metrics.Loc
		pkgEl.Metrics.Ncloc += fileEl.Metrics.Ncloc
	}
	// Aggregate metrics for final package element
	covEl.Project.Metrics.Complexity += pkgEl.Metrics.Complexity
	covEl.Project.Metrics.Statements += pkgEl.Metrics.Statements
	covEl.Project.Metrics.CoveredStatements += pkgEl.Metrics.CoveredStatements
	covEl.Project.Metrics.Elements += pkgEl.Metrics.Elements
	covEl.Project.Metrics.CoveredElements += pkgEl.Metrics.CoveredElements
	covEl.Project.Metrics.Classes += pkgEl.Metrics.Classes
	covEl.Project.Metrics.Loc += pkgEl.Metrics.Loc
	covEl.Project.Metrics.Ncloc += pkgEl.Metrics.Ncloc

	var out *os.File
	if outfile == "" {
		var dir string
		dir, err = ioutil.TempDir("", "cover")
		if err != nil {
			return err
		}
		out, err = os.Create(filepath.Join(dir, "coverage.xml"))
	} else {
		out, err = os.Create(outfile)
	}
	exitOnError(err, 1)

	if xmlstring, err := xml.MarshalIndent(covEl, "", "    "); err == nil {
		xmlstring = []byte(xml.Header + string(xmlstring))
		//		dbg.Printf("%s\n", xmlstring)
		out.Write(xmlstring)
	}
	if err == nil {
		err = out.Close()
	}
	if err != nil {
		return err
	}

	if outfile == "" {
		if !startBrowser("file://" + out.Name()) {
			fmt.Fprintf(os.Stderr, "XML output written to %s\n", out.Name())
		}
	}
	return nil
}

func exitOnError(err error, code int) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(code)
	}
}

func xmlGen(file string, profile *Profile, pkgEl *PackageElement, fileEl *FileElement) error {
	funcs, err := findFuncs(file)
	if err != nil {
		return err
	}

	// Now correlate functions and profile blocks.
	for _, f := range funcs {

		fnEl := &ClassElement{Name: f.name}
		fileEl.Class = append(fileEl.Class, fnEl)
		covered, total := f.coverage2(profile, fnEl)
		fnEl.Metrics.Complexity = f.complexity
		fnEl.Metrics.Statements = total
		fnEl.Metrics.CoveredStatements = covered
		fnEl.Metrics.Elements = total
		fnEl.Metrics.CoveredElements = covered

		fileEl.Metrics.Complexity += f.complexity
		fileEl.Metrics.Statements += total
		fileEl.Metrics.CoveredStatements += covered
		fileEl.Metrics.Elements += total
		fileEl.Metrics.CoveredElements += covered
		fileEl.Metrics.Classes++

		// Use source lines of code (LOC) == total statements
		fileEl.Metrics.Loc += total
		fileEl.Metrics.Ncloc += total

		dbg.Printf("%s: start=%d.%d end=%d.%d\n", f.name, f.startLine, f.startCol, f.endLine, f.endCol)
	}
	// Avg function complexity per file
	//    fileEl.Metrics.Complexity = fileEl.Metrics.Complexity / fileEl.Metrics.Classes
	return nil
}

// coverage2 returns the fraction of the statements in the function that were covered, as a numerator and denominator (similar to coverage() in func.go).
func (f *FuncExtent) coverage2(profile *Profile, fnEl *ClassElement) (num, den int) {
	// We could avoid making this n^2 overall by doing a single scan and annotating the functions,
	// but the sizes of the data structures is never very large and the scan is almost instantaneous.
	var covered, total int
	// The blocks are sorted, so we can stop counting as soon as we reach the end of the relevant block.
	for _, b := range profile.Blocks {
		if b.StartLine > f.endLine || (b.StartLine == f.endLine && b.StartCol >= f.endCol) {
			// Past the end of the function.
			break
		}
		if b.EndLine < f.startLine || (b.EndLine == f.startLine && b.EndCol <= f.startCol) {
			// Before the beginning of the function
			continue
		}
		total += b.NumStmt
		if b.Count > 0 {
			covered += b.NumStmt
		}
		for i := 0; i < b.NumStmt; i++ {
			lnEl := &LineElement{}
			fnEl.Line = append(fnEl.Line, lnEl)
			lnEl.Num = b.StartLine + i
			lnEl.Type = "stmt"
			lnEl.Count = b.Count
		}
	}
	if total == 0 {
		total = 1 // Avoid zero denominator.
	}
	return covered, total
}

type CoverageElement struct {
	XMLName   xml.Name       `xml:"coverage"`
	Clover    string         `xml:"clover,attr"`
	Generated int64          `xml:"generated,attr"`
	Project   ProjectElement `xml:"project"`
	//    TestProject         TestProjectElement  `xml:"testproject"`
}

type ProjectElement struct {
	XMLName   xml.Name          `xml:"project"`
	Name      string            `xml:"name,attr,omitempty"`
	Timestamp int64             `xml:"timestamp,attr"`
	Metrics   ProjectMetrics    `xml:"metrics"`
	Package   []*PackageElement `xml:"package"`
}

type PackageElement struct {
	XMLName xml.Name       `xml:"package"`
	Name    string         `xml:"name,attr"`
	Metrics PackageMetrics `xml:"metrics"`
	File    []*FileElement `xml:"file"`
}

type FileElement struct {
	XMLName xml.Name        `xml:"file"`
	Name    string          `xml:"name,attr"`
	Path    string          `xml:"path,attr"`
	Metrics FileMetrics     `xml:"metrics"`
	Class   []*ClassElement `xml:"class"`
}

type ClassElement struct {
	XMLName xml.Name       `xml:"class"`
	Name    string         `xml:"name,attr"`
	Metrics ClassMetrics   `xml:"metrics"`
	Line    []*LineElement `xml:"line"`
}

type LineElement struct {
	XMLName xml.Name `xml:"line"`
	Num     int      `xml:"num,attr"`
	Type    string   `xml:"type,attr"`
	Count   int      `xml:"count,attr"`
}

type ProjectMetrics struct {
	PackageMetrics
	Packages int `xml:"packages,attr"`
}

type PackageMetrics struct {
	FileMetrics
	Files int `xml:"files,attr"`
}

type FileMetrics struct {
	ClassMetrics
	Classes int `xml:"classes,attr"`
	Loc     int `xml:"loc,attr"`
	Ncloc   int `xml:"ncloc,attr"`
}

type ClassMetrics struct {
	XMLName             xml.Name `xml:"metrics"`
	Complexity          int      `xml:"complexity,attr"`
	Elements            int      `xml:"elements,attr"`
	CoveredElements     int      `xml:"coveredelements,attr"`
	Conditionals        int      `xml:"conditionals,attr"`
	CoveredConditionals int      `xml:"coveredconditionals,attr"`
	Statements          int      `xml:"statements,attr"`
	CoveredStatements   int      `xml:"coveredstatements,attr"`
	Methods             int      `xml:"methods,attr"`
	CoveredMethods      int      `xml:"coveredmethods,attr"`
}
