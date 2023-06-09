package ecode

import "errors"

var (
	ErrNullPointer             = errors.New("null reference instance")
	ErrNoCommand               = errors.New("no command")
	ErrLogLevel                = errors.New("invalid log level")
	ErrNoBlueskyHandle         = errors.New("no Bluesky handle")
	ErrInvalidBlueskyRecordURI = errors.New("invalid Bluesky record URI")
	ErrInvalidMastodonUserId   = errors.New("invalid Mastodon user ID")
	ErrNoContent               = errors.New("no content")
	ErrTooLargeImage           = errors.New("too large image (>1MB)")
	ErrNoAPODImage             = errors.New("no APOD image")
	ErrExistAPODData           = errors.New("exist APOD data")
	ErrNoFeed                  = errors.New("no feed")
)

/* Copyright 2023 Spiegel
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
