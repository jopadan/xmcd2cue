/* xmcd to .cue sheet converter © 2001-2015 Viktor Szakáts

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE. */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type album struct {
	filename  string
	id        string
	artist    string
	title     string
	year      string
	genre     string
	ext       string
	preOffset int
	lengthSec int
	trkList   []*trk
}

type trk struct {
	offset int
	artist string
	title  string
	ext    string
}

var swapalb = flag.Bool("t", false, "swap album artist and title")
var swaptrk = flag.Bool("s", false, "swap artist and title in track titles")
var rename = flag.Bool("n", false, "rename the .cue sheet file according to album title")

func main() {
	flag.Parse()

	mask := flag.Arg(0)

	if mask != "" {
		fns, _ := filepath.Glob(mask)
		if fns != nil {
			for _, fn := range fns {
				alb, ok := albReadFreeDb(fn)
				if ok {

					if *rename {
						if alb.artist != "" {
							fn = "(" + alb.artist + ") "
						} else {
							fn = ""
						}
						fn += alb.title
						fn = fNameFilter(fn)
					}

					fn = strings.TrimSuffix(fn, filepath.Ext(fn)) + ".cue"

					if albWriteCue(alb, fn) {
						fmt.Printf("Successfully converted to '%s'", fn)
					} else {
						fmt.Printf("Error writing '%s'", fn)
					}
				}
			}
		}
	} else {
		fmt.Println("xmcd to .cue sheet converter 2.0.1 © 2001-2015 Viktor Szakáts")
		fmt.Println("https://github.com/vszakats/xmcd2cue/")
		fmt.Println()
		fmt.Println("Syntax: xmcd2cue [options[s]] xmcd-file[mask]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println()
		fmt.Println(" -t  Swap album artist and title")
		fmt.Println(" -s  Swap artist and title in track titles")
		fmt.Println(" -n  Rename the .cue sheet file according to album title")
	}
}

func albReadFreeDb(fn string) (album, bool) {
	alb := *new(album)
	alb.filename = fn
	alb.id = ""
	alb.title = ""
	alb.year = ""
	alb.genre = ""
	alb.ext = ""
	alb.preOffset = 0
	alb.lengthSec = 0
	alb.trkList = make([]*trk, 100)
	f, err := os.Open(fn)
	if err != nil {
		return alb, false
	}
	offset := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) <= 256 && line != "" {
			if strings.HasPrefix(line, "#") {
				line = strings.TrimSpace(line[1:])
				switch {
				case strings.HasPrefix(line, "xmcd"):
					// Valid input file
				case strings.HasPrefix(line, "Disc length: "):
					line = strings.Replace(line, "Disc length: ", "", 1)
					line = strings.Replace(line, " seconds", "", 1)
					alb.lengthSec, _ = strconv.Atoi(line)
				default:
					preOffset, err := strconv.Atoi(line)
					if err == nil {
						if offset == 0 {
							alb.preOffset = preOffset
						} else {
							alb.trkList[offset-1] = new(trk)
							alb.trkList[offset-1].offset = preOffset
						}
						offset++
					}
				}
			} else {
				value := strings.SplitN(line, "=", 2)
				if len(value) == 2 {
					switch {
					case strings.HasPrefix(value[0], "DISCID"):
						// Concatenating may not be standard
						alb.id += value[1]
					case strings.HasPrefix(value[0], "DTITLE"):
						alb.title += value[1]
					case strings.HasPrefix(value[0], "DYEAR"):
						// Concatenating may not be standard
						alb.year += value[1]
					case strings.HasPrefix(value[0], "DGENRE"):
						// Concatenating may not be standard
						alb.genre += value[1]
					case strings.HasPrefix(value[0], "DEXT"):
						alb.ext += value[1]
					case strings.HasPrefix(value[0], "TTITLE"):
						value[0] = strings.Replace(value[0], "TTITLE", "", 1)
						track, err := strconv.Atoi(value[0])
						if err == nil && track >= 0 && track <= 99 {
							if alb.trkList[track] == nil {
								alb.trkList[track] = new(trk)
							}
							alb.trkList[track].title += value[1]
						}
					case strings.HasPrefix(value[0], "EXTT"):
						value[0] = strings.Replace(value[0], "EXTT", "", 1)
						track, err := strconv.Atoi(value[0])
						if err == nil && track >= 0 && track <= 99 {
							if alb.trkList[track] == nil {
								alb.trkList[track] = new(trk)
							}
							alb.trkList[track].ext += value[1]
						}
					}
				}
			}
		}
	}
	f.Close()

	// Post-process album title

	var artist string
	var title string

	if value := strings.SplitN(alb.title, " / ", 2); len(value) == 2 {
		artist = value[0]
		title = value[1]
	} else if value := strings.SplitN(alb.title, " - ", 2); len(value) == 2 { // Non-standard
		artist = value[0]
		title = value[1]
	} else {
		artist = ""
		title = alb.title
	}

	title = strings.TrimSpace(title)
	artist = strings.TrimSpace(artist)

	if *swapalb {
		alb.artist = title
		alb.title = artist
	} else {
		alb.artist = artist
		alb.title = title
	}

	// Post-process track titles

	for _, trk := range alb.trkList {
		if trk != nil {
			if value := strings.SplitN(trk.title, " / ", 2); len(value) == 2 {
				artist = value[0]
				title = value[1]
			} else if value := strings.SplitN(trk.title, " - ", 2); len(value) == 2 { // Non-standard
				artist = value[0]
				title = value[1]
			} else {
				artist = ""
				title = trk.title
			}

			title = strings.TrimSpace(title)
			artist = strings.TrimSpace(artist)

			if *swaptrk {
				trk.artist = title
				trk.title = artist
			} else {
				trk.artist = artist
				trk.title = title
			}
		}
	}

	return alb, true
}

func albWriteCue(alb album, fn string) bool {
	f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return false
	}

	if alb.artist != "" {
		io.WriteString(f, fmt.Sprintf("PERFORMER \"%.80s\"\n", alb.artist))
	}

	io.WriteString(f, fmt.Sprintf("TITLE \"%.80s\"\n", alb.title))

	if alb.lengthSec != 0 {
		sec := alb.lengthSec
		hour := int(sec / 3600)
		sec -= hour * 3600
		min := int(sec / 60)
		sec -= int(min * 60)
		io.WriteString(f, fmt.Sprintf("REM Total length: %02d:%02d:%02d\n", hour, min, sec))
	}

	if alb.id != "" {
		io.WriteString(f, fmt.Sprintf("REM ID: %s\n", alb.id))
	}
	if *swapalb {
		io.WriteString(f, fmt.Sprintf("REM Command line option: %s\n", "-t"))
	}
	if *swaptrk {
		io.WriteString(f, fmt.Sprintf("REM Command line option: %s\n", "-s"))
	}
	if *rename {
		io.WriteString(f, fmt.Sprintf("REM Command line option: %s\n", "-n"))
	}

	io.WriteString(f, fmt.Sprintf("FILE \"%s\" MP3\n", strings.TrimSuffix(fn, filepath.Ext(fn))+".mp3"))

	prevOffset := alb.preOffset

	for nr, trk := range alb.trkList {
		if trk != nil {

			io.WriteString(f, fmt.Sprintf("  TRACK %02d AUDIO\n", nr+1))

			if trk.title != "" {
				io.WriteString(f, fmt.Sprintf("    TITLE \"%.80s\"\n", trk.title))
			}
			if trk.artist != "" {
				io.WriteString(f, fmt.Sprintf("    PERFORMER \"%.80s\"\n", trk.artist))
			}

			// Add PREGAP for CD writing
			if nr == 0 {
				io.WriteString(f, "    PREGAP 00:02:00\n")
			}

			frg := int((prevOffset - alb.preOffset) / 75)
			io.WriteString(f, fmt.Sprintf("    INDEX 01 %02d:%02d:%02d\n", int(frg/60), frg-(int(frg/60)*60), int((prevOffset-alb.preOffset)-(frg*75))))

			prevOffset = trk.offset
		}
	}

	f.Close()

	return true
}

func fNameFilter(fn string) string {
	out := ""
	for _, chr := range fn {
		if unicode.IsPrint(chr) && !strings.ContainsRune(":/?*\\\"", chr) {
			out += string(chr)
		} else {
			out += "_"
		}
	}
	return out
}

func xmcdDataConv(str string) string {
	str = strings.Replace(str, "\\n", "\n", -1)
	str = strings.Replace(str, "\\t", "\t", -1)
	str = strings.Replace(str, "\\\\", "\\", -1)
	return str
}
