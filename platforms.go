package crossbot

import "strconv"

type Platform uint8

const (
	PlatformUndefined Platform = iota
	PlatformDiscord
	PlatformTelegram
)

func GetStringPlatform(s string) (Platform, error) {
	res, err := strconv.Atoi(s)
	if err != nil {
		return PlatformUndefined, err
	}

	return Platform(res), nil
}
