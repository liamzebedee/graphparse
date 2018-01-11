import * as d3 from 'd3'
import graphdata from '../graph.json';
import './graph.css';

// Cheers https://bl.ocks.org/puzzler10/4438752bb93f45dc5ad5214efaa12e4a

document.addEventListener('DOMContentLoaded', function() {
  var svg = d3.select("#graph").append('svg')
  var width = 800;
  var height = 500;
  svg.attr('width', '100%').attr('height', '100%')

  
  

  // Setup simulation
  // ----------------

  var simulation = d3.forceSimulation()
            .nodes(graphdata.nodes);
                                
  var link_force =  d3.forceLink(graphdata.edges)
                      .id(function(d) { return d.id; });            
  
  var charge_force = d3.forceManyBody().strength(-400)

                        
  simulation
      .force("links", link_force)
      .force('charge', charge_force)
      .force('center', d3.forceCenter(width / 2, height / 2))
      .force('collision', d3.forceCollide().radius(function(d) {
        return 5 + radius(d)
      }))
   ;
  
          
  //add tick instructions: 
  simulation.on("tick", tickActions );
  
  //add encompassing group for the zoom 
  var g = svg.append("g")
             .attr("class", "everything");
  

  // var info = g.append('text')
  
  //draw lines for the links 
  var link = g.append("g")
              .attr("class", "links")
              .selectAll("line")
              .data(graphdata.edges)
              .enter().append("line")
              .attr("stroke-width", 2)
              .style("stroke", linkColour);
  
  //draw circles for the nodes 
  var node = g.append("g")
          .attr("class", "nodes") 
          .selectAll("g")
          .data(graphdata.nodes)
          .enter()
          .append('g')
          // .on("click", function(d){

          // })
  
function radius(d) {
  const cradius = 18; 
  return cradius * d.rank;
}

  let circle = node.append("circle")
              .attr("r", radius)
  
  let label = node.append('text')
              .text(function(d) { return d.label })
              .style("font-size", function(d) {
                let size =  Math.min(2 * radius(d), (2 * radius(d) - 8) / this.getComputedTextLength() * 16);
                return `${size}px`; 
              })
              .attr("dy", ".35em")

  //add drag capabilities  
  var drag_handler = d3.drag()
    .on("start", drag_start)
    .on("drag", drag_drag)
    .on("end", drag_end);	
    
  drag_handler(circle);
  
  
  //add zoom capabilities 
  var zoom_handler = d3.zoom()
      .on("zoom", zoom_actions);
  
  zoom_handler(svg);     
  
  /** Functions **/
  
  //Function to choose what color circle we have
  //Let's return blue for males and red for females
  function circleColour(d){
    return "blue";
  }
  
  //Function to choose the line colour and thickness 
  //If the link type is "A" return green 
  //If the link type is "E" return red 
  function linkColour(d){
    return "green";
  }
  
  //Drag functions 
  //d is the node 
  function drag_start(d) {
   if (!d3.event.active) simulation.alphaTarget(0.3).restart();
      d.fx = d.x;
      d.fy = d.y;
  }
  
  //make sure you can't drag the circle outside the box
  function drag_drag(d) {
    d.fx = d3.event.x;
    d.fy = d3.event.y;
  }
  
  function drag_end(d) {
    if (!d3.event.active) simulation.alphaTarget(0);
    d.fx = null;
    d.fy = null;
  }
  
  //Zoom functions 
  function zoom_actions(){
      g.attr("transform", d3.event.transform)
  }
  

  function tickActions() {
      // update circle positions each tick of the simulation 
         node
          .attr('transform', (d) => `translate(${d.x} ${d.y})`)
          
      // update link positions 
      link
          .attr("x1", function(d) { return d.source.x; })
          .attr("y1", function(d) { return d.source.y; })
          .attr("x2", function(d) { return d.target.x; })
          .attr("y2", function(d) { return d.target.y; });
  } 



});
