/*
 * (c) Oleg Puchinin
 * puchininolegigorevich@gmail.com
 */

package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//go:embed go_brut_x.tmpl
var templ string

type Brut struct {
	f         string
	templ     string
	main_code string
	diap_high string
	diap_low  string
	step      string
}

func (brut *Brut) compile() error {
	var (
		f   *os.File
		e   error
		cmd *exec.Cmd
	)
	f, e = os.Create("/tmp/go_brut_x.c")
	if e != nil {
		return e
	}
	f.WriteString(brut.main_code)
	f.Close()
	defer os.Remove("/tmp/go_brut_x.c")

	os.Chdir("/tmp")
	cmd = exec.Command("gcc", "go_brut_x.c", "-lm")
	e = cmd.Run()
	if e != nil {
		os.Stderr.WriteString("Can't compile.\n")
		return e
	}
	cmd = exec.Command("./a.out")
	cmd.Stdout = os.Stdout
	cmd.Run()
	os.Remove("/tmp/a.out")
	return nil
}

func (brut *Brut) process_template(f string) {
	brut.main_code = brut.templ
	brut.main_code = strings.ReplaceAll(brut.main_code, "@f", brut.f)
	brut.main_code = strings.ReplaceAll(brut.main_code, "@step", brut.step)
	brut.main_code = strings.ReplaceAll(brut.main_code, "@diap_low", brut.diap_low)
	brut.main_code = strings.ReplaceAll(brut.main_code, "@diap_high", brut.diap_high)
}

func (brut *Brut) main_loop() {
	var (
		e   error
		buf *bufio.Reader
		s   string
	)
	buf = bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Enter equation (eq zero): ")
		s, e = buf.ReadString('\n')
		if e != nil {
			return
		}
		s = strings.TrimRight(s, "\n")
		brut.f = s
		brut.process_template(s)
		brut.compile()
	}
}

func main() {
	var (
		brut *Brut
		low  *string
		high *string
		step *string
	)

	step = flag.String("step", "0.001", "step. stupid step")
	low = flag.String("diap_low", "-1000", "low")
	high = flag.String("diap_high", "1000", "high")
	flag.Parse()

	brut = &Brut{}
	brut.diap_low = *low
	brut.diap_high = *high
	brut.step = *step

	brut.templ = templ
	brut.main_loop()
}
