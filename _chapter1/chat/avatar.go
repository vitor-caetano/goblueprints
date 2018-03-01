package main

import (
	"errors"
)

// ErrorNoAvatar is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrorNoAvatar = errors.New("chat: Unable to get an avatar URL.")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client,
	// or returns an error if something goes wrong.
	// ErrorNoAvatar is returned if the object is unable to get
	// a URL for the specified client.
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrorNoAvatar
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userId, ok := c.userData["userid"]; ok {
		if userIdStr, ok := userId.(string); ok {
			return "//www.gravatar.com/avatar/" + userIdStr, nil
		}
	}
	return "", ErrorNoAvatar
}
