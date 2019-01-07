Downloads: https://github.com/jopadan/xmcd2cue/releases/tag/2.0.1

or

    go get github.com/jopadan/xmcd2cue

xmcd2cue
========
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE.md) [![Build Status](https://travis-ci.org/jopadan/xmcd2cue.svg)](https://travis-ci.org/jopadan/xmcd2cue) [![Go Report](https://goreportcard.com/badge/github.com/jopadan/xmcd2cue)](https://goreportcard.com/report/github.com/jopadan/xmcd2cue)

previous project included in GoLibs [gitee.com/GoLibs/xmcd2cue](https://gitee.com/GoLibs/xmcd2cue)

Original developers project links:
[github.com/vszakats/xmcd2cue](https://github.com/vszakats/xmcd2cue)
[Viktor Szakats](https://vszakats.net/)

Purpose
-------

This program will convert [xmcd compatible](http://www.freedb.org/pub/freedb/misc/freedb_database_format_specs.zip)
CD descriptor files (which can be downloaded from [freedb.org](http://www.freedb.org/)) into
[.cue](https://en.wikipedia.org/wiki/Cue_sheet_(computing)#Cue_sheet_syntax)
[sheet](https://web.archive.org/web/digitalx.org/cue-sheet/syntax/index.html) files,
so that you can use them players to see which song is played from an all-in-one CD rip.

There exists a freeware Windows GUI tool named [CueMaster](http://cuemaster.org/)
to handle similar tasks.

Content
-------

The package contains full source code, and a precompiled binary for
the Windows platform.

License
-------

This program is &copy;&nbsp;2001&ndash;2014 Viktor Szakáts. It is [free](https://www.gnu.org/philosophy/free-sw.html),
[open source](https://opensource.org/docs/definition.php) and it
is licensed under the [MIT License](LICENSE.md).

Usage
-----

<pre>
xmcd2cue [options[s]] xmcd-file[mask]

Options:

 -t  Swap album artist and title
 -s  Swap artist and title in track titles
 -n  Rename the .cue sheet file according to album title
</pre>

Notes
-----

The downloadable xmcd files do not always conform to the standard.
Sometimes wrong artist / title separator is used, so xmcd2cue will
recognize `" - "` as a valid separator as well, but it cannot do anything
with plain `"-"` and `"/"` (without surrounding spaces), if this is the
case you should fix the xmdb file manually. Sometimes the order of
the artist / title is reversed, in xmcd2cue you should use the `-t` and
`-s` command line options to correct this.

You can use wildcards to convert more than one file at once.

History
-------
<div>
<pre>
   2.0.1 (2019-01-07)

   	* forked from gitee.com to github.com
	* updated project locations in README.md

   2.0.0 (2013-07-16)

        + rewritten in Go
        + added -t option
        - deleted -L option

   1.0.9 (2012-10-02)

        * moved to GitHub
        * internal cleanups

   1.0.8 (2012-01-27)

        * hb_osnewline() -> hb_eol()
        * dropped Clipper compatibility
        * removed spaces at EOL
        * version update
        * updated comment on usage

   1.0.7 (2008-05-05)

        * Copyright date updates.
        * Converted spaces to tabs.
        * Removed @ from email addresses.
        % Reduced binary size.

   1.0.6 (2003-02-24)

        * URL change.

   1.0.5 (2002-05-26)

        % Some FUNCs converted to PROCs.
        % ExtGet() optimized.
        ! Chr( 13 ) + Chr( 10 ) -> hb_OSNewLine()
        ! TITLE and PERFORMER lines swapped again for the
          album.

   1.0.4 (2002-05-24)

        + Puts album ID to a REM.
        * Puts album length in a separate REM line instead
          of cluttering the TITLE.
        + Added PREGAP for the first track.
        + TITLE and PERFORMER lines swapped.
        ; Thanks to Ynte for the above ideas.
        + Puts command line switches in a REM.
        ! readme.html updated.
        ! Never write more than 99 tracks.
        ! Limit the TITLE and PERFORMER command parameters
          to a maximum of 80 characters.

   1.0.3 (2002-05-22)

        ! Fixed time format where the third part of the time was
          not correctly calculated. Thanks to Ynte for reporting
          this.
        * Updated email address.

   1.0.2 (2001-11-21)

        + Generates separate PERFORMER line for the album header instead
          of combining the artist name with the album title.

   1.0.1 (2001-11-13)

        ! Fixed handling of artist/title when more than one space
          was used as a separator between them.

   1.0.0 (2001-11-09)

        ; Initial version
</pre>
</div>

Internals
---------

The program has a complete read algorithm for the xmcd files, so it
can be easily extended to support other output formats.

The program is written in [Go](https://golang.org/) language.
You can compile it to binary executable on most platforms using:

   `go build xmcd2cue.go`

or run directly:

   `go run xmcd2cue.go <arguments>`

Author:<br />
Viktor Szakáts
<hr />

[Sample xmcd file](http://www.freedb.org/freedb/rock/7e09f20a)
---

Sample .cue file (generated from the file above)
---

<pre>
PERFORMER "Pink Floyd"
TITLE "Dark Side of The Moon"
REM Total length: 00:42:28
REM ID: 7e09f20a
REM Command line option: -t
FILE "7e09f20a.mp3" MP3
  TRACK 01 AUDIO
    TITLE "Speak to Me"
    PREGAP 00:02:00
    INDEX 01 00:00:00
  TRACK 02 AUDIO
    TITLE "Breathe"
    INDEX 01 01:13:18
  TRACK 03 AUDIO
    TITLE "On the Run"
    INDEX 01 04:00:41
  TRACK 04 AUDIO
    TITLE "Time"
    INDEX 01 07:33:41
  TRACK 05 AUDIO
    TITLE "The Great Gig in the Sky"
    INDEX 01 14:38:24
  TRACK 06 AUDIO
    TITLE "Money"
    INDEX 01 18:43:32
  TRACK 07 AUDIO
    TITLE "Us and Them"
    INDEX 01 25:15:31
  TRACK 08 AUDIO
    TITLE "Any Colour You Like"
    INDEX 01 33:05:48
  TRACK 09 AUDIO
    TITLE "Brain Damage"
    INDEX 01 36:31:18
  TRACK 10 AUDIO
    TITLE "Eclipse"
    INDEX 01 40:21:55
</pre>
