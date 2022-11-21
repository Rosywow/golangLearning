//go:build ignore

//annotation:service
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"w2w.io/cmn"
)

var missingAuthorPrompt = `
Missing service developer information, please specify 
  it in %s by bellow format

    //author:{"name":"tom sawyer","tel":"13580452503", "email":"KManager@GMail.com"}

`

var invalidAuthorPrompt = `
[%s: %d] %s

Invalid service developer information, please specify 
  it in %s by bellow JSON format

    //author:{"name":"tom sawyer","tel":"13580452503", "email":"KManager@GMail.com"}

`

func main() {

	type service struct {
		ServiceName string `json:"serviceName"`

		Author     *cmn.ModuleAuthor `json:"author"`
		sourceFile string
	}

	type templateArg struct {
		Services   map[string]service
		ModuleName string
	}

	s, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	moduleName := strings.ReplaceAll(strings.ReplaceAll(string(s), "\n", ""), "\r", "")

	//-----------------
	directories, err := ioutil.ReadDir("../serve")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	rDeveloperPattern := regexp.MustCompile(`(?i)^//\s*author\s*:\s*(?P<info>{.*})\s*$`)
	rServiceNamePattern := regexp.MustCompile(`^\s*//\s*annotation\s*:\s*(?P<name>.*)-service\s*$`)
	rFileNamePattern := regexp.MustCompile(`(?i).*\.go$`)
	services := make(map[string]service)
	authors := make(map[string]*cmn.ModuleAuthor)

	for _, serviceDir := range directories {
		if !serviceDir.IsDir() {
			continue
		}

		serviceName := ""
		var devAuthor *cmn.ModuleAuthor
		serviceFilePath := "../serve/" + serviceDir.Name()
		err = filepath.Walk(serviceFilePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return err
			}

			if info.IsDir() {
				return nil
			}

			if !rFileNamePattern.MatchString(info.Name()) {
				return nil
			}

			// it's a goLang source file
			fileName := fmt.Sprintf("../serve/%s/%s", serviceDir.Name(), info.Name())

			f, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
			if err != nil {
				fmt.Println("open file error: " + err.Error())
				return err
			}

			defer f.Close()

			rd := bufio.NewReader(f)
			lineNum := 0
			//read every line
			for {
				lineNum++
				line, err := rd.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						break
					}

					fmt.Println("read file line error: " + err.Error())
					return err
				}

				//find service name
				match := rServiceNamePattern.FindStringSubmatch(line)
				if len(match) > 0 {
					for i, name := range rServiceNamePattern.SubexpNames() {
						if i != 0 && name != "" {
							serviceName = match[i]
							services[serviceDir.Name()] = service{ServiceName: serviceName, sourceFile: serviceFilePath + "/" + info.Name()}
						}
					}
				}

				//find author
				match = rDeveloperPattern.FindStringSubmatch(line)
				if len(match) > 0 {
					for i, name := range rDeveloperPattern.SubexpNames() {
						if i != 0 && name != "" {
							devAuthor = &cmn.ModuleAuthor{}
							err = json.Unmarshal([]byte(match[i]), devAuthor)
							if err != nil {
								fmt.Printf(invalidAuthorPrompt, fileName, lineNum, err.Error(), fileName)
								os.Exit(-1)
							}
							devAuthor.Addi = match[i]
							authors[serviceDir.Name()] = devAuthor
						}
					}
				}

			}

			return err
		})

		if err != nil {
			fmt.Println(err.Error())
		}

		if serviceName == "" {
			continue
		}
	}

	if len(services) == 0 {
		return
	}

	for k, v := range services {
		d, ok := authors[k]
		if !ok {
			fmt.Printf(missingAuthorPrompt, strings.ReplaceAll(v.sourceFile, "../", ""))
			os.Exit(-1)
		}
		services[k] = service{
			ServiceName: v.ServiceName, Author: d,
		}
	}

	tmplArg := templateArg{
		ModuleName: string(moduleName),
		Services:   services,
	}

	t, err := template.ParseFiles("services.go.tmpl")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	f, err := os.Create("services-generated.go")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()

	err = t.ExecuteTemplate(f, "services.go.tmpl", tmplArg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
