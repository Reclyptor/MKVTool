package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const VERSION = "0.0.0"

func printUsage() {
	println("Usage: ripmkv [options]")
	println("Version:", VERSION)
	println("Example: ripmkv -d /dev/sr0 -n Title -o /path/to/output -t 0 1 2 -a eng jpn -s eng")
	println("Options:")
	println("  -l, --list                   List available tracks")
	println("  --minsize <size>             Filter tracks of at least this size, used with -l, e.g. 100M, 1.5G")
	println("  --minlength <seconds>        Filter tracks of at least this length, used whenever -t is omitted, e.g. 3600")
	println("  -d, --drive <path>           Specify the drive path, e.g. /dev/sr0")
	println("  -t, --track <track>          Specify the tracks to rip, e.g. 0 1 2 ..., or all if none specified")
	println("  -a, --audio <lang>           Specify the audio languages to keep, e.g. eng jpn")
	println("  -s, --subtitle <lang>        Specify the subtitle languages to keep, e.g. eng jpn")
	println("  -n, --name <name>            Specify the output title name prefix, also used as segment title")
	println("  -o, --outdir <output dir>    Specify the output directory, default is current directory")
	println("  -v, --version                Show version information")
	println("  -h, --help                   Show this help message")
}

type Arguments struct {
	List      bool
	MinSize   string
	MinLength string
	Drive     string
	Tracks    []int64
	Audio     []string
	Subtitle  []string
	Name      string
	OutDir    string
	Version   bool
	Help      bool
}

func parseArgs() Arguments {
	var arguments Arguments
	for idx := 1; idx < len(os.Args); idx++ {
		switch os.Args[idx] {
		case "-l", "--list":
			arguments.List = true
		case "--minsize":
			arguments.MinSize = os.Args[idx+1]
			idx++
		case "--minlength":
			arguments.MinLength = os.Args[idx+1]
			idx++
		case "-d", "--drive":
			arguments.Drive = os.Args[idx+1]
			idx++
		case "-t", "--track":
			for subIdx := idx + 1; subIdx < len(os.Args); subIdx++ {
				if matched, _ := regexp.MatchString(`^-`, os.Args[subIdx]); matched {
					break
				}
				track, err := strconv.ParseInt(os.Args[subIdx], 10, 0)
				if err != nil {
					fmt.Println("Invalid track number:", os.Args[subIdx])
					printUsage()
					os.Exit(1)
				}
				arguments.Tracks = append(arguments.Tracks, track)
			}
			idx += len(arguments.Tracks)
		case "-a", "--audio":
			for subIdx := idx + 1; subIdx < len(os.Args); subIdx++ {
				if matched, _ := regexp.MatchString(`^-`, os.Args[subIdx]); matched {
					break
				}
				arguments.Audio = append(arguments.Audio, os.Args[subIdx])
			}
		case "-s", "--subtitle":
			for subIdx := idx + 1; subIdx < len(os.Args); subIdx++ {
				if matched, _ := regexp.MatchString(`^-`, os.Args[subIdx]); matched {
					break
				}
				arguments.Subtitle = append(arguments.Subtitle, os.Args[subIdx])
			}
		case "-n", "--name":
			arguments.Name = os.Args[idx+1]
			idx++
		case "-o", "--outdir":
			arguments.OutDir = os.Args[idx+1]
			idx++
		case "-v", "--version":
			arguments.Version = true
		case "-h", "--help":
			arguments.Help = true
		}
	}
	return arguments
}

func main() {
	args := parseArgs()

	if args.Help {
		printUsage()
		os.Exit(0)
	}

	if args.Version {
		fmt.Println("Version:", VERSION)
		os.Exit(0)
	}

	if args.List {
		ListTracks(args)
		os.Exit(0)
	}

	RipTracks(args)
	os.Exit(0)
}
