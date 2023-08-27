package main

import (
	"github.com/joho/godotenv"
)

type App interface {
	ReadEnv() (map[string]string,  error)
}

type app struct {}

func (app) ReadEnv() (map[string]string, error) {
	envs, err := godotenv.Read(".env")

	envMaps := make(map[string]string)
	if err != nil {
		return envMaps, err
	}

	envMaps["fbClientId"] = envs["FB_CLIENT_ID"]
	envMaps["fbClientSecret"] = envs["FB_CLIENT_SECRET"]
	envMaps["fbCallback"] = envs["FB_CALLBACK"]

	envMaps["githubClientId"] = envs["GITHUB_CLIENT_ID"]
	envMaps["githubClientSecret"] = envs["GITHUB_CLIENT_SECRET"]
	envMaps["githubCallback"] = envs["GITHUB_CALLBACK"]

	envMaps["googleClientId"] = envs["GOOGLE_CLIENT_ID"]
	envMaps["googleClientSecret"] = envs["GOOGLE_CLIENT_SECRET"]
	envMaps["googleCallback"] = envs["GOOGLE_CALLBACK"]
	
	return envMaps, nil
}