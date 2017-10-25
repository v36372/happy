import React from 'react'
import ReactDOM from 'react-dom'

class Song extends React.Component {
	render() {
		const {name, thumbnail, onClickHandler, index } = this.props
		return (
			<li className='track' onClick={() => onClickHandler(index)}>
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
