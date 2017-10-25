import React from 'react'
import ReactDOM from 'react-dom'
import Song from './song'

class Playlist extends React.Component {

	constructor(props) {
		super(props)
		this.onClickHandler = this.onClickHandler.bind(this)

		this.state = {
			songArr: []
		}

	}

	receiveSong(json) {
		window.Happy = new HappyPlayer(json)
		this.setState({
			songArr: json
		})
	}

	componentDidMount() {
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

	onClickHandler() {
		console.log("123123312")
	}

	renderPlaylist() {
		const { songArr } = this.state;
		if (songArr) {
		console.log("123123312")
			return songArr.map((song, index) => 
				<Song key={index} onClick={this.onClickHandler} 
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

export default Playlist
var node = document.querySelector('#playlist-scroll')
ReactDOM.render(<Playlist />, node);
