package main

const (
	CiDiscType     = 1  // e.g. "Blu-ray disc", "DVD disc"
	CiDiscName     = 2  // e.g. "Up (Disc 1)"
	CiLangCode     = 28 // ISO code, e.g. "eng"
	CiLangName     = 29 // e.g. "English"
	CiTitle        = 30 // often same as DiscName
	CiUiHeaderHTML = 31 // HTML header, e.g. "<b>Source information</b><br>"
	CiVolumeLabel  = 32 // e.g. "UP_USA"
	CiLayerInfo    = 33 // side/layer/region-ish (often "0")
)

const (
	TiName           = 2  // Human name (movie/episode), e.g. "Up (Disc 1)"
	TiChapters       = 8  // "5"
	TiDuration       = 9  // "HH:MM:SS"
	TiSizeHuman      = 10 // "920.1 MB"
	TiSizeBytes      = 11 // "964829184"
	TiPlaylist       = 16 // "00038.mpls" or source file
	TiVideoTracks    = 25 // count
	TiAudioTracks    = 26 // count
	TiDefaultOutName = 27 // "Up (Disc 1)_t00.mkv"
	TiLangCode       = 28 // "eng"
	TiLangName       = 29 // "English"
	TiLongDesc       = 30 // e.g. "Up (Disc 1) - 5 chapter(s) , 920.1 MB"
	TiUiHeaderHTML   = 31 // HTML header
	TiUnknown33      = 33 // observed "0"
)

const (
	SiTypeName      = 1  // "Video" | "Audio" | "Subtitles"
	SiAttr          = 2  // e.g. "Surround 5.1", "Stereo"
	SiLangCode      = 3  // "eng"
	SiLangName      = 4  // "English"
	SiCodecID       = 5  // "A_AC3", "V_MPEG4/ISO/AVC", "S_HDMV/PGS"
	SiCodecShort    = 6  // "DD", "DTS", "PGS"…
	SiCodecLong     = 7  // "Dolby Digital", "DTS-HD MA", "HDMV PGS Subtitles"
	SiBitrate       = 13 // "640 Kb/s"
	SiChannels      = 14 // "6"
	SiSampleRate    = 17 // "48000"
	SiBitsPerSample = 18 // audio bit depth, e.g. "24"
	SiResolution    = 19 // "1920x1080"
	SiAspectRatio   = 20 // "16:9"
	SiFrameRate     = 21 // "23.976 (24000/1001)"
	SiPgSizeOrFlg   = 22 // seen "0" / "6144" (PGS packet size/flag)
	SiLangCode2     = 28 // duplicate lang code
	SiLangName2     = 29 // duplicate lang text
	SiLongDesc      = 30 // "DD Surround 5.1 English"
	SiUiHeaderHTML  = 31 // HTML header
	SiTrackID       = 33 // numeric track id/order
	SiUnknown38     = 38 // seen "", "d"
	SiDefaultFlag   = 39 // "Default"
	SiChannelLayout = 40 // "5.1(side)", "stereo"
	SiNotes         = 42 // "( Lossless conversion )"
)
