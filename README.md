# Docfs

Docfs is a file service, which translates pdf and epub documents into Altid-flavored markup.

`go install github.com/altid/docfs`

## Usage

`docfs [-p <dir>]`

 - <dir> will default to /tmp/altid/docs if none is provided.

## Configuration

```
# altid/config - Place this in your operating systems' default configuration directory
service=docs
	log=/usr/halfwit/log
	#listen_address=192.168.0.4
```
 
 - log is a location to store the body of markdown from parsed documents. A special value of `none` disables logging.
 - listen_address is a more advanced topic, explained here: [Using listen_address](https://altid.github.io/using-listen-address.html)

> See [altid configuration](https://altid.github.io/altid-configurations.html) for more information

## PDF

Currently, pdfs are unsupported. You can toy around with opening them, and will see some structural elements such as a ToC and a title; but the main body parsing is incomplete.

## Epub

Currently, epubs are in an alpha state of support. They are fully converted to Altid-compatible directory structures, but there may be elements which are missing, incorrectly formatted, or invalid when read by particular clients.

EPUB3 are currently not supported. Any help here is greatly appreciated!

