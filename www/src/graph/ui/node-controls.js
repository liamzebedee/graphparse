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
} from '../actions';


import {
    getSelectedNode
} from '../selectors';

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
                            shownNodeTypes={selectedNode.selection.ins.shownNodeTypes}
                            depth={selectedNode.selection.ins.maxDepth}
                            changeDepth={(depth) => changeDepth(selectedNode.id, 'ins', depth)}
                            toggleFilter={(variant) => toggleFilter(selectedNode.id, 'ins', variant)}
                            />
                        </section>
                        
                        <section>
                        <h5>Out</h5>
                        <RelationControl 
                            node={selectedNode} 
                            shownNodeTypes={selectedNode.selection.ins.shownNodeTypes}
                            depth={selectedNode.selection.ins.maxDepth}
                            changeDepth={(depth) => changeDepth(selectedNode.id, 'outs', depth)}
                            toggleFilter={(variant) => toggleFilter(selectedNode.id, 'outs', variant)}/>
                        </section>
                    </div>
                </div>
            </div>
        }
    </div>
}

const RelationControl = ({ node, depth, changeDepth, toggleFilter, shownNodeTypes }) => {
    return <div>
        <Filters node={node} shownNodeTypes={shownNodeTypes} toggleFilter={toggleFilter}/>
        <Depth node={node} depth={depth} changeDepth={changeDepth}/>
    </div>
}

const mapStateToProps = state => {
    return {
        selectedNode: getSelectedNode(state.graph),
    }
}

const mapDispatchToProps = dispatch => {
    return {
        changeDepth: (id, relationships, depth) => dispatch(changeDepth(id, relationships, depth)),
        toggleFilter: (id, relationships, variant) => dispatch(toggleFilter(id, relationships, variant)),
    }
}
 
export default connect(
    mapStateToProps,
    mapDispatchToProps
)(NodeControls);