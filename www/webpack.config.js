const path = require('path');

module.exports = {
  context: path.join(__dirname, 'src'),
  entry: {
    ast: './ast',
    graph: './graph'
  },
  output: {
    path: path.join(__dirname, 'www'),
    filename: '[name].bundle.js',
  },
  module: {
    rules: [
      {
        test: /\.css$/,
        use: [ 'style-loader', 'css-loader' ]
      },
      {
        test: /\.js$/,
        exclude: /(node_modules|vendor)/,
        use: ['babel-loader'],
      },
    ],
  },
};
