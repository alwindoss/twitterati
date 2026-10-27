# twitterati

`twitterati` is a command-line tool that downloads the first MP4 video (or animated GIF) attached to a tweet from its URL.

## Requirements

- Go 1.27.3 or later
- Network access to X/Twitter's syndication API and the selected video's media URL

## Build

From the repository root, build the executable with:

```sh
go build -o twitterati .
```

Alternatively, run the tool directly with `go run .` without building an executable.

## Usage

```text
twitterati download --url <tweet-url> [--output <filename>]
```

The tweet URL is required. `--url` can also be written as `-u`. The output filename is optional and can be supplied with `--output` or `-o`.

### Examples

Run directly from the source:

```sh
go run . download --url https://x.com/user/status/1234567890123456789
```

Save to a specific filename:

```sh
./twitterati download -u https://x.com/user/status/1234567890123456789 -o clip.mp4
```

Show general or download-specific help:

```sh
./twitterati --help
./twitterati download --help
```

When no output filename is given, the tool uses a timestamped `output_*.mp4` filename. If the requested output file already exists, it tries to create a timestamp-prefixed filename instead of overwriting it.

## Limitations

- The tweet URL must contain a numeric tweet status ID.
- The tweet must be available to the syndication API. Private, deleted, or age-restricted tweets may not be downloadable.
- A tweet without a supported video MP4 variant cannot be downloaded.
