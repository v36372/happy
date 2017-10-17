(function(){
	'use strict';
	var HappyPlayer = function() {
		this.init;
	};

	const MODE_NORMAL = 0, MODE_SHUFFLE = 1, MODE_REPEAT_ALL = 2, MODE_REPEAT_ONE = 3;

	HappyPlayer.prototype = {
		init: function(){
			var self = this || Happy;

			self._playMode = MODE_NORMAL;
		};
	};
})());
