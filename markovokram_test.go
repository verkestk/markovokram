package markovokram

import (
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Test_chainPrefix_string(t *testing.T) {
	prefix := []string{"what", "noise", "annoys"}
	expected := "what noise annoys"
	actual := chainPrefix(prefix).string()
	if expected != actual {
		t.Logf("prefix: %v\n", prefix)
		t.Errorf("expected string() \"%s\", got \"%s\"", expected, actual)
	}
}

func Test_chainPrefix_shift(t *testing.T) {
	prefix := []string{"what"}
	expected := []string{"noise"}
	chainPrefix(prefix).shift("noise")
	if !reflect.DeepEqual(expected, prefix) {
		t.Logf("expected: %v", expected)
		t.Logf("actual: %v", prefix)
		t.Errorf("unexpected chainPrefix values")
	}

	prefix = []string{"what", "noise"}
	expected = []string{"noise", "annoys"}
	chainPrefix(prefix).shift("annoys")
	if !reflect.DeepEqual(expected, prefix) {
		t.Logf("expected: %v", expected)
		t.Logf("actual: %v", prefix)
		t.Errorf("unexpected chainPrefix values")
	}
}

func Test_NewChain(t *testing.T) {
	chain := NewChain(1)
	if chain.forwards == nil {
		t.Errorf("new chain has nil forwards map")
	}
	if chain.backwards == nil {
		t.Errorf("new chain has nil backwards map")
	}
	if chain.prefixLength != 1 {
		t.Errorf("expected prefixLength for new chain is %d, gott %d", 1, chain.prefixLength)
	}
}

func Test_Chain_Build(t *testing.T) {
	sentence1 := "What noise annoys a noisy oyster?"
	sentence2 := "A noisy noise annoys a noisy oyster."

	chain := NewChain(1)
	chain.Build(strings.Fields(sentence1))
	chain.Build(strings.Fields(sentence2))
	if len(chain.forwards) != 7 {
		t.Errorf("expected forwards map length 7, got %d", len(chain.forwards))
	}
	if len(chain.backwards) != 7 {
		t.Errorf("expected backwards map length 7, got %d", len(chain.forwards))
	}
}

func Test_Chain_GenerateForward(t *testing.T) {
	sentence1 := "What noise annoys a noisy oyster?"
	sentence2 := "A noisy noise annoys a noisy oyster."

	chain := NewChain(1)
	chain.Build(strings.Fields(sentence1))
	chain.Build(strings.Fields(sentence2))
	generation := chain.GenerateForward()
	next := generation.Next()
	if next != "What" && next != "A" {
		t.Errorf("expected \"What\" or \"A\", got \"%s\"", next)
	}
}

func Test_Chain_GenerateForwardFromPrefix(t *testing.T) {
	sentence1 := "What noise annoys a noisy oyster?"
	sentence2 := "A noisy noise annoys a noisy oyster."

	chain := NewChain(1)
	chain.Build(strings.Fields(sentence1))
	chain.Build(strings.Fields(sentence2))
	generation := chain.GenerateForwardFromPrefix([]string{"noisy"})
	next := generation.Next()
	if next != "noise" && next != "oyster?" && next != "oyster." {
		t.Errorf("expected \"noise\" or \"oyster?\" or \"oyster.\", got \"%s\"", next)
	}

	generation = chain.GenerateForwardFromPrefix([]string{"a", "noisy"})
	next = generation.Next()
	if next != "noise" && next != "oyster?" && next != "oyster." {
		t.Errorf("expected \"noise\" or \"oyster?\" or \"oyster.\", got \"%s\"", next)
	}

	generation = chain.GenerateForwardFromPrefix([]string{})
	next = generation.Next()
	if next == "" {
		t.Errorf("expected non empty string, got \"\"")
	}
}

func Test_Chain_GenerateBackward(t *testing.T) {
	sentence1 := "What noise annoys a noisy oyster?"
	sentence2 := "A noisy noise annoys a noisy oyster."

	chain := NewChain(1)
	chain.Build(strings.Fields(sentence1))
	chain.Build(strings.Fields(sentence2))
	generation := chain.GenerateBackward()
	next := generation.Next()
	if next != "oyster?" && next != "oyster." {
		t.Errorf("expected \"oyster?\" or \"oyster.\", got \"%s\"", next)
	}
}

func Test_Chain_GenerateBackwardFromPrefix(t *testing.T) {
	sentence1 := "What noise annoys a noisy oyster?"
	sentence2 := "A noisy noise annoys a noisy oyster."

	chain := NewChain(1)
	chain.Build(strings.Fields(sentence1))
	chain.Build(strings.Fields(sentence2))
	generation := chain.GenerateBackwardFromPrefix([]string{"noise"})
	next := generation.Next()
	if next != "What" && next != "noisy" {
		t.Errorf("expected \"What\" or \"noisy\", got \"%s\"", next)
	}

	generation = chain.GenerateBackwardFromPrefix([]string{"oyster.", "noisy"})
	next = generation.Next()
	if next != "a" && next != "A" {
		t.Errorf("expected \"a\" or \"A\", got \"%s\"", next)
	}

	generation = chain.GenerateBackwardFromPrefix([]string{})
	next = generation.Next()
	if next == "" {
		t.Errorf("expected non empty string, got \"\"")
	}

}

func Test_Generation_Next(t *testing.T) {
	// mostly tested in the Test_Chain_Generate* functions above

	// what about the end of the chain?
	sentence1 := "What noise annoys a noisy oyster?"
	sentence2 := "A noisy noise annoys a noisy oyster."

	chain := NewChain(1)
	chain.Build(strings.Fields(sentence1))
	chain.Build(strings.Fields(sentence2))
	generation := chain.GenerateForwardFromPrefix([]string{"oyster."})
	next := generation.Next()

	if next != "" {
		t.Errorf("expected empty next string, got \"%s\"", next)
	}
}

func Test_Generation_NextWith(t *testing.T) {
	sentence1 := "What noise annoys a noisy oyster?"
	sentence2 := "A noisy noise annoys a noisy oyster."

	// use a valid "next" value that
	chain := NewChain(1)
	chain.Build(strings.Fields(sentence1))
	chain.Build(strings.Fields(sentence2))
	generation := chain.GenerateForwardFromPrefix([]string{"What"})
	generation.NextWith("noise")
	next := generation.Next()
	expected := "annoys"
	if next != expected {
		t.Errorf("expected next token \"%s\" - got \"%s\"", expected, next)
	}

	// use an invalid "next" value that should still work
	// this should result in a generated sequence that couldn't have been
	// generated "naturally"
	generation = chain.GenerateForwardFromPrefix([]string{"What"})
	generation.NextWith("noisy")
	next = generation.Next()
	if next != "oyster?" && next != "noise" && next != "oyster." {
		t.Errorf("expected next token \"oyster?\" or \"noise\" or \"oyster.\" - got \"%s\"", next)
	}

	// use an invalid "next" value that should force the end of the sequence
	generation = chain.GenerateForwardFromPrefix([]string{"What"})
	generation.NextWith("impossible")
	next = generation.Next()
	if next != "" {
		t.Errorf("expected empty next token - got \"%s\"", next)
	}
}

func Test_Generation_Options(t *testing.T) {
	sentence1 := "What noise annoys a noisy oyster?"
	sentence2 := "A noisy noise annoys a noisy oyster."

	chain := NewChain(1)
	chain.Build(strings.Fields(sentence1))
	chain.Build(strings.Fields(sentence2))
	generation := chain.GenerateForwardFromPrefix([]string{"noisy"})
	options := generation.Options()
	expectedOptions := []string{"oyster?", "noise", "oyster."}
	if len(options) != len(expectedOptions) {
		t.Errorf("expected %d options = got %d", len(expectedOptions), len(options))
	}
	if !reflect.DeepEqual(expectedOptions, options) {
		t.Logf("expected: %v", expectedOptions)
		t.Logf("actual: %v", options)
		t.Errorf("unexpected options")
	}
}
