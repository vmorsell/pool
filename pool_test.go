package pool

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccessful(t *testing.T) {
	jobsDone := 0
	job := func() error {
		jobsDone++
		return nil
	}

	pool := New(2)
	err := pool.Run(job, job, job, job)
	require.Nil(t, err)
	require.Equal(t, jobsDone, 4)
}

func TestWithErrors(t *testing.T) {
	err1 := fmt.Errorf("err1")

	jobsDone := 0
	jobsWithErrors := 0
	jobOK := func() error {
		jobsDone++
		return nil
	}
	jobWithError := func() error {
		jobsWithErrors++
		return err1
	}

	pool := New(2)
	err := pool.Run(jobOK, jobOK, jobWithError, jobWithError)
	require.Equal(t, 2, jobsDone)
	require.Equal(t, 2, jobsWithErrors)
	require.Equal(t, err1, err)
}

func TestMoreWorkersThanJobs(t *testing.T) {
	jobsDone := 0
	pool := New(100)
	err := pool.Run(func() error {
		jobsDone++
		return nil
	})
	require.Equal(t, 1, jobsDone)
	require.Nil(t, err)
}
