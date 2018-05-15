import App from './app';
import { history, store } from './store';
import ReactDOM from 'react-dom'
import React from 'react'
import { AppContainer } from 'react-hot-loader';

const rootEl = document.getElementById('root');

const render = Component => {
  ReactDOM.render(
    <AppContainer>
      <Component store={store} history={history} />
    </AppContainer>,
    rootEl,
  )
}

render(App)

// if (module.hot) {
//   module.hot.accept([
//     './app',
//   ], () => {
//     render(App)
//   });
// }