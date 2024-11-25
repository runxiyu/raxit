package main

// TODO: Fix error handling

import (
	"errors"
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var err_video_driver = errors.New("cannot set video driver hint")

func main() {
	if err := main_internal(); err != nil {
		log.Fatalln(err)
	}
}

func main_internal() error {
	var font *ttf.Font
	var err error
	var window *sdl.Window
	var renderer *sdl.Renderer
	var w, h int32
	var bg sdl.Rect

	if err = ttf.Init(); err != nil {
		return err
	}
	defer ttf.Quit()


	if !sdl.SetHint(sdl.HINT_VIDEODRIVER, "wayland") {
		log.Println(err_video_driver)
	}

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}
	defer sdl.Quit()

	window, err = sdl.CreateWindow(
		"test",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600,
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE|sdl.WINDOW_ALLOW_HIGHDPI,
	)
	if err != nil {
		return err
	}
	defer window.Destroy()


	font, err = ttf.OpenFont("/usr/share/fonts/rsms-inter-vf-fonts/InterVariable.ttf", 20);
	if  err != nil {
		return err
	}
	defer font.Close()


	if renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED); err != nil {
		return err
	}

	if w, h, err = renderer.GetOutputSize(); err != nil {
		return err
	}

	bg = sdl.Rect{X: 0, Y: 0, W: w, H: h}
	renderer.FillRect(&bg)
	renderer.Present()

event_loop:
	for {
		if event := sdl.PollEvent(); event != nil {
			switch event.(type) {
			case *sdl.QuitEvent:
				break event_loop
			case *sdl.WindowEvent:
				window_event := event.(*sdl.WindowEvent)
				switch window_event.Event {
				case sdl.WINDOWEVENT_RESIZED:
					w, h, err = renderer.GetOutputSize()
					if err != nil {
						return err
					}

					bg = sdl.Rect{X: 0, Y: 0, W: w, H: h}
					renderer.SetDrawColor(0, 0, 0, 255)
					renderer.FillRect(&bg)

					var text *sdl.Surface
					if text, err = font.RenderUTF8Blended("Hello, World!", sdl.Color{R: 255, G: 255, B: 255, A: 255}); err != nil {
						return err
					}
					defer text.Free()

					var texture *sdl.Texture
					if texture, err = renderer.CreateTextureFromSurface(text); err != nil {
						return err
					}

					var text_rect sdl.Rect
					text_rect.W = text.W
					text_rect.H = text.H
					text_rect.X = 0
					text_rect.Y = 0

					renderer.Copy(texture, nil, &text_rect)


					renderer.Present()
					window.SetSize(window_event.Data1, window_event.Data2)
				}
			}
		}
	}

	return nil
}
