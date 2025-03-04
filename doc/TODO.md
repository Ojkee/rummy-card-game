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
- [ ] Dragging a sequence into another to get a joker/merge  
- [ ] Handling reshuffle when the draw pile gets empty  

### Further work  
- [ ] AI Bots
- [ ] Other players' number of cards should be visible  

### Debug mode  
- [x] Meld-ready draw  
- [x] Reset table client call  

### Extra features  
- [x] Wrong card sequence candidate UI  
- [ ] Chat  
- [ ] Nicer UI  
- [ ] Fast emulation using 4 bots with data gathering  
- [ ] Highlight locked cards that don't construct a sequence  
- [ ] Player can enter a nickname  
- [ ] End game screen  
- [ ] Continue playing for 2nd (and so on) places  
- [ ] Reset table server call  

## Bugs  
### FIX ASAP  
- [x] Discarding a card with an exact copy discards both cards  
- [x] Meld button doesn't disappear after meld  
- [x] Joker appends to the end of the sequence after an Ace  
- [x] After clicking meld when locked sequences don't meet requirements, it resets lock color and messes locks  
- [x] Replacing needs better refresh of `jokerImitation` when appending the first card  
- [x] Player shouldn't append the last card in hand to a table sequence  
- [x] Available spot rectangle should update after appending a new card to the sequence  
- [ ] Can't append a joker at the beginning  
- [ ] Append at the beginning/end if the first/last card is a joker  

### FIX PLEASE  
- [ ] Ascending sequence ACE -> TWO -> THREE doesn't work
- [ ] Rearranging sequences or drawing a card shouldn't unlock sequences in hand


