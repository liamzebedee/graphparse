import React from 'react';
import { connect } from 'react-redux'
import _ from 'underscore';

import WatchFilledIcon from '@atlaskit/icon/glyph/watch-filled';
import MediaServicesPreselectedIcon from '@atlaskit/icon/glyph/media-services/preselected';
import BitbucketBranchesIcon from '@atlaskit/icon/glyph/bitbucket/branches';

import './controls.css';

import {

} from '../actions';

import {
    ClickActions
} from '../types';

import ClickControl from './click-control';

const Controls = ({  }) => {
    return <div>
        <ClickControl clickAction={ClickActions.select} keybinding={'1'}>
            <MediaServicesPreselectedIcon/> Select
        </ClickControl>
        <ClickControl clickAction={ClickActions.relationships} keybinding={'2'} >
            <BitbucketBranchesIcon/> Show relationships
        </ClickControl>
        <ClickControl clickAction={ClickActions.visibility} keybinding={'3'}>
            <WatchFilledIcon/> Hide/Show
        </ClickControl>
    </div>
}

const mapStateToProps = state => {
    return {

    }
}

const mapDispatchToProps = dispatch => {
    return {
        
    }
}
 
export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Controls);