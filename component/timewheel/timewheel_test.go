package timewheel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTimeWheel(t *testing.T) {
	tw := New(5*time.Millisecond, 12)
	tw.Start()
	defer tw.Stop()
	durations := []time.Duration{
		10 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
		2 * time.Second,
		3 * time.Second,
	}

	wg := sync.WaitGroup{}
	wg.Add(len(durations))
	for _, duration := range durations {
		d := duration
		tw.AfterFunc(d, func() {
			fmt.Println("task run", d, time.Now().UnixMilli())
			wg.Done()
		})
	}

	// 等待任务执行完
	wg.Wait()

}

func TestReset(t *testing.T) {
	tw := New(5*time.Millisecond, 12)
	tw.Start()
	defer tw.Stop()

	wg := sync.WaitGroup{}
	wg.Add(1)
	timer := tw.AfterFunc(time.Second, func() {
		fmt.Println("task run", time.Now().UnixMilli())
		wg.Done()
	})

	timer.Reset(time.Second * 3)

	wg.Wait()
}

func TestStandardReset(t *testing.T) {
	timer := time.AfterFunc(time.Second, func() {
		fmt.Println("task run", time.Now().UnixMilli())
	})

	time.Sleep(time.Second * 2)
	timer.Reset(time.Second * 3)

	time.Sleep(time.Second * 5)
}

func BenchmarkTimeWheel_StartStop(b *testing.B) {
	tw := New(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()

	cases := []struct {
		name string
		N    int // the data size (i.e. number of existing timers)
	}{
		{"N-1m", 1000000},
		{"N-5m", 5000000},
		{"N-10m", 10000000},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			base := make([]*Timer, c.N)
			for i := 0; i < len(base); i++ {
				base[i] = tw.AfterFunc(genD(i), func() {})
			}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				tw.AfterFunc(time.Second, func() {}).Stop()
			}

			b.StopTimer()
			for i := 0; i < len(base); i++ {
				base[i].Stop()
			}
		})
	}
}

func BenchmarkStandardTimer(b *testing.B) {
	cases := []struct {
		name string
		N    int // the data size (i.e. number of existing timers)
	}{
		{"N-1m", 1000000},
		{"N-5m", 5000000},
		{"N-10m", 10000000},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			base := make([]*time.Timer, c.N)
			for i := 0; i < len(base); i++ {
				base[i] = time.AfterFunc(genD(i), func() {})
			}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				time.AfterFunc(time.Second, func() {}).Stop()
			}

			b.StopTimer()
			for i := 0; i < len(base); i++ {
				base[i].Stop()
			}
		})
	}
}

func genD(i int) time.Duration {
	return time.Duration(i%10000) * time.Millisecond
}
