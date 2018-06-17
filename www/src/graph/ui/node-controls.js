import React from 'react';
import { connect } from 'react-redux'
import _ from 'underscore';

import './node-controls.css';

import Filters from './filters';
import Depth from './depth';

import { ToggleStateless } from '@atlaskit/toggle';

import {
    changeDepth,
    toggleFilter
} from './actions';

const NodeControls = ({ 
    selectedNode, 
    changeDepth 
}) => {
    return <div styleName='ui-pad'>
        <h3>Current node</h3>
        { selectedNode == null 
            ? <i>No node selected</i>
            : <div>
                <div styleName='current-node'>
                    <u>{selectedNode.label}</u>
                    {/* <div>Shown <ToggleStateless/></div> */}

                    <h4>Relations</h4>
                    <div styleName='relations'>
                        <section>
                        <h5>In</h5>
                        <RelationControl 
                            node={selectedNode} 
                            changeDepth={(depth) => changeDepth(node, 'ins', depth)}
                            toggleFilter={(variant) => toggleFilter(node, 'ins', variant)}
                            />
                        </section>
                        
                        <section>
                        <h5>Out</h5>
                        <RelationControl 
                            node={selectedNode} 
                            changeDepth={(depth) => changeDepth(node, 'outs', depth)}
                            toggleFilter={(variant) => toggleFilter(node, 'outs', variant)}/>
                        </section>
                    </div>
                </div>
            </div>
        }
    </div>
}

const RelationControl = ({ node, depth, changeDepth, toggleFilter, shownNodeTypes }) => {
    return <div>
        <Filters shownNodeTypes={shownNodeTypes} toggleFilter={toggleFilter}/>
        <Depth depth={depth} changeDepth={changeDepth}/>
    </div>
}

const mapStateToProps = state => {
    let selectedNode = null;
    if(state.graph.selectedNode) {
        selectedNode = _.find(state.graph.nodes, { id: state.graph.selectedNode });
    }

    return {
        selectedNode,
    }
}

const mapDispatchToProps = dispatch => {
    return {
        changeDepth: (node, relationships, depth) => dispatch(changeDepth(node, relationships, depth)),
        toggleFilter: (node, relationships, variant) => dispatch(toggleFilter(node, relationships, variant)),
    }
}
 
export default connect(
    mapStateToProps,
    mapDispatchToProps
)(NodeControls);