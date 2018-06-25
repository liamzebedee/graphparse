import React from 'react';
import classNames from 'classnames';
import { connect } from 'react-redux'
import shortcut from 'keyboard-shortcut';
import copy from 'copy-to-clipboard';
import _ from 'underscore';
import * as d3 from 'd3v4';

import {
    hoverNode,
    clickNode,
    setGrabbing,
    clearSelection,
    toggleNodeTypeFilter,
} from './actions'
import {
    hexToRgb
} from '../util'
import nodeColor from './colours';

import TypesOverview from './types-overview';

const graphCSS = require('!!raw-loader!./graph.css');
import './graph.css';

import Worker from './graph-logic.worker';
const worker = new Worker();

import Blanket from '@atlaskit/blanket';

class D3Graph extends React.Component {
    state = {
        grabbing: false,

        zoom: {
            transform: {
                x: 0,
                y: 0,
                z: 0
            }
        },

        shiftKey: false,
    }

    constructor() {
        super()
        this.graphDOT = null;
    }


    componentDidMount() {
        this.addZoom();

        shortcut('ctrl c', {}, () => {
            copy(this.state.graphDOT);
            // copy(this.svg.outerHTML)
        });

        [1,2,3,4].map((num) => {
            shortcut(`${num}`, {}, () => {
                this.props.toggleNodeTypeFilter(num)
            })
        })

        document.addEventListener("keydown", this.onKeyDown);
        document.addEventListener("keyup", this.onKeyUp);
    }

    componentWillUnmount(){
        document.removeEventListener("keydown", this.onKeyDown, false);
        document.removeEventListener("keyup", this.onKeyUp, false);
    }

    addZoom = () => {
        var zoomHandler = d3.zoom()
        .on("zoom", () => {
            this.setState({ zoom: d3.event.transform })
        });
  
        zoomHandler(d3.select(this.svg));
    }

    onKeyDown = (ev) => {
        this.setState({ shiftKey: ev.shiftKey })
    }

    onKeyUp = (ev) => {
        this.setState({ shiftKey: ev.shiftKey })
    }

    render() {
        let zoom = this.state.zoom;
        let { uiView, clearSelection, clickNode, nodes, edges, generating } = this.props;

        return <div styleName='graph-ctn' onKeyDown={this.onKeyDown} onKeyUp={this.onKeyUp}>
            <Blanket isTinted={generating} canClickThrough={!generating}/>

            <svg styleName='svg'
            ref={(ref) => this.svg = ref}
            onClick={clearSelection}
            >
            <defs>
                <filter id="shadow" x="0" y="0" width="200%" height="200%">
                <feOffset result="offOut" in="SourceAlpha" dx="10" dy="10" />
                <feGaussianBlur result="blurOut" in="offOut" stdDeviation="3" />
                <feBlend in="SourceGraphic" in2="blurOut" mode="normal" />
                </filter>

                <style type="text/css">{`<![CDATA[
                ${graphCSS}
                ]]>`}</style>
            </defs>

            <g 
                style={{
                    transform: `translate3d(${zoom.x}px, ${zoom.y}px, 0px) scale(${zoom.k})`
                }} 
                ref={ref => this.everything = ref}
                >

                <g>
                {nodes.map(node => {
                    // return null;
                    return <Node 
                        key={node.id} 
                        clickNode={clickNode} 
                        {...node}/>
                })}
                </g>

                <g styleName='edges'>
                {edges.map((edge, i) => {
                    // return null;
                    return <Edge key={edge.id} {...edge}/>
                })}
                </g>

            </g>
        </svg>
        </div>;
    }
}

const Node = ({ id, interesting, layout, variant, label, clickNode, selected }) => {
    let { cx, cy, rx, ry } = layout;
    
    return <g 
        styleName='node'
        onMouseOver={() => hoverNode(id)}
        onClick={(ev) => {
            clickNode(id, ev.shiftKey);
            ev.stopPropagation();
            return false;
        }}
        transform={`translate(${cx}, ${cy})`}
        >
        <ellipse 
            stroke='#000000'
            rx={rx}
            ry={ry}
            fill={nodeColor(variant)}
            styleName={classNames({
                'interesting': interesting,
                'active': selected
            })}
        >
        </ellipse>
        <text 
            textAnchor='middle'
            x="0" y="0" alignmentBaseline="middle" fontSize="12" strokeWidth="0" stroke="#000" textAnchor="middle"
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

    return <g styleName={`${edgeVariantStr(edge.variant)}`}>
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
    let g = state.graph;

    return {
        nodes: g.layout.nodes.map(({ id, layout }) => {
            return {
                ...(_.findWhere(g.nodes, { id } )),
                layout,
            }
        }),

        edges: g.layout.edges.map(({ id, layout }) => {
            return {
                ...(_.findWhere(g.edges, { id } )),
                layout,
            }
        })
    }
}

const mapDispatchToProps = dispatch => {
    return {
        clickNode: (id, shiftKey) => {
            dispatch(clickNode(id, shiftKey))
        },
        hoverNode: id => {
            dispatch(hoverNode(id))
        },
        setGrabbing: (grabbing) => dispatch(setGrabbing(grabbing)),
        clearSelection: () => dispatch(clearSelection()),
        toggleNodeTypeFilter: (i) => dispatch(toggleNodeTypeFilter(i))
    }
}

 
export default connect(
    mapStateToProps,
    mapDispatchToProps
)(D3Graph);