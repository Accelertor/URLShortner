package main

import "math/rand"

type URLShortener struct {
	urls map[string]string
}

func encode() string {
	const set = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
	sk := make([]byte, 5)
	for i := range sk {
		sk[i] = set[rand.Intn(len(set))]
	}
	return string(sk)
}

func (url *URLShortener) ShortUrl(input string) string {
	if url.urls == nil {
		url.urls = make(map[string]string)
	}

	code := encode()

	// Ensure uniqueness
	for {
		if _, exists := url.urls[code]; !exists {
			break
		}
		code = encode()
	}

	url.urls[code] = input // ✅ Store mapping
	return code
}

func (url *URLShortener) FindURL(short string) string {
	if url.urls == nil {
		return ""
	}
	return url.urls[short]
}
