package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func PrintDiscTree(disc Disc, args Arguments) {
	fmt.Printf("Name:   %s\n", disc.Name)
	fmt.Printf("Type:   %s\n", disc.Type)
	fmt.Printf("Volume: %s\n", disc.Volume)
	fmt.Printf("Titles: %d\n", len(disc.Titles))
	fmt.Println()
	fmt.Printf("%-7s  %-30s %-8s %-2s %-8s  %-28s %-40s %-24s\n", "TrackID", "Name", "Duration", "Ch", "Size", "Video", "Audio", "Subtitles")
	fmt.Printf("%-7s  %-30s %-8s %-2s %-8s  %-28s %-40s %-24s\n", "-------", "------------------------------", "--------", "--", "--------", "----------------------------", "----------------------------------------", "------------------------")

	titles := append([]Title(nil), disc.Titles...)
	if (args.MinSize != "") && (args.MinSize != "0") {
		minBytes := sizeToBytes(args.MinSize)
		titles = filter(titles, func(t Title) bool { return t.Bytes >= minBytes })
	}
	sort.Slice(titles, func(i, j int) bool { return titles[i].ID < titles[j].ID })
	for _, t := range titles {
		videoStr := "—"
		if len(t.Video) > 0 {
			v := t.Video[0]
			codec := firstNonEmpty(v.CodecShort, v.CodecLong, v.CodecID)
			res := normalizeResolution(v.Resolution)
			fps := shortFrameRate(v.FrameRate)
			videoStr = fmt.Sprintf("%s • %s • %s", codec, res, fps)
		}

		audioStr := formatAudioGrouped(t.Audio)
		subsStr := formatSubsDeduped(t.Subtitles)

		fmt.Printf("%-7s  %-30.30s %-8.8s %02d %8s  %-28.28s %-40.40s %-24.24s\n",
			fmt.Sprintf("%02d", t.ID),
			t.Name,
			t.Duration,
			t.Chapters,
			t.Size,
			videoStr,
			audioStr,
			subsStr,
		)
	}
}

func normalizeResolution(res string) string {
	switch res {
	case "1920x1080":
		return "1080p"
	case "3840x2160":
		return "2160p"
	case "1280x720":
		return "720p"
	default:
		return res
	}
}

func shortFrameRate(fr string) string {
	if i := strings.Index(fr, " "); i > 0 {
		return fr[:i]
	}
	return fr
}

func formatChannels(ch int) string {
	switch ch {
	case 1:
		return "1.0"
	case 2:
		return "2.0"
	case 6:
		return "5.1"
	case 8:
		return "7.1"
	default:
		return strconv.Itoa(ch)
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return "?"
}

type audioKey struct {
	Channels string
	Codec    string
	Layout   string
}

func formatAudioGrouped(list []Audio) string {
	if len(list) == 0 {
		return "—"
	}
	perLang := map[string]map[audioKey]bool{}
	for _, a := range list {
		lang := strings.ToLower(firstNonEmpty(a.LanguageCode, a.Language))
		if lang == "" || lang == "?" {
			lang = "und"
		}
		key := audioKey{
			Channels: formatChannels(a.Channels),
			Codec:    firstNonEmpty(a.CodecShort, a.CodecLong, a.CodecID),
			Layout:   strings.TrimSpace(a.Layout),
		}
		if _, ok := perLang[lang]; !ok {
			perLang[lang] = map[audioKey]bool{}
		}
		perLang[lang][key] = perLang[lang][key] || a.Default
	}

	langs := make([]string, 0, len(perLang))
	for l := range perLang {
		langs = append(langs, l)
	}
	sort.Strings(langs)

	langChunks := make([]string, 0, len(langs))
	for _, l := range langs {
		items := make([]audioKey, 0, len(perLang[l]))
		for k := range perLang[l] {
			items = append(items, k)
		}
		sort.Slice(items, func(i, j int) bool {
			w := func(s string) int {
				switch s {
				case "7.1":
					return 3
				case "5.1":
					return 2
				case "2.0":
					return 1
				default:
					return 0
				}
			}
			wi, wj := w(items[i].Channels), w(items[j].Channels)
			if wi != wj {
				return wi > wj
			}
			if items[i].Codec != items[j].Codec {
				return items[i].Codec < items[j].Codec
			}
			return items[i].Layout < items[j].Layout
		})

		var parts []string
		for _, k := range items {
			entry := fmt.Sprintf("%s %s", k.Channels, k.Codec)
			if perLang[l][k] {
				entry += "*"
			}
			parts = append(parts, entry)
		}
		langChunks = append(langChunks, fmt.Sprintf("%s: %s", l, strings.Join(parts, " • ")))
	}

	return strings.Join(langChunks, " / ")
}

type subKey struct {
	Lang  string
	Flags string
}

func formatSubsDeduped(list []Subtitles) string {
	if len(list) == 0 {
		return "—"
	}
	m := map[subKey]bool{}
	for _, s := range list {
		lang := strings.ToLower(firstNonEmpty(s.LanguageCode, s.Language))
		if lang == "" || lang == "?" {
			lang = "und"
		}
		flags := subFlags(s.Description)
		k := subKey{Lang: lang, Flags: flags}
		m[k] = m[k] || s.Default
	}

	type outRow struct {
		Lang  string
		Flags string
		Def   bool
	}
	rows := make([]outRow, 0, len(m))
	for k, def := range m {
		rows = append(rows, outRow{k.Lang, k.Flags, def})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Lang != rows[j].Lang {
			return rows[i].Lang < rows[j].Lang
		}
		pri := func(fl string) int {
			switch fl {
			case "⚑":
				return 1
			case "Ⓢ":
				return 2
			case "⚑Ⓢ":
				return 3
			default:
				return 0
			}
		}
		return pri(rows[i].Flags) < pri(rows[j].Flags)
	})

	out := make([]string, 0, len(rows))
	for _, r := range rows {
		entry := r.Lang + r.Flags
		if r.Def {
			entry += "*"
		}
		out = append(out, entry)
	}
	return strings.Join(out, ", ")
}

func subFlags(desc string) string {
	d := strings.ToLower(desc)
	hasForced := strings.Contains(d, "forced")
	hasSDH := strings.Contains(d, "sdh") || strings.Contains(d, "hoh")
	var flags []string
	if hasForced {
		flags = append(flags, "⚑")
	}
	if hasSDH {
		flags = append(flags, "Ⓢ")
	}
	return strings.Join(flags, "")
}
