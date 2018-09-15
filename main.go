package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
	"errors"

	"github.com/djhworld/gomeboycolor/cartridge"
	"github.com/djhworld/gomeboycolor/config"
	"github.com/djhworld/gomeboycolor/gbc"
	"github.com/djhworld/gomeboycolor-glfw/saves"
	"gopkg.in/urfave/cli.v1"
)

const TITLE string = "gomeboycolor"

var VERSION string = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "gomeboycolor"
	app.Usage = "Gameboy Color emulator"
	app.ArgsUsage = "<path-to-ROM>"
	app.Version=VERSION
	app.UsageText= "gomeboycolor [flags] <path-to-ROM-file>"
	app.Action = run

	app.Flags = []cli.Flag {
		cli.StringFlag{
		  Name: "title",
		  Value: TITLE,
		  Usage: "Title to use",
		},
		cli.BoolFlag{
		  Name: "showfps",
		  Usage: "Calculate and display frames per second",
		},
		cli.BoolFlag{
		  Name: "skipboot",
		  Usage: "Skip boot sequence",
		},
		cli.BoolFlag{
		  Name: "no-color",
		  Usage: "Disable Gameboy Color Hardware",
		},
		cli.BoolFlag{
		  Name: "headless",
		  Usage: "Run emulator without output",
		},
		cli.Int64Flag{
		  Name: "fpslock",
		  Value: 58,
		  Usage: "Lock framerate to this. Going higher than default might be unstable!",
		},
		cli.IntFlag{
		  Name: "size",
		  Value: 1,
		  Usage: "Screen size multiplier",
		},
		cli.BoolFlag{
		  Name: "debug",
		  Usage: "Enable debugger",
		},
		cli.BoolFlag{
		  Name: "dump",
		  Usage: "Print state of machine after each cycle (WARNING - WILL RUN SLOW)",
		},
		cli.StringFlag{
		  Name: "b",
		  Value: "0x0000",
		  Usage: "Break into debugger when PC equals a given value between 0x0000 and 0xFFFF",
		},
	  }
  
	err := app.Run(os.Args)
	if err != nil {
	  log.Fatal(err)
	}
}


func run(c *cli.Context) error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("%s. %s\n", TITLE, VERSION)
	fmt.Println("Copyright (c) 2018. Daniel James Harper.")
	fmt.Println("http://djhworld.github.io/gomeboycolor")
	fmt.Println(strings.Repeat("*", 120))

	if c.NArg() != 1 {
		return errors.New("Please specify the location of a ROM to boot")
	}

	var colorMode bool = true
	if c.Bool("no-color") {
		colorMode = false
	}

	//Parse and validate settings file (if found)
	conf := &config.Config{
		Title:         TITLE,
		ScreenSize:    c.Int("size"),
		SkipBoot:      c.Bool("skipboot"),
		DisplayFPS:    c.Bool("showfps"),
		ColorMode:     colorMode,
		Debug:         c.Bool("debug"),
		BreakOn:       c.String("b"),
		DumpState:     c.Bool("dump"),
		Headless:      c.Bool("headless"),
		FrameRateLock: c.Int64("fpslock"),
	}
	fmt.Println(conf)

	cart, err := createCartridge(c.Args().Get(0))
	if err != nil {
		return err
	}

	log.Println("Starting emulator")

	emulator, err := gbc.Init(cart, getSaveStore(), conf, NewGlfwIO(conf.FrameRateLock, conf.Headless, conf.DisplayFPS))
	if err != nil {
		return err
	}

	//Starts emulator code in a goroutine
	go emulator.Run()

	//lock the OS thread here
	runtime.LockOSThread()

	//set the IO controller to run indefinitely (it waits for screen updates)
	emulator.RunIO()

	log.Println("Goodbye!")
	return nil
}

func getSaveStore() *saves.FileSystemStore {
	user, _ := user.Current()
	saveDir := user.HomeDir + "/.gomeboycolor/saves"

	os.MkdirAll(saveDir, os.ModeDir)

	return saves.NewFileSystemStore(saveDir)
}

func createCartridge(romFilename string) (*cartridge.Cartridge, error) {
	romContents, err := retrieveROM(romFilename)
	if err != nil {
		return nil, err
	}

	return cartridge.NewCartridge(romFilename, romContents)
}

func retrieveROM(filename string) ([]byte, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err
}
