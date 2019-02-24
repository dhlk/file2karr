package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

type tmplArg struct {
	Name string
	Data []byte
}

var tmpl = template.Must(template.New("file2karr").Parse(`
// parsed from {{.Name}}
const char output[] = { {{range .Data}}0x{{printf "%x" .}}, {{end}}};
unsigned int output_len = {{len .Data}};
`))

func output(path string) (err error) {
	var (
		arg tmplArg
		out *os.File
	)

	arg.Name = path
	arg.Data, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	out, err = os.Create(path + ".sr.c")
	if err != nil {
		return
	}
	defer out.Close()

	return tmpl.Execute(out, arg)
}

func main() {
	defer fmt.Scanln()

	if len(os.Args) < 2 {
		fmt.Println("missing required argument(s): files to work on")
		return
	}

	for _, path := range os.Args[1:] {
		err := output(path)
		if err != nil {
			fmt.Printf("error processing %s: %v\n", path, err)
		} else {
			fmt.Printf("successfully ran %s\n", path)
		}
	}
}
