package main

import "strings"

func filter_1(s string) bool {
	s_ := strings.Split(s, ">")
	_s := s_[len(s_)-1]
	return !strings.Contains(s, "<script") &&
		!strings.Contains(s, "a new window") &&
		!strings.Contains(s, "<style") &&
		!strings.Contains(s, ".js") &&
		!strings.Contains(s, ".css") &&
		!strings.Contains(s, "wp-admin") &&
		!strings.Contains(s, "wp-admin") &&
		!strings.Contains(s, ".php") &&
		!strings.Contains(s, ".json") &&
		!strings.Contains(s, "__") &&
		!strings.Contains(s, "#shopify") &&
		!strings.Contains(s, "@media") &&
		!strings.Contains(s, "@font") &&
		!strings.Contains(s, ":{}") &&
		!strings.Contains(s, "::") &&
		!strings.Contains(s, "~=") &&
		!strings.Contains(s, "});") &&
		!strings.Contains(s, "//@") &&
		!strings.Contains(s, ".preventDefault") &&
		!strings.Contains(s, "function(") &&
		!strings.Contains(strings.ToLower(s), "jquery") && strings.ToLower(_s) != "youtube" &&
		strings.ToLower(_s) != "google" &&
		strings.ToLower(_s) != "facebook" &&
		strings.ToLower(_s) != "instagram" &&
		strings.ToLower(_s) != "spotify" &&
		strings.ToLower(_s) != "shopify" &&
		strings.ToLower(_s) != "pinterest" &&
		strings.ToLower(_s) != "tiktok" &&
		strings.ToLower(_s) != "tik tok" &&
		strings.ToLower(_s) != "about" &&
		strings.ToLower(_s) != "contact" &&
		strings.ToLower(_s) != "home" &&
		strings.ToLower(_s) != "twitter" &&
		strings.ToLower(_s) != "face book" &&
		strings.ToLower(_s) != "you tube" &&
		strings.ToLower(_s) != "snapchat" &&
		strings.ToLower(_s) != "snap chat" &&
		strings.Count(_s, ";") < 6
}

func fixback(s string) string {
	chars := "abcdefghijklmnopqrstuvwxyz"
	if len(s) > 1 {
		if !strings.Contains(chars, string(s[len(s)-1])) {
			fixfront(s[:len(s)-1])
		}
	}
	return s
}

func fixfront(s string) string {
	chars := "abcdefghijklmnopqrstuvwxyz"
	if len(s) > 1 {
		if !strings.Contains(chars, string(s[0])) {
			fixfront(s[1:])
		}
	}
	return s
}

func isMostlyLetters(s string) bool {
	chars := "abcdefghijklmnopqrstuvwxyz"
	charcount := 0
	for i := 0; i <= len(s)-1; i++ {
		if strings.Contains(chars, strings.ToLower(string(s[i]))) {
			charcount = charcount + 1
		}
	}
	if charcount > 0 {
		return float64(len(s))/float64(charcount) < 2.00
	}
	return false
}
