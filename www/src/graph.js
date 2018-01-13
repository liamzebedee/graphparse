import * as d3 from 'd3'
import graphdata from '../graph.json';
import './graph.css';


import React from 'react';
import ReactDOM from 'react-dom';
import mobx, { observable, computed } from 'mobx';
import { observer } from "mobx-react";
import Fuse from 'fuse.js';
import _ from 'underscore';

// Cheers https://bl.ocks.org/puzzler10/4438752bb93f45dc5ad5214efaa12e4a

// https://bl.ocks.org/mbostock/01ab2e85e8727d6529d20391c0fd9a16
// Fisheye https://bost.ocks.org/mike/fisheye/

// Use donut chart for import visualisation https://bl.ocks.org/mbostock/3887193
// Use circular segment for contribution viz https://bl.ocks.org/mbostock/3422480

// Zoomable sub circles https://bl.ocks.org/mbostock/7607535

class GraphStore {
  @observable hoverId = null;
  @observable highlightedNodes = null;
  @observable selectedNodeIds = [];

  @computed get selectionInfo() {
    let key = this.hoverId || this.selectionId;
    return graphdata.nodesLookup[key]
  }

  @computed.struct get selectedNodes() {
    return Array.from(this.selectedNodeIds).map(id => { console.log(id); return graphdata.nodesLookup[id] })
  }

  mouseoverNode(id) {
    this.hoverId = id;
  }

  mouseoutNode(id) {
    this.hoverId = null;
  }

  selectNode(id) {
    if(!_.contains(this.selectedNodeIds, id))
      this.selectedNodeIds.push(id)
  }
  // toggleSelection(id) {
    // if(this.selectionId) this.selectionId = null;
    // else this.selectionId = id;
  // }

  clearSelection(id) {
    this.selectionId = null;
  }

  searchAndHighlightNodes(query) {
    if(query == "") {
      this.highlightedNodes = null;
      return
    }
    let matches = [];

    function splitId(id) {
      // return id.split(/(?=[A-Z])/g).reverse();
      return [id]
    }

    let searchData = graphdata.nodes.map(node => {
      return {
        id: node.id, 
        label: splitId(node.label),
      }
    })

    var options = {
      shouldSort: true,
      findAllMatches: true,
      includeScore: true,
      threshold: 0.3,
      location: 0,
      distance: 100,
      maxPatternLength: 32,
      minMatchCharLength: 2,
      keys: [
        "label"
      ]
    };

    let fuse = new Fuse(searchData, options);
    let results = fuse.search(query)
    this.highlightedNodes = results.filter(el => el.score < 0.7)
  }

  getNodeHighlight(id) {
    return this.highlightedNodes.filter(el => el.item.id == id)[0];
  }
}

const graphStore = new GraphStore()

const LoadingWrap = (props) => {
  if(props.data != null) return props.children;
  else return "";
}

@observer
class InfoView extends React.Component {
  searchNodes = (ev) => {
    this.props.store.searchAndHighlightNodes(ev.target.value)
  }

  render() {
    const store = this.props.store;

    let selectedNodes = store.selectedNodes;

    return <div>
      <input type="text" placeholder="&#x1F50E; Find structs, funcs, methods..." className='search' onChange={this.searchNodes}/>
      <SelectedNodes selectedNodes={selectedNodes}/>
      <LoadingWrap data={store.selectionInfo}><SelectionInfo sel={store.selectionInfo}/></LoadingWrap>
    </div>
  }
}

const SelectedNodes = (props) => {
  return <div>
    {props.selectedNodes.map((node, i) => <SelectedNode key={i} node={node}/>)}
  </div>
}

const SelectedNode = (props) => {
  return <div>
    <div style={{ backgroundColor: nodeColor(props.node.variant) }}>{props.node.label}</div>
  </div>
}

const SelectionInfo = (props) => {
  return <div style={{ padding:'1em' }}>
    {graphdata.nodeTypes[props.sel.variant]} <br/>
    <b>{props.sel.label}</b>
  </div>
}


var nodeColor = d3.scaleOrdinal(d3.schemeCategory20);

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

  

  // Setup simulation
  // ----------------

  var simulation = d3.forceSimulation()
                                
  var link_force =  d3.forceLink(graphdata.edges)
                      .id(function(d) { return d.id; })
                      .strength(1)
                      .distance(10)
  
  var charge_force = d3.forceManyBody()
                       .strength(-200)

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
            // graphStore.toggleSelection(d.id)
            graphStore.selectNode(d.id)
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
  
  mobx.autorun(() => {
    if(graphStore.highlightedNodes == null) {
      node.attr('opacity', 1)
      link.attr('opacity', 1)
    } else {
      node
      .attr('opacity', 0.2)
      .filter(function(d) {
        return graphStore.getNodeHighlight(d.id) != null;
      })
      .attr('opacity', function(d) {
        return 1 - graphStore.getNodeHighlight(d.id).score;
      });

      link
      .attr('opacity', 0)
    }
  })

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
