package main

import (
	"course/config"
	_ "course/config"
	"course/exps"
	_ "course/pkg/clock"
	"course/presenter"
	"course/scene"
	"course/stats"
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	runExp()
}

func runExp() {
	pres := new(presenter.Presenter)

	st := stats.NewDriversStats()

	for name := range exps.Optimizers() {
		err := os.Mkdir(fmt.Sprintf("exps/output/%s", name), 0777)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				continue
			}
			log.Fatal(err)
		}
	}

	for expCount := 0; expCount < config.C().ExperimentsCount; expCount++ {
		ttBuilder, dhBuilder, bsBuilder := scene.GenScene()
		for k, opt := range exps.Optimizers() {
			tt, dh, bs := ttBuilder.Build(), dhBuilder.Build(), bsBuilder.Build()
			opt.Optimize(tt, bs, dh)
			st.Collect(tt, dh, bs, k)

			pres.Present(fmt.Sprintf("exps/output/%s/%d", k, expCount+1), tt, dh, bs)
		}

	}

	err := st.SaveStatistics("exps/output/statistics.pdf")
	if err != nil {
		log.Fatal(err)
	}
}
