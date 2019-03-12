# Docfs

Docfs is a file service, which translates pdf and epub documents into ubqt-flavored markup.

## Usage

`docfs [-p <dir>]`

 - <dir> will default to /tmp/ubqt/docs if none is provided.

## Typcial Installation

To add docfs to your PATH:
`go install github.com/ubqt-systems/docfs`

Alternatively, you can build it in your current working directory
`go build -o docfs github.com/ubqt-systems/docfs`

## Configuration

```
# ubqt.cfg - Place this in your distros' default configuration directory
service=docs
	log=/usr/halfwit/log
	#listen_address=192.168.0.4
```
 
 - log is a location to store the body of markdown from parsed documents. A special value of `none` disables logging.
 - listen_address is an advanced feature, see [Using listen_address](https://ubqt-systems.github.io/using-listen-address.html)

> See [ubqt configuration](https://ubqt-systems.github.io/ubqt-configurations.html) for more information

## PDF

Currently, pdfs are unsupported. You can toy around with opening them, and will see some structural elements such as a ToC and a title; but the main body parsing is incomplete.

## Epub

Currently, epubs are in an alpha state of support. They are fully converted to ubqt-compatible directory structures, but there may be elements which are missing, incorrectly formatted, or invalid when read by particular clients. 
Image assets for a given document will be made available, but currently are not.
