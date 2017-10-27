import React from 'react'
import ReactDOM from 'react-dom'
import Song from './song'
import MyHappyPlayer from './happy'

class Playlist extends React.Component {

	constructor(props) {
		super(props)
		this.onClickHandler = this.onClickHandler.bind(this)
		this.onClickDeleteHandler = this.onClickDeleteHandler.bind(this)

		this.state = {
			songArr: []
		}

	}

	broadcastEvent(event) {
		switch (event.target.id){
			case "button-pause":
				window.MyPlayer.thinkAndDo({type:"onPauseButtonClick"});
				break;
			case "button-play":
				window.MyPlayer.thinkAndDo({type:"onPlayButtonClick"});
				break;
			case "button-next":
				window.MyPlayer.thinkAndDo({type:"onNextButtonClick"});
				break;
			case "button-prev":
				window.MyPlayer.thinkAndDo({type:"onPrevButtonClick"});
				break;
			case "button-shuffle":
				window.MyPlayer.thinkAndDo({type:"onShuffleButtonClick"});
				break;
			case "button-repeat":
				window.MyPlayer.thinkAndDo({type:"onRepeatButtonClick"});
				break;
		}
	}

	receiveSong(json) {
		window.MyPlayer = new MyHappyPlayer(json)

		$("#button-next").click(this.broadcastEvent);

		$("#button-prev").click(this.broadcastEvent);

		$("#button-pause").click(this.broadcastEvent);

		$("#button-play").click(this.broadcastEvent);

		$("#button-shuffle").click(this.broadcastEvent);

		$("#button-repeat").click(this.broadcastEvent);
		this.setState({
			songArr: json
		})
	}

	getSong() {
		fetch(`http://localhost:3000/song`, {
			method: 'GET',
			headers: {
				'Accept': 'application/json',
				'Content-Type':'application/json',
			},
		})
			.then(response => response.json())
			.then(json => this.receiveSong(json))
	}

	delSong(id) {
		fetch(`http://localhost:3000/song/` + id, {
			method: 'DELETE',
			headers: {
				'Accept': 'application/json',
				'Content-Type':'application/json',
			},
		})
	}

	componentDidMount() {
		this.getSong()
	}

	onClickDeleteHandler(id, index) {
		console.log("vo delete roi ne")
		this.delSong(id)
		window.MyPlayer.thinkAndDo({data:index,type:"onTrackDelete"})
		var temp = this.state.songArr.slice()
		temp = temp.filter(function(_, i){
				return i != index;
		})
		this.setState({
			songArr: temp
		})
	}

	onClickHandler(index) {
		window.MyPlayer.thinkAndDo({data:index,type:"onTrackClick"})
	}

	renderPlaylist() {
		const { songArr } = this.state;
		if (songArr) {
			return songArr.map((song, index) => 
				<Song 
					key={index} 
					index={index} 
					id={song.id}
					onClickHandler={this.onClickHandler} 
					onClickDeleteHandler={this.onClickDeleteHandler} 
					name={song.name}
					thumbnail={song.thumbnail}
				/>);
		}

		return ""
	}

	render() {
		const playlist = this.renderPlaylist()
		return (
			<ol>
			{playlist}
			</ol>
		)
	}
}

$(".button-back").click(function(){
	$(".player").toggleClass("playlist");
});

$("#hamburger-menu").click(function(){
	$(".player").toggleClass("playlist");
});

$("#button-pause").hide();

export default Playlist
var node = document.querySelector('#playlist-scroll')
ReactDOM.render(<Playlist />, node);
