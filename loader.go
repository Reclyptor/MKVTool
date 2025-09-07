package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Video struct {
	CodecID    string
	CodecShort string
	CodecLong  string
	Bitrate    string
	Resolution string
	Aspect     string
	FrameRate  string
}

type Audio struct {
	CodecID       string
	CodecShort    string
	CodecLong     string
	Language      string
	LanguageCode  string
	Description   string
	Channels      int
	Layout        string
	SampleRate    int
	BitsPerSample int
	Default       bool
}

type Subtitles struct {
	CodecID      string
	CodecShort   string
	CodecLong    string
	Language     string
	LanguageCode string
	Description  string
	Default      bool
}

type Title struct {
	ID        int
	Name      string
	Chapters  int
	Duration  string
	Playlist  string
	Bytes     int64
	Size      string
	Video     []Video
	Audio     []Audio
	Subtitles []Subtitles
}

type Disc struct {
	Type   string
	Name   string
	Volume string
	Titles []Title
}

func buildTitles(tracks map[int]Track, streams map[int]map[int]Stream) []Title {
	var titles []Title
	for trackID, track := range tracks {
		title := Title{}
		title.ID = trackID
		title.Name = track.Name
		title.Chapters = track.Chapters
		title.Duration = track.Duration
		title.Playlist = track.Playlist
		title.Bytes = track.SizeBytes
		title.Size = track.SizeHuman
		for _, stream := range streams[trackID] {
			switch stream.TypeName {
			case "Video":
				video := Video{}
				video.CodecID = stream.CodecID
				video.CodecShort = stream.CodecShort
				video.CodecLong = stream.CodecLong
				video.Bitrate = stream.Bitrate
				video.Resolution = stream.Resolution
				video.Aspect = stream.AspectRatio
				video.FrameRate = stream.FrameRate
				title.Video = append(title.Video, video)
			case "Audio":
				audio := Audio{}
				audio.CodecID = stream.CodecID
				audio.CodecShort = stream.CodecShort
				audio.CodecLong = stream.CodecLong
				audio.Language = stream.LangName
				audio.LanguageCode = stream.LangCode
				audio.Description = stream.Attr
				audio.Channels = stream.Channels
				audio.Layout = stream.ChannelLayout
				audio.SampleRate = stream.SampleRate
				audio.BitsPerSample = stream.BitsPerSample
				audio.Default = stream.DefaultFlag
				title.Audio = append(title.Audio, audio)
			case "Subtitles":
				subtitles := Subtitles{}
				subtitles.CodecID = stream.CodecID
				subtitles.CodecShort = stream.CodecShort
				subtitles.CodecLong = stream.CodecLong
				subtitles.Language = stream.LangName
				subtitles.LanguageCode = stream.LangCode
				subtitles.Description = stream.LongDesc
				subtitles.Default = stream.DefaultFlag
				title.Subtitles = append(title.Subtitles, subtitles)
			}
		}
		titles = append(titles, title)
	}
	return titles
}

func LoadDisc(args Arguments) Disc {
	if args.Drive == "" {
		fmt.Println("Drive not specified. Use -d or --drive to specify the drive.")
		printUsage()
		os.Exit(1)
	}

	var argv []string
	argv = append(argv, "-r")
	argv = append(argv, "info")
	argv = append(argv, "dev:"+args.Drive)
	cmd := exec.Command("makemkvcon", argv...)

	var output, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &errb
	if err := cmd.Run(); err != nil {
		fmt.Println(errb)
		os.Exit(1)
	}

	cinfo, tinfo, sinfo := Tokenize(output.String())
	container := ParseCInfo(cinfo)
	tracks := ParseTInfo(tinfo)
	streams := ParseSInfo(sinfo)

	disc := Disc{}
	disc.Type = container.DiscType
	disc.Name = container.DiscName
	disc.Volume = container.VolumeLabel
	disc.Titles = buildTitles(tracks, streams)

	return disc
}

func RipDisc(args Arguments) {
	if args.Drive == "" {
		fmt.Println("Drive not specified. Use -d or --drive to specify the drive.")
		printUsage()
		os.Exit(1)
	}
	if args.OutDir == "" {
		fmt.Println("Output directory not specified. Use -o or --outdir to specify the output directory.")
		printUsage()
		os.Exit(1)
	}

	if err := os.MkdirAll(args.OutDir, 0o755); err != nil {
		fmt.Println("Failed to create output directory:", err)
		os.Exit(1)
	}

	var argv []string
	argv = append(argv, "mkv")
	argv = append(argv, "--progress")
	argv = append(argv, "--noscan")
	argv = append(argv, "--directio=true")
	if (args.MinLength != "") && (args.MinLength != "0") {
		argv = append(argv, "--minlength="+args.MinLength)
	}
	if (args.Audio != nil) && (len(args.Audio) > 0) {
		argv = append(argv, "--audio="+fmt.Sprintf("%s", strings.Join(args.Audio, ",")))
	}
	if (args.Subtitle != nil) && (len(args.Subtitle) > 0) {
		argv = append(argv, "--subtitle="+fmt.Sprintf("%s", strings.Join(args.Subtitle, ",")))
	}
	if args.Drive != "" {
		argv = append(argv, "dev:"+args.Drive)
	}
	if (args.Tracks != nil) && (len(args.Tracks) > 0) {
		for _, track := range args.Tracks {
			argv = append(argv, strconv.FormatInt(track, 10))
		}
	} else {
		argv = append(argv, "all")
	}

	tmpDir, err := os.MkdirTemp("", "*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}
	argv = append(argv, tmpDir)

	cmd := exec.Command("makemkvcon", argv...)

	var errb bytes.Buffer
	cmd.Stderr = &errb
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println(errb)
		os.Exit(1)
	}

	glob := filepath.Join(tmpDir, "*.mkv")
	files, err := filepath.Glob(glob)
	if err != nil {
		fmt.Println("Error reading temporary directory:", err)
		os.Exit(1)
	}
	if files == nil {
		fmt.Println("No MKVs produced. Nothing to do.")
		os.Exit(0)
	}
	slices.Sort(files)

	regex := regexp.MustCompile(`^.*?(?P<id>\d+)\.mkv$`)
	for _, file := range files {
		if args.Name != "" {
			mkvpropedit := exec.Command("mkvpropedit", file, "--edit", "info", "--set", "title="+args.Name)
			err := mkvpropedit.Run()
			if err != nil {
				fmt.Printf("mkvpropedit failed for %s\n", file)
			}
		}

		base := filepath.Base(file)
		matches := regex.FindStringSubmatch(base)
		trackID := ""
		if matches != nil {
			trackID = "_" + matches[1]
		}

		dest := filepath.Join(args.OutDir, fmt.Sprintf("%s%s.mkv", args.Name, trackID))
		fmt.Printf("→ %s  ==>  %s\n", base, filepath.Base(dest))
		if err = copyFile(file, dest); err != nil {
			fmt.Println("Error renaming file:", err)
		}
	}

	fmt.Printf("✓ Done. Wrote %d file(s) to: %s\n", len(files), args.OutDir)
}
