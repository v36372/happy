var webpack = require('webpack');
var path = require('path');

var BUILD_DIR = path.resolve(__dirname, '../static/js');
var APP_DIR = path.resolve(__dirname, 'src/');

var config = {
	entry: {
		index: APP_DIR + '/index.js'
	},
	cache: false,
	output: {
		path: BUILD_DIR,
		filename: "[name].bundle.js"
	},
	module : {
		loaders : [
			{
				test : /\.js?/,
				include : APP_DIR,
				loader : 'babel-loader'
			}
		]
	},
	plugins: [
		new webpack.DefinePlugin({
			'process.env': {
				// This has effect on the react lib size
				'NODE_ENV': JSON.stringify('production'),
			}
		}),
		new webpack.optimize.UglifyJsPlugin({
			mangle: true,
			compress: {
				warnings: false, // Suppress uglification warnings
				pure_getters: true,
				unsafe: true,
				unsafe_comps: true,
				screw_ie8: true
			},
			output: {
				comments: false,
			},
			exclude: [/\.min\.js$/gi] // skip pre-minified libs
		})
	]
};

module.exports = config;
