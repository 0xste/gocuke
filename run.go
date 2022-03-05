package gocuke

import (
	"github.com/cucumber/gherkin-go/v19"
	_ "github.com/cucumber/gherkin-go/v19"
	"github.com/cucumber/messages-go/v16"
	_ "github.com/cucumber/messages-go/v16"
	"gotest.tools/v3/assert"
	"os"
	"path/filepath"
	"testing"
)

type ScenarioContext struct {
	stepDefs []*stepDef
	t        *testing.T
}

func Run(t *testing.T, opts Options, setupScenario func(t *testing.T, ctx *ScenarioContext)) {
	r := &runner{
		opts:          opts,
		setupScenario: setupScenario,
		topLevelT:     t,
		incr:          &messages.Incrementing{},
	}
	r.run()
}

type runner struct {
	opts          Options
	setupScenario func(*testing.T, *ScenarioContext)
	topLevelT     *testing.T
	incr          *messages.Incrementing
}

type Options struct {
	Paths    []string
	Parallel bool
}

func (r *runner) run() {
	paths := r.opts.Paths
	if len(paths) == 0 {
		paths = []string{"features/*.feature"}
	}

	for _, path := range paths {
		files, err := filepath.Glob(path)
		assert.NilError(r.topLevelT, err)

		for _, file := range files {
			r.topLevelT.Run(file, func(t *testing.T) {
				if r.opts.Parallel {
					t.Parallel()
				}

				f, err := os.Open(file)
				assert.NilError(t, err)
				defer func() {
					err := f.Close()
					if err != nil {
						panic(err)
					}
				}()

				doc, err := gherkin.ParseGherkinDocument(f, r.incr.NewId)
				assert.NilError(t, err)
				r.runDoc(t, doc)
			})
		}
	}
}
