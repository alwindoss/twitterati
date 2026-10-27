package twitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var idRe = regexp.MustCompile(`/status(?:es)?/(\d+)`)
var stripRe = regexp.MustCompile(`0+|\.`)

type tweet struct {
	MediaDetails []struct {
		Type      string `json:"type"`
		VideoInfo struct {
			Variants []struct {
				ContentType string `json:"content_type"`
				Bitrate     int    `json:"bitrate"`
				URL         string `json:"url"`
			} `json:"variants"`
		} `json:"video_info"`
	} `json:"mediaDetails"`
}

// Mirrors the token the X embed widget computes from the tweet ID.
func token(id string) string {
	n, _ := strconv.ParseFloat(id, 64)
	f := n / 1e15 * math.Pi
	ip := math.Floor(f)
	frac := f - ip
	s := strconv.FormatInt(int64(ip), 36)
	if frac > 0 {
		s += "."
		const digits = "0123456789abcdefghijklmnopqrstuvwxyz"
		for i := 0; i < 11 && frac > 0; i++ {
			frac *= 36
			d := int(frac)
			s += string(digits[d])
			frac -= float64(d)
		}
	}
	return stripRe.ReplaceAllString(s, "")
}

func get(url string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/124.0 Safari/537.36")
	return http.DefaultClient.Do(req)
}

type Config struct {
	OutputWriter io.Writer
	URL          string
}

func Run(cfg *Config) error {
	m := idRe.FindStringSubmatch(cfg.URL)
	if m == nil {
		return errors.New("could not find a tweet ID in the URL")
	}
	id := m[1]

	api := fmt.Sprintf("https://cdn.syndication.twimg.com/tweet-result?id=%s&lang=en&token=%s", id, token(id))
	resp, err := get(api)
	if err != nil {
		fatal(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("API returned " + resp.Status + " (tweet may be private, deleted, or age-restricted)")
	}

	var t tweet
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return errors.New("failed to decode API response: " + err.Error())
	}

	best, bestRate := "", -1
	for _, md := range t.MediaDetails {
		if md.Type != "video" && md.Type != "animated_gif" {
			continue
		}
		for _, v := range md.VideoInfo.Variants {
			if v.ContentType == "video/mp4" && v.Bitrate >= bestRate {
				best, bestRate = v.URL, v.Bitrate
			}
		}
		if best != "" {
			break // first video in the tweet
		}
	}
	if best == "" {
		return errors.New("no video found in the tweet")
	}

	if err := download(best, cfg.OutputWriter); err != nil {
		return err
	}
	return nil
}

func download(url string, out io.Writer) error {
	resp, err := get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	total := resp.ContentLength
	var done int64
	buf := make([]byte, 32*1024)
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			if _, err := out.Write(buf[:n]); err != nil {
				return err
			}
			done += int64(n)
			if total > 0 {
				fmt.Printf("\r%s %.1f%%", strings.Repeat("#", int(done*30/total)), float64(done)*100/float64(total))
			}
		}
		if rerr == io.EOF {
			return nil
		}
		if rerr != nil {
			return rerr
		}
	}
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, "error:", msg)
	os.Exit(1)
}
