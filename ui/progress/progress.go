package progress

import (
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"io"
	"strings"
	"sync"
	"time"
)

// Bar represents a progress bar
type Bar struct {
	b       *mpb.Bar
	Name    string
	Desc    string
	Start   time.Time
	Current time.Time
}

// Progress is a progress object
type Progress struct {
	p       *mpb.Progress
	bars    sync.Map
	builder strings.Builder
	lock    sync.Mutex
}

func ShowDesc(bar *Bar, wcc ...decor.WC) decor.Decorator {
	producer := func(bar *Bar, wcc ...decor.WC) decor.DecorFunc {
		return func(s decor.Statistics) string {
			return bar.Desc
		}
	}
	return decor.Any(producer(bar), wcc...)
}

func (p *Progress) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) (body io.ReadCloser) {
	p.IOBar(src, stream, totalSize)
	return io.NopCloser(strings.NewReader(p.builder.String()))
}

// CreateProgress creates a new progress object
func CreateProgress() *Progress {
	p := &Progress{
		p:    mpb.New(),
		bars: sync.Map{},
	}
	return p
}
func (p *Progress) IOBar(name string, reader io.Reader, total int64) {

	bar := p.p.New(total,
		mpb.BarStyle().Rbound("|"),
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
	)
	// create proxy reader
	proxyReader := bar.ProxyReader(reader)
	defer proxyReader.Close()

	// copy from proxyReader, ignoring errors
	_, _ = io.Copy(&p.builder, proxyReader)
	p.p.Wait()
}

// Add adds a new bar to the progress
func (p *Progress) Add(name string, total int64) {

	_, ok := p.bars.Load(name)
	if ok {
		return
	}
	var bar Bar
	bar.b = p.p.AddBar(
		total,
		mpb.BarWidth(100),
		mpb.PrependDecorators(
			decor.Name(name, decor.WCSyncSpaceR),
			decor.CountersNoUnit("[%d/%d]", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.Elapsed(decor.ET_STYLE_GO),
				"Success",
			),
		),
	)
	bar.Start = time.Now()
	bar.Current = time.Now()
	bar.Name = name
	p.bars.Store(name, &bar)
}

func (p *Progress) Increment(name string, n int64) {

	bar, ok := p.bars.Load(name)
	if !ok {
		return
	}
	bar.(*Bar).b.IncrInt64(n)
	bar.(*Bar).Current = time.Now()
}

func (p *Progress) Current(name string, n int64, desc ...string) {
	bar, ok := p.bars.Load(name)
	if !ok {
		return
	}
	bar.(*Bar).b.SetCurrent(n)
	bar.(*Bar).Current = time.Now()
	bar.(*Bar).Desc = desc[0]
}

func (p *Progress) SetTotal(name string, n int64) {
	bar, ok := p.bars.Load(name)
	if !ok {
		return
	}
	bar.(*Bar).b.SetTotal(n, false)
	bar.(*Bar).Current = time.Now()
}

func (p *Progress) Next(name string) {
	p.Increment(name, 1)
}

func (p *Progress) Done(name string) {
	bar, ok := p.bars.Load(name)
	if !ok {
		return
	}
	bar.(*Bar).b.EnableTriggerComplete()
}

func (p *Progress) Wait(name string) {
	bar, ok := p.bars.Load(name)
	if !ok {
		return
	}
	bar.(*Bar).b.Wait()
}
