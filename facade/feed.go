package facade

import (
	"github.com/goark/errs"
	"github.com/goark/gocli/rwi"
	"github.com/goark/toolbox/ecode"
	"github.com/goark/toolbox/webpage"
	"github.com/spf13/cobra"
)

// newFeedCmd returns cobra.Command instance for show sub-command
func newFeedCmd(ui *rwi.RWI) *cobra.Command {
	webpageCmd := &cobra.Command{
		Use:     "feed",
		Aliases: []string{"rss"},
		Short:   "Handling information for Web feed",
		Long:    "Handling information for Web feed.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return debugPrint(ui, errs.Wrap(ecode.ErrNoCommand))
		},
	}
	webpageCmd.PersistentFlags().StringP("url", "u", "", "Feed URL")
	webpageCmd.PersistentFlags().StringP("flickr-id", "", "", "Flickr ID")
	webpageCmd.PersistentFlags().StringP("feed-list-file", "f", "", "path of Feed list file")
	webpageCmd.PersistentFlags().BoolP("save", "", false, "Save webpage data to cache")

	webpageCmd.AddCommand(
		newFeedLookupCmd(ui),
		newFeedPostCmd(ui),
	)
	return webpageCmd
}

func getFeedAll(cmd *cobra.Command, cfg *webpage.Webpage) ([]*webpage.Info, error) {
	urlStr, err := cmd.Flags().GetString("url")
	if err != nil {
		return nil, errs.Wrap(err)
	}
	flickrID, err := cmd.Flags().GetString("flickr-id")
	if err != nil {
		return nil, errs.Wrap(err)
	}
	feedListPath, err := cmd.Flags().GetString("feed-list-file")
	if err != nil {
		return nil, errs.Wrap(err)
	}
	if len(urlStr) > 0 {
		if err := cfg.Feed(cmd.Context(), urlStr); err != nil {
			return nil, errs.Wrap(err, errs.WithContext("url", urlStr))
		}
	}
	if len(flickrID) > 0 {
		if err := cfg.FeedFlickr(cmd.Context(), flickrID); err != nil {
			return nil, errs.Wrap(err, errs.WithContext("flickr_id", flickrID))
		}
	}
	if len(feedListPath) > 0 {
		fl, err := webpage.NewFeedList(feedListPath)
		if err != nil {
			return nil, errs.Wrap(err, errs.WithContext("feed_list_file", feedListPath))
		}
		if err := fl.Parse(cmd.Context(), cfg); err != nil {
			return nil, errs.Wrap(err, errs.WithContext("feed_list_file", feedListPath))
		}
	}
	cfg.StopPool()
	if err := cfg.GetErrorInPool(); err != nil {
		return nil, errs.Wrap(err, errs.WithContext("feed_list_file", feedListPath))
	}
	return cfg.SaveInfoToCache(), nil
}

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
