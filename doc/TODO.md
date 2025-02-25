# TODO

## Features
### Important
- [x] Placing sequences on the screen
- [x] Propper turn handling (draw -> action -> discard)
- [x] Player should be able to rearange cards
- [x] Initial meld (51 points and at least one pure sequence)
- [x] Drawing card from discard pile
- [x] Replacing joker with propper card
- [x] Adding cards to existing sequences
- [ ] Dragging sequence into another to get Joker/Merge
- [ ] Handling reshuffle when draw pile gets empty

### Further work
- [ ] AI Bots
- [ ] Other players number of cards should be visible

### Debug mode
- [x] Meld ready draw 

### Extra features 
- [ ] Chat
- [ ] Nicer UI
- [ ] Fast emulation using 4 bots with data gathering
- [ ] Highlight locked cards that don't construct sequence
- [ ] Player can enter nickname
- [ ] End game screen

## Bugs
### FIX ASAP
- [x] Discarding card with exact copy discards both cards
- [x] Meld button doesn't disappear after meld
- [ ] Append begin/end if first/last is joker
- [ ] Joker appends to the end of sequence after Ace
- [ ] After clicking meld when lock seqs don't meet requirements restarts lock color and messes locks 
- [ ] Replacing needs better refresh of jokImitation when appending first card


### FIX PLEASE
- [ ] Ascending sequence ACE -> TWO -> TRHEE doesn't work
- [ ] Rearranging sequences or drawing card shouldn't unlock sequences on hand
