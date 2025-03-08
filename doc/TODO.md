# TODO

## Features
### Important
- [x] Placing sequences on the screen
- [x] Proper turn handling (draw -> action -> discard)
- [x] Player should be able to rearrange cards
- [x] Initial meld (51 points and at least one pure sequence)
- [x] Drawing a card from the discard pile
- [x] Replacing a joker with the proper card
- [x] Adding cards to existing sequences
- [x] Handling reshuffle when the draw pile gets empty 
- [x] `Enter ip` window before connecting

### Further work
- [ ] Other players' number of cards should be visible
- [ ] Dragging a sequence into another to get a joker/merge
- [ ] AI Bots
- [ ] Refactor functions in table file into game_manager file

### Debug mode
- [x] Meld-ready draw
- [x] Reset table client call

### Extra features
- [x] Wrong card sequence candidate UI
- [x] Highlight locked cards that don't construct a sequence
- [x] Copy-paste ip mechanism 
- [ ] Player can enter a nickname 
- [ ] Points counter
- [ ] Turn and connection info
- [ ] End game screen
- [ ] Continue playing for 2nd and 3rd places
- [ ] Fast emulation using 4 bots with data gathering
- [ ] Reset table server call
- [ ] Chat
- [ ] Nicer UI

## Bugs
### FIX ASAP
- [x] Discarding a card with an exact copy discards both cards
- [x] Meld button doesn't disappear after meld
- [x] Joker appends to the end of the sequence after an Ace
- [x] After clicking meld when locked sequences don't meet requirements, it resets lock color and messes locks
- [x] Replacing needs better refresh of `jokerImitation` when appending the first card
- [x] Player shouldn't append the last card in hand to a table sequence
- [x] Available spot rectangle should update after appending a new card to the sequence
- [x] Can't append a joker at the beginning
- [x] Append at the beginning/end if the first/last card is a joker
- [x] Most left joker placement should appear as lowest ranked card
- [x] Same rank sequence on table shows available spots for wrong cards and couses server crash upon adding card
- [x] Can't replace joker in same rank sequence

### FIX PLEASE
- [x] Ascending sequence ACE -> TWO -> THREE doesn't work
- [x] Rearranging sequences on hand within other player turn shouldn't show `not your turn` message
- [ ] Rearranging sequences or drawing a card shouldn't unlock sequences in hand
