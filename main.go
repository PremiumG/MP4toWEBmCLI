package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	// Get the user's home directory

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return
	}

	// Construct the path to the Downloads folder
	downloadsPath := filepath.Join(homeDir, "Downloads")

	// Get cutoff date (one week ago)
	cutoffDate := time.Now().AddDate(0, 0, -1)

	// Create a slice to store MP4 filenames
	var mp4Files []string

	// Walk through the Downloads folder
	err = filepath.Walk(downloadsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if it's a regular file (not a directory)
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".mp4") {
			if info.ModTime().After(cutoffDate) {
				mp4Files = append(mp4Files, info.Name())
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	//--------------------------------------------------------------------------------
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	g4 := widgets.NewGauge()
	g4.Title = "Loading"
	g4.SetRect(0, 11, 100, 14)
	g4.Percent = 0
	g4.Label = ""
	g4.BarColor = ui.ColorCyan
	g4.LabelStyle = ui.NewStyle(ui.ColorYellow)
	ui.Render(g4)

	l := widgets.NewList()
	l.Title = "MP4 Files - 0"
	l.Rows = mp4Files
	l.TextStyle = ui.NewStyle(ui.ColorCyan)
	l.WrapText = false

	l.SetRect(0, 0, 100, 10)
	ui.Render(l)

	downloadFFMPEG(g4)
	g4.Title = "Waiting for start"
	g4.Percent = 0
	ui.Render(g4)

	previousKey := ""
	selectedIndex := 0

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
			if selectedIndex < len(mp4Files)-1 {
				selectedIndex++
			}
		case "k", "<Up>":
			l.ScrollUp()
			if selectedIndex > 0 {
				selectedIndex--
			}
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				l.ScrollTop()
			}
		case "<Home>":
			l.ScrollTop()
		case "G", "<End>":
			l.ScrollBottom()

		case "<Enter>":
			convertThatBitch(mp4Files[selectedIndex], downloadsPath, g4)
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}
		l.Title = "MP4 Files - " + strconv.Itoa(selectedIndex)

		ui.Render(l, g4)
	}
}

func downloadFFMPEG(g4 *widgets.Gauge) {
	url := "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip"
	outputFile := "ffmpeg.zip"
	g4.Title = "Loading"
	g4.Percent = 2
	ui.Render(g4)
	info, err := os.Stat("./ffmpeg")
	if os.IsNotExist(err) {

	} else {
		if info.IsDir() {
			g4.Percent = 100
			ui.Render(g4)
			return
		}
	}
	g4.Percent = 10
	ui.Render(g4)

	// Use wget or curl to download
	cmd := exec.Command("curl", "-L", url, "-o", outputFile) // Download FFmpeg to a file
	g4.Percent = 20
	ui.Render(g4)
	cmd.Stdout = nil
	cmd.Stderr = nil

	err = cmd.Run()
	g4.Percent = 50
	ui.Render(g4)
	if err != nil {
		fmt.Println("Error downloading FFmpeg:", err)
		return
	}

	// Untar or unzip based on the OS
	fmt.Println("Extracting FFmpeg...")
	extractCmd := exec.Command("powershell", "-Command", fmt.Sprintf("Expand-Archive -Path %s -DestinationPath ffmpeg", outputFile))
	extractCmd.Stdout = nil

	extractCmd.Stderr = nil
	g4.Percent = 55
	ui.Render(g4)
	err = extractCmd.Run()
	g4.Percent = 90
	ui.Render(g4)
	if err != nil {
		return
	}

}
func replaceExtension(filename, oldExt, newExt string) string {
	if strings.HasSuffix(filename, oldExt) {
		return filename[:len(filename)-len(oldExt)] + newExt
	}
	return filename
}
func convertThatBitch(path string, dFolder string, g4 *widgets.Gauge) error {
	g4.Title = "Running (Don't panic if it looks stuck)"
	g4.Percent = 2
	ui.Render(g4)
	dir, _ := os.Getwd()
	ffmpegDir := dir + "\\" + "ffmpeg\\ffmpeg-master-latest-win64-gpl\\bin\\ffmpeg.exe" // Adjust this to your FFmpeg directory
	inputFile := dFolder + "\\" + path                                                  // Replace with your input file
	outputFile := dFolder + "\\" + replaceExtension(path, ".mp4", ".webm")
	g4.Percent = 40
	ui.Render(g4)
	cmd := exec.Command(ffmpegDir, "-i", inputFile, outputFile)
	g4.Percent = 50
	ui.Render(g4)
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Run()
	g4.Title = "Done (select next)"
	g4.Percent = 100
	ui.Render(g4)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil

}
