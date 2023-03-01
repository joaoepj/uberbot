package core

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

// response.go
// This file contains structures and functions that make it easier to create and send response embeds

// ResponseComponents
// Stores the components for response
// allows for functions to add data.
type ResponseComponents struct {
	Components        []discordgo.MessageComponent
	SelectMenuOptions []discordgo.SelectMenuOption
	Rows              int
}

// Response
// The Response type, can be build and sent to a given guild.
type Response struct {
	Ctx                *CmdContext
	Reply              bool
	Success            bool
	Deferred           bool
	Ephemeral          bool
	Embeds             []*discordgo.MessageEmbed
	ResponseComponents ResponseComponents
}

// -- Main Response handlers --

func NewResponse(ctx *CmdContext, deferred bool, ephemeral bool, rows int) *Response {
	r := &Response{
		Ctx:                ctx,
		Deferred:           deferred,
		Ephemeral:          ephemeral,
		ResponseComponents: newResponseComponents(rows),
		Embeds: []*discordgo.MessageEmbed{
			CreateEmbed(0, "", "", nil),
		},
	}
	if r.Deferred && ctx.Interaction != nil {
		if ephemeral {
			_ = Session.InteractionRespond(r.Ctx.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					// Ephemeral is type 64 don't ask why
					Flags: 1 << 6,
				},
			})
		}
		_ = Session.InteractionRespond(r.Ctx.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
	}
	if ctx.Cmd.Trigger != "" {
		r.AppendCommand()
	}
	if ctx.Message != nil {
		r.AppendField(0, "Invoked by:", r.Ctx.Message.Author.Mention(), false)
	}
	return r
}

// Send a compiled response.
func (r *Response) Send(success bool, title, description string, color int) {
	// Determine what color to use based on the success state and if the embed has an existing color
	if success && color == 0 {
		color = ColorSuccess
	} else if !success && color == 0 {
		color = ColorFailure
	}
	// Fill the main embed
	r.Embeds[0].Title = title
	r.Embeds[0].Description = description
	r.Embeds[0].Color = color

	// if the guild is nil, this is supposed to be sent to bot admins
	if r.Ctx.Guild == nil {
		for admin := range botAdmins {
			dmChannel, dmCreateErr := Session.UserChannelCreate(admin)
			if dmCreateErr != nil {
				// Since error reports also use DMs, sending this as an error report would be redundant
				// Just log the error
				Log.Errorf("failed to send response dm to admin: %s; Response title: %s", admin, r.Embeds[0].Title)
				return
			}
			_, dmSendErr := Session.ChannelMessageSendComplex(dmChannel.ID, &discordgo.MessageSend{
				Embeds:     r.Embeds,
				Components: r.ResponseComponents.Components,
			})
			if dmSendErr != nil {
				// Since error reports also use DMs, sending this as an error report would be redundant
				// Just log the error
				Log.Errorf("failed to send response dm to admin: %s; Response title %s", admin, r.Embeds[0].Title)
				return
			}
			return
		}
	}
	// If this is an interaction (Application Command or MessageComponent)
	// Handle as an interaction response
	if r.Ctx.Interaction != nil {
		// Some commands might take a bit to process information
		// Slash commands expect a response in three seconds or the interaction becomes invalidated
		// So we check to see if the command has been deferred
		if r.Deferred {
			r.handleDeferredResponse()
			return
		}
		r.handleInteractionResponse()
	}

	// Try sending the response in the configured output channel
	// If that fails, try sending the response in the current channel
	// If THAT fails, send an error report
	_, err := Session.ChannelMessageSendComplex(r.Ctx.Guild.Info.ResponseChannelID, &discordgo.MessageSend{
		Embeds:     r.Embeds,
		Components: r.ResponseComponents.Components,
	})
	if err != nil && r.Reply {
		// Reply to user if no output channel
		_, err = ReplyToUser(r.Ctx.Message.ChannelID, &discordgo.MessageSend{
			Embeds:     r.Embeds,
			Components: r.ResponseComponents.Components,
			Reference: &discordgo.MessageReference{
				MessageID: r.Ctx.Message.ID,
				ChannelID: r.Ctx.Message.ChannelID,
				GuildID:   r.Ctx.Guild.ID,
			},
			AllowedMentions: &discordgo.MessageAllowedMentions{
				RepliedUser: false,
			},
		})
		if err != nil {
			SendErrorReport(r.Ctx.Guild.ID, r.Ctx.Message.ChannelID, r.Ctx.Message.Author.ID, "Ultimately failed to send bot response", err)
		}
	} else if !r.Reply {
		// If the command does not want to reply lets just send it to the channel the command was invoked
		_, err = Session.ChannelMessageSendComplex(r.Ctx.Message.ChannelID, &discordgo.MessageSend{
			Embeds:     r.Embeds,
			Components: r.ResponseComponents.Components,
		})
	}
}

// handleInteractionResponse
// handles interaction responses.
func (r *Response) handleInteractionResponse() {
	// Check to see if the command is ephemeral (only shown to the user)
	if r.Ephemeral {
		err := Session.InteractionRespond(r.Ctx.Interaction, &discordgo.InteractionResponse{
			// Ephemeral is type 64 don't ask why
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:      1 << 6,
				Embeds:     r.Embeds,
				Components: r.ResponseComponents.Components,
			},
		})
		if err != nil {
			Log.Errorf("unable to respond to interaction %s in %s (%s) ", r.Ctx.Interaction.ID, r.Ctx.Guild.Name, r.Ctx.Guild.ID)
		}
		return
	}

	// Default response for interaction
	err := Session.InteractionRespond(r.Ctx.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     r.Embeds,
			Components: r.ResponseComponents.Components,
		},
	})
	if err != nil {
		if err != nil {
			SendErrorReport(r.Ctx.Guild.ID, r.Ctx.Interaction.ChannelID, r.Ctx.Message.Author.ID, "Unable to send interaction messages", err)
		}
		if r.Ctx.Guild.Info.ResponseChannelID != "" {
			_, err = Session.ChannelMessageSendComplex(r.Ctx.Guild.Info.ResponseChannelID, &discordgo.MessageSend{
				Embeds:     r.Embeds,
				Components: r.ResponseComponents.Components,
			})
		} else {
			_, err = Session.ChannelMessageSendComplex(r.Ctx.Message.ChannelID, &discordgo.MessageSend{
				Embeds:     r.Embeds,
				Components: r.ResponseComponents.Components,
			})
		}

		if err != nil {
			SendErrorReport(r.Ctx.Guild.ID, r.Ctx.Interaction.ChannelID, r.Ctx.Message.Author.ID, "Unable to send message", err)
		}
	}
	return
}

// handleDeferredResponse
// handles responses that have been deferred.
func (r *Response) handleDeferredResponse() {
	// Check to see if the command is ephemeral (only shown to the interaction initiator)
	if r.Ephemeral {
		_, err := Session.InteractionResponseEdit(r.Ctx.Interaction, &discordgo.WebhookEdit{
			Components: &r.ResponseComponents.Components,
			Embeds:     &r.Embeds,
		})
		// Just in case the interaction has been removed
		// Just in case the interaction gets removed.
		if err != nil {
			if err != nil {
				SendErrorReport(r.Ctx.Guild.ID, r.Ctx.Interaction.ChannelID, r.Ctx.Message.Author.ID, "Unable to send interaction messages", err)
			}
			if r.Ctx.Guild.Info.ResponseChannelID != "" {
				_, err = Session.ChannelMessageSendComplex(r.Ctx.Guild.Info.ResponseChannelID, &discordgo.MessageSend{
					Embeds:     r.Embeds,
					Components: r.ResponseComponents.Components,
				})
			} else {
				_, err = Session.ChannelMessageSendComplex(r.Ctx.Message.ChannelID, &discordgo.MessageSend{
					Embeds:     r.Embeds,
					Components: r.ResponseComponents.Components,
				})
			}

			if err != nil {
				SendErrorReport(r.Ctx.Guild.ID, r.Ctx.Interaction.ChannelID, r.Ctx.Message.Author.ID, "Unable to send message", err)
			}
		}
		r.Deferred = false
		return
	}

	// Just respond normally
	_, err := Session.InteractionResponseEdit(r.Ctx.Interaction, &discordgo.WebhookEdit{
		Embeds:     &r.Embeds,
		Components: &r.ResponseComponents.Components,
	})
	// Just in case the interaction gets removed.
	if err != nil {
		_, err := Session.ChannelMessageSendComplex(r.Ctx.Guild.Info.ResponseChannelID, &discordgo.MessageSend{
			Embeds:     r.Embeds,
			Components: r.ResponseComponents.Components,
		})
		if err != nil {
			_, err = Session.ChannelMessageSendComplex(r.Ctx.Message.ChannelID, &discordgo.MessageSend{
				Embeds:     r.Embeds,
				Components: r.ResponseComponents.Components,
			})
		}
	}
	r.Deferred = false
	return
}

// -- Embeds --

// CreateEmbed
// Creates a MessageEmbed struct with basic information.
func CreateEmbed(color int, title, description string, fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
		Fields:      fields,
	}
}

// AppendEmbed
// Creates a new embed and appends it to the existing Embed slice.
func (r *Response) AppendEmbed(color int, title, description string, fields []*discordgo.MessageEmbedField) {
	r.Embeds = append(r.Embeds, CreateEmbed(color, title, description, fields))
}

// -- Embed Fields --

// CreateField
// Creates a MessageEmbedField struct with given information.
func CreateField(name string, value string, inline bool) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	}
}

// AppendField
// Create a new basic field and append it to an existing Response.
func (r *Response) AppendField(index int, name, value string, inline bool) {
	if !(index >= 0 && index < len(r.Embeds)) {
		Log.Errorf("embed index is out of range!")
		return
	}
	r.Embeds[index].Fields = append(r.Embeds[index].Fields, CreateField(name, value, inline))
}

// PrependField
// Create a new basic field and prepend it to an existing Response.
func (r *Response) PrependField(index int, name, value string, inline bool) {
	if !(index >= 0 && index < len(r.Embeds)) {
		Log.Errorf("embed index is out of range!")
		return
	}
	fields := []*discordgo.MessageEmbedField{CreateField(name, value, inline)}
	r.Embeds[index].Fields = append(fields, r.Embeds[index].Fields...)
}

// AppendUsage
// Add the command usage to the response. Intended for syntax error responses.
func (r *Response) AppendUsage() {

}

// PrependAuthor
// Adds the author field to the specified embed.
func (r *Response) PrependAuthor(index int, name, url, iconUrl string) {
	if !(index >= 0 && index < len(r.Embeds)) {
		Log.Errorf("embed index is out of range!")
		return
	}
	r.Embeds[index].Author = &discordgo.MessageEmbedAuthor{
		URL:     url,
		Name:    name,
		IconURL: iconUrl,
	}
}

// AppendFooter
// Adds the footer field to the specified embed.
func (r *Response) AppendFooter(index int, text, icon string, timestamp bool) {
	if !(index >= 0 && index < len(r.Embeds)) {
		Log.Errorf("embed index is out of range!")
		return
	}
	if timestamp {
		r.Embeds[index].Timestamp = time.Now().Format(time.RFC3339)
	}
	r.Embeds[index].Footer = &discordgo.MessageEmbedFooter{
		Text:    text,
		IconURL: icon,
	}
}

// -- Message Components --

// CreateComponentFields
// Returns a slice of a Message Component, containing a singular ActionsRow.
func createComponentFields(rows int) []discordgo.MessageComponent {
	if rows == 1 {
		return []discordgo.MessageComponent{
			discordgo.ActionsRow{},
		}
	} else if rows < 1 {
		return nil
	}
	msgComponent := make([]discordgo.MessageComponent, 0, rows)
	for i := 0; i < rows; i++ {
		msgComponent = append(msgComponent, discordgo.ActionsRow{})
	}
	return msgComponent
}

// newResponseComponents.
func newResponseComponents(rows int) ResponseComponents {
	resp := ResponseComponents{
		Components: createComponentFields(rows),
		Rows:       rows,
	}
	return resp
}

// CreateButton creates a new button struct
// label: the button label
// style: the button style
// customID: the id for handling the button
// url: an optional url for the button
// disabled: if the button is disabled
func CreateButton(label string, style discordgo.ButtonStyle, customID string, url string, disabled bool) *discordgo.Button {
	button := &discordgo.Button{
		Label:    label,
		Style:    style,
		Disabled: disabled,
		Emoji:    discordgo.ComponentEmoji{},
		URL:      url,
		CustomID: customID,
	}
	return button
}

// CreateSelect TODO add doc.
func CreateSelect(customID string, placeholder string, options []discordgo.SelectMenuOption) discordgo.SelectMenu {
	dropDown := discordgo.SelectMenu{
		CustomID:    customID,
		Placeholder: placeholder,
		Options:     options,
	}
	return dropDown
}

// AppendButton
// Appends a button.
func (r *Response) AppendButton(label string, style discordgo.ButtonStyle, url string, customID string, rowID int) {
	row := r.ResponseComponents.Components[rowID].(discordgo.ActionsRow)
	row.Components = append(row.Components, CreateButton(label, style, customID, url, false))
	r.ResponseComponents.Components[rowID] = row
}

// -- Util -- //

// AppendCommand
// Appends command information for logging purposes.
func (r *Response) AppendCommand() {
	// Get the command used as a string, and all interpreted arguments, so it can be a part of the output
	commandUsed := ""
	if r.Ctx.Cmd.IsChild {
		commandUsed = fmt.Sprintf("%s%s %s", r.Ctx.Guild.Info.Prefix, r.Ctx.Cmd.ParentID, r.Ctx.Cmd.Trigger)
	} else {
		commandUsed = r.Ctx.Guild.Info.Prefix + r.Ctx.Cmd.Trigger
	}
	// Just makes the thing prettier
	if r.Ctx.Interaction != nil {
		commandUsed = "/" + r.Ctx.Cmd.Trigger
	}
	for _, k := range r.Ctx.Cmd.Arguments.Keys() {
		arg := r.Ctx.Args[k]
		if arg.StringValue() == "" {
			continue
		}
		vv, ok := r.Ctx.Cmd.Arguments.Get(k)

		if ok {
			argInfo := vv.(*ArgInfo)
			switch argInfo.TypeGuard {
			case Int:
				fallthrough
			case Boolean:
				fallthrough
			case String:
				commandUsed += " " + arg.StringValue()
				break
			case User:
				user, err := arg.UserValue(Session)
				if err != nil {
					commandUsed += " " + arg.StringValue()
				} else {
					commandUsed += " " + user.Mention()
				}
			case Role:
				role, err := arg.RoleValue(Session, r.Ctx.Guild.ID)
				if err != nil {
					commandUsed += " " + arg.StringValue()
				} else {
					commandUsed += " " + role.Mention()
				}
			case Channel:
				channel, err := arg.ChannelValue(Session)
				if err != nil {
					commandUsed += " " + arg.StringValue()
				} else {
					commandUsed += " " + channel.Mention()
				}
			}
		} else {
			commandUsed += " " + arg.StringValue()
		}
	}

	commandUsed = "```\n" + commandUsed + "\n```"

	r.AppendField(0, "Command used:", commandUsed, false)
}
func ReplyToUser(channelID string, messageSend *discordgo.MessageSend) (*discordgo.Message, error) {
	return Session.ChannelMessageSendComplex(channelID, messageSend)
}
