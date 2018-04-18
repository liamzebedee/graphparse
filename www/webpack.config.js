const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack');

module.exports = {
  mode: 'development',
  target: 'web',

  context: path.join(__dirname, 'src'),

  entry: {
    // ast: './ast',
    ast: ['./ast', 'webpack/hot/only-dev-server', 'webpack-dev-server/client?http://0.0.0.0:3001'],
    graph: ['./graph', 'webpack/hot/only-dev-server', 'webpack-dev-server/client?http://0.0.0.0:3001'],
  },

  output: {
    path: path.join(__dirname, 'dist'),
    filename: '[name].bundle.js',
    
    // devtoolModuleFilenameTemplate: '[absolute-resource-path]',
    // devtoolFallbackModuleFilenameTemplate: '[absolute-resource-path]?[hash]'
  },

  // devtool: "inline-cheap-module-source-map",
  devtool: "inline-source-map",
  
  plugins: [
    new HtmlWebpackPlugin({
      title: 'APP'
    }),
    new webpack.NamedModulesPlugin(),
    new webpack.HotModuleReplacementPlugin(),
    new webpack.SourceMapDevToolPlugin({
      filename: "[file].map"
    }),
    new webpack.DefinePlugin({
      "process.env": { 
         NODE_ENV: JSON.stringify("production") 
       }
    })
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
      }
    ],
  },
};
