package main

type Container struct {
	DiscType     string // CiDiscType
	DiscName     string // CiDiscName
	LangCode     string // CiLangCode
	LangName     string // CiLangName
	Title        string // CiTitle
	UIHeaderHTML string // CiUiHeaderHTML
	VolumeLabel  string // CiVolumeLabel
	LayerInfo    string // CiLayerInfo
	Raw          map[int]string
}

type Track struct {
	TrackID      int            // the TINFO <title_id> (ordinal)
	Name         string         // TiName
	Chapters     int            // TiChapters
	Duration     string         // TiDuration "HH:MM:SS"
	SizeHuman    string         // TiSizeHuman
	SizeBytes    int64          // TiSizeBytes
	Playlist     string         // TiPlaylist (e.g., "00038.mpls")
	VideoTracks  int            // TiVideoTracks
	AudioTracks  int            // TiAudioTracks
	DefaultOut   string         // TiDefaultOutName
	LangCode     string         // TiLangCode
	LangName     string         // TiLangName
	LongDesc     string         // TiLongDesc
	UIHeaderHTML string         // TiUiHeaderHTML
	Unknown33    string         // TiUnknown33 (string to avoid lossy coercion)
	Raw          map[int]string // all TINFO fieldID -> value
}

type Stream struct {
	TrackID       int            // SINFO <title_id> (ordinal)
	StreamID      int            // SINFO <stream_id> (ordinal)
	TypeName      string         // SiTypeName
	Attr          string         // SiAttr
	LangCode      string         // SiLangCode
	LangName      string         // SiLangName
	CodecID       string         // SiCodecID
	CodecShort    string         // SiCodecShort
	CodecLong     string         // SiCodecLong
	Bitrate       string         // SiBitrate
	Channels      int            // SiChannels
	SampleRate    int            // SiSampleRate
	BitsPerSample int            // SiBitsPerSample
	Resolution    string         // SiResolution
	AspectRatio   string         // SiAspectRatio
	FrameRate     string         // SiFrameRate
	PgSizeOrFlag  string         // SiPgSizeOrFlg (string to keep raw)
	LangCode2     string         // SiLangCode2
	LangName2     string         // SiLangName2
	LongDesc      string         // SiLongDesc
	UIHeaderHTML  string         // SiUiHeaderHTML
	StreamTrackID string         // SiTrackID
	Unknown38     string         // SiUnknown38
	DefaultFlag   bool           // SiDefaultFlag
	ChannelLayout string         // SiChannelLayout
	Notes         string         // SiNotes
	Raw           map[int]string // all SINFO fieldID -> value
}

func ParseCInfo(cInfo []CInfo) Container {
	var container Container
	for _, info := range cInfo {
		switch info.Field {
		case CiDiscType:
			container.DiscType = info.Value
		case CiDiscName:
			container.DiscName = info.Value
		case CiLangCode:
			container.LangCode = info.Value
		case CiLangName:
			container.LangName = info.Value
		case CiTitle:
			container.Title = info.Value
		case CiUiHeaderHTML:
			container.UIHeaderHTML = info.Value
		case CiVolumeLabel:
			container.VolumeLabel = info.Value
		case CiLayerInfo:
			container.LayerInfo = info.Value
		default:
			if container.Raw == nil {
				container.Raw = make(map[int]string)
			}
			container.Raw[info.Field] = info.Value
		}
	}
	return container
}

func ParseTInfo(tInfo []TInfo) map[int]Track {
	tracks := make(map[int]Track)
	for _, info := range tInfo {
		track, exists := tracks[info.Track]
		if !exists {
			track = Track{}
			track.TrackID = info.Track
		}
		switch info.Field {
		case TiName:
			track.Name = info.Value
		case TiChapters:
			track.Chapters = atoi(info.Value)
		case TiDuration:
			track.Duration = info.Value
		case TiSizeHuman:
			track.SizeHuman = info.Value
		case TiSizeBytes:
			track.SizeBytes = atoi64(info.Value)
		case TiPlaylist:
			track.Playlist = info.Value
		case TiVideoTracks:
			track.VideoTracks = atoi(info.Value)
		case TiAudioTracks:
			track.AudioTracks = atoi(info.Value)
		case TiDefaultOutName:
			track.DefaultOut = info.Value
		case TiLangCode:
			track.LangCode = info.Value
		case TiLangName:
			track.LangName = info.Value
		case TiLongDesc:
			track.LongDesc = info.Value
		case TiUiHeaderHTML:
			track.UIHeaderHTML = info.Value
		case TiUnknown33:
			track.Unknown33 = info.Value
		default:
			if track.Raw == nil {
				track.Raw = make(map[int]string)
			}
			track.Raw[info.Field] = info.Value
		}
		tracks[info.Track] = track
	}
	return tracks
}

func ParseSInfo(sInfo []SInfo) map[int]map[int]Stream {
	streams := make(map[int]map[int]Stream)
	for _, info := range sInfo {
		streamMap, exists := streams[info.Track]
		if !exists {
			streamMap = make(map[int]Stream)
		}
		stream, exists := streamMap[info.Stream]
		if !exists {
			stream = Stream{}
			stream.TrackID = info.Track
			stream.StreamID = info.Stream
		}
		switch info.Field {
		case SiTypeName:
			stream.TypeName = info.Value
		case SiAttr:
			stream.Attr = info.Value
		case SiLangCode:
			stream.LangCode = info.Value
		case SiLangName:
			stream.LangName = info.Value
		case SiCodecID:
			stream.CodecID = info.Value
		case SiCodecShort:
			stream.CodecShort = info.Value
		case SiCodecLong:
			stream.CodecLong = info.Value
		case SiBitrate:
			stream.Bitrate = info.Value
		case SiChannels:
			stream.Channels = atoi(info.Value)
		case SiSampleRate:
			stream.SampleRate = atoi(info.Value)
		case SiBitsPerSample:
			stream.BitsPerSample = atoi(info.Value)
		case SiResolution:
			stream.Resolution = info.Value
		case SiAspectRatio:
			stream.AspectRatio = info.Value
		case SiFrameRate:
			stream.FrameRate = info.Value
		case SiPgSizeOrFlg:
			stream.PgSizeOrFlag = info.Value
		case SiLangCode2:
			stream.LangCode2 = info.Value
		case SiLangName2:
			stream.LangName2 = info.Value
		case SiLongDesc:
			stream.LongDesc = info.Value
		case SiUiHeaderHTML:
			stream.UIHeaderHTML = info.Value
		case SiTrackID:
			stream.StreamTrackID = info.Value
		case SiUnknown38:
			stream.Unknown38 = info.Value
		case SiDefaultFlag:
			stream.DefaultFlag = info.Value == "Default"
		case SiChannelLayout:
			stream.ChannelLayout = info.Value
		case SiNotes:
			stream.Notes = info.Value
		default:
			if stream.Raw == nil {
				stream.Raw = make(map[int]string)
			}
			stream.Raw[info.Field] = info.Value
		}
		streamMap[info.Stream] = stream
		streams[info.Track] = streamMap
	}
	return streams
}
