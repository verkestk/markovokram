# markovokram
golang markov chain and text generation - builds forwards and backwards markov chains from input text

# Usage

Before you generate text, don't forget to:
```
rand.Seed(time.Now().UnixNano())
```

## Creat a chain

```
chain := markovokram.NewChain(2)
```

The parameter here is prefix length. The shorter the prefix, the more variety you'll see in the generated text. The longer the prefix, the more "natural" the generated language will be. Be careful - a long prefix combined with a short corpus will leave you with little to no variety - it will just parrot examples word-for-word, from your corpus.

For very small corpora, you may need to use a chain length of 1. For medium sized corpora, a chain length of 2 is likely best. Only for very large corpora would you find good results with a value of 3 or more.

## Build the chain

```
sentence1 := "What noise annoys a noisy oyster?"
sentence2 := "A noisy noise annoys a noisy oyster."
chain.Build(strings.Fields(sentence1))
chain.Build(strings.Fields(sentence2))
```

Pass a slice of strings to the `Build` method. It's up to you how you divide up your corpus. By paragraph or sentence perhaps. If you pass your entire corpus in one invocation of `Build`, and then you consequently generate text based on no specific prefix, then your generated text will always start with the same word. It's a good idea to break up your corpus.

## Generate text

You have 4 options for generating text:
```
generation := chain.GenerateForward()
```
```
generation := chain.GenerateForwardFromPrefix([]string{...})
```
```
generation := chain.GenerateBackward()
```
```
generation := chain.GenerateBackwardFromPrefix([]string{...})
```

Forward generation is the most typical. If you want to start a random sentence, use `GenerateForward`. If you want text generated from a specific prefix, use `GenerateForwardFromPrefix`. You can pass a prefix of any length - however if the prefix you pass is longer that the chain's prefix length, some of your prefix will be ignored.

Backward generation is useful for some specific appliciations - for example poetry generation. If you need a statement that ends in a specific word for rhyming purposes, you can use `GenerateBackwardFromPrefix`, passing the single word as the prefix. `GenerateBackward` is also available.

Once you've created a `generation`, then you call:

```
next := generation.Next()
```

This will select a random next word in the sequence. This method will return an empty string if there are no available next words.

## Examples

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

### Forward generation from a prefix

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

### Basic backward generation

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

### Basic backward generation from a prefix

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
