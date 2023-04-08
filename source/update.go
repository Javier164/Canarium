package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func Update(renderer *sdl.Renderer, common []Common, data *Data, config Config, rss *RSS, nid int, rid int) {
	Text(renderer, "Current Conditions", "assets/fonts/StarJR.ttf", 44, 50, 15, false, 0)

	Text(renderer, fmt.Sprintf("Temperature: %d\u00b0F", data.Observation.Imperial.Temp), "assets/fonts/StarJR.ttf", 32, 50, 75, false, 0)
	Text(renderer, fmt.Sprintf("Cloud Cover: %s", data.Observation.Cover), "assets/fonts/StarJR.ttf", 32, 50, 105, false, 0)
	Text(renderer, fmt.Sprintf("High: %d\u00b0F", data.Observation.Imperial.High), "assets/fonts/StarJR.ttf", 32, 50, 135, false, 0)
	Text(renderer, fmt.Sprintf("Low: %d\u00b0F", data.Observation.Imperial.Low), "assets/fonts/StarJR.ttf", 32, 50, 165, false, 0)
	Text(renderer, fmt.Sprintf("UV Index: %d (%s)", data.Observation.UV, data.Observation.Desc), "assets/fonts/StarJR.ttf", 32, 50, 195, false, 0)
	Text(renderer, fmt.Sprintf("Wind Speed: %dmph", data.Observation.Imperial.Wind), "assets/fonts/StarJR.ttf", 32, 50, 225, false, 0)

	clock := time.Now()
	format := clock.Format("3:04 PM")

	Text(renderer, format, "assets/fonts/StarJR.ttf", 48, 650, 15, false, 0)

	t, err := time.Parse("2006-01-02T15:04:05-0700", common[0].V3WxForecastDaily5Day.Moonrise[nid])
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return
	}

	dformat := t.Format("01/02/2006")
	tformat := t.Format("3:04:05 AM")

	t, err = time.Parse("2006-01-02T15:04:05-0700", common[0].V3WxForecastDaily5Day.Sunrise[nid])
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return
	}

	rformat := t.Format("3:04:05 AM")

	t, err = time.Parse("2006-01-02T15:04:05-0700", common[0].V3WxForecastDaily5Day.Sunset[nid])
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return
	}

	sformat := t.Format("3:04:05 PM")

	Text(renderer, fmt.Sprintf("Sunrise: %s", rformat), "assets/fonts/StarJR.ttf", 22, 50, 615, false, 0)
	Text(renderer, fmt.Sprintf("Sunset: %s", sformat), "assets/fonts/StarJR.ttf", 22, 275, 615, false, 0)
	Text(renderer, fmt.Sprintf("Moon Phase: %s, on %s at %s.", common[0].V3WxForecastDaily5Day.MoonPhase[nid], dformat, tformat), "assets/fonts/StarJR.ttf", 22, 50, 635, false, 0)

	str := WordWrap(fmt.Sprintf("%s: %s", common[0].V3WxForecastDaily5Day.DayOfWeek[nid], common[0].V3WxForecastDaily5Day.Narrative[nid]), 72)

	base := 655 // Base X coordinate used when word wrapping.
	for _, chunk := range str {
		Text(renderer, chunk, "assets/fonts/StarJR.ttf", 22, 50, int32(base), false, 0)
		base += 20
	}

	base = 325 // Now changed for the RSS feed.

	Text(renderer, "Today's News", "assets/fonts/StarJR.ttf", 44, 50, 265, false, 0)

	str = WordWrap(rss.Channel.Items[rid].Title, 48)
	for _, chunk := range str {
		Text(renderer, chunk, "assets/fonts/StarJR.ttf", 32, 50, int32(base), false, 0)
		base += 30
	}

	str = WordWrap(rss.Channel.Items[rid].Description, 72)
	for _, chunk := range str {
		Text(renderer, chunk, "assets/fonts/StarJR.ttf", 22, 50, int32(base)+20, false, 0)
		base += 20
	}

	renderer.Present()
}
