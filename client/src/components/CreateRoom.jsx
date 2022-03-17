import React from 'react'
import configData from './config.json'

const CreateRoom = (props) => {
    const create = async event => {
        event.preventDefault();

        const resp = await fetch(configData.REACT_APP_URL + "create");
        const { room_id } = await resp.json();

        props.history.push(`/room/${room_id}`)
    }

    return (
        <div>
            <button onClick={create}>Create Room</button>
        </div>
    )
}

export default CreateRoom
