package disgord

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/andersfylling/disgord/discordws"
	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/guild"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Token      string
	HTTPClient *http.Client
	Debug      bool
}

// NewDisgord creates a new default disgord instance
func NewClient(conf *Config) (*Client, error) {

	if conf.HTTPClient == nil {
		// http client configuration
		conf.HTTPClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	// Use discordws to keep the socket connection going
	dws, err := discordws.NewClient(&discordws.Config{
		// user settings
		Token:      conf.Token,
		HTTPClient: conf.HTTPClient,
		Debug:      conf.Debug,

		// lib specific
		DAPIVersion:  endpoint.APIVersion,
		DAPIEncoding: endpoint.APIComEncoding,
	})
	if err != nil {
		return nil, err
	}

	// create a disgord instance
	d := &Client{
		HTTPClient: conf.HTTPClient,
		ws:         dws,
		EventChan:  dws.GetEventChannel(),
		Token:      conf.Token,
		Dispatcher: event.NewDispatcher(),
		State:      NewStateCache(),
	}

	return d, nil
}

// NewRequiredDisgord same as NewDisgord, but exits program if an error occours
func NewRequiredClient(conf *Config) *Client {
	dg, err := NewClient(conf)
	if err != nil {
		logrus.Fatal(err)
	}

	return dg
}

type Client struct {
	sync.RWMutex

	Token string

	ws *discordws.Client

	HTTPClient *http.Client

	EventChan <-chan discordws.EventInterface

	// register listeners for events
	*event.Dispatcher

	// cache
	State StateCacher
}

func (c *Client) String() string {
	return c.ws.String()
}

func (c *Client) logInfo(msg string) {
	logrus.WithFields(logrus.Fields{
		"lib": c.ws.String(),
	}).Info(msg)
}

func (c *Client) logErr(msg string) {
	logrus.WithFields(logrus.Fields{
		"lib": c.ws.String(),
	}).Error(msg)
}

// Connect establishes a websocket connection to the discord API
func (c *Client) Connect() (err error) {
	c.logInfo("Connecting to discord Gateway")
	err = c.ws.Connect()
	if err != nil {
		c.logErr(err.Error())
		return
	}
	c.logInfo("Connected")

	// setup event observer
	go c.eventObserver()

	return nil
}

// Disconnect closes the discord websocket connection
func (c *Client) Disconnect() (err error) {
	fmt.Println()
	c.logInfo("Closing Discord gateway connection")
	err = c.ws.Disconnect()
	if err != nil {
		c.logErr(err.Error())
		return
	}
	c.logInfo("Disconnected")

	return nil
}

func (c *Client) eventObserver() {
	for {
		select {
		case evt, ok := <-c.EventChan:
			if !ok {
				logrus.Error("Event channel is dead!")
				break
			}

			ctx := context.Background()

			// TODO: parsing JSON uses panic and not logging on issues..

			eventName := evt.Name()
			switch eventName {
			//case event.Ready:
			//case event.Resumed:
			//case event.ChannelCreate:
			//case event.ChannelUpdate:
			//case event.ChannelDelete:
			//case event.ChannelPinsUpdate:
			case event.GuildCreate, event.GuildUpdate, event.GuildDelete:
				g := &guild.Guild{}
				err := json.Unmarshal(evt.Data(), g)
				if err != nil {
					panic(err)
				}

				switch eventName { // internal switch statement for guild events
				case event.GuildCreate:
					// add to cache
					c.State.AddGuild(g)
					// notifify listeners
					go c.GuildCreateEvent.Trigger(ctx, g)
				case event.GuildUpdate:
					// update cache
					c.State.UpdateGuild(g)
					// notifify listeners
					go c.GuildUpdateEvent.Trigger(ctx, g)
				case event.GuildDelete:
					cachedGuild, err := c.State.Guild(g.ID)
					if err != nil {
						// guild has not been cached earlier for some reason..
					} else {
						// update instance with complete info.
						// Assumption: The cached version has no outdated information.
						g = nil
						g = cachedGuild
						// delete the guild object from the cache
						c.State.DeleteGuild(g)
					}

					// notify listeners
					go c.GuildDeleteEvent.Trigger(ctx, g)
				} // END internal switch statement for guild events
			//case event.GuildBanAdd:
			//case event.GuildBanRemove:
			//case event.GuildEmojisUpdate:
			//case event.GuildIntegrationsUpdate:
			//case event.GuildMemberAdd:
			//case event.GuildMemberRemove:
			//case event.GuildMemberUpdate:
			//case event.GuildMemberChunk:
			//case event.GuildRoleCreate:
			//case event.GuildRoleUpdate:
			//case event.GuildRoleDelete:
			//case event.MessageCreate:
			//case event.MessageUpdate:
			//case event.MessageDelete:
			//case event.MessageDeleteBulk:
			//case event.MessageReactionAdd:
			//case event.MessageReactionRemove:
			//case event.MessageReactionRemoveAll:
			//case event.PresenceUpdate:
			//case event.TypingStart:
			//case event.UserUpdate:
			//case event.VoiceStateUpdate:
			//case event.VoiceServerUpdate:
			//case event.WebhooksUpdate:

			default:
				fmt.Printf("------\nTODO\nImplement event handler for `%s`, data: \n%+v\n------\n\n", evt.Name(), string(evt.Data()))
			}
		}
	}
}
