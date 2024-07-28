package bot

import (
	"errors"
	"os"
	"strconv"
)

type Environment struct {
	LOOP_INTERVAL int
	DISCORD_WEBHOOK string
}

func GetEnv() (Environment, error) {
	env := Environment{}

	loopInterval, err := strconv.Atoi(os.Getenv("LOOP_INTERVAL"))
	if err != nil {
			return env, errors.New("LOOP_INTERVAL is not set")
	}
	env.LOOP_INTERVAL = loopInterval

	discordWebhook := os.Getenv("DISCORD_WEBHOOK")
	if discordWebhook == "" {
			return env, errors.New("DISCORD_WEBHOOK is not set")
	}
	env.DISCORD_WEBHOOK = discordWebhook

	return env, nil
}