# docfs
Temporary implementation for a document -- html server

## Usage
Usage: Modify DOCPATH to match where your documents are stored

`printf '%s\n' "$DOCPATH/mydoc.pdf" > "$XDG_RUNTIME_DIR/doc/ctl"`

This will create a folder named mydoc.pdf, currently containing the html representation of the document passed in. In time, it will contain simply the markdown and image resources. 
