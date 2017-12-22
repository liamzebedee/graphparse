PageRank AST
============

![generated AST graph - 20 Nov 2017](https://i.imgur.com/SgXxWeX.png)

WIP implementation of PageRank on the Abstract Syntax Trees of Go source code.

 - parse and 'link' Go source
 - walk AST and extract relationships between nodes (subject to design experimentation)
 - build a new graph from this
 - run PageRank on this graph
 - convert the graph to .dot GRAPHVIZ format, with node sizes normalised according to their importance
 - visualise in browser

## Try
 - `./run.sh` to generate .dot graph file.
 - `cd www && npm run start` to view the graph.
 - Bonus: tool for quickly getting details on the AST representation of source code