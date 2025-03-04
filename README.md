# Rummy card game 

Learning websocket programming with Go. 


# Game Rules

## Table of Contents
1. [Introduction](#introduction)
2. [Game Components](#game-components)
3. [Setup](#setup)
4. [Gameplay](#gameplay)
5. [Valid Melds](#valid-melds)
6. [Scoring](#scoring)
7. [Special Rules](#special-rules)
8. [Winning](#winning)

---

## Introduction
Rummy is a card game in which players form valid combinations of cards called **melds**. The objective is to arrange all your cards into sets and sequences before your opponents.

### Game Components
- **Deck:** The game is played with one or two standard 52-card decks, including Jokers.
- **Players:** 2-4 players can participate.

### Setup
1. Each player is dealt **13** cards.
2. The remaining cards form the **draw pile**.
3. The top card is placed face up as the **discard pile**.

### Gameplay
Each turn consists of (order required):
1. **Drawing a card**
2. **Forming melds** (if possible)
3. **Discarding a card**

## Valid Melds
- **Sequences:** Three or more consecutive cards of the same suit.  
#### Valid
![Ex](https://raw.githubusercontent.com/Ojkee/rummy-card-game/master/doc/imgs/seqJQK.png)
![ExAce](https://raw.githubusercontent.com/Ojkee/rummy-card-game/master/doc/imgs/seqQKA.png)
![ExAce2](https://raw.githubusercontent.com/Ojkee/rummy-card-game/master/doc/imgs/seqA23.png)
#### Wrong
![ExWrong](https://raw.githubusercontent.com/Ojkee/rummy-card-game/master/doc/imgs/seqWrongAscend.png)
- **Sets:** Three or four cards of the same rank but different suits.  
#### Valid
![ExSuit](https://raw.githubusercontent.com/Ojkee/rummy-card-game/master/doc/imgs/seqSuit.png)
![ExSuitJok](https://raw.githubusercontent.com/Ojkee/rummy-card-game/master/doc/imgs/seqSuitJok.png)
![ExSuitFull](https://raw.githubusercontent.com/Ojkee/rummy-card-game/master/doc/imgs/seqSuitFull.png)
![ExSuitFullJok](https://raw.githubusercontent.com/Ojkee/rummy-card-game/master/doc/imgs/seqSuitFullJok.png)
#### Wrong
![ExWrongSuit](https://raw.githubusercontent.com/Ojkee/rummy-card-game/master/doc/imgs/seqWrongSuitFull.png)
![ExWrongSuit](https://raw.githubusercontent.com/Ojkee/rummy-card-game/master/doc/imgs/seqWrongSuit.png)
- **Jokers:** Can substitute any missing card in a meld.
![ExJok](https://raw.githubusercontent.com/Ojkee/rummy-card-game/master/doc/imgs/seq10JJok.png)
![ExJok2](https://raw.githubusercontent.com/Ojkee/rummy-card-game/master/doc/imgs/seqJok10J.png)

- **Placement** Its not required to set cards in hand into sequences to lock and meld, but if joker appears to be first in ascending sequence, it is considered to be lowest cards
![ExOrder](https://raw.githubusercontent.com/Ojkee/rummy-card-game/master/doc/imgs/seqDistantShuffle.png)

### Scoring
- **Face cards (K, Q, J):** 10 points each.
- **Number cards:** Face value.
- **Ace:** 1 or 11 points.
- **Joker:** 0 points.

### Winning
- The first player to arrange all their cards into valid melds wins.

---

*Have fun!* 🃏
