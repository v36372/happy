(function(){
	'use strict';
	var HappyPlayer = function(playList){
		this.init(playList);
	}

	HappyPlayer.prototype = {
		init: function(playList){
			console.log("parent init")
			var self = this

			self.currentIndex = 0;
			self.playerState = "wait";
			if (playList) {
				self.playListLength = playList.length;
				self.setPlaylist(playList);
			}
		},

		setPlaylist: function(playList) {
			var self = this;
			self.playList = playList.map(function(song, index){
				return {
					song: song,
					next: index == self.playListLength-1 ? 0 : index+1,
					prev: index == 0 ? self.playListLength-1 : index-1,
				}
			});
		},

		updateUI: function(){

		},

		seekToBegin: function(){

		},

		thinkAndDo: function(event){
			var self=this;
			if (event.type == "onSongEnded") {
				if (self.repeatOne) {
					self.seekToBegin();
				}
				self.next();
			}
			if (event.type == "onPauseButtonClick") {
				self.pauseSong();
			}
			if (event.type == "onPlayButtonClick") {
				if (self.playerState == "pause") {
					self.resumeSong();
				} else {
					self.playSong(event);
				}
			}
			if (event.type == "onNextButtonClick") {
				self.next();
			}
			if (event.type == "onPrevButtonClick") {
				self.prev();
			}
			if (event == "onShuffleButtonClick") {
				self.shuffleToggle();
			}
			if (event.type == "onRepeatButtonClick") {
				self.repeatOneToggle();
			}
			if (event.type == "onTrackClick") {
				self.playSong(event);
			}
			self.updateUI();
		},

		repeatOneToggle: function() {
			var self = this;
			self.repeatOne = !self.repeatOne;
		},

		shuffleToggle: function() {
			var self = this;
			if (self.shuffle) {
				self.playList = self.originalPlayList;
			} else {
				self.playList = self.originalPlayList.slice();
				var temp=null,j=0;
				self.playList.forEach(function(e,i){
					var j = Math.floor(Math.random() * (i + 1))
					var temp = self.playList[i].song;
					self.playList[i].song = self.playList[j].song;
					self.playList[j].song = temp;
				});
			}
			self.shuffle = !self.shuffle;
		},

		playSong: function(){
			var self = this;
			self.playerState = "playing";
			console.log("parent run")
		},

		resumeSong: function(){
			var self = this;
			self.playerState = "playing";
		},

		pauseSong: function(){
			var self = this;
			self.playerState = "pause";
		},

		prev: function() {
			var self = this;
			self.currentIndex = self.playList[self.currentIndex].prev;
		},

		next: function() {
			var self=this;
			self.currentIndex = self.playList[self.currentIndex].next;
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
