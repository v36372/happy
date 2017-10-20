(function(){
	'use strict';
	var HappyPlayer = function(playList){
		this.init(playList);
	}

	HappyPlayer.prototype = {
		init: function(playList){
			var self = this

			self._currentIndex = 0;
			self._playListLength = playList.length;
			self._playerState = "wait";
			self._playList = playList.map(function(song, index){
				return {
					song: song,
					next: index == self._playListLength-1 ? 0 : index+1,
					prev: index == 0 ? self._playListLength-1 : index-1,
				}
			})
		},

		updateUI: function(){

		},

		seekToBegin: function(){

		},

		thinkAndDo: function(event){
			var self=this;
			if (event == "onSongEnded") {
				if (self._repeatOne) {
					self.seekToBegin();
				}
				self.next();
			}
			if (event == "onPauseButtonClick") {
				self.pauseSong();
			}
			if (event == "onPlayButtonClick") {
				if (self._playerState == "pause") {
					self.resumeSong();
				} else {
					self.playSong();
				}
			}
			if (event == "onNextButtonClick") {
				self.next();
			}
			if (event == "onPrevButtonClick") {
				self.prev();
			}
			if (event == "onShuffleButtonClick") {
				self.shuffleToggle();
			}
			if (event == "onRepeatButtonClick") {
				self.repeatOneToggle();
			}
			self.updateUI();
		},

		repeatOneToggle: function() {
			var self = this;
			self._repeatOne = !self._repeatOne;
		},

		shuffleToggle: function() {
			var self = this;
			if (self._shuffle) {
				self._playList = self._originalPlayList;
			} else {
				self._playList = self._originalPlayList.slice();
				var temp=null,j=0;
				self._playList.forEach(function(e,i){
					var j = Math.floor(Math.random() * (i + 1))
					var temp = self._playList[i].song;
					self._playList[i].song = self._playList[j].song;
					self._playList[j].song = temp;
				});
			}
			self._shuffle = !self._shuffle;
		},

		playSong: function(){
			var self = this;
			self._playerstate = "playing";
		},

		resumeSong: function(){
			var self = this;
			self._playerstate = "playing";
		},

		pauseSong: function(){
			var self = this;
			self._playerstate = "pause";
		},

		prev: function() {
			var self = this;
			self._currentIndex = self.playList[self._currentIndex].prev;
		},

		next: function() {
			var self=this;
			self._currentIndex = self.playList[self._currentIndex].next;
		},

		onSongEnded: function(){
			var self = this;
			self.thinkAndDo("onSongEnded");
		},
	};

	if (typeof window !== 'undefined'){
		window.HappyPlayer = HappyPlayer;
	}
}());
