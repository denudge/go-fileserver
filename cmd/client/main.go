package main

import (
	"flag"
	"fmt"
	"github.com/denudge/go-fileserver/pkg/fileserver"
	"os"
)

func main() {
	flag.StringVar(&fileserver.ServerAddress, "server-address", "http://localhost:8080", "HTTP host to stream files to/from")
	flag.Parse()

	// parse command
	if len(os.Args) < 4 {
		fmt.Println("missing arguments")
		printUsage(os.Args[0])
		os.Exit(1)
	}

	progName := os.Args[0]
	command := os.Args[1]
	folder := os.Args[2]
	file := os.Args[3]

	c := fileserver.NewClient()

	switch command {
	case CommandUpload:
		if len(os.Args) < 5 {
			fmt.Println("missing local file argument")
			printUsage(progName)
			os.Exit(1)
		}

		err := c.UploadFile(folder, file, os.Args[4])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("file successfully uploaded")
	case CommandDownload:
		if len(os.Args) < 5 {
			fmt.Println("missing local file argument")
			printUsage(progName)
			os.Exit(1)
		}

		err := c.DownloadFile(folder, file, os.Args[4])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("file successfully downloaded")
	case CommandDelete:
		err := c.DeleteFile(folder, file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("file successfully deleted")
	default:
		fmt.Printf("unknown command %q\n", command)
		os.Exit(1)
	}
}
