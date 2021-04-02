# glom

A code editor for the visual thinkers

Glom is a text editor makes structured documents easy to
create and change. Many text editors now allow you to select text and drag it
around, but almost without exception, a couple of stray keypresses can turn a
neat 20 level hierarchy of brackets into spaghetti covered in nasty red syntax
error markers.

Just for fun, select a block of code in a typical modern IDE, and and drag it
around and see how long it takes you to turn it into a 3000 piece jigsaw puzzle.

Glom understands structure and makes it easy to create documents with coherent
structure at the macro as well as micro scale.

When you move your mouse around text, it faintly highlights the segmentation
that will be enforced if you click and drag, you can fold sections in, and
easily rearrange complex bracketed structures as found frequently in graphical
user interface code.

Not only that, but in addition to code, Glom makes it simple to annotate 
documents, allowing detailed side-by-side documentation and code.

Lastly, the undo buffer branches instead of truncating previous edits off the
tree, it provides a cue to indicate that you can fork and get back on a 
different timeline entirely, without losing the old, as a temporarily visible 
represetation of the graph of edits.

Instead of storing the document formatted, it continuously records input and 
replays the chains of events at the will of the user, it is assembled only for 
display and transmission.
