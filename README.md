# metro - terminal metronome

Terminal metronome, written in Go.
Prints status on 1 line, using filled/hollow emojis.
Uses Sonic Pi's collection of samples.

## Usage

```
-b int
    beats (default 4)
-f string
    file (default "samples/tabla_te2.flac")
-t int
    tempo (default 120)
```

## Planned features

+ "Big" mode - TUI cells
+ On-the-fly tempo changing
+ Multiple tempo modes - currently only resampling available
+ Strong and weak beats
+ Polyrhythms - e.g. `-b 3/2`
+ Embedding the samples/finding Sonic Pi samples on the system
