package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/diamondburned/arikawa/session"
	"github.com/diamondburned/arikawa/state"
	"github.com/joho/godotenv"
)

func main() {
	cfgPath := flag.String("c", "", "Path to the config (.env) file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s [flags...] mountpoint\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *cfgPath != "" {
		if err := godotenv.Load(*cfgPath); err != nil {
			log.Fatalln("Can't parse config file at", *cfgPath)
		}
	}

	var (
		token      = os.Getenv("TOKEN")
		username   = os.Getenv("USERNAME")
		password   = os.Getenv("PASSWORD")
		mountpoint = flag.Arg(0)
	)

	if mountpoint == "" {
		flag.Usage()
		os.Exit(2)
	}

	var ses *session.Session
	var err error

	switch {
	case token != "":
		ses, err = session.New(token)
	case username != "" && password != "":
		ses, err = session.Login(username, password, "")
	default:
		log.Fatalln("No token or username and password given.")
	}

	if err != nil {
		log.Fatalln("Failed to authenticate:", err)
	}

	s, err := state.NewFromSession(ses, state.NewDefaultStore(nil))
	if err != nil {
		log.Fatalln("Failed to create a Discord state:", err)
	}

	if err := s.Open(); err != nil {
		log.Fatalln("Failed to open a Discord connection:", err)
	}
	defer s.Close()

	FS, err := NewFS(s)
	if err != nil {
		log.Fatalln("Failed to create a filesystem:", err)
	}

	log.Println("Created a filesystem")

	c, err := fuse.Mount(mountpoint)
	if err != nil {
		log.Fatalln("Failed to mount FUSE:", err)
	}

	log.Println("Mount point created at", mountpoint)

	u, err := s.Me()
	if err != nil {
		log.Fatalln("Failed to get myself:", err)
	}

	log.Println("Serving for user", u.Username)

	if err := fs.Serve(c, FS); err != nil {
		log.Fatalln("Failed to serve filesystem:", err)
	}

	// Block until there's an error
	<-c.Ready

	if err := c.MountError; err != nil {
		log.Fatalln("Mount error:", err)
	}
}