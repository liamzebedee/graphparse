const history = require('connect-history-api-fallback');
const convert = require('koa-connect');

module.exports = {
  reload: false,
  hot: false,

  add: (app, middleware, options) => {
    const historyOptions = {};
    app.use(convert(history(historyOptions)));
  }
};