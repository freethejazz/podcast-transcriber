# A podcast transcription tool

This repository was created for a go-lang Sydney talk.

### What it does
A simple API that:
1. Accepts a url to a podcast mp3
1. Downloads the mp3 locally
1. Transcribes the mp3 using OpenAI's `whisper` command-line tool
1. Lightly processes and indexes the SRT (subtitle) file to elasticsearch

You can then use full-text search to find specific locations where
topics are mentioned in the podcasts indexed.

### Basic usage

#### Pre-requisites
The codebase relies on a command-line version of OpenAI's whisper
to be installed locally. Assuming you have python installed, you can
simply run:
```
pip install -U openai-whisper
```
Whisper relies on `ffmpeg`, which is usually available via package
managers:

```
# on Ubuntu or Debian
sudo apt update && sudo apt install ffmpeg

# on Arch Linux
sudo pacman -S ffmpeg

# on MacOS using Homebrew (https://brew.sh/)
brew install ffmpeg

# on Windows using Chocolatey (https://chocolatey.org/)
choco install ffmpeg

# on Windows using Scoop (https://scoop.sh/)
scoop install ffmpeg
```

#### Running it
Make sure elasticsearch is running (`docker-compose up`), then
build and run the application using your method of choice, e.g.,

```
go build -o transcriber
./transcriber
```
or for live-reloading during development
```
air
```

Once Gin finishes starting up, you can interact with the API.

#### Creating indexing jobs
The following is an example request to kick off a transcription job for
a particular URL.
```
curl --location 'http://localhost:8080/job' \
--header 'Content-Type: application/json' \
--data '{"url": "https://dts.podtrac.com/redirect.mp3/feeds.soundcloud.com/stream/1181362132-ideo_u-well-said-strategy-is-a-set-of-choices.mp3"}'
```

Example response:
```
{
    "jobId": "c6175d1b-744e-494d-8360-f2a88a943686",
    "status": "started"
}
```

#### Searching
The following is an example request to search for words or phrases.
```
curl --location 'http://localhost:8080/search' \
--header 'Content-Type: application/json' \
--data '{"query": "choices"}'
```

Example response:
```
{
    "results": [
        {
            "url": "https://dts.podtrac.com/redirect.mp3/feeds.soundcloud.com/stream/1181362132-ideo_u-well-said-strategy-is-a-set-of-choices.mp3",
            "index": 7,
            "text": "And so that's why I want to focus on the problem as the starting point because that should",
            "context": "set of choices. And so that's why I want to focus on the problem as the starting point because that should inform what kind of choices would make it go away.",
            "timestamp_from": 30920000000,
            "timestamp_to": 37000000000,
            "clip_length": 6080000000
        },
        {
            "url": "https://dts.podtrac.com/redirect.mp3/feeds.soundcloud.com/stream/1181362132-ideo_u-well-said-strategy-is-a-set-of-choices.mp3",
            "index": 5,
            "text": "that problem is going to stay around if not get worse until such time as you make a different",
            "context": "And if there's a problem, it's a problem of your current choices and generally speaking that problem is going to stay around if not get worse until such time as you make a different set of choices.",
            "timestamp_from": 23960000000,
            "timestamp_to": 29320000000,
            "clip_length": 5360000000
        },
        // etc
    ]
}
```

### What is an SRT file?
A structured text file that represents subtitles and the timestamps they
are associated with. As an example:

```
1
00:00:00,000 --> 00:00:06,640
Strategy should be thought first and foremost as a problem-solving tool.

2
00:00:06,640 --> 00:00:12,640
That is to say you have some problem because whatever problem you have now is the result

3
00:00:12,640 --> 00:00:17,840
of all the choices you've made that cause you to be doing what you're doing now.

4
00:00:17,840 --> 00:00:23,960
And if there's a problem, it's a problem of your current choices and generally speaking

5
00:00:23,960 --> 00:00:29,320
that problem is going to stay around if not get worse until such time as you make a different

6
00:00:29,320 --> 00:00:30,920
set of choices.

7
00:00:30,920 --> 00:00:37,000
And so that's why I want to focus on the problem as the starting point because that should

8
00:00:37,000 --> 00:00:41,960
inform what kind of choices would make it go away.
```

#### What's with the post processing thing for SRT captions?
I couldn't have indexed the whole podcast transcript as a single
document. `whisper` writes out a standard text file that has everything
in a human readable form without timestamps. This would have been good
enough to pinpoint a particular podcast, but then I'd be hunting around
trying to find the snippet I actually wanted to hear. Indexing the SRT
data means I can both search for the podcast, but also find timestamp
ranges pretty close to the search results.

The problem that arises here is that separate documents don't know about
each other, so if a particular phrase falls across two subtitle lines, it
won't match. By concatenating the subtitle text from the line before and
the line after the current subtitle, it ensures we don't lose that sort
of searchability.

Thanks goes to [this SO post](https://stackoverflow.com/questions/28431583/searching-subtitle-data-in-elasticsearch).
