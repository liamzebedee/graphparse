
space - toggle node visibility

change options in sidebar -> affects currently selected node

click node - select
click top half - show ins
click bottom half - show outs

click graph and no node - clear selection


http://rich-iannone.github.io/DiagrammeR/graphviz_and_mermaid.html
https://github.com/noflo/klayjs-noflo
https://rtsys.informatik.uni-kiel.de/confluence/display/KIELER/KLay+Layered

http://typefox.io/

https://github.com/OpenKieler/elkjs
https://github.com/noflo/klayjs-noflo

https://rtsys.informatik.uni-kiel.de/confluence/display/KIELER/KLay+Layered

https://ptolemy.berkeley.edu/ptolemyII/index.htm#4
https://ptolemy.berkeley.edu/ptolemyII/ptIIlatest/ptII/ptolemy/domains/modal/doc/index.htm

https://github.com/OpenKieler/debukviz

What is the pain I'm trying to solve? 

https://rtsys.informatik.uni-kiel.de/elklive/elkgraph.html
https://github.com/eclipse/elk/pull/106#issuecomment-268772479
https://www.eclipse.org/elk/documentation/tooldevelopers/graphdatastructure/jsonformat.html



1) visualise graphparse pre and post-redesign, to illuminate how much you could've estimated the cost.
2) try and understand btcd


try and put as much information in there as possible as a design sprint


empathise - hard to see codebase relationships
define - understanding code from a birds-eye view is not possible really
ideate - 
prototype - 
validate




1. Write down two use cases
2. Iterate and explore until I solve the interface issue for those two cases.
3. Then distribute to 3 friends to get feedback / ideas


Include only useful filters
be able to explore the btcd codebase's use of witness (up and down the hierarchy)
don't draw links back to originating node if it's a call within itself? 

I want to see how Graph is related up this fucking call stack
Is it that difficult? 
you have a relatio

https://datavizcatalogue.com/methods/sunburst_diagram.html
https://datavizcatalogue.com/methods/network_diagram.html
https://www.visualcinnamon.com/2015/08/stretched-chord.html


build initial context for graph
calc stats of + filter graph
build layout for graph
render graph


webapi:
    graph
        endpoints...
    



load graph from repo
repo = 




Make selecting a node independent from clicking
Make clicking a default action to expand
Make parent filter apply to children (makes sense?)
Fix the error with cx (when layout is null)



find all files that use this method
ins:Method depth:3


https://datavizcatalogue.com/methods/sunburst_diagram.html




Build UX for refactoring all the global variables of Graph
    User story
        I want to see all of the variables I must encapsulate that interact with Graph
    tech
        node = find(Graph)
        node.uses => map 'use
            switch(node.type)
            case func:
                map func.uses
            case var:
                return true if var.def is child of file

use GraphQL?



there are defs and uses
issue- how will we manage having multiple nodes in the same place in the layout?
this is a place where understanding the schema of what I'm building would help so I could have greater working memory

(nodes, edges) -> <Graph/> ->  
generateGraphDOT ->
    node1 -> node2
    node2 -> node2_body
    edge1

    subgraph node2_body {
        node1;

        node1_clone -> node1;
    }
generateLayout ->
    nodesLayout
    edgesLayout
getLayout -> 
    layout.nodes.map(<Node/>)
    layout.edges.map(<Edge>)


we use graphdot purely for dot engine layout
layout is then used to generate various components of the UI
but we have to make sure that the lay



https://dreampuf.github.io/GraphvizOnline/


digraph structs {
node [shape=record];
    struct1 [shape=record,label="<f0> left|<f1> middle|<f2> right"];
    struct2 [shape=record,label="<f0> one|<f1> two"];
    struct3 [shape=record,label="hello\nworld |{ b |{c|<here> d|e}| f}| g | h"];
    struct1:f1 -> struct2:f0;
    struct1:f2 -> struct3:here;
}

tree
digraph g {
node [shape = record,height=.1];
node0[label = "<f0> |<f1> G|<f2> "];
node1[label = "<f0> |<f1> E|<f2> "];
node2[label = "<f0> |<f1> B|<f2> "];
node3[label = "<f0> |<f1> F|<f2> "];
node4[label = "<f0> |<f1> R|<f2> "];
node5[label = "<f0> |<f1> H|<f2> "];
node6[label = "<f0> |<f1> Y|<f2> "];
node7[label = "<f0> |<f1> A|<f2> "];
node8[label = "<f0> |<f1> C|<f2> "];
"node0":f2 -> "node4":f1;
"node0":f0 -> "node1":f1;
"node1":f0 -> "node2":f1;
"node1":f2 -> "node3":f1;
"node2":f2 -> "node8":f1;
"node2":f0 -> "node7":f1;
"node4":f2 -> "node6":f1;
"node4":f0 -> "node5":f1;
}


digraph G {
size="6,6";

	a -> b -> c;

	subgraph cluster0 {
		x0 -> y0;
		x0 -> z0;
	}

	

	subgraph cluster2 {
		x2 -> y2;
		x2 -> z2;
	}

	a -> x0;
	b -> x1;
	b -> x2;
	a -> z2;
	c -> z1;
}







idea-
show links to other nodes as you see fit
simply build iterative understanding for now


https://graphviz.gitlab.io/_pages/Gallery/directed/cluster.html



what do I want to do? 
run a bunch of experiments on different designs
and then choose the best




simple mode for editing diagram:
    cut edges / children

Think about the best way to do this
    query language a la sql
    path from Graph -> usages
Edit node view:
    hidden
    only show certain types
Show better 'packed' view for function body's





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




//Fix bug.
//Make the visualisation make sense
Deploy for other people.


//add more debug information to node
take existing colours and make them slightly better and static.



center the frame



// idea: actually model it like webpages
//       looking for a piece of code? use current scope as starting page
//       use the type system to autofill the vars? 

// other thing with VR:
// need a visual shape-based constraint/design language