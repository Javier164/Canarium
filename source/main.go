package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Config struct {
	Color string `json:"color"`
	City  string `json:"city"`
	State string `json:"state"`
	ZIP   string `json:"zip"`
	Key   string `json:"wkey"`
	Feed  string `json:"feed"`
}

type Data struct {
	Observation struct {
		UV       byte   `json:"uv_index"`
		Cover    string `json:"sky_cover"`
		Phrase   string `json:"phrase_32char"`
		Desc     string `json:"uv_desc"`
		Imperial struct {
			Temp int16 `json:"temp"`
			High int16 `json:"temp_max_24hour"`
			Low  int8  `json:"temp_min_24hour"`
			Dew  int8  `json:"dewpt"`
			Wind uint8 `json:"wspd"`
		} `json:"imperial"`
	} `json:"observation"`
}

type CommonData struct {
	MoonPhase []string `json:"moonPhase"`
	Moonrise  []string `json:"moonriseTimeLocal"`
	Narrative []string `json:"narrative"`
	DayOfWeek []string `json:"dayOfWeek"`
	Sunrise   []string `json:"sunriseTimeLocal"`
	Sunset    []string `json:"sunsetTimeLocal"`
}

type Common struct {
	V3WxForecastDaily5Day CommonData `json:"v3-wx-forecast-daily-5day"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

func main() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Failed to read configuration file: %s", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalf("Failed to parse configuration file: %s", err)
	}

	r, g, b, err := ParseHexColor(config.Color)
	if err != nil {
		log.Fatalf("Failed to parse hexadecimal code: %s", err)
	}

	if err := sdl.Init(sdl.INIT_AUDIO | sdl.INIT_VIDEO); err != nil {
		log.Fatalf("Failed to initialize SDL: %s", err)
	}
	defer sdl.Quit()

	if err := mix.Init(mix.INIT_MP3 | mix.INIT_FLAC | mix.INIT_OGG); err != nil {
		log.Fatalf("Failed to initialize mixer: %s", err)
	}
	defer mix.Quit()

	if err := mix.OpenAudio(48000, sdl.AUDIO_S16, 2, 4096); err != nil {
		log.Fatalf("Failed to open mixer: %s", err)
	}
	defer mix.CloseAudio()

	if err := ttf.Init(); err != nil {
		log.Fatalf("Failed to initialize TTF: %s", err)
	}
	defer ttf.Quit()

	window, err := sdl.CreateWindow(fmt.Sprintf("Weather in %s, %s", config.City, config.State), sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 960, 720, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("Failed to create window: %s", err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("Failed to create renderer: %s", err)
	}
	defer renderer.Destroy()

	renderer.SetDrawColor(r, g, b, 255) // Default Color: #190160
	renderer.Clear()

	data, err := ParseObservationData(config)
	if err != nil {
		log.Fatalf("Failed to get observation data: %s", err)
	}

	common, err := ParseCommonData(config)
	if err != nil {
		log.Fatalf("Failed to get observation data: %s", err)
	}

	rss, err := ParseRSSFeed(config)
	if err != nil {
		log.Fatalf("Failed to parse RSS feed: %s", err)
	}

	nid := 0 // Narrative ID
	rid := 0 // RSS ID

	Update(renderer, common, data, config, rss, 0, 0)

	var music []string
	err = filepath.Walk("assets/music", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (filepath.Ext(path) == ".flac" ||
			filepath.Ext(path) == ".ogg" ||
			filepath.Ext(path) == ".mp3") {
			music = append(music, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to queue music: %s", err)
		return
	}

	var queue []string

	for _, file := range music {
		queue = append(queue, file)
	}

	rand.Shuffle(len(queue), func(i, j int) {
		queue[i], queue[j] = queue[j], queue[i]
	})

	go func() {
		for {
			if mix.Playing(-1) == 0 {
				if len(queue) == 0 {
					for _, file := range music {
						queue = append(queue, file)
					}

					rand.Shuffle(len(queue), func(i, j int) {
						queue[i], queue[j] = queue[j], queue[i]
					})
				}

				file := queue[0]
				queue = queue[1:]

				music, err := mix.LoadMUS(file)
				if err != nil {
					log.Fatalf("Failed to load music: %s", err)
				}
				defer music.Free()

				if err := music.Play(0); err != nil {
					log.Fatalf("Failed to play music: %s", err)
				} else {
					log.Println("Playing song:", file)
					for mix.PlayingMusic() {
						time.Sleep(time.Millisecond * 10)
					}
				}
			}

			time.Sleep(time.Millisecond * 10)
		}
	}()

	go func() { // Weather Data Goroutine
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ndata, err := ParseObservationData(config)
				if err != nil {
					log.Printf("Failed to get observation data: %s", err)
				} else {
					ncdata, err := ParseCommonData(config)
					if err != nil {
						log.Printf("Failed to get common data: %s", err)
					} else {
						nrdata, err := ParseRSSFeed(config)
						if err != nil {
							log.Printf("Failed to parse RSS feed: %s", err)
						} else {
							data = ndata
							common = ncdata
							rss = nrdata

							nid = nid + 1
							rid = rid + 1

							if nid > 5 {
								nid = 0
							}

							if rid > len(rss.Channel.Items)-1 {
								rid = 0
							}

							renderer.SetDrawColor(r, g, b, 255)
							renderer.Clear()

							Update(renderer, common, data, config, rss, nid, rid)
							log.Printf("Updated!")
						}
					}
				}
			}
		}
	}()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			case *sdl.KeyboardEvent:
				if event.(*sdl.KeyboardEvent).Keysym.Sym == sdl.K_q ||
					event.(*sdl.KeyboardEvent).Keysym.Sym == sdl.K_ESCAPE {
					running = false
					break
				}
			}
		}
		sdl.Delay(5)
	}
}
