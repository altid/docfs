# Docfs
This doesn't even build right now. 

## Overview

docfs converts pdf files, such that they can be read through ubqt clients.
Currently, a simple text representation of a given pdf is provided only.

## Usage

`docfs -p <dir>`

Currently, only pdfs are supported. To load a pdf for reading via ubqt, you would do as follows.

```
echo 'open /path/to/a/normal/document.pdf' > <dir>/ctrl
ls <dir>/document/
	doc
	title
	status
	tabs
```

## Persistence

ubqt servers are meant to live in temporary storage, with anything long lasting relying on persistent caches. docfs uses the log= config parameter to dictate where these logs live - it stores the document, translated to markdown, at <log>/doc/thatfile.md

