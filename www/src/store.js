
import { createStore, combineReducers, applyMiddleware, compose } from 'redux'
import createHistory from 'history/createBrowserHistory'
import thunk from 'redux-thunk';
import { 
    routerReducer, 
    routerMiddleware, 
} from 'react-router-redux'
import graph from './graph/reducers';

export const history = createHistory()
const enhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;

export const store = createStore(
    combineReducers({
        router: routerReducer,
        graph,
    }),
    enhancers(
        applyMiddleware(
            routerMiddleware(history),
            thunk,
        )
    ),
)