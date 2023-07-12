package pool

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	t.Run("no errors", func(t *testing.T) {
		ok := 0
		fn := func() error {
			ok++
			return nil
		}

		jobs := []Job{fn, fn, fn, fn, fn, fn}
		p := New(3)
		err := p.Run(jobs...)

		require.Nil(t, err)
		require.Equal(t, len(jobs), ok)
	})

	t.Run("with errors", func(t *testing.T) {
		processed := 0
		ok := 0
		wantErr := fmt.Errorf("err")

		fn := func() error {
			processed++
			ok++
			return nil
		}

		fnErr := func() error {
			processed++
			return wantErr
		}

		okJobs := []Job{fn, fn, fn, fn, fn}
		jobs := append(okJobs, fnErr)
		jobs = append(jobs, okJobs...)

		size := 2
		p := New(size)
		err := p.Run(jobs...)

		// The number of jobs processed differs depending on how many of the
		// error jobs the workers start processing before the first error
		// is returned.
		require.GreaterOrEqual(t, processed, len(okJobs)+1)
		require.Equal(t, processed-1, ok)
		require.Equal(t, err, wantErr)
	})
}
