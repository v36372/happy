import React from 'react'
import ReactDOM from 'react-dom'

class Playlist extends React.Component {
	constructor(props) {
		super(props)
	}

	componentDidMount() {
		this.fetchSong()
	}

	fetchSong() {
		fetch(`http://localhost:3000/song`, {
			method: 'GET',
			headers: {
				'Accept': 'application/json',
				'Content-Type': 'application/json',
			}
		})
			.then(response => response.json())
			.then(json => {
				console.log(json)
			})
			.catch(function(error){
				console.log(error)
			})
	}

	render() {
		return (
			<div>
				asdasd
			</div>
		)
	}
}


var node = document.querySelector('#playlist-scroll')
ReactDOM.render(<Playlist />, node);
export default Playlist
