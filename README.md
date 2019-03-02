# Docfs
This doesn't even build right now. 

## Overview

docfs converts pdf files, such that they can be read through ubqt clients.
Currently, a simple text representation of a given pdf is provided only.

## General Usage

`docfs [-p <dir>]`
<dir> will default to /tmp/ubqt/docs if none is provided.

To load a document for reading via ubqt, you would do as follows.

```
echo 'open /path/to/a/normal/document.pdf' > <dir>/docs/ctrl
ls <dir>/docs/document.pdf/
	document
	title
	status
	tabs
```


## Typcial Installation

To add docfs to your PATH:
`go install github.com/ubqt-systems/docfs`

Alternatively, you can build it in your current working directory
`go build -o docfs github.com/ubqt-systems/docfs`

## Configuration

docfs uses the general ubqt.cfg file, used by many of the file servers in ubqt.
The only available option at the time of writing is setting the directory for cached 
documents (documents you've previously opened)

```
service=docs
	log=/usr/halfwit/log

```

## PDF

Currently, pdfs are unsupported. You can toy around with opening them, and will see some structural elements such as a ToC and a title; but the main body parsing is incomplete.

## Epub

Currently, epubs are in an alpha state of support. They are fully converted to ubqt-compatible directory structures, but there may be elements which are missing, incorrectly formatted, or invalid when read by particular clients. 

## Persistence

ubqt servers are meant to live in temporary storage, with anything long lasting relying on persistent caches. docfs uses the log= config parameter to dictate where these logs live - it stores the document, translated to markdown, at <log>/doc/thatfile

## With 9p-server, Clients

When 9p-server is ran in the parent directory of docfs's path, ie /tmp/ubqt, it will find and integrate any documents you open. Normal buffer management, including open, and close work as usual. Closing a document will remove it from docfs's path, ie /tmp/ubqt/docs/, but the persistent copy that lives in your log directory will not be destroyed. Subsequent calls to open on the same document will skip parsing the main body of the document, and instead only need to populate the sidebar (ToC), the title, and any applicable status messages will be written to status.

Clients connecting over 9p will be able to view a tab list of open buffers from docfs if authenticated to do so, and the read state on any given document will be stashed - meaning you can migrate clients and not lose your place in a document!

It is important to note, read states will *not* be stored across restarts of docfs, and opening on a previously closed document, likewise will not synchronize read state. 
