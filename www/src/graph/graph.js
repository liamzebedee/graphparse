import React from 'react';
// import 'script-loader!../vendor/d3.v4.min.js';
import classNames from 'classnames';
import { connect } from 'react-redux'
import shortcut from 'keyboard-shortcut';
import copy from 'copy-to-clipboard';
import _ from 'underscore';
import * as d3 from 'd3v4';

import {
    hoverNode,
    clickNode,
    setGrabbing
} from './actions'
import {
    hexToRgb
} from '../util'
import nodeColor from './colours';

import TypesOverview from './types-overview';
import {

} from './graph-logic';

import './graph.css';

import {
    // preFilterNodesAndEdges,
    // generateLayout,
    GraphLogic
} from './graph-logic';

const logic = new GraphLogic();

class D3Graph extends React.Component {
    constructor() {
        super()
        this.graphDOT = null;
    }

    state = {
        grabbing: false,

        zoom: {
            transform: {
                x: 0,
                y: 0,
                z: 0
            }
        },

        nodes: [],
        edges: [],
        graphDOT: "",
    }

    componentDidMount() {
        this.addZoom()

        shortcut('ctrl c', {}, () => {
            copy(this.state.graphDOT);
        })
    }

    addZoom = () => {
        var zoomHandler = d3.zoom()
        .on("zoom", () => {
            this.setState({ zoom: d3.event.transform })
        });
  
        zoomHandler(d3.select(this.svg));
    }

    static getDerivedStateFromProps(nextProps, prevState) {
        let { 
            nodes, edges,
            currentNode,
            selection,

            showDefinitions,
            maxDepth
        } = nextProps;

        logic.refresh(nodes, edges, currentNode, selection, maxDepth);

        let layout = logic.getLayout();

        return {
            ...layout
        }
    }

    render() {
        let zoom = this.state.zoom;
        let { uiView } = this.props;

        return <svg
                ref={(ref) => this.svg = ref}
                >
            <defs>
                <filter id="shadow" x="0" y="0" width="200%" height="200%">
                <feOffset result="offOut" in="SourceAlpha" dx="10" dy="10" />
                <feGaussianBlur result="blurOut" in="offOut" stdDeviation="3" />
                <feBlend in="SourceGraphic" in2="blurOut" mode="normal" />
                </filter>
            </defs>

            <g 
                style={{
                    transform: `translate3d(${zoom.x}px, ${zoom.y}px, 0px) scale(${zoom.k})`
                }} 
                ref={ref => this.everything = ref}
                >

                { uiView == "show relationships" ? [
                    <g>
                    {this.state.nodes.map(node => {
                        return <Node 
                            key={node.id} clickNode={this.props.clickNode} 
                            {...node}/>
                    })}
                    </g>,
    
                    <g class='edges'>
                    {this.state.edges.map((edge, i) => {
                        return <Edge key={edge.id} {...edge}/>
                    })}
                    </g>
                ] : <TypesOverview/> }

            </g>
        </svg>
    }
}

const Node = ({ id, interesting, layout, variant, label, clickNode }) => {
    let { cx, cy, rx, ry } = layout;
    
    return <g 
        class='node'
        onMouseOver={() => hoverNode(id)}
        onClick={() => clickNode(id)}
        transform={`translate(${cx}, ${cy})`}
        >
        <ellipse 
            stroke='#000000'
            rx={rx}
            ry={ry}
            fill={nodeColor(variant)}
            class={classNames({
                'interesting': interesting
            })}
        >
        </ellipse>
        <text 
            textAnchor='middle'
            x="0" y="0" alignment-baseline="middle" font-size="12" stroke-width="0" stroke="#000" text-anchor="middle"
            >
            {label}
        </text>
    </g>
}


const edgeVariantStr = (variant) => {
    switch(variant) {
        case 0: return 'use';
        case 1: return 'def';
    }
}
const Edge = (edge) => {
    let layout = edge.layout;
    let { points, arrowPts } = layout;

    let computeD = () => {
        return points.map((point, i) => {
            if(i == 0) return `M${point.join(',')}C`;
            return `${point.join(',')} `;
        }).join('')
    }

    return <g class={`${edgeVariantStr(edge.variant)}`}>
        <path 
            
            d={computeD()}/>
        <polygon
            
            points={`${arrowPts.join(' ')} ${arrowPts[0]}`}/>
    </g>
}


// function zoomToBoundingBox(node) {
//     var bounds = path.bounds(d),
//       dx = bounds[1][0] - bounds[0][0],
//       dy = bounds[1][1] - bounds[0][1],
//       x = (bounds[0][0] + bounds[1][0]) / 2,
//       y = (bounds[0][1] + bounds[1][1]) / 2,
//       scale = Math.max(1, Math.min(8, 0.9 / Math.max(dx / width, dy / height))),
//       translate = [width / 2 - scale * x, height / 2 - scale * y];

//     let transform = d3.zoomIdentity.translate(translate[0],translate[1]).scale(scale);
//     return transform;
// }


const mapStateToProps = state => {
    return {
        ...state.graph,
    }
}

const mapDispatchToProps = dispatch => {
    return {
        clickNode: id => {
            dispatch(clickNode(id))
        },
        hoverNode: id => {
            dispatch(hoverNode(id))
        },
        setGrabbing: (grabbing) => dispatch(setGrabbing(grabbing)),
    }
}

 
export default connect(
    mapStateToProps,
    mapDispatchToProps
)(D3Graph);