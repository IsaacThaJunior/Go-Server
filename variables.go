package main

import "time"

const InternalServerMessage = "An error occured from our end"
const lengthOfChirp = 140

var badWords = map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}

const tokenTime = 1 * time.Hour
const RefreshTokenExipires = 60 * 24 * time.Hour
