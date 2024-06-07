package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Velocidex/go-ext4/parser"
	kingpin "github.com/alecthomas/kingpin/v2"
)

var (
	stat_command = app.Command(
		"stat", "Stat a FAT file.")

	stat_command_file_arg = stat_command.Arg(
		"file", "The image file to inspect",
	).Required().OpenFile(os.O_RDONLY, os.FileMode(0666))

	stat_command_arg = stat_command.Arg(
		"path", "The first cluster to read from.",
	).Required().String()

	stat_command_image_offset = stat_command.Flag(
		"image_offset", "An offset into the file.",
	).Default("0").Int64()
)

func doStat() {
	reader, _ := parser.NewPagedReader(
		getReader(*stat_command_file_arg), 1024, 10000)

	ext4, err := parser.GetEXT4Context(reader)
	kingpin.FatalIfError(err, "Can not open filesystem")

	components := strings.Split(*stat_command_arg, "/")
	inode, err := ext4.OpenInodeWithPath(components)
	kingpin.FatalIfError(err, "Can not open inode for %v", *stat_command_arg)

	stat := inode.Stat()
	Dump(stat)

	// Get the runlist
	runs := inode.Runs(ext4)

	fmt.Printf("\n\nExtents (%v)\n", len(runs))
	fmt.Printf("FileOffset      Length    DiskOffset\n")
	for _, r := range runs {
		fmt.Printf("% #10x % #11x % #13x\n",
			r.FileOffset, r.Length, r.DiskOffset)
	}
	fmt.Println("")
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case "stat":
			doStat()
		default:
			return false
		}
		return true
	})
}
