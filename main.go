package main

import (
	"fmt"
	"image"
	"image/png"
	"os/exec"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/oliamb/cutter"
)

func main() {
	// Get the bounds of the screen
	bounds := screenshot.GetDisplayBounds(0)

	// Create a pipe to write frames to FFmpeg
	ffmpegCmd := exec.Command("ffmpeg",
		"-f", "image2pipe",
		"-framerate", "10", // record at 10 fps
		"-i", "-", // read frames from stdin
		"-pix_fmt", "yuv420p", // set pixel format for compatibility with most players
		"-y", // overwrite output file without asking
		"recording.mp4",
	)
	pipeIn, err := ffmpegCmd.StdinPipe()
	if err != nil {
		fmt.Println("Failed to create FFmpeg input pipe:", err)
		return
	}
	defer pipeIn.Close()

	// Start FFmpeg process
	err = ffmpegCmd.Start()
	if err != nil {
		fmt.Println("Failed to start FFmpeg process:", err)
		return
	}
	defer ffmpegCmd.Process.Kill()

	// Record for 10 seconds
	for i := 0; i < 100; i++ {
		// Capture the screen image
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			fmt.Println("Failed to capture screen:", err)
			return
		}

		// Crop the image to remove the menu bar, task bar, etc.
		croppedImg, err := cutter.Crop(img, cutter.Config{
			Width:   bounds.Dx(),
			Height:  bounds.Dy() - 50, // adjust height to remove 50 pixels from the bottom of the screen
			Anchor:  image.Point{0, 0},
			Options: cutter.Copy,
		})
		if err != nil {
			fmt.Println("Failed to crop image:", err)
			return
		}

		// Write the cropped image to the FFmpeg input pipe
		err = png.Encode(pipeIn, croppedImg)
		if err != nil {
			fmt.Println("Failed to write image to FFmpeg pipe:", err)
			return
		}

		// Wait for 100 milliseconds before capturing the next frame
		time.Sleep(time.Millisecond * 100)
	}

	fmt.Println("Screen recording completed")
}
