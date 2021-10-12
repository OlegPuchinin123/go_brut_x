/*
 * (c) Oleg Puchinin 2021
 * puchininolegigorevich@gmail.com
 */

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Brut struct {
	f         string
	templ     string
	main_code string
	diap      *string
	step      *float64
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
	os.Chdir("/tmp")
	cmd = exec.Command("gcc", "go_brut_x.c", "-lm", "-O3")
	cmd.Run()
	os.Remove("/tmp/go_brut_x.c")
	if cmd.ProcessState.ExitCode() != 0 {
		fmt.Fprintf(os.Stderr, "Can't compile the equation.\n")
		return errors.New("Can't compile.")
	}
	cmd = exec.Command("./a.out")
	cmd.Stdout = os.Stdout
	cmd.Run()
	os.Remove("/tmp/a.out")
	return nil
}

func (brut *Brut) process_template(f string) {
	var (
		diap_spl []string
	)
	brut.main_code = brut.templ
	brut.main_code = strings.ReplaceAll(brut.main_code, "@f", brut.f)
	brut.main_code = strings.ReplaceAll(brut.main_code, "@step", fmt.Sprintf("%f", *brut.step))

	diap_spl = strings.Split(*brut.diap, "..")
	if len(diap_spl) == 2 {
		brut.main_code = strings.ReplaceAll(brut.main_code, "@diap_low", diap_spl[0])
		brut.main_code = strings.ReplaceAll(brut.main_code, "@diap_high", diap_spl[1])
	}
	//println(brut.main_code)
}

func (brut *Brut) main_loop() {
	var (
		e   error
		buf *bufio.Reader
		s   string
	)
	buf = bufio.NewReader(os.Stdout)
	for {
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
		buf  []byte
		e    error
		s    string
	)

	brut = new(Brut)
	brut.diap = flag.String("diap", "-1000..1000", "")
	brut.step = flag.Float64("step", 0.001, "")
	flag.Parse()
	buf, e = os.ReadFile("./go_brut_x.tmpl")
	if e != nil {
		fmt.Printf("Can't read file ./go_brut_x.tmpl")
		return
	}
	s = string(buf)
	brut.templ = s
	brut.main_loop()
}
