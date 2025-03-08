package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	dm "rummy-card-game/src/game_logic/deck_manager"
)

type CardModel struct {
	srcCard    *dm.Card
	rect       rl.Rectangle
	innerRect  rl.Rectangle
	isSelected bool
	sequenceId int
	isSmall    bool
}

func NewCardModel(card *dm.Card, rect rl.Rectangle, isSmall bool) *CardModel {
	gap := float32(CARD_GAP)
	innerWidth := float32(CARD_INNER_WIDTH)
	innerHeight := float32(CARD_INNER_HEIGHT)
	if isSmall {
		gap /= 2
		innerWidth = SEQUENCE_INNER_CARD_WIDTH
		innerHeight = SEQUENCE_INNER_CARD_HEIGHT
	}
	_innerRect := rl.NewRectangle(
		rect.X+gap,
		rect.Y+gap,
		innerWidth,
		innerHeight,
	)
	return &CardModel{
		srcCard:    card,
		rect:       rect,
		innerRect:  _innerRect,
		isSelected: false,
		sequenceId: -1,
		isSmall:    isSmall,
	}
}

func (card *CardModel) SetSrcCard(srcCard *dm.Card) {
	card.srcCard = srcCard
}

func (card *CardModel) SetSequenceId(sequenceId int) {
	card.sequenceId = sequenceId
}

func (card *CardModel) Reset() {
	card.isSelected = false
	card.sequenceId = -1
}

func (card *CardModel) IsSelectedSequence() bool {
	return card.sequenceId == -1
}

func (card *CardModel) MoveX(x float32) {
	card.innerRect.X = x + float32(CARD_GAP)
	card.rect.X = x
}

func (card *CardModel) Draw() {
	var selectedOffset float32 = 0
	if card.isSelected {
		selectedOffset = -CARD_SELECTED_OFFSET
	}
	card.drawFrame(selectedOffset)
	if card.srcCard == nil {
		return
	}
	card.drawSuitTexture(selectedOffset)
	card.drawRank(selectedOffset)
	if card.sequenceId != -1 {
		card.drawLockSequence(selectedOffset)
	}
}

func (card *CardModel) drawFrame(selectedOffset float32) {
	rl.DrawRectangle(
		card.rect.ToInt32().X,
		card.rect.ToInt32().Y+int32(selectedOffset),
		card.rect.ToInt32().Width,
		card.rect.ToInt32().Height,
		COLOR_BEIGE,
	)

	innerColor := COLOR_TAUPE
	if card.isSelected {
		innerColor = COLOR_WALNUT_BROWN
	}
	rl.DrawRectangleRec(
		rl.NewRectangle(
			card.innerRect.X,
			card.innerRect.Y+selectedOffset,
			card.innerRect.Width,
			card.innerRect.Height,
		),
		innerColor,
	)
}

func (card *CardModel) drawSuitTexture(selectedOffset float32) {
	var rotation float32 = 0
	var scale float32 = 1
	img := RANK_IMGS[card.srcCard.Suit]
	if card.isSmall {
		img = RANK_IMGS_SMALL[card.srcCard.Suit]
	}
	rl.DrawTextureEx(
		img,
		rl.NewVector2(
			card.rect.X+float32(CARD_GAP),
			card.rect.Y+float32(CARD_GAP)+selectedOffset,
		),
		rotation,
		scale,
		COLOR_DARK_GRAY,
	)
}

func (card *CardModel) drawRank(selectedOffset float32) {
	rankString := card.srcCard.Rank.String()
	textVec := GetTextVec(rankString)
	rl.DrawTextEx(
		FONT,
		rankString,
		rl.NewVector2(
			card.rect.X+(card.rect.Width-textVec.X)/2,
			card.rect.Y+(card.rect.Height-textVec.Y)+selectedOffset,
		),
		FONT_SIZE,
		FONT_SPACING,
		COLOR_BEIGE,
	)
}

func (card *CardModel) drawLockSequence(selectedOffset float32) {
	rl.DrawRectangle(
		card.rect.ToInt32().X,
		card.rect.ToInt32().Y+int32(selectedOffset),
		card.rect.ToInt32().Width,
		4,
		LOCK_COLORS[card.sequenceId],
	)
}

func (card *CardModel) InRect(mousePos rl.Vector2) bool {
	return rl.CheckCollisionPointRec(mousePos, card.rect)
}
