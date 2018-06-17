import React from 'react'
import ReactDOM from 'react-dom'
import { Route } from 'react-router'
import { hot } from 'react-hot-loader'
import { ConnectedRouter } from 'react-router-redux';
import Home from './home'
import GraphUI from './graph';
import ASTView from './ast';
import { Provider } from 'react-redux';

import '@atlaskit/css-reset';

export default ({ history, store }) => <Provider store={store}>
    <ConnectedRouter history={history}>
        <div>
            <Route exact path="/" component={Home}/>
            <Route path="/repo/" component={GraphUI}/>
            {/* <Route path="/ast" component={ASTView}/> */}
        </div>
    </ConnectedRouter>
</Provider>;