package main

import (
	"flag"
	"fmt"
	"go-ddd-generator/strcase"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/rakyll/statik/fs"
	_ "go-ddd-generator/statik"
)

type Relation struct {
	TemplateFilePath string
	OutputFilePath string
}

func main() {

	// 引数を取得
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatalln("Not found args for module name")
	}

	moduleName := args[0]

	fmt.Println("Processing...\n")

	relations := createRelations(moduleName)

	for _, relation := range relations {
		outputFileName, err := generateFile(moduleName, relation)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Generated: %s\n", outputFileName)
	}
}

func createRelations(moduleName string) []*Relation {
	return []*Relation{
		{
			TemplateFilePath: "templates/service.template",
			OutputFilePath: fmt.Sprintf("domain/service/%s/interface.go", moduleName),
		},
		{
			TemplateFilePath: "templates/application.template",
			OutputFilePath: fmt.Sprintf("application/%s.go", moduleName),
		},
		{
			TemplateFilePath: "templates/repository.template",
			OutputFilePath: fmt.Sprintf("domain/repository/%s.go", moduleName),
		},
		{
			TemplateFilePath: "templates/infrastructure.template",
			OutputFilePath: fmt.Sprintf("infrastructure/%s.go", moduleName),
		},
	}
}

func generateFile(moduleName string, relation *Relation) (string, error) {
	template, err := readTemplate(relation.TemplateFilePath)

	replacedCode := convertCode(template, moduleName)

	err = writeGeneratedCodeToFile(relation.OutputFilePath, replacedCode)
	if err != nil{
		return "", err
	}

	return relation.OutputFilePath, nil
}

func readTemplate(filename string) (string, error) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	f, err := statikFS.Open(filename)
	if err != nil{
		return "", err
	}
	defer f.Close()

	template, err := ioutil.ReadAll(f)
	if err != nil{
		return "", err
	}

	return string(template), nil
}

func convertCode(template string, moduleName string) string {
	code := strings.Replace(template, "${project}", moduleName, -1)

	// user -> User
	//moduleNameFirstCapitalLetter := strings.ToUpper(moduleName[:1]) + moduleName[1:]
	camel := strcase.ToCamel(moduleName)
	code = strings.Replace(code, "${Project}", camel, -1)

	lowerCamel := strcase.ToLowerCamel(moduleName)
	code = strings.Replace(code, "${projectName}", lowerCamel, -1)

	workingDir := getWorkingDir()
	code = strings.Replace(code, "${working_dir}", workingDir, -1)

	return code
}

func writeGeneratedCodeToFile(outputFileName string, code string) error {
	pathElements := strings.Split(outputFileName, "/")

	// 存在しない場合、ディレクトリも作成する
	dicName := ""
	for _, e := range pathElements {
		dicName += e + "/"

		if strings.Contains(e, ".") {
			continue
		}

		if !isExist(dicName) {
			err := os.Mkdir(dicName, 0755)
			if err != nil {
				return err
			}
		}
	}


	outpuFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil{
		return err
	}

	fmt.Fprintf(outpuFile, code)

	return nil
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func getWorkingDir() string {
	p, _ := os.Getwd()
	path := strings.Split(p, "/")
	return path[len(path)-1]
}
