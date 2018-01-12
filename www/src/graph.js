import * as d3 from 'd3'
import graphdata from '../graph.json';
import './graph.css';


import React from 'react';
import ReactDOM from 'react-dom';
import { observable, computed } from 'mobx';
import { observer } from "mobx-react";

// Cheers https://bl.ocks.org/puzzler10/4438752bb93f45dc5ad5214efaa12e4a

// https://bl.ocks.org/mbostock/01ab2e85e8727d6529d20391c0fd9a16
// Fisheye https://bost.ocks.org/mike/fisheye/

// Use donut chart for import visualisation https://bl.ocks.org/mbostock/3887193
// Use circular segment for contribution viz https://bl.ocks.org/mbostock/3422480

class GraphStore {
  @observable selectionId = null;
  @observable hoverId = null;

  @computed get selectionInfo() {
    let key = this.hoverId || this.selectionId;
    return graphdata.nodesLookup[key]
  }

  mouseoverNode(id) {
    this.hoverId = id;
  }

  mouseoutNode(id) {
    this.hoverId = null;
  }

  toggleSelection(id) {
    if(this.selectionId) this.selectionId = null;
    else this.selectionId = id;
  }

  clearSelection(id) {
    this.selectionId = null;
  }
}

const graphStore = new GraphStore()

@observer
class InfoView extends React.Component {
  render() {
    const store = this.props.store;
    return store.selectionInfo ? <SelectionInfo sel={store.selectionInfo}/> : null;
  }
}

const SelectionInfo = (props) => {
  return <div>
    {graphdata.nodeTypes[props.sel.variant]} <br/>
    <b>{props.sel.label}</b>
  </div>
}

document.addEventListener('DOMContentLoaded', function() {
  ReactDOM.render(
    <InfoView store={graphStore}/>,
    document.getElementById('react-mount')
  );


  var svg = d3.select("#graph").append('svg')
  var width = 900;
  var height = 700;
  svg.attr('width', '100%').attr('height', '100%')

  svg.append("defs").selectAll("marker")
    .data(["end"])
  .enter().append("svg:marker")
    .attr("id", String)
    .attr("viewBox", "0 -5 10 10")
    .attr("refX", 15)
    .attr("refY", -1.5)
    .attr("markerWidth", 6)
    .attr("markerHeight", 6)
    .attr("orient", "auto")
  .append("svg:path")
    .attr("d", "M0,-5L10,0L0,5");

  var nodeColor = d3.scaleOrdinal(d3.schemeCategory20);
  

  // Setup simulation
  // ----------------

  var simulation = d3.forceSimulation()
                                
  var link_force =  d3.forceLink(graphdata.edges)
                      .id(function(d) { return d.id; })
                      .strength(2)
                      .distance(10)
  
  var charge_force = d3.forceManyBody()
                      //  .strength(-200)


  simulation
      .nodes(graphdata.nodes)
      .force("links", link_force)
      .force('charge', charge_force)
      .force('center', d3.forceCenter(width / 2, height / 2))
      .force('collision', d3.forceCollide().radius(function(d) {
        return 15 + radius(d)
      }))
   ;
  
          
  //add tick instructions: 
  simulation.on("tick", tickActions );
  
  //add encompassing group for the zoom 
  var g = svg.append("g")
             .attr("class", "everything");

  //draw lines for the links 
  var link = g.append("g")
              .attr("class", "links")
              .selectAll("path")
              .data(graphdata.edges)
              .enter()
              .append("path")
              .attr("class", "link")
              .attr("marker-end", "url(#end)");
  

  //draw circles for the nodes 
  var node = g.append("g")
          .attr("class", "nodes") 
          .selectAll("g")
          .data(graphdata.nodes)
          .enter()
          .append('g')
          .on("mouseover", function(d){
            graphStore.mouseoverNode(d.id)
          })
          .on("mouseout", function(d){
            graphStore.mouseoutNode(d.id)
          })
          .on('click', function(d) {
            graphStore.toggleSelection(d.id)
          })
  
function radius(d) {
  const cradius = 18; 
  return cradius * d.rank;
}

  let circle = node.append("circle")
              .attr("r", function(d) {
                return radius(d)
              })
              .attr("fill", function(d) { return nodeColor(d.variant); })
  
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

  //   link.attr("d", function(d) {
  //     var dx = d.target.x - d.source.x,
  //         dy = d.target.y - d.source.y,
  //         dr = Math.sqrt(dx * dx + dy * dy);
  //     return "M" + 
  //         d.source.x + "," + 
  //         d.source.y + "A" + 
  //         dr + "," + dr + " 0 0,1 " + 
  //         d.target.x + "," + 
  //         d.target.y;
  // });
  link.attr("d", function(d) {
    // Total difference in x and y from source to target
    let diffX = d.target.x - d.source.x;
    let diffY = d.target.y - d.source.y;

    // Length of path from center of source node to center of target node
    let pathLength = Math.sqrt((diffX * diffX) + (diffY * diffY));

    // x and y distances from center to outside edge of target node

    let offsetX = (diffX * radius(d.target)) / pathLength;
    let offsetY = (diffY * radius(d.target)) / pathLength;

    return "M" + d.source.x + "," + d.source.y + "L" + (d.target.x - offsetX) + "," + (d.target.y - offsetY);
});
  
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
