'use strict';
const webpack = require('webpack');
const merge = require('webpack-merge');
const baseConfig = require('./webpack.base');
const config = require('./config');
const utils = require('./utils');
const TerserJSPlugin = require('terser-webpack-plugin');
const OptimizeCSSAssetsPlugin = require('optimize-css-assets-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');

const env = process.env.NODE_ENV === 'testing'
    ? require('./config/test.env')
    : require('./config/prod.env');

const webpackConfig = {
    mode: 'production',

    devtool: config.build.productionSourceMap ? config.build.devtool : false,
    output: {
        path: config.build.assetsRoot,
        filename: utils.assetsPath('js/[name].[chunkhash].js'),
        chunkFilename: utils.assetsPath('js/[id].[chunkhash].js')
    },

    module: {
        rules: [
            {
                test: /\.(sa|sc|c)ss$/,
                use: [
                    {
                        loader: MiniCssExtractPlugin.loader,
                        options: {
                            // you can specify a publicPath here
                            // by default it uses publicPath in webpackOptions.output
                            publicPath: '../'
                        },
                    },
                    'css-loader',
                    'postcss-loader',
                    'sass-loader'
                ]
            }
        ]
    },

    plugins: [
        new webpack.DefinePlugin({
            'process.env': env
        }),
        // generate dist index.html with correct asset hash for caching.
        // you can customize output by editing /index.html
        // see https://github.com/ampedandwired/html-webpack-plugin
        new HtmlWebpackPlugin({
            filename: process.env.NODE_ENV === 'testing'
                ? 'index.html'
                : config.build.index,
            template: './frontend/index.html',
            inject: true,
            minify: {
                removeComments: true,
                collapseWhitespace: true,
                removeAttributeQuotes: true
                // more options:
                // https://github.com/kangax/html-minifier#options-quick-reference
            },
            // necessary to consistently work with multiple chunks via CommonsChunkPlugin
            chunksSortMode: 'dependency'
        }),
        new MiniCssExtractPlugin({
            filename: utils.assetsPath('css/[name].[chunkhash].css'),
            chunkFilename: utils.assetsPath('css/[id].[chunkhash].css'),
            ignoreOrder: false, // Enable to remove warnings about conflicting order
        })
    ],

    optimization: {
        minimizer: [new TerserJSPlugin({}), new OptimizeCSSAssetsPlugin({})]
    }
};

module.exports = merge(baseConfig, webpackConfig);
