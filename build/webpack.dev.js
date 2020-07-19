'use strict';
const webpack = require('webpack');
const { merge } = require('webpack-merge');
const baseConfig = require('./webpack.base');
const config = require('./config');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const FriendlyErrorsPlugin = require('friendly-errors-webpack-plugin');
const portfinder = require('portfinder');

const HOST = process.env.HOST;
const PORT = process.env.PORT && Number(process.env.PORT);

const webpackConfig = {
    mode: 'development',

    // cheap-module-eval-source-map is faster for development
    devtool: config.dev.devtool,

    devServer: {
        contentBase: './static',
        inline: true,
        hot: true,
        progress: true,
        host: HOST || config.dev.host,
        port: PORT || config.dev.port,
        quiet: true, // necessary for FriendlyErrorsPlugin
        overlay: true
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
                            publicPath: '../',
                            hmr: true
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
            'process.env': require('./config/dev.env')
        }),
        new webpack.HotModuleReplacementPlugin(),
        new webpack.NamedModulesPlugin(),
        new HtmlWebpackPlugin({
            template: './frontend/index.html',
            filename: 'index.html',
            chunks: ['app']
        }),
        new MiniCssExtractPlugin({
            filename: 'css/[name].css',
            chunkFilename: 'css/[id].css',
            ignoreOrder: false, // Enable to remove warnings about conflicting order
        })
    ]
};

const devWebpackConfig = merge(baseConfig, webpackConfig);

module.exports = new Promise((resolve, reject) => {
    portfinder.basePort = process.env.PORT || config.dev.port;
    portfinder.getPort((err, port) => {
        if (err) {
            reject(err);
        } else {
            // publish the new Port, necessary for e2e tests
            process.env.PORT = port;
            // add port to devServer config
            devWebpackConfig.devServer.port = port;

            // Add FriendlyErrorsPlugin
            devWebpackConfig.plugins.push(new FriendlyErrorsPlugin({
                compilationSuccessInfo: {
                    messages: [`Your application is running here: http://${devWebpackConfig.devServer.host}:${port}`]
                }
            }));

            resolve(devWebpackConfig);
        }
    })
});
