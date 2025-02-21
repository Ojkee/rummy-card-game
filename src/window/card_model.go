package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	dm "rummy-card-game/src/game_logic/deck_manager"
)

type CardModel struct {
	srcCard    *dm.Card
	isSelected bool
	rect       rl.Rectangle
}

func NewCardModel(card *dm.Card, rect rl.Rectangle) *CardModel {
	return &CardModel{
		srcCard:    card,
		isSelected: false,
		rect:       rect,
	}
}

func (card *CardModel) SetSrcCard(srcCard *dm.Card) {
	card.srcCard = srcCard
}

func (card *CardModel) Draw() {
	var selectedOffset float32 = 0
	if card.isSelected {
		selectedOffset = -20
	}
	card.drawFrame(selectedOffset)
	if card.srcCard == nil {
		return
	}
	card.drawSuitTexture(selectedOffset)
	card.drawRank(selectedOffset)
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
			card.rect.X+float32(CARD_GAP),
			float32(card.rect.ToInt32().Y+CARD_GAP)+selectedOffset,
			float32(CARD_INNER_WIDTH),
			float32(CARD_INNER_HEIGHT),
		),
		innerColor,
	)
}

func (card *CardModel) drawSuitTexture(selectedOffset float32) {
	var rotation float32 = 0
	var scale float32 = 1
	rl.DrawTextureEx(
		RANK_IMGS[card.srcCard.Suit],
		rl.NewVector2(card.rect.X+float32(CARD_GAP), card.rect.Y+float32(CARD_GAP)+selectedOffset),
		rotation,
		scale,
		COLOR_DARK_GRAY,
	)
}

func (card *CardModel) drawRank(selectedOffset float32) {
	randString := card.srcCard.Rank.String()
	textVec := GetTextVec(randString)
	rl.DrawTextEx(
		FONT,
		randString,
		rl.NewVector2(
			card.rect.X+float32(SUIT_WIDTH-int32(textVec.X))/2,
			card.rect.Y+float32(SUIT_HEIGHT*3/2)+selectedOffset,
		),
		float32(FONT_SIZE),
		FONT_SPACING,
		COLOR_BEIGE,
	)
}

func (card *CardModel) IsClicked(mousePos rl.Vector2) bool {
	return rl.CheckCollisionPointRec(mousePos, card.rect)
}
