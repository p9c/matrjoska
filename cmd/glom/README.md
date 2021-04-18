# glom

#### Plan 9 Crypto's Answer to Acme

One of the most notable and visible elements of Plan 9 from Bell Labs operating system is the editor, 'Acme', which
combines the functionalities of a code editor with a terminal.

Glom is not an attempt to copy Acme as we consider it to be a cut-down copy of Oberon 2's GUI/OS anyway, and yes that
hints at the kind of seamless integration between the language, and the system that it will provide.

## Glom is not a text editor, it is not an IDE, and it is not a shell, it is all of those things in one.

Commands can be:

- executed one by one, with all declared variables in scope up to the session boundary

- declarations of symbols can be referred to by later statements by content hash or a symbol at the end of a path, with
  an automatic session-based global scope

- statements previously declared can be grouped into a sequence, the variables assigned to be parameters, return values,
  or internal values, tagged with versions

- the version at the HEAD of a path is automatically stored in a statement that invokes it, and on the receiving side,
  this version hash is carried with the message, and the function then retrieved to handle it, if not already cached.

- values that are not exported from the internal set are automatically exported as return values in a refactored child
  fork of the function

- locally defined symbols and packages can be automatically exported to a path at a URL the user has permission to run a
  Glom Cache, which can be imported then by other users, forked and combined into new packages

- simply by adding a main function to a package, Glom will replicate the currently focused symbol to a forked package
  that is now an entrypoint and whose parameters become an API

- comments can include both text.Template or Printf formatting substrings, including embedding function/API calls, which
  are collated into a function

- documentation can be generated from these automatically, and a filter set defined to output these documentation/logs
  as a structured log, tied to the executing code via content hashes

- with a full log available, the execution be replayed, the values automatically mocked in at the given times they are
  received from outside the package.

- the command history is a directed acyclic graph that forks on modification, that is completely addressable via hash
  chains, can be tied to paths and symbols and (subjective) timestamps
  
