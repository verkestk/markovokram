/*
Package markovokram is a simple string Markov chain building package

It builds forwards and backwards chains from text and can generate both
text both forwards and backwards.
*/
package markovokram

import (
	"math/rand"
	"strings"
)

// Chain contains the forwards and backwards markov chains. A Markov chain is a
// map of prefixes to suffixes. A prefix can be one or multiple tokens. While a
// prefix can map to multiple suffixes, a single suffix is always a single
// token.
type Chain struct {
	forwards     map[string][]string
	backwards    map[string][]string
	prefixLength int
}

// prefix is a Markov chain prefix of one or more words. Convenience wrapper for
// string slice that provides a couple of reusable methods.
type chainPrefix []string

// string returns the Prefix as a string (for use as a map key).
func (p chainPrefix) string() string {
	return strings.Join(p, " ")
}

// shift removes the first word from the Prefix and appends the given word.
func (p chainPrefix) shift(word string) {
	if len(p) == 1 {
		p[0] = word
	} else {
		copy(p, p[1:])
		p[len(p)-1] = word
	}
}

// Generation keeps track of a specific prefix and allows the consumer to
// continue generating text but randomly selecing a suffix and shifting the
// prefix
type Generation struct {
	chainMap map[string][]string // from *Chain either forwards or backwards
	prefix   chainPrefix
}

// NewChain returns a new Chain with prefixes of prefixLength words.
func NewChain(prefixLength int) *Chain {
	return &Chain{
		forwards:     make(map[string][]string),
		backwards:    make(map[string][]string),
		prefixLength: prefixLength,
	}
}

// Build reads text and parses it into prefixes and suffixes that are stored in
// Chain.
func (c *Chain) Build(tokens []string) {
	forwardsPrefix := make(chainPrefix, c.prefixLength)
	backwardsPrefix := make(chainPrefix, c.prefixLength)
	for i := range tokens {
		// build forwards chain
		forwardsKey := forwardsPrefix.string()
		c.forwards[forwardsKey] = append(c.forwards[forwardsKey], tokens[i])
		forwardsPrefix.shift(tokens[i])

		// build backwards chain
		backwardsKey := backwardsPrefix.string()
		c.backwards[backwardsKey] = append(c.backwards[backwardsKey], tokens[len(tokens)-i-1])
		backwardsPrefix.shift(tokens[len(tokens)-i-1])
	}
}

// GenerateForward generates fowards text based on empty prefix
func (c *Chain) GenerateForward() *Generation {
	return &Generation{chainMap: c.forwards, prefix: make(chainPrefix, c.prefixLength)}
}

// GenerateForwardFromPrefix generates forwads text based on a specified prefix.
// If the specified prefix is longer that the chain's prefixLength, the items
// are removed from the front of the prefix to accomodate. If the specified
// prefix is shorter than the chain's prefixLength, the beginning of the prefix
// is padded.
func (c *Chain) GenerateForwardFromPrefix(prefix []string) *Generation {
	if len(prefix) > c.prefixLength {
		return &Generation{chainMap: c.forwards, prefix: chainPrefix(prefix[len(prefix)-c.prefixLength:])}
	}

	if len(prefix) < c.prefixLength {
		padding := make([]string, c.prefixLength-len(prefix))
		return &Generation{chainMap: c.forwards, prefix: chainPrefix(append(padding, prefix...))}
	}

	return &Generation{chainMap: c.forwards, prefix: chainPrefix(prefix)}
}

// GenerateBackward generates backwards text based on empty prefix.
func (c *Chain) GenerateBackward() *Generation {
	return &Generation{chainMap: c.backwards, prefix: make(chainPrefix, c.prefixLength)}
}

// GenerateBackwardFromPrefix generates forwads text based on a specified
// prefix. If the specified prefix is longer that the chain's prefixLength, the
// items are removed from the back of the prefix to accomodate. If the specified
// prefix is shorter than the chain's prefixLength, the end of the prefix is
// padded.
func (c *Chain) GenerateBackwardFromPrefix(prefix []string) *Generation {
	if len(prefix) > c.prefixLength {
		return &Generation{chainMap: c.backwards, prefix: chainPrefix(prefix[len(prefix)-c.prefixLength:])}
	}

	if len(prefix) < c.prefixLength {
		padding := make([]string, c.prefixLength-len(prefix))
		return &Generation{chainMap: c.backwards, prefix: chainPrefix(append(padding, prefix...))}
	}

	return &Generation{chainMap: c.backwards, prefix: chainPrefix(prefix)}
}

// Next generates a new token for the sequence.
func (g *Generation) Next() string {
	suffixes := g.chainMap[g.prefix.string()]
	if len(suffixes) == 0 {
		return ""
	}

	next := suffixes[rand.Intn(len(suffixes))]
	g.prefix.shift(next)
	return next
}

// NextWith Uses the specified token to shift tthe chain.
func (g *Generation) NextWith(str string) {
	g.prefix.shift(str)
}

// Options returns all of the possible next tokens.
func (g *Generation) Options() []string {
	options := make([]string, len(g.chainMap[g.prefix.string()]))
	copy(options, g.chainMap[g.prefix.string()])
	return options
}
