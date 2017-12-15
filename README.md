# docfs
Temporary implementation for a document -- html server
Eventually this will be done in a dedicated binary

## Usage
Usage: Modify DOCPATH to match where your documents are stored

`printf 'open %s\n' "$DOCPATH/mydoc.pdf" > "$XDG_RUNTIME_DIR/doc/ctl"`

This will create a folder named mydoc.pdf, currently containing the html representation of the document passed in. In time, it will contain simply the markdown and image resources. 

## Future
This will allow a filesystem representation of arbitrary documents, such as pdf files or djvu files, via our ubqt-flavored markdown and supplemental resources (images, links)

This will allow a user to fluidly handle reading a single document, over the span of a day, on potentially many systems. 
