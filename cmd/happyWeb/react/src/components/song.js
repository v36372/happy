import React from 'react'
import ReactDOM from 'react-dom'

class Song extends React.Component {
	render() {
		const {name, thumbnail, onClickHandler, onClickDeleteHandler, id, index } = this.props
		return (
			<li className='track' >
				<a>
					<img src={thumbnail} onClick={() => onClickHandler(index)}/>
					<div className="track-name" onClick={() => onClickHandler(index)}>
						<h1>{name}</h1>
					</div>
					<i onClick={() => onClickDeleteHandler(id, index)} className="fa fa-times fa-2x button-delete" aria-hidden="true"></i>
				</a>
				<hr/>
			</li>
		)
	}
}

export default Song
