// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

// Package fmtp implements per codec parsing of fmtp lines
package fmtp

import (
	"strings"
)

func parseParameters(line string) map[string]string {
	parameters := make(map[string]string)

	for _, p := range strings.Split(line, ";") {
		pp := strings.SplitN(strings.TrimSpace(p), "=", 2)
		key := strings.ToLower(pp[0])
		var value string
		if len(pp) > 1 {
			value = pp[1]
		}
		parameters[key] = value
	}

	return parameters
}

// FMTP interface for implementing custom
// FMTP parsers based on MimeType.
type FMTP interface {
	// MimeType returns the MimeType associated with
	// the fmtp
	MimeType() string
	// Match compares two fmtp descriptions for
	// compatibility based on the MimeType
	Match(f FMTP) bool
	// Parameter returns a value for the associated key
	// if contained in the parsed fmtp string
	Parameter(key string) (string, bool)
}

// Parse parses an fmtp string based on the MimeType.
func Parse(mimeType, line string) FMTP {
	var fmtp FMTP

	parameters := parseParameters(line)

	switch {
	case strings.EqualFold(mimeType, "video/h264"):
		fmtp = &h264FMTP{
			parameters: parameters,
		}

	case strings.EqualFold(mimeType, "video/vp9"):
		fmtp = &vp9FMTP{
			parameters: parameters,
		}

	case strings.EqualFold(mimeType, "video/av1"):
		fmtp = &av1FMTP{
			parameters: parameters,
		}

	default:
		fmtp = &genericFMTP{
			mimeType:   mimeType,
			parameters: parameters,
		}
	}

	return fmtp
}

type genericFMTP struct {
	mimeType   string
	parameters map[string]string
}

func (g *genericFMTP) MimeType() string {
	return g.mimeType
}

// Match returns true if g and b are compatible fmtp descriptions
// The generic implementation is used for MimeTypes that are not defined.
func (g *genericFMTP) Match(b FMTP) bool {
	fmtp, ok := b.(*genericFMTP)
	if !ok {
		return false
	}

	if !strings.EqualFold(g.mimeType, fmtp.MimeType()) {
		return false
	}

	for k, v := range g.parameters {
		if vb, ok := fmtp.parameters[k]; ok && !strings.EqualFold(vb, v) {
			return false
		}
	}

	for k, v := range fmtp.parameters {
		if va, ok := g.parameters[k]; ok && !strings.EqualFold(va, v) {
			return false
		}
	}

	return true
}

func (g *genericFMTP) Parameter(key string) (string, bool) {
	v, ok := g.parameters[key]

	return v, ok
}
