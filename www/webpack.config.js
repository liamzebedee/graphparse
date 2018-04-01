const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack');

module.exports = {
  target: 'web',  

  context: path.join(__dirname, 'src'),
  entry: {
    // ast: './ast',
    ast: ['./ast', 'webpack/hot/only-dev-server', 'webpack-dev-server/client?http://0.0.0.0:3000'],
    graph: ['./graph', 'webpack/hot/only-dev-server', 'webpack-dev-server/client?http://0.0.0.0:3000'],
  },

  output: {
    path: path.join(__dirname, 'dist'),
    filename: '[name].bundle.js',
  },

  devtool: 'inline-source-map',
  
  plugins: [
    // new CleanWebpackPlugin(['dist']),
    new HtmlWebpackPlugin({
      title: 'APP'
    }),
    new webpack.NamedModulesPlugin(),
    new webpack.HotModuleReplacementPlugin(),
    new webpack.SourceMapDevToolPlugin({
      filename: "[file].map"
    }),
    new webpack.EnvironmentPlugin(['NODE_ENV', 'DEBUG'])
  ],

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
