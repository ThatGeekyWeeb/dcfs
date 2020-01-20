package main

import (
	"log"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/pkg/errors"
)

func (fs *Filesystem) UpdateGuilds() error {
	guilds, err := fs.State.Guilds()
	if err != nil {
		log.Fatalln("Failed to get guilds:", err)
	}

	newGuilds := guilds[:0]

	fs.mu.Lock()
	defer fs.mu.Unlock()

Main:
	for _, guild := range guilds {
		for _, g := range fs.Guilds {
			if g.ID == guild.ID {
				continue Main
			}
		}

		newGuilds = append(newGuilds, guild)
	}

	for _, g := range newGuilds {
		guild := &Guild{
			ID:    g.ID,
			FS:    fs,
			Inode: NewInode(),
		}

		if err := guild.UpdateChannels(); err != nil {
			return errors.Wrap(err, "Failed to update guild "+g.ID.String())
		}

		// Subscribe to guilds
		fs.State.Gateway.GuildSubscribe(gateway.GuildSubscribeData{
			GuildID:    g.ID,
			Typing:     true,
			Activities: true,
		})

		fs.Guilds = append(fs.Guilds, guild)
	}

	return nil
}

func (g *Guild) UpdateChannels() error {
	channels, err := g.FS.State.Channels(g.ID)
	if err != nil {
		return errors.Wrap(err, "Failed to get channels")
	}

	newChs := channels[:0]

	g.mu.Lock()
	defer g.mu.Unlock()

Main:
	for _, channel := range channels {
		for _, ch := range g.Channels {
			if ch.ID == channel.ID {
				continue Main
			}
		}

		newChs = append(newChs, channel)
	}

	for _, ch := range newChs {
		g.Channels = append(g.Channels, &Channel{
			ID:       ch.ID,
			FS:       g.FS,
			Inode:    NewInode(),
			Category: ch.CategoryID,
			Position: ch.Position,
		})
	}

	return nil
}

func (ch *Channel) Messages() ([]discord.Message, error) {
	return ch.FS.State.Messages(ch.ID)
}