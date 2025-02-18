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
	card.drawFrame()
	if card.srcCard == nil {
		return
	}
	card.drawSuitTexture()
	card.drawRank()
}

func (card *CardModel) drawFrame() {
	rl.DrawRectangleRec(
		card.rect,
		COLOR_BEIGE,
	)

	innerColor := COLOR_TAUPE
	if card.isSelected {
		innerColor = COLOR_WALNUT_BROWN
	}
	var selectedOffset float32 = 0
	if card.isSelected {
		selectedOffset = -20
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

func (card *CardModel) drawSuitTexture() {
	var rotation float32 = 0
	var scale float32 = 1
	rl.DrawTextureEx(
		RANK_IMGS[card.srcCard.Suit],
		rl.NewVector2(card.rect.X+float32(CARD_GAP), card.rect.Y+float32(CARD_GAP)),
		rotation,
		scale,
		COLOR_DARK_GRAY,
	)
}

func (card *CardModel) drawRank() {
	randString := card.srcCard.Rank.String()
	textVec := GetTextVec(randString)
	rl.DrawTextEx(
		FONT,
		randString,
		rl.NewVector2(
			card.rect.X+float32(SUIT_WIDTH-int32(textVec.X))/2,
			float32(CARD_POS_Y+SUIT_HEIGHT*3/2),
		),
		float32(FONT_SIZE),
		FONT_SPACING,
		COLOR_BEIGE,
	)
}

func (card *CardModel) IsClicked(mousePos rl.Vector2) bool {
	return rl.CheckCollisionPointRec(mousePos, card.rect)
}
