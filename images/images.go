package images

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"io"
	"os"

	"github.com/goark/errs"
	"github.com/goark/fetch"
	"github.com/goark/toolbox/ecode"
	"golang.org/x/image/draw"
)

// FetchFromURL returns binary image from URL.
func FetchFromURL(ctx context.Context, urlStr string) ([]byte, error) {
	u, err := fetch.URL(urlStr)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", urlStr))
	}
	resp, err := fetch.New().GetWithContext(ctx, u)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", u.String()))
	}
	defer resp.Close()

	b, err := io.ReadAll(resp.Body())
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("url", u.String()))
	}
	return b, nil
}

// FetchFromFile returns binary image from local file.
func FetchFromFile(fn string) ([]byte, error) {
	b, err := os.ReadFile(fn)
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("file_name", fn))
	}
	return b, nil
}

const (
	imageMaxSize     = 1024
	imageFileMaxSize = 1024 * 1024
)

func AjustImage(src []byte) (io.Reader, error) {
	// check file size
	if len(src) < imageFileMaxSize {
		return bytes.NewReader(src), nil
	}

	// decode image
	imgSrc, t, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		return nil, errs.Wrap(err)
	}
	// convert JPEG
	if t != "jpeg" {
		b, err := convertJPEG(imgSrc, 100)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		if len(b) < imageFileMaxSize {
			return bytes.NewReader(b), nil
		}
		src = b
		imgSrc, _, err = image.Decode(bytes.NewReader(src))
		if err != nil {
			return nil, errs.Wrap(err)
		}
	}
	// quality down
	quality := 100
	for _, q := range []int{75, 50, 25} {
		b, err := convertJPEG(imgSrc, q)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		quality = q
		if len(b) < imageFileMaxSize {
			return bytes.NewReader(b), nil
		}
	}

	// rectange of image
	rctSrc := imgSrc.Bounds()
	rate := 1.0
	if rctSrc.Dx() > rctSrc.Dy() {
		if rctSrc.Dx() > imageMaxSize {
			rate = imageMaxSize / float64(rctSrc.Dx())
		}
	} else {
		if rctSrc.Dy() > imageMaxSize {
			rate = imageMaxSize / float64(rctSrc.Dy())
		}
	}
	if rate >= 1.0 {
		return nil, errs.Wrap(ecode.ErrTooLargeImage)
	}

	// scale down
	dstX := int(float64(rctSrc.Dx()) * rate)
	dstY := int(float64(rctSrc.Dy()) * rate)
	imgDst := image.NewRGBA(image.Rect(0, 0, dstX, dstY))
	draw.BiLinear.Scale(imgDst, imgDst.Bounds(), imgSrc, rctSrc, draw.Over, nil)
	b, err := convertJPEG(imgDst, quality)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	if len(b) > imageFileMaxSize {
		return nil, errs.Wrap(ecode.ErrTooLargeImage)
	}
	return bytes.NewReader(b), nil
}

func convertJPEG(src image.Image, quality int) ([]byte, error) {
	dst := &bytes.Buffer{}
	if err := jpeg.Encode(dst, src, &jpeg.Options{Quality: quality}); err != nil {
		return nil, errs.Wrap(err)
	}
	return dst.Bytes(), nil
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
