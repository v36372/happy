(function(){
	'use strict';
	var HappyPlayer = function(playList){
		this.init(playList);
	}

	HappyPlayer.prototype = {
		init: function(playList){
			var self = this

			self.currentIndex = 0;
			self.playerState = "wait";
			self.shuffle = false;
			self.repeatOne = false;
			if (playList) {
				self.playListLength = playList.length;
				self.setPlaylist(playList);
			}
		},

		setPlaylist: function(playList) {
			var self = this;
			self.originalPlayList = self.playList = playList.map(function(song, index){
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
			console.log("think and do", event)
			if (event.type == "onSongEnded") {
				if (self.repeatOne) {
					self.seekToBegin();
					return;
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
					event.data = self.currentIndex;
					self.playSong(event);
				}
			}
			if (event.type == "onNextButtonClick") {
				self.next();
			}
			if (event.type == "onPrevButtonClick") {
				self.prev();
			}
			if (event.type == "onShuffleButtonClick") {
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
			console.log("parent shuffle")
			var self = this;
			if (self.shuffle) {
				var currentPlayingID = self.playList[self.currentIndex].song.id;
				self.playList = self.originalPlayList;
				self.currentIndex = self.playList.findIndex(function(el){return el.song.id == currentPlayingID});

			} else {
				var currentPlayingID = self.playList[self.currentIndex].song.id;
				self.playList = self.originalPlayList.map(function(el){
					return {
						next: el.next,
						prev: el.prev,
						song: {
							id: el.song.id,
							name: el.song.name,
							link: el.song.link,
							provider: el.song.provider,
							thumbnail: el.song.thumbnail
						}
					}
				});
				var temp=null,j=0;
				self.playList.forEach(function(e,i){
					var j = Math.floor(Math.random() * (i + 1))
					var temp = self.playList[i].song;
					self.playList[i].song = self.playList[j].song;
					self.playList[j].song = temp;
				});
				self.currentIndex = self.playList.findIndex(function(el){return el.song.id == currentPlayingID});
			}
			self.shuffle = !self.shuffle;
		},

		playSong: function(event){
			var self = this;
			self.playerState = "playing";
			self.currentIndex = event.data;
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
			self.playSong({data:self.playList[self.currentIndex].prev})
		},

		next: function() {
			var self=this;
			self.playSong({data:self.playList[self.currentIndex].next})
		},

		onSongEnded: function(){
			var self = this;
			self.thinkAndDo({type:"onSongEnded"});
		},
	};

	if (typeof window !== 'undefined'){
		window.HappyPlayer = HappyPlayer;
	}
}());
