'use strict';
const path = require('path');
const Dotenv = require('dotenv-webpack');

module.exports = {
    entry: {
        app: './frontend/js/app.js'
    },
    output: {
        filename: 'js/[name].js',
        chunkFilename: 'js/[name].chunk.js',
        path: path.resolve(__dirname, '../public')
    },
    module: {
        rules: [
            {
                test: /\.(png|jpe?g|gif)$/,
                use: {
                    loader: 'url-loader',
                    options: {
                        name: '[name].[ext]',
                        outputPath: 'images/',
                        limit: 8192
                    }
                }
            },
            {
                test: /\.js$/,
                include: /src/,
                exclude: /node_modules/,
                use: [
                    {
                        loader: 'expose-loader'
                    },
                    {
                        loader: 'babel-loader'
                    }
                ]
            },
            {
                test: require.resolve('jquery'),
                loader: 'expose-loader',
                options: {
                    exposes: ['$', 'jQuery'],
                },
            },
        ]
    },

    plugins: [
        new Dotenv({
            path: process.env.NODE_ENV === 'production' ? './.env.prod' : './.env'
        })
    ]
};
