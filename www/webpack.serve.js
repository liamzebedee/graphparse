const history = require('connect-history-api-fallback');
const convert = require('koa-connect');

module.exports = {
  reload: false,
  hot: true,

  // host: "0.0.0.0",
  // port: 8080,

  add: (app, middleware, options) => {
    const historyOptions = {};
    app.use(convert(history(historyOptions)));
  }
};