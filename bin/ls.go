package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Velocidex/go-ext4/parser"
	kingpin "github.com/alecthomas/kingpin/v2"
)

var (
	ls_command = app.Command(
		"ls", "List files.")

	ls_command_file_arg = ls_command.Arg(
		"file", "The image file to inspect",
	).Required().OpenFile(os.O_RDONLY, os.FileMode(0666))

	ls_command_arg = ls_command.Arg(
		"path", "The path to list separated by /",
	).Default("/").String()

	ls_command_image_offset = ls_command.Flag(
		"image_offset", "An offset into the file.",
	).Default("0").Int64()

	ls_command_verbose = ls_command.Flag(
		"vebode", "More information").
		Short('v').Default("false").Bool()

	ls_command_recursive = ls_command.Flag(
		"recurse", "List directories recursively").
		Short('r').Default("false").Bool()
)

func doLS() {
	reader, _ := parser.NewPagedReader(
		getReader(*ls_command_file_arg), 1024, 10000)

	ext4, err := parser.GetEXT4Context(reader)
	kingpin.FatalIfError(err, "Can not open filesystem")

	components := strings.Split(*ls_command_arg, "/")
	inode, err := ext4.OpenInodeWithPath(components)
	kingpin.FatalIfError(err, "Can not open inode for %v", *ls_command_arg)

	kingpin.FatalIfError(listDir(ext4, inode),
		"Can not open inode for %v", *ls_command_arg)
}

func listDir(ext4 *parser.EXT4Context, inode *parser.Inode) error {
	entries, err := inode.Dir(ext4)
	if err != nil {
		return err
	}

	if *ls_command_verbose {
		for _, e := range entries {
			if e.Name() == "." || e.Name() == ".." {
				continue
			}

			fmt.Printf("%v\n", e.Dict())
		}

	} else {
		for _, e := range entries {
			if e.Name() == "." || e.Name() == ".." {
				continue
			}

			fmt.Printf("%v % 5d %d %d % 20d %v %v\n",
				e.Mode(), e.Inode(), e.Uid(), e.Gid(),
				e.Size(), e.ModTime().Format(time.RFC3339), e.FullPath())

			if e.Mode().IsDir() && *ls_command_recursive {
				inode, err := ext4.OpenInode(e.Inode(), e.Components())
				if err != nil {
					continue
				}
				listDir(ext4, inode)
			}

		}
	}

	return nil
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case "ls":
			doLS()
		default:
			return false
		}
		return true
	})
}
