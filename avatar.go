package main

import (
	"errors"
	"fmt"
	"os"
	"path"
)

type AuthAvatar struct{}
var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()

	if len(url) > 0 {
		return url, nil
	}

	return "", ErrNoAvatarURL
}

// ErrNoAvatar is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client,
	// or returns an error if something goes wrong.
	// ErrNoAvatarURL is returned if the object is unable to get
	// a URL for the specified client.
	GetAvatarURL(ChatUser) (string, error)
}

type TryAvatars []Avatar

type GravatarAvatar struct{}
var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return fmt.Sprintf("//www.gravatar.com/avatar/%x", u.UniqueID()), nil

}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := os.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
				return fmt.Sprintf("/avatars/%s", file.Name()), nil
			}
		}
	}

	return "", ErrNoAvatarURL
}

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatars := range a {
		if url, err := avatars.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}
