package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aquasecurity/tracee/tracee"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Tracee",
		Usage: "Trace OS events and syscalls using eBPF",
		Action: func(c *cli.Context) error {
			if c.Bool("list") {
				printList()
				return nil
			}
			cfg, err := tracee.NewConfig(
				c.StringSlice("event"),
				c.Bool("container"),
				c.Bool("detect-original-syscall"),
				c.Bool("show-exec-env"),
				c.String("output"),
				c.Int("perf-buffer-size"),
				c.Bool("extract-files"),
				c.String("output-path"),
				c.StringSlice("filter-file-write"),
			)
			if c.Bool("show-all-syscalls") {
				cfg.EventsToTrace = append(cfg.EventsToTrace, tracee.EventsNameToID["raw_syscalls"])
			}
			if err != nil {
				return fmt.Errorf("error creating Tracee config: %v", err)
			}
			t, err := tracee.New(*cfg)
			if err != nil {
				// t is being closed internally
				return fmt.Errorf("error creating Tracee: %v", err)
			}
			return t.Run()
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "table",
				Usage:   "output format: table (default)/json",
			},
			&cli.StringSliceFlag{
				Name:    "event",
				Aliases: []string{"e"},
				Value:   nil,
				Usage:   "trace only the specified event or syscall. use this flag multiple times to choose multiple events",
			},
			&cli.BoolFlag{
				Name:    "list",
				Aliases: []string{"l"},
				Value:   false,
				Usage:   "just list tracable events",
			},
			&cli.BoolFlag{
				Name:    "container",
				Aliases: []string{"c"},
				Value:   false,
				Usage:   "trace only containers",
			},
			&cli.BoolFlag{
				Name:  "detect-original-syscall",
				Value: false,
				Usage: "when tracing kernel functions which are not syscalls (such as cap_capable), detect and show the original syscall that called that function",
			},
			&cli.BoolFlag{
				Name:  "show-exec-env",
				Value: false,
				Usage: "when tracing execve/execveat, show environment variables",
			},
			&cli.IntFlag{
				Name:    "perf-buffer-size",
				Aliases: []string{"b"},
				Value:   64,
				Usage:   "size, in pages, of the internal perf ring buffer used to submit events from the kernel",
			},
			&cli.BoolFlag{
				Name:    "show-all-syscalls",
				Aliases: []string{"a"},
				Value:   false,
				Usage:   "log all syscalls invocations, including syscalls which were not fully traced by tracee (shortcut to -e raw_syscalls)",
			},
			&cli.BoolFlag{
				Name:  "extract-files",
				Value: false,
				Usage: "extract file writes to output path. \"filter-file-write\" must be specified",
			},
			&cli.StringFlag{
				Name:  "output-path",
				Value: "/tmp/tracee",
				Usage: "output path for extracted files, i.e. when \"extract-files\" is used",
			},
			&cli.StringSliceFlag{
				Name:  "filter-file-write",
				Value: nil,
				Usage: "filter path of file writes. Up to 3 filters can be specified, 16 characters each",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func printList() {
	const sep = ", "
	var b strings.Builder
	for _, name := range tracee.EventsSyscalls {
		b.WriteString(name)
		b.WriteString(sep)
	}
	fmt.Println("System calls:")
	fmt.Println(strings.TrimSuffix(b.String(), sep))
	b.Reset()
	fmt.Println()

	for _, name := range tracee.EventsTracepoints {
		b.WriteString(name)
		b.WriteString(sep)
	}
	fmt.Println("Tracepoints:")
	fmt.Println(strings.TrimSuffix(b.String(), sep))
	b.Reset()
	fmt.Println()

	for _, name := range tracee.EventsKprobes {
		b.WriteString(name)
		b.WriteString(sep)
	}
	fmt.Println("System events:")
	fmt.Println(strings.TrimSuffix(b.String(), sep))
}
