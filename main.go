package main

import (
	"flag"
	"fmt"
	"github.com/figoxu/Figo"
	"github.com/quexer/utee"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

// ./Translate -input=/home/zmy/newsda/workspace/Translate/java -output=/home/zmy/newsda/workspace/Translate/ts/crowd.model.ts
func main() {
	javaDir := flag.String("input", "", "java dir for translate")
	outfile := flag.String("output", "", "output typescript filename")
	flag.Parse()
	files, _ := ioutil.ReadDir(*javaDir)
	content := ""
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			fname := fmt.Sprint(*javaDir, "/", file.Name())
			fmt.Println("fname:", fname)
			b, err := ioutil.ReadFile(fname)
			utee.Chk(err)
			name := strings.Replace(file.Name(), ".java", "", -1)
			code := parse2ts(string(b), name)
			log.Println(code)
			content = fmt.Sprintln(content, code)
		}
	}
	log.Println(content)
	ioutil.WriteFile(*outfile, []byte(content), os.ModePerm)
}

func parse2ts(content, name string) string {
	parser := Figo.Parser{
		PrepareReg: []string{"private.+;"},
		ProcessReg: []string{".+serialVersionUID.+", "private ", ";"},
	}
	results := parser.Exe(content)
	result := fmt.Sprintln("import {BaseModel} from \"../base.model\";\n", "export class ", name, " extends BaseModel{")
	for _, v := range results {
		vs := strings.Split(v, " ")
		result = fmt.Sprintln(result, "\t\t", vs[1], "\t:", parseJava2TsTp(vs[0]), ";")
	}
	result = fmt.Sprintln(result, "}")
	return result
}

func parseJava2TsTp(tp string) string {
	tp = strings.ToLower(tp)
	checkTp := func(reg string) bool {
		log.Println("tp:", tp)
		log.Println("reg:", reg)
		rs := regexp.MustCompile(reg).FindAllString(tp, -1)
		return len(rs) > 0
	}

	if checkTp("bool") {
		return "boolean"
	}
	if checkTp("string") {
		return "string"
	}
	if checkTp("int") || checkTp("long") || checkTp("double") || checkTp("float") {
		return "number"
	}
	if checkTp("date") {
		return "Date"
	}
	if checkTp("list") {
		replace := strings.Replace(tp, "list", "Array", -1)
		index := strings.Index(replace, "<")
		upper := strings.ToUpper(replace[(index + 1):(index + 2)])
		sprint := fmt.Sprint(replace[0:index+1], upper, replace[(index+2):len(replace)])
		return sprint
	}
	return tp
}
