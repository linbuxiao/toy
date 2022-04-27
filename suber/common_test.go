package suber_test

import (
	"github.com/linbuxiao/toy/suber"
	"github.com/stretchr/testify/require"
	"testing"
)

func checkContents(t *testing.T, c chan interface{}, val []string) {
	var result []string
	for k := range c {
		result = append(result, k.(string))
	}
	require.NotNil(t, result)
	require.ElementsMatch(t, val, result)
}

func TestSuber(t *testing.T) {
	s := suber.New()
	c := s.Sub("test")
	s.Pub("test msg", "test")
	// you must shut down all channel initiative. Unless you will be blocked in loop.
	s.Shutdown()
	checkContents(t, c, []string{"test msg"})
}
