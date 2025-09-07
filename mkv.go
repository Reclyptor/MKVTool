package main

func ListTracks(args Arguments) {
	disc := LoadDisc(args)
	PrintDiscTree(disc, args)
}

func RipTracks(args Arguments) {
	RipDisc(args)
}
