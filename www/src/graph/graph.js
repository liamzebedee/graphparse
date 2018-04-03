
import React from 'react';
import 'script-loader!../vendor/d3.v4.min.js';

import graphJSON from '../../graph.json';
import graphDOT from 'raw-loader!../../graph.dot';

import Viz from 'viz.js';
import _ from 'underscore';


import { observer } from 'mobx-react';

import {
    hoverNode,
    unhoverNode,
    selectNode
} from './actions'

import {
    hexToRgb
} from '../util'


@observer
export default class D3Graph extends React.Component {
    constructor() {
        super()
    }

    componentDidMount() {
        this.svg = d3.select(this.svgCtn).append('svg')
                     .attr('width', '100%')
                     .attr('height', '100%')
        this.setupSvg()
    }

    setupSvg = () => {
        let svg = this.svg;
        let g = svg.append("g").attr("class", "everything");

        let zoom = handleZoom(svg, g)
        let layout = loadGraphVizLayout(g)
        renderGraph(g, layout)
        focusOnPackageNode(svg, zoom, layout)
    }

    componentWillReceiveProps() {
        console.log(arguments)
    }

    render() {
        return <div ref={(ref) => this.svgCtn = ref}></div>
    }
}

function defineMarker(svg) {
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
}

function handleZoom(svg, g) {
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
    return zoomHandler;
}


let nodeColor = d3.scaleOrdinal(d3.schemeCategory20);

const toSvgPointSpace = point => [ point[0], -point[1] ];

function loadGraphVizLayout(g) {
    var graphvizData = JSON.parse(Viz(graphDOT, { format: 'json' }));
    
    let nodes = graphvizData.objects.map(obj => {
        let pos = obj.pos.split(',').map(Number);
        let id = new Number(obj.name)
        let info = graphJSON.nodeLookup[id + ""];
        if(!info) return;

        return {
            cx: pos[0],
            cy: -pos[1],
            rx: obj._draw_[1].rect[2],
            ry: obj._draw_[1].rect[3],
            id,
            ...obj,
            ...graphJSON.nodeLookup[id]
        }
    })

    let edges = graphvizData.edges.map(edge => {
        let points = edge._draw_[1].points.map(toSvgPointSpace);
        
        let { head, tail } = edge;
        function findNodeForGvid(id) {
            let obj = _.find(graphvizData.objects, obj => obj._gvid == id)
            if(!obj) throw new Error()
            return new Number(obj.name)
        }

        let source = findNodeForGvid(head)
        let target = findNodeForGvid(tail)

        let arrowPts = edge._hdraw_[3].points.map(toSvgPointSpace)

        return {
            points,
            source,
            target,
            arrowPts
        }
    })

    return {
        nodes,
        edges
    }
}

function renderGraph(g, layout) {
    let nodes = g.append("g").attr("class", "nodes").selectAll('g').data(layout.nodes, (d) => `${d.id}`)
    
    function buildNode(g) {
        g
        .append('ellipse')
        .attr('stroke', "#000000")
        .attr('cx', d => d.cx)
        .attr('cy', d => d.cy)
        .attr('rx', d => d._draw_[1].rect[2])
        .attr('ry', d => d._draw_[1].rect[3])
        .attr("fill", (d) => nodeColor(d.variant))

        g
        .append('text')
        .attr('text-anchor', 'middle')
        .attr('x', d => d._ldraw_[2].pt[0])
        .attr('y', d => -d._ldraw_[2].pt[1])
        .text(d => d.label)
        // .attr('color', d => {
        //     let {r,g,b} = hexToRgb(nodeColor(d.variant))
        //     return contrastColour(r,g,b)
        // })

        g
        .on("mouseover", d => {
            hoverNode(d.id)
        })
        .on("mouseout", function(d){
            unhoverNode(d.id)
        })
        .on('click', function(d) {
            selectNode(d.id)
        })
    }
    buildNode(nodes.enter().append('g'))

    let edges = g.append("g").attr("class", "edges").selectAll('g').data(layout.edges, (d) => `${d.source.id}${d.target.id}`)

    function buildEdge(g) {
        g.append('path')
        .attr('fill', 'none')
        .attr('stroke', '#000000')
        .attr('d', d => {
            let a = d.points[0].join(',')
            let b = d.points.splice(1).map(point => point.join(',')).join(' ')
            return `M${a}C${b}`;
        })

        g.append('polygon')
        .attr('fill', '#000000')
        .attr('stroke', '#000000')
        .attr('points', d => {
            return `${d.arrowPts.join(' ')} ${d.arrowPts[0]}`;
        })
    }
    buildEdge(edges.enter().append('g'))
}

function renderGraphviz(g) {
    var graphvizSvg = Viz(graphDOT, { format: 'svg' });
    g.html(graphvizSvg)
}


function focusOnPackageNode(svg, zoom, layout) {
    let node = _.find(layout.nodes, node => {
        return node.variant == graphJSON.nodeTypes.indexOf('RootPackage')
    })
    // let transform = to_bounding_box(getCenter(node.cx, node.cy), node.rx, node.ry, 0)
    // g.transition().duration(200).call(zoom.transform, transform);
    zoom.translateTo(svg, node.cx / 2, node.cy / 2)
}

function contrastColour(r,g,b) {
    let d = 0;
    // Counting the perceptive luminance - human eye favors green color... 
    let a = 1 - ( 0.299 * r + 0.587 * g + 0.114 * b)/255;
    if (a < 0.5)
        return "black";
    else
        return "white";
}

  
//     var t = d3.transition()
//       .duration(350);
  
//     function radius(d) {
//       const cradius = 18; 
//       return cradius * d.rank;
//     }
  
//     var circle;
  
//     var $nodes = g.append("g")
//       .attr("class", "nodes");
  
//     var links = g.append("g")
//       .attr("class", "links");
  
//     d3RenderData = function(nodes, edges) {
//       // Links
//       // -----
  
//       link = links
//         .selectAll("path")
//         .data(edges, (d) => `${d.source.id}${d.target.id}`)
  
//       link.exit().remove()
  
//       link
//         .enter()
//         .append("path")
//         .attr("class", "link")
//         .attr("marker-end", "url(#end)")
//         .merge(link)
  
  
//       // Node
//       // ----
  
//       node = $nodes
//         .selectAll("g")
//         .data(nodes, (d) => `${d.id}`)
      
//       node
//         .exit()
//         .remove();
  
//       node = node
//         .enter()
//         .append('g')
//         .on("mouseover", function(d){
//           graphStore.mouseoverNode(d.id)
//         })
//         .on("mouseout", function(d){
//           graphStore.mouseoutNode(d.id)
//         })
//         .on('click', function(d) {
//           graphStore.selectNode(d.id)
//         })
  
//       circle = node
//         .append("circle")
//         .attr("r", radius)
//         .attr("fill", (d) => nodeColor(d.variant))
  
//       let label = node
//         .append('text')
//         .text(function(d) { return d.label })
//         .style("font-size", function(d) {
//           let size =  Math.min(2 * radius(d), (2 * radius(d) - 8) / this.getComputedTextLength() * 16);
//           return `${size}px`; 
//         })
//         .attr("dy", ".35em")
  
//       node
//         .merge(node)
//     }
  
//     function contrastColour(r,g,b) {
//       let d = 0;
//       // Counting the perceptive luminance - human eye favors green color... 
//       let a = 1 - ( 0.299 * r + 0.587 * g + 0.114 * b)/255;
//       if (a < 0.5)
//          return "black";
//       else
//          return "white"
//     }
//     d3RenderData(nodes, edges)
  
//     mobx.autorun(() => {
//       if(graphStore.highlightedNodes == null) {
//         node
//           .attr('opacity', 1)
        
//         circle
//           .classed('selected', false)
        
//         link
//           .attr('opacity', 1)
        
//       } else {
//         node
//           .attr('opacity', 0.2)
//           .filter(function(d) {
//             return graphStore.getNodeHighlight(d.id) != null;
//           })
//           .attr('opacity', function(d) {
//             return 1 - graphStore.getNodeHighlight(d.id).score;
//           });
        
//         circle
//           .classed('selected', function(d) {
//             return _.has(graphStore.selectedNodeIds, d.id)
//           })
  
//         link
//           .attr('opacity', 0)
//       }
//     })
  
    // var zoomHandler = d3.zoom()
    //   // .scaleExtent([1, Infinity])
    //   // .translateExtent([[0, 0], [width, height]])
    //   // .extent([[0, 0], [width, height]])
    //   .on("zoom", () => {
    //     let transform = d3.event.transform;
    //     // g.attr("transform", d3.event.transform)
    //     g.style("transform", () => `translate3d(${transform.x}px, ${transform.y}px, 0px) scale(${transform.k})`)
    //   });
    
    // zoomHandler(svg);
  
  
//     var dragHandler = d3.drag()
//       .on("start", (d) => {
//         if (!d3.event.active) simulation.alphaTarget(0.3).restart();
//           d.fx = d.x;
//           d.fy = d.y;
//       })
//       .on("drag", (d) => {
//         d.fx = d3.event.x;
//         d.fy = d3.event.y;
//       })
//       .on("end", (d) => {
//         if (!d3.event.active) simulation.alphaTarget(0);
//         d.fx = null;
//         d.fy = null;
//       })
  
//     dragHandler(node);
  
//     tickActions = function() {
//       link.attr("d", function(d) {
//         // Total difference in x and y from source to target
//         let diffX = d.target.x - d.source.x;
//         let diffY = d.target.y - d.source.y;
  
//         // Length of path from center of source node to center of target node
//         let pathLength = Math.sqrt((diffX * diffX) + (diffY * diffY));
  
//         // x and y distances from center to outside edge of target node
//         let offsetX = (diffX * radius(d.target)) / pathLength;
//         let offsetY = (diffY * radius(d.target)) / pathLength;
  
//         return "M" + d.source.x + "," + d.source.y + "L" + (d.target.x - offsetX) + "," + (d.target.y - offsetY);
//       });
    
//       node
//         // .attr("style", function (d) { return `transform: translate3d(${d.x}px, ${d.y}px, 0px)` })
//         .attr('transform', (d) => `translate(${d.x} ${d.y})`)
            
//       link
//         .attr("x1", function(d) { return d.source.x; })
//         .attr("y1", function(d) { return d.source.y; })
//         .attr("x2", function(d) { return d.target.x; })
//         .attr("y2", function(d) { return d.target.y; });
//     } 
//   }


/*

function getCenter(x, y) {
    return {
        x: x / 2,
        y: y / 2
      };
}

/*
Returns a transform for center a bounding box in the browser viewport
    - W and H are the witdh and height of the window
    - w and h are the witdh and height of the bounding box
    - center cointains the coordinates of the bounding box center
    - margin defines the margin of the bounding box once zoomed
    

function to_bounding_box(center, w, h, margin) {
    const W = window.innerWidth;
    const H = window.innerHeight;
    
    var k, kh, kw, x, y;
    kw = (W - margin) / w;
    kh = (H - margin) / h;
    k = d3.min([kw, kh]);
    x = W / 2 - center.x * k;
    y = H / 2 - center.y * k;
    return d3.zoomIdentity.translate(x, y).scale(k);
};

*/