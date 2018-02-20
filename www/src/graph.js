import './graph.css';
import * as d3old from 'd3'
import 'script-loader!./vendor/d3.v4.min.js';
import 'script-loader!./vendor/viz-lite.js';
import 'script-loader!./vendor/d3-graphviz.min.js';
import React from 'react';
import ReactDOM from 'react-dom';
import mobx, { observable, computed } from 'mobx';
import { observer } from "mobx-react";
import _ from 'underscore';
import GraphStore from './graph/store';
var graphStore = new GraphStore()

import graphdataFromFile from '../graph.json';
var graphdata = graphdataFromFile;
import graphDot from 'raw-loader!../graph.dot';

// Cheers https://bl.ocks.org/puzzler10/4438752bb93f45dc5ad5214efaa12e4a

// https://bl.ocks.org/mbostock/01ab2e85e8727d6529d20391c0fd9a16
// Fisheye https://bost.ocks.org/mike/fisheye/

// Use donut chart for import visualisation https://bl.ocks.org/mbostock/3887193
// Use circular segment for contribution viz https://bl.ocks.org/mbostock/3422480

// Zoomable sub circles https://bl.ocks.org/mbostock/7607535




var svg, simulation;
var linkForce;
var link, node;
var tickActions;

var d3RenderData;

var nodeColor = d3.scaleOrdinal(d3.schemeCategory20);

const LoadingWrap = (props) => {
  if(props.data != null) return props.children;
  else return "";
}

@observer
class InfoView extends React.Component {
  searchNodes = (ev) => {
    this.props.store.searchAndHighlightNodes(ev.target.value)
  }

  toggleThread = () => {
    this.props.store.getCodeThread()
  }

  reloadGraph = () => {
    this.props.store.loadGraph()
  }

  toggleNodeTypeFilter = (nodeTypeIdx) => {
    this.props.store.toggleNodeTypeFilter(nodeTypeIdx)
  }

  render() {
    const store = this.props.store;

    let selectedNodes = store.selectedNodes;

    return <div>
      <GraphUI 
        nodes={store.nodes} 
        edges={store.edges} 
        updateGraph={store.updateGraph}/>

      { false ? 
      <div className="info-view">
        <div className="ui-pane">
          <input type="text" placeholder="&#x1F50E; Find structs, funcs, methods..." className='search' onChange={this.searchNodes}/>
          <LoadingWrap data={store.selectionInfo}>
            <SelectionInfo sel={store.selectionInfo}/>
          </LoadingWrap>

          <SelectedNodes selectedNodes={selectedNodes}/>

          <button onClick={this.toggleThread}>Show thread</button>
          <button onClick={this.reloadGraph}>Reload</button>
        </div>

        <div className="ui-pane">
          <h3>Filters</h3>
          {graphdata.nodeTypes.map((nodeType, i) => {
            return <div onClick={() => this.toggleNodeTypeFilter(i)} key={i} style={{ 
              backgroundColor: nodeColor(i),
              padding: '0.25em',
              color: 'black',
              fontWeight: 800,
              display: 'inline-block',
              cursor: 'pointer',
              opacity: _.contains(this.props.store.nodeTypesHidden, i) ? 0.3 : 1,
            }}>{nodeType}</div>
          })}

        </div>
      </div> : null }
    </div> 
  }
}

export function renderGraph(nodes, edges) {
  renderGraphVizJs(nodes, edges)
}

class GraphUI extends React.Component {
  componentDidMount() {
    svg = d3.select(this.ui).append('svg')
      .attr('width', '100%')
      .attr('height', '100%')

    // renderGraphD3(this.props.nodes, this.props.edges)
    renderGraphVizJs(this.props.nodes, this.props.edges)
  }

  componentDidUpdate(prevProps, prevState) {
    if(this.props.updateGraph === prevProps.updateGraph) return;

    simulation.stop();

    // d3RenderData(this.props.nodes, this.props.edges)

    linkForce = d3.forceLink(this.props.edges)
                    .id(function(d) { return d.id; })
                    .strength(1)
                    .distance(10)
    
    simulation
      .nodes(this.props.nodes)
      .force("links", linkForce)
    
    d3.timeout(function() {  
      // See https://github.com/d3/d3-force/blob/master/README.md#simulation_tick
      for (var i = 0, n = Math.ceil(Math.log(simulation.alphaMin()) / Math.log(1 - simulation.alphaDecay())); i < n; ++i) {
        simulation.tick();
        tickActions()
      }
    });
  }

  updateGraph = () => {

  }

  render() {
    return <div ref={(ref) => this.ui = ref}></div>
  }
}

const SelectedNodes = (props) => {
  return <div>
    {props.selectedNodes.map((node, i) => <SelectedNode key={i} node={node}/>)}
  </div>
}

const SelectedNode = (props) => {
  return <div>
    <div style={{ 
      backgroundColor: nodeColor(props.node.variant),
      padding: '0.25em',
      color: 'black',
      fontWeight: 800
    }}>{props.node.label}</div>
  </div>
}

const SelectionInfo = (props) => {
  return <div style={{ padding:'1em' }}>
    {graphdata.nodeTypes[props.sel.variant]} <br/>
    <b>{props.sel.label}</b>
  </div>
}


document.addEventListener('DOMContentLoaded', function() {
  ReactDOM.render(
    <InfoView store={graphStore}/>,
    document.getElementById('react-mount')
  );
});




function renderGraphVizJs(nodes, edges) {
  svg
    .graphviz()
    .renderDot(graphDot)
    .totalMemory(16777216 * 2)
}

function renderGraphD3(nodes, edges) {
  var width = 900;
  var height = 700;

  svg.append("defs").selectAll("marker")
    .data(["end"])
    .enter()
    .append("svg:marker")
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

  simulation = d3.forceSimulation()
                                
  linkForce = d3.forceLink(edges)
                  .id(function(d) { return d.id; })
                  .strength(1)
                  .distance(10)
  
  var chargeForce = d3.forceManyBody()
                       .strength(-200)

  simulation
    .nodes(nodes)
    .force("links", linkForce)
    .force('charge', chargeForce)
    .force('center', d3.forceCenter(width / 2, height / 2))
    .force('collision', d3.forceCollide().radius(function(d) {
      return 1 + radius(d)
    }))

  simulation.stop();
  // simulation.on("tick", tickActions );

  d3.timeout(function() {  
    // See https://github.com/d3/d3-force/blob/master/README.md#simulation_tick
    for (var i = 0, n = Math.ceil(Math.log(simulation.alphaMin()) / Math.log(1 - simulation.alphaDecay())); i < n; ++i) {
      simulation.tick();
      tickActions()
    }
  });


  var g = svg.append("g")
    .attr("class", "everything");

  var t = d3.transition()
    .duration(350);

  function radius(d) {
    const cradius = 18; 
    return cradius * d.rank;
  }

  var circle;

  var $nodes = g.append("g")
    .attr("class", "nodes");

  var links = g.append("g")
    .attr("class", "links");

  d3RenderData = function(nodes, edges) {
    // Links
    // -----

    link = links
      .selectAll("path")
      .data(edges, (d) => `${d.source.id}${d.target.id}`)

    link.exit().remove()

    link
      .enter()
      .append("path")
      .attr("class", "link")
      .attr("marker-end", "url(#end)")
      .merge(link)


    // Node
    // ----

    node = $nodes
      .selectAll("g")
      .data(nodes, (d) => `${d.id}`)
    
    node
      .exit()
      .remove();

    node = node
      .enter()
      .append('g')
      .on("mouseover", function(d){
        graphStore.mouseoverNode(d.id)
      })
      .on("mouseout", function(d){
        graphStore.mouseoutNode(d.id)
      })
      .on('click', function(d) {
        graphStore.selectNode(d.id)
      })

    circle = node
      .append("circle")
      .attr("r", radius)
      .attr("fill", (d) => nodeColor(d.variant))

    let label = node
      .append('text')
      .text(function(d) { return d.label })
      .style("font-size", function(d) {
        let size =  Math.min(2 * radius(d), (2 * radius(d) - 8) / this.getComputedTextLength() * 16);
        return `${size}px`; 
      })
      .attr("dy", ".35em")

    node
      .merge(node)
  }

  function contrastColour(r,g,b) {
    let d = 0;
    // Counting the perceptive luminance - human eye favors green color... 
    let a = 1 - ( 0.299 * r + 0.587 * g + 0.114 * b)/255;
    if (a < 0.5)
       return "black";
    else
       return "white"
  }
  d3RenderData(nodes, edges)

  mobx.autorun(() => {
    if(graphStore.highlightedNodes == null) {
      node
        .attr('opacity', 1)
      
      circle
        .classed('selected', false)
      
      link
        .attr('opacity', 1)
      
    } else {
      node
        .attr('opacity', 0.2)
        .filter(function(d) {
          return graphStore.getNodeHighlight(d.id) != null;
        })
        .attr('opacity', function(d) {
          return 1 - graphStore.getNodeHighlight(d.id).score;
        });
      
      circle
        .classed('selected', function(d) {
          return _.has(graphStore.selectedNodeIds, d.id)
        })

      link
        .attr('opacity', 0)
    }
  })

  var zoomHandler = d3.zoom()
    // .scaleExtent([1, Infinity])
    // .translateExtent([[0, 0], [width, height]])
    // .extent([[0, 0], [width, height]])
    .on("zoom", () => {
      let transform = d3.event.transform;
      // g.attr("transform", d3.event.transform)
      g.style("transform", () => `translate3d(${transform.x}px, ${transform.y}px, 0px) scale(${transform.k})`)
    });
  
  zoomHandler(svg);


  var dragHandler = d3.drag()
    .on("start", (d) => {
      if (!d3.event.active) simulation.alphaTarget(0.3).restart();
        d.fx = d.x;
        d.fy = d.y;
    })
    .on("drag", (d) => {
      d.fx = d3.event.x;
      d.fy = d3.event.y;
    })
    .on("end", (d) => {
      if (!d3.event.active) simulation.alphaTarget(0);
      d.fx = null;
      d.fy = null;
    })

  dragHandler(node);

  tickActions = function() {
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
  
    node
      // .attr("style", function (d) { return `transform: translate3d(${d.x}px, ${d.y}px, 0px)` })
      .attr('transform', (d) => `translate(${d.x} ${d.y})`)
          
    link
      .attr("x1", function(d) { return d.source.x; })
      .attr("y1", function(d) { return d.source.y; })
      .attr("x2", function(d) { return d.target.x; })
      .attr("y2", function(d) { return d.target.y; });
  } 
}