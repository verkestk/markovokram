# markovokram
golang markov chain and text generation - builds forwards and backwards markov chains from input text

# Usage

Before you generate text, don't forget to:
```
rand.Seed(time.Now().UnixNano())
```

### Basic forward generation
`GenerateForward` will always start with an empty prefix. That means its first generated word will always be the first token from one of the calls to `Build`.

```
sentence1 := "What noise annoys a noisy oyster?"
sentence2 := "A noisy noise annoys a noisy oyster."

chain := markovokram.NewChain(1)
chain.Build(strings.Fields(sentence1))
chain.Build(strings.Fields(sentence2))

for i := 0; i < 10; i++ {
  generation := chain.GenerateForward()
  words := []string{}
  next := generation.Next()
  for next != "" {
    words = append(words, next)
    next = generation.Next()
  }

  fmt.Println(strings.Join(words, " "))
}
```

Example Output
```
What noise annoys a noisy oyster.
What noise annoys a noisy oyster.
A noisy oyster.
What noise annoys a noisy oyster?
A noisy oyster?
A noisy oyster?
What noise annoys a noisy noise annoys a noisy noise annoys a noisy noise annoys a noisy oyster?
A noisy noise annoys a noisy oyster.
A noisy oyster?
What noise annoys a noisy oyster.
```

# Forward generation from a prefix

Swap out the `GenerateForward` call for:
```
generation := chain.GenerateForwardFromPrefix([]string{"annoys"})
```

Example Output
```
a noisy oyster?
a noisy oyster.
a noisy oyster.
a noisy oyster?
a noisy noise annoys a noisy oyster?
a noisy noise annoys a noisy noise annoys a noisy oyster.
a noisy oyster.
a noisy noise annoys a noisy oyster.
a noisy oyster?
a noisy oyster.
```

# Basic backward generation

`GenerateBackward` will always start with an empty prefix. That means its first generated word will always be the final token from one of the calls to `Build`.

```
sentence1 := "What noise annoys a noisy oyster?"
sentence2 := "A noisy noise annoys a noisy oyster."

chain := markovokram.NewChain(1)
chain.Build(strings.Fields(sentence1))
chain.Build(strings.Fields(sentence2))

for i := 0; i < 10; i++ {
  generation := chain.GenerateBackward()
  words := []string{}
  next := generation.Next()
  for next != "" {
    words = append([]string{next}, words...)
    next = generation.Next()
  }

  fmt.Println(strings.Join(words, " "))
}
```

Example Output
```
What noise annoys a noisy oyster?
What noise annoys a noisy noise annoys a noisy oyster?
A noisy oyster.
What noise annoys a noisy oyster?
A noisy oyster?
What noise annoys a noisy oyster.
What noise annoys a noisy oyster.
A noisy oyster?
A noisy oyster?
What noise annoys a noisy oyster?
```

# Basic backward generation from a prefix

Swap out the `GenerateBackward` call above with:

```
generation := chain.GenerateBackwardFromPrefix([]string{"annoys"})
```

Example Output
```
A noisy noise
What noise
What noise
What noise annoys a noisy noise annoys a noisy noise annoys a noisy noise
What noise annoys a noisy noise annoys a noisy noise
What noise
What noise
A noisy noise
What noise
A noisy noise
```


# Thank you

Big thanks to the authors of [https://golang.org/doc/codewalk/markov/](https://golang.org/doc/codewalk/markov/). This implementation helped me understand how Markov chains work and how they could be implemented in golang.
