package main

import (
	"github.com/figoxu/Figo"
	"io/ioutil"
	"fmt"
	"github.com/quexer/utee"
	"strings"
	"regexp"
	"flag"
	"os"
	"log"
)

func main(){
	javaDir := *flag.String("input", "", "java dir for translate")
	outfile := *flag.String("output", "", "output typescript filename")
	files,_ := ioutil.ReadDir(javaDir)
	content := ""
	for _,file := range files {
		if file.IsDir() {
			continue
		}else{
			fname := fmt.Sprint(javaDir,"/",file.Name())
			b,err:=ioutil.ReadFile(fname)
			utee.Chk(err)
			name:=strings.Replace(file.Name(),".java","",-1)
			code := parse2ts(string(b),name)
			log.Println(code)
			content = fmt.Sprintln(content,code)
		}
	}
	log.Println(content)
	ioutil.WriteFile(outfile,[]byte(content),os.ModePerm)
}


func parse2ts(content,name string) string{
	parser := Figo.Parser{
		PrepareReg:[]string{"private.+;"},
		ProcessReg:[]string{".+serialVersionUID.+","private ",";"},
	}
	results := parser.Exe(content)
	result := fmt.Sprintln("class ",name,"{")
	for _,v := range results {
		vs := strings.Split(v," ");
		result = fmt.Sprintln(result,"\t\t",vs[1],"\t:",parseJava2TsTp(vs[0]),";")
	}
	result = fmt.Sprintln(result,"}")
	return result
}

func parseJava2TsTp(tp string)string{
	tp =strings.ToLower(tp)
	checkTp := func(reg string) bool {
		rs := regexp.MustCompile(reg).FindAllString(tp, -1)
		return len(rs)>0
	}

	if checkTp("bool"){
		return "boolean"
	}
	if checkTp("string"){
		return "string"
	}
	if checkTp("int")||checkTp("long")||checkTp("double")||checkTp("float"){
		return "string"
	}
	if checkTp("date"){
		return "Date"
	}
	return tp;
}