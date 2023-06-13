package helpers

import (
	"time"

	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
)

// function to simulate a progress bar
func ShowProgressBar(total int) *mpb.Progress {
	p := mpb.New(mpb.WithWidth(64))

	bar := p.Add(int64(total),
		mpb.NewBarFiller("[=>-|"),
		mpb.PrependDecorators(
			decor.CountersNoUnit("%d / %d", decor.WCSyncSpaceR),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.AverageETA(decor.ET_STYLE_GO, decor.WCSyncSpace),
				"done!",
			),
		),
	)

	go func() {
		max := 100 * time.Millisecond
		for i := 0; i < total; i++ {
			time.Sleep(max)
			bar.Increment()
		}
		p.Wait()
	}()

	return p
}
