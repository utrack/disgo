package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/disgo/rest/route"
)

type User struct {
	discord.User
	Bot *Bot
}

// AvatarURL returns the Avatar URL of this User
func (u *User) AvatarURL(size int) *string {
	return u.getAssetURL(route.UserAvatar, u.Avatar, size)
}

// DefaultAvatarURL returns the default avatar URL of this User
func (u *User) DefaultAvatarURL(size int) string {
	discriminator, _ := strconv.Atoi(u.Discriminator)
	compiledRoute, err := route.DefaultUserAvatar.Compile(nil, route.PNG, size, discriminator%5)
	if err != nil {
		return ""
	}
	return compiledRoute.URL()
}

// EffectiveAvatarURL returns either this User avatar or default avatar depending on if this User has one
func (u *User) EffectiveAvatarURL(size int) string {
	if u.Avatar == nil {
		return u.DefaultAvatarURL(size)
	}
	return *u.AvatarURL(size)
}

// BannerURL returns the Banner URL of this User
func (u *User) BannerURL(size int) *string {
	return u.getAssetURL(route.UserBanner, u.Banner, size)
}

// String returns this User as a mention
func (u *User) String() string {
	return "<@" + u.ID.String() + ">"
}

// Tag returns this User's Username and Discriminator in Username#Discriminator format
func (u *User) Tag() string {
	return fmt.Sprintf("%s#%s", u.Username, u.Discriminator)
}

func (u *User) getAssetURL(cdnRoute *route.CDNRoute, assetId *string, size int) *string {
	if assetId == nil {
		return nil
	}
	format := route.PNG
	if strings.HasPrefix(*assetId, "a_") {
		format = route.GIF
	}
	compiledRoute, err := cdnRoute.Compile(nil, format, size, u.ID, *assetId)
	if err != nil {
		return nil
	}
	url := compiledRoute.URL()
	return &url
}

// OpenDMChannel creates a DMChannel between this User and the Bot
func (u *User) OpenDMChannel(opts ...rest.RequestOpt) (*Channel, rest.Error) {
	channel, err := u.Bot.RestServices.UserService().CreateDMChannel(u.ID, opts...)
	if err != nil {
		return nil, err
	}
	// TODO: should we caches it here? or do we get a gateway event?
	return u.Bot.EntityBuilder.CreateChannel(*channel, CacheStrategyYes), nil
}
