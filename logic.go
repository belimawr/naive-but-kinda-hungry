package main

import (
	"context"

	"github.com/rs/zerolog"
)

// info returns the Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
func info(ctx context.Context) BattlesnakeInfoResponse {
	zerolog.Ctx(ctx).Info().Msg("Sending snake info")
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "Tiago Queiroz",
		Color:      "#b30000",
		Head:       "smart-caterpillar",
		Tail:       "round-bum",
	}
}

// start is called everytime the Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
func start(ctx context.Context, state GameState) {
	zerolog.Ctx(ctx).Info().Msg("Starting game!")
}

// end is called when a a game has ended.
func end(ctx context.Context, state GameState) {
	zerolog.Ctx(ctx).Info().Msg("The game has ended")
}

// mode is called on every turn and returns the Battlesnake's next move
func move(ctx context.Context, state GameState) BattlesnakeMoveResponse {
	logger := zerolog.Ctx(ctx)

	nextMove := hungry(ctx, state)

	logger.Info().Msgf("MOVE: %s", nextMove)
	return BattlesnakeMoveResponse{
		Move: nextMove,
	}
}
