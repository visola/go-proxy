const CopyPlugin = require('copy-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const path = require('path');
const VueLoaderPlugin = require('vue-loader/lib/plugin')

module.exports = {
  devtool: 'inline-source-map', 
  entry: './src/js/main.js',

  module: {
    rules: [
      { test: /\.vue$/, loader: 'vue-loader' },
    ]
  },

  output: {
    filename: 'main.js',
    path: path.resolve(__dirname, 'dist'),
  },

  plugins: [
    new CopyPlugin([
      { from: 'src/css', to: path.resolve(__dirname, 'dist/css') },
    ]),

    new HtmlWebpackPlugin({
      template: 'src/html/index.html',
      title: 'go-proxy admin',
    }),

    new VueLoaderPlugin(),
  ]
};
