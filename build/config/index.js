'use strict';
const path = require('path');

module.exports = {
  dev: {
    // Paths
    assetsSubDirectory: 'static',
    assetsPublicPath: '/',

    // Various Dev Server settings
    host: 'localhost', // can be overwritten by process.env.HOST
    port: 7171, // can be overwritten by process.env.PORT, if port is in use, a free one will be determined

    // https://webpack.js.org/configuration/devtool/#development
    devtool: 'eval-cheap-module-source-map'
  },

  build: {
    // Template for index.html
    index: path.resolve(__dirname, '../../public/index.html'),

    // Paths
    assetsRoot: path.resolve(__dirname, '../../public'),
    assetsSubDirectory: 'static',
    assetsPublicPath: '/',

    /**
     * Source Maps
     */

    productionSourceMap: false,
    // https://webpack.js.org/configuration/devtool/#production
    devtool: '#source-map'
  }
};
