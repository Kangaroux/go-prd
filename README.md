# go-prd

This is a simple implementation for approximating the C-value used in [Pseudo Random Distribution](https://dota2.fandom.com/wiki/Random_Distribution) (PRD).

The script is designed to generate a dataset of C-values. It runs simulations for each to approximate the average probability. Simulations are run in parallel and the precision of the C-value is controlled to keep the run time reasonable.

## Usage

Running the script outputs a CSV with 3 columns: the C-value, the probability, and the expected value. There are some constants at the top of the file that can be adjusted.

You can simply run the script with:

```sh
go run .
```

Or, probably more useful, dump the results to a CSV:

```sh
go run . > out.csv
```

## A Brief Overview of PRD

PRD seeks to address the problem of randomness in multiplayer games where a player can be very lucky or very unlucky. This can lead to an overwhelming advantage in some cases, and is generally not fun for the other players. PRD was first introduced in Warcraft 3 and later adopted in Dota 2.

Unlike uniform distributions, where each coin flip or dice roll is independent, PRD uses dependent randomness. Every time an event doesn't occur, the probability increases linearly. The probability of the first attempt is low, and increases until the probability eventually exceeds 100%, which forces the event to occur.

One important feature of PRD is that -- averaged over a large sample set -- the probability of an event occuring is the same as a uniform distribution. The differences are that:

1. Lucky streaks are unlikely due to the initial probability being lower
2. Unlucky streaks are reduced/impossible. An event _must_ occur after a certain number of attempts. The probability of it occuring is also skewed towards the lower end of the distribution

Consider an event with a 25% chance of occuring. In both distributions the expected number of attempts is 4. However, with PRD, the upper limit for attempts is 11, whereas a uniform distribution has no upper limit -- the event could simply never occur.

![Bar graph comparing a PRD and a uniform distribution.](https://static.wikia.nocookie.net/dota2_gamepedia/images/8/8b/AttacksUntilNextProc25.jpg/revision/latest?cb=20130505045408)

Image source: [Dota2 Wiki](https://dota2.fandom.com/wiki/Random_Distribution)
