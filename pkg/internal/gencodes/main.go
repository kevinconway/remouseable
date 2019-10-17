// Package main generates a Go source file that contains a mapping of all
// evdev codes by extracting them from the linux source code files.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dave/jennifer/jen"
	flag "github.com/spf13/pflag"
	sys "golang.org/x/sys/unix"
)

var sourceFiles = []string{
	"/usr/include/linux/input.h",
	"/usr/include/linux/input-event-codes.h",
}

// pattern is copied from github.com/gvalkov/golang-evdev
const pattern = `#define +((?:KEY|ABS|REL|SW|MSC|LED|BTN|REP|SND|ID|EV|BUS|SYN|FF)_\w+)\s+(\w+)`

func main() {
	fs := flag.NewFlagSet("gencodes", flag.ExitOnError)
	srcs := fs.StringArray("sources", sourceFiles, "Linux header source files to process.")
	dst := fs.String("destination", "evdevcodes.go", "The destination file path. Use - for stdout.")
	_ = fs.Parse(os.Args[1:])

	reg, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}

	codes := make([][]string, 0)
	for _, src := range *srcs {
		f, fErr := os.Open(src)
		if fErr != nil {
			panic(fErr)
		}
		defer f.Close()

		bf := bufio.NewScanner(f)
		for bf.Scan() {
			line := bf.Text()
			matches := reg.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				codes = append(codes, match[1:])
			}
		}
		if bf.Err() != nil {
			panic(bf.Err())
		}
	}

	un := &sys.Utsname{}
	if err = sys.Uname(un); err != nil {
		panic(err)
	}

	f := jen.NewFile("remouse")
	f.Comment("// Code generated DO NOT EDIT").Line()
	f.Comment(
		fmt.Sprintf(
			"// Generated using %s %s %s.",
			string(bytes.Trim(un.Sysname[:], "\x00")),
			string(bytes.Trim(un.Release[:], "\x00")),
			string(bytes.Trim(un.Machine[:], "\x00")),
		),
	)
	f.Comment(
		fmt.Sprintf(
			"// Generated at %s.",
			time.Now().Format(time.RFC3339),
		),
	)
	f.Comment(
		fmt.Sprintf(
			"// Generated from %s.",
			strings.Join(*srcs, ", "),
		),
	)

	keyMaps := make([]jen.Code, 0)
	absMaps := make([]jen.Code, 0)
	relMaps := make([]jen.Code, 0)
	swMaps := make([]jen.Code, 0)
	mscMaps := make([]jen.Code, 0)
	ledMaps := make([]jen.Code, 0)
	btnMaps := make([]jen.Code, 0)
	repMaps := make([]jen.Code, 0)
	sndMaps := make([]jen.Code, 0)
	idMaps := make([]jen.Code, 0)
	evMaps := make([]jen.Code, 0)
	busMaps := make([]jen.Code, 0)
	synMaps := make([]jen.Code, 0)
	ffMaps := make([]jen.Code, 0)
	defs := make([]jen.Code, 0, len(codes))
	for _, code := range codes {
		name := code[0]
		value := code[1]
		defs = append(defs, jen.Id(name).Op("=").Id(value))

		if name == "EV_VERSION" || strings.HasSuffix(name, "_MAX") {
			// EV_VERSION is not a uint16 value and is also not an event type.
			// *_MAX are often duplicated by other named values.
			continue
		}
		if name == "BTN_TRIGGER" || name == "BTN_SOUTH" || name == "BTN_DIGI" || name == "BTN_WHEEL" || name == "BTN_TRIGGER_HAPPY" || name == "BTN_MISC" || name == "BTN_MOUSE" {
			// BTN_TRIGGER is a duplicate of BTN_TASK
			// BTN_SOUTH is a duplicate of BTN_GAMEPAD
			// BTN_DIGI is a duplicate of BTN_TOOL_PEN
			// BTN_WHEEL is a duplicate of BTN_GEAR_DOWN
			// BTN_TRIGGER_HAPPY is a duplicate of BTN_TRIGGER_HAPPY1
			// BTN_MISC is a duplicate of BTN_0
			// BTN_MOUSE is a duplicate of BTN_LEFT
			continue
		}
		switch {
		case strings.HasPrefix(name, "KEY") && !strings.HasPrefix(value, "KEY"):
			keyMaps = append(keyMaps, jen.Id(value).Op(":").Lit(name))
		case strings.HasPrefix(name, "ABS") && !strings.HasPrefix(value, "ABS"):
			absMaps = append(absMaps, jen.Id(value).Op(":").Lit(name))
		case strings.HasPrefix(name, "REL") && !strings.HasPrefix(value, "REL"):
			relMaps = append(relMaps, jen.Id(value).Op(":").Lit(name))
		case strings.HasPrefix(name, "SW") && !strings.HasPrefix(value, "SW"):
			swMaps = append(swMaps, jen.Id(value).Op(":").Lit(name))
		case strings.HasPrefix(name, "MSC") && !strings.HasPrefix(value, "MSC"):
			mscMaps = append(mscMaps, jen.Id(value).Op(":").Lit(name))
		case strings.HasPrefix(name, "LED") && !strings.HasPrefix(value, "LED"):
			ledMaps = append(ledMaps, jen.Id(value).Op(":").Lit(name))
		case strings.HasPrefix(name, "BTN") && !strings.HasPrefix(value, "BTN"):
			btnMaps = append(btnMaps, jen.Id(value).Op(":").Lit(name))
		case strings.HasPrefix(name, "REP") && !strings.HasPrefix(value, "REP"):
			repMaps = append(repMaps, jen.Id(value).Op(":").Lit(name))
		case strings.HasPrefix(name, "SND") && !strings.HasPrefix(value, "SND"):
			sndMaps = append(sndMaps, jen.Id(value).Op(":").Lit(name))
		case strings.HasPrefix(name, "ID") && !strings.HasPrefix(value, "ID"):
			idMaps = append(idMaps, jen.Id(value).Op(":").Lit(name))
		case strings.HasPrefix(name, "EV") && !strings.HasPrefix(value, "EV"):
			evMaps = append(evMaps, jen.Id(value).Op(":").Lit(name))
		case strings.HasPrefix(name, "BUS") && !strings.HasPrefix(value, "BUS"):
			busMaps = append(busMaps, jen.Id(value).Op(":").Lit(name))
		case strings.HasPrefix(name, "SYN") && !strings.HasPrefix(value, "SYN"):
			synMaps = append(synMaps, jen.Id(value).Op(":").Lit(name))
		case strings.HasPrefix(name, "FF") && !strings.HasPrefix(value, "FF"):
			ffMaps = append(ffMaps, jen.Id(value).Op(":").Lit(name))
		}
	}
	f.Const().Defs(defs...)
	f.Var().Id("KEYMap").Op("=").Map(jen.Uint16()).String().Values(keyMaps...)
	f.Var().Id("ABSMap").Op("=").Map(jen.Uint16()).String().Values(absMaps...)
	f.Var().Id("RELMap").Op("=").Map(jen.Uint16()).String().Values(relMaps...)
	f.Var().Id("SWMap").Op("=").Map(jen.Uint16()).String().Values(swMaps...)
	f.Var().Id("MSCMap").Op("=").Map(jen.Uint16()).String().Values(mscMaps...)
	f.Var().Id("LEDMap").Op("=").Map(jen.Uint16()).String().Values(ledMaps...)
	f.Var().Id("BTNMap").Op("=").Map(jen.Uint16()).String().Values(btnMaps...)
	f.Var().Id("REPMap").Op("=").Map(jen.Uint16()).String().Values(repMaps...)
	f.Var().Id("SNDMap").Op("=").Map(jen.Uint16()).String().Values(sndMaps...)
	f.Var().Id("IDMap").Op("=").Map(jen.Uint16()).String().Values(idMaps...)
	f.Var().Id("EVMap").Op("=").Map(jen.Uint16()).String().Values(evMaps...)
	f.Var().Id("BUSMap").Op("=").Map(jen.Uint16()).String().Values(busMaps...)
	f.Var().Id("SYNMap").Op("=").Map(jen.Uint16()).String().Values(synMaps...)
	f.Var().Id("FFMap").Op("=").Map(jen.Uint16()).String().Values(ffMaps...)

	f.Func().Id("CodeString").Params(jen.Id("etype").Uint16(), jen.Id("code").Uint16()).String().Block(
		jen.Var().Id("stype").Op("=").Id("EVMap").Index(jen.Id("etype")),
		jen.Switch(jen.Id("stype")).Block(
			jen.Case(jen.Lit("EV_SYN")).Block(
				jen.Return(jen.Id("SYNMap").Index(jen.Id("code"))),
			),
			jen.Case(jen.Lit("EV_KEY")).Block(
				jen.Return(jen.Id("KEYMap").Index(jen.Id("code"))),
			),
			jen.Case(jen.Lit("EV_ABS")).Block(
				jen.Return(jen.Id("ABSMap").Index(jen.Id("code"))),
			),
			jen.Case(jen.Lit("EV_REL")).Block(
				jen.Return(jen.Id("RELMap").Index(jen.Id("code"))),
			),
			jen.Case(jen.Lit("EV_SW")).Block(
				jen.Return(jen.Id("SWMap").Index(jen.Id("code"))),
			),
			jen.Case(jen.Lit("EV_MSC")).Block(
				jen.Return(jen.Id("MSCMap").Index(jen.Id("code"))),
			),
			jen.Case(jen.Lit("EV_LED")).Block(
				jen.Return(jen.Id("LEDMap").Index(jen.Id("code"))),
			),
			jen.Case(jen.Lit("EV_SND")).Block(
				jen.Return(jen.Id("SNDMap").Index(jen.Id("code"))),
			),
			jen.Case(jen.Lit("EV_REP")).Block(
				jen.Return(jen.Id("REPMap").Index(jen.Id("code"))),
			),
			jen.Case(jen.Lit("EV_FF")).Block(
				jen.Return(jen.Id("FFMap").Index(jen.Id("code"))),
			),
			jen.Default().Block(
				jen.Return(jen.Lit("")),
			),
		),
	)

	if *dst == "-" {
		if err := f.Render(os.Stdout); err != nil {
			panic(err)
		}
		return
	}
	if err = f.Save(*dst); err != nil {
		panic(err)
	}
}
