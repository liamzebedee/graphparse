import React from 'react';
import 'script-loader!../vendor/d3.v4.min.js';

import graphDOT from 'raw-loader!../../graph.dot';
import graphJSON from '../../graph.json';

import Viz from 'viz.js';
import _ from 'underscore';

import {
    hoverNode,
} from './actions'

import {
    hexToRgb
} from '../util'


import { connect } from 'react-redux'


let dispatch;

class D3Graph extends React.Component {
    constructor() {
        super()
        this.graphDOT = null;
    }

    componentDidMount() {
        dispatch = this.props.dispatch;
        this.svg = d3.select(this.svgCtn).append('svg')
                     .attr('width', '100%')
                     .attr('height', '100%')
        this.g = this.svg.append("g").attr("class", "everything");
        
        this.setupSvg()
        this.renderD3()
    }

    setupSvg = () => {
        let svg = this.svg;
        let g = this.g;
        let zoom = handleZoom(svg, g)
        g.append('g').attr("class", "nodes")
        g.append('g').attr("class", "edges")
    }

    componentDidUpdate(props) {
        // focusOnPackageNode(svg, zoom, layout)
        this.renderD3()
    }

    renderD3 = () => {
        let props = this.props;

        let dot = generateGraphDOT(props.nodes, props.edges);
        if(this.graphDOT != dot) {
            console.log("Regenerating layout...")

            this.graphDOT = dot;
            this.layout = generateLayout(this.graphDOT)
        }

        // Merge layout with node info for D3
        this.layout.nodes = _.map(this.layout.nodes, obj => {
            return _.assign(
                obj, 
                // _.find(props.nodeLookup, { id: obj.id }),
                props.nodeLookup[obj.id],
                { interesting: _.contains(this.props.interested, obj.id), }
            )
        })

        this.layout.edges = this.layout.edges.map((obj, i) => {
            obj.interesting = _.contains(this.props.interested, obj.source);
            obj.id = i;
            return obj;
        })

        renderGraph(this.g, this.layout)
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
  
    zoomHandler(d3.select(svg));
    return zoomHandler;
}


let nodeColor = d3.scaleOrdinal(d3.schemeCategory20);

const toSvgPointSpace = point => [ point[0], -point[1] ];


function generateGraphDOT(nodes, edges) {
    // TODO stub out.
    return graphDOT
}

// Passes DOT to Graphviz, generates layout of nodes and edges in JSON, merges with node data to be bound to D3
function generateLayout(graphDOT) {
    let graphvizData = JSON.parse(Viz(graphDOT, { format: 'json' }));
    
    let nodes = graphvizData.objects.map(obj => {
        let pos = obj.pos.split(',').map(Number);
        let id = new Number(obj.name)

        return {
            cx: pos[0],
            cy: -pos[1],
            rx: obj._draw_[1].rect[2],
            ry: obj._draw_[1].rect[3],
            ...obj,
            id,
            // ...nodeLookup[id]
        }
    })

    let edges = graphvizData.edges.map(edge => {
        let points = edge._draw_[1].points.map(toSvgPointSpace);
        
        let { head, tail } = edge;
        function findNodeForGvid(id) {
            let obj = _.find(graphvizData.objects, obj => obj._gvid == id)
            if(!obj) throw new Error()
            return Number(obj.name)
        }

        let source = findNodeForGvid(head)
        let target = findNodeForGvid(tail)

        let arrowPts = edge._hdraw_[3].points.map(toSvgPointSpace)

        return {
            points,
            arrowPts,
            source,
            target,
        }
    })

    return {
        nodes,
        edges,
    }
}

// D3 pattern:
// select
// join
// update
// enter
// exit

function renderGraph(g, layout) {
    let nodes = g.select('.nodes').selectAll('g').data(layout.nodes, (d) => `${d.id}`)
    
    function buildNode(g) {
        g
        .classed('interesting', d => d.interesting)

        g
        .append('ellipse')
        .attr('stroke', "#000000")
        .attr('cx', d => d.cx)
        .attr('cy', d => d.cy)
        .attr('rx', d => d._draw_[1].rect[2])
        .attr('ry', d => d._draw_[1].rect[3])
        .attr("fill", (d) => nodeColor(d.variant))
        .classed('interesting', d => d.interesting)

        g
        .append('text')
        .attr('text-anchor', 'middle')
        .attr('x', d => d._ldraw_[2].pt[0])
        .attr('y', d => -d._ldraw_[2].pt[1])
        .text(d => d.label)

        g
        .on("mouseover", d => {
            // hoverNode(d.id)
            dispatch(hoverNode(d.id))
        })
        // .on("mouseout", function(d){
        //     unhoverNode(d.id)
        // })
        // .on('click', function(d) {
        //     selectNode(d.id)
        // })

        return g
    }
    function updateNode(g) {
        g
        .classed('interesting', d => d.interesting)
    }
    
    updateNode(buildNode(nodes.enter().append('g')).merge(nodes))
    nodes.exit().remove()

    let edges = g.select('.edges').selectAll('g').data(layout.edges, (d) => `${d.id}`)
    // ===== update selection
    // let edges = g.select('.edges').selectAll('g').data(layout.edges)

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

        return g
    }
    function updateEdge(g) {
        g
        .classed('interesting', d => d.interesting)
    }
    edges.exit().remove()
    updateEdge(buildEdge(edges.enter().append('g')).merge(edges))
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
    // zoom.translateTo(svg, node.cx / 2, node.cy / 2)
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




const mapStateToProps = state => {
    return {
        ...state.graph,
    }
}

const mapDispatchToProps = dispatch => {
    return {
        // HACK for d3.
        dispatch, 

        // hoverNode: id => {
        //     dispatch(hoverNode(id))
        // }
    }
}

const D3GraphCtn = connect(
  mapStateToProps,
  mapDispatchToProps
)(D3Graph)
 
export default D3GraphCtn;