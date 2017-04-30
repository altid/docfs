# docfs
Temporary implementation for a document -- markdown server

Usage: Modify DOCPATH to match where your documents are stored
mkdir "$XDG_RUNTIME_DIR/docs/nameofadocument"

inotifywait will pick up the creation of the directory, and parse the document. Currently, pdf and epub are supported.
