import React from 'react'
import ReactDOM from 'react-dom'

class Song extends React.Component {
	render() {
		const {name, thumbnail } = this.props
		return (
			<li className='track'>
				<a>
					<img src={thumbnail}/>
					<div>
						<h1>{name}</h1>
					</div>
				</a>
				<hr/>
			</li>
		)
	}
}

export default Song
