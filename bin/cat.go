package main

import (
	"io"
	"os"
	"strings"

	"github.com/Velocidex/go-ext4/parser"
	kingpin "github.com/alecthomas/kingpin/v2"
)

var (
	cat_command = app.Command(
		"cat", "Cat a FAT Cluster.")

	cat_command_file_arg = cat_command.Arg(
		"file", "The image file to inspect",
	).Required().OpenFile(os.O_RDONLY, os.FileMode(0666))

	cat_command_arg = cat_command.Arg(
		"path", "The first cluster to read from.",
	).Required().String()

	cat_command_image_offset = cat_command.Flag(
		"image_offset", "An offset into the file.",
	).Default("0").Int64()
)

func doCAT() {
	reader, _ := parser.NewPagedReader(
		getReader(*cat_command_file_arg), 1024, 10000)

	ext4, err := parser.GetEXT4Context(reader)
	kingpin.FatalIfError(err, "Can not open filesystem")

	components := strings.Split(*cat_command_arg, "/")
	inode, err := ext4.OpenInodeWithPath(components)
	kingpin.FatalIfError(err, "Can not read files")

	stream, err := inode.GetReader(ext4)
	kingpin.FatalIfError(err, "Can not read files")

	var fd io.WriteCloser = os.Stdout
	buf := make([]byte, 4096)
	var offset int64
	for {
		n, err := stream.ReadAt(buf, offset)
		if err == io.EOF {
			break
		}

		if err != nil {
			kingpin.FatalIfError(err, "Can not read files")
		}

		fd.Write(buf[:n])
		offset += int64(n)
	}
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case "cat":
			doCAT()
		default:
			return false
		}
		return true
	})
}
