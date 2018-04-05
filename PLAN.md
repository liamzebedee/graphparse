Load codebase from different URL/path

 - Clone($CODEBASE)
 - graph = GraphparseEngine($CODEBASE)
 - ServeAPI(graph)

Web client
    Load for repo


What's the biggest pain points in the process?
    Adjusting visualisation and reloading - have to reload page everytime viz code is changed, or parser is changed and needs to be vizzed
    Adjusting parsing code and reloading - have to restart parser everytime make a change to the engine
    Changing various aspects of parsing on a small scale - have to scroll through text to inspect parser

Hot Module Replacement
Unit tests for different elements
Abstract engine from parse tree




Pain point
    How to simply identify a node throughout traversing the AST?

Could use universal AST -
    HYPERFOCUS
    Very different techincal requirements.
    Solve a problem, don't build a platform




Goal:
    build a good name graph for subnet.
        ^^first hyperfocussed pain point



Integration or unit tests?
Integration I guess.




Pain: how to understand subnet quickly


The real pain? 
Being able to refactor well and see relationships



okay now I have to refactor loadGraphVizLayout
