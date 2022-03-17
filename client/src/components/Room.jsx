import React, { useEffect, useRef } from 'react'
import configData from './config.json'
const Room = (props) => {
    const userVideo = useRef()
    const userStream = useRef()
    const partnerVideo = useRef()
    const peerRef = useRef()
    const webSocketRef = useRef()

    const openCamera = async () => {
        const allDevices = await navigator.mediaDevices.enumerateDevices()
        const cameras = allDevices.filter(
            (device) => device.kind == "videoinput"
        )

        const constraints = {
            audio: true,
            video: {
                deviceId: cameras[0].deviceId,
            },
        }

        try {
            return await navigator.mediaDevices.getUserMedia(constraints)
        } catch (error) {
            console.log(error)
        }
    }

    useEffect(() => {
        openCamera().then((stream) => {
            userVideo.current.srcObject = stream
            userStream.current = stream

            webSocketRef.current = new WebSocket(
                configData.REACT_APP_WEBSOCKET_URL + `join?roomID=${props.match.params.roomID}`
            )

            webSocketRef.current.addEventListener("open", () => {
                webSocketRef.current.send(JSON.stringify({ join: true }))
            })

            webSocketRef.current.addEventListener("message", async (e) => {
                const message = JSON.parse(e.data)

                if (message.join) {
                    callUser();
                }

                if (message.offer) {
                    handleOffer(message.offer)
                }

                if (message.answer) {
                    console.log("Receiving asnwer")
                    peerRef.current.setRemoteDescription(
                        new RTCSessionDescription(message.answer)
                    )
                }

                if (message.iceCandidate) {
                    console.log("Receiving and adding ICE candidate")
                    try {
                        await peerRef.current.addIceCandidate(message.iceCandidate)
                    } catch (error) {
                        console.log("Error receiving ICE candidate", error)
                    }
                }
            })
        })
    })

    const handleOffer = async (offer) => {
        console.log("Received offer. Creating answer")
        peerRef.current = createPeer()

        await peerRef.current.setRemoteDescription(
            new RTCSessionDescription(offer)
        )

        userStream.current.getTracks().forEach(track => {
            peerRef.current.addTrack(track, userStream.current)
        })

        const answer = await peerRef.current.createAnswer()
        await peerRef.current.setLocalDescription(answer)

        webSocketRef.current.send(
            JSON.stringify({ answer: peerRef.current.localDescription })
        )
    }

    const callUser = () => {
        console.log("Calling other user")
        peerRef.current = createPeer()

        userStream.current.getTracks().forEach(track => {
            peerRef.current.addTrack(track, userStream.current)
        })
    }

    const createPeer = () => {
        console.log("Creating peer connection")
        const peer = new RTCPeerConnection({
            iceServers: [{ urls: configData.REACT_APP_RTC_PEER_CONNECTION_URL }]
        })

        peer.onnegotiationneeded = handleNegotiationNeeded
        peer.onicecandidate = handleIceCandidateEvent
        peer.ontrack = handleTrackEvent

        return peer
    }

    const handleNegotiationNeeded = async () => {
        console.log("Creating offer")

        try {
            const myOffer = await peerRef.current.createOffer()
            await peerRef.current.setLocalDescription(myOffer)

            webSocketRef.current.send(
                JSON.stringify({ offer: peerRef.current.localDescription })
            )
        } catch (error) {
            console.log(error)
        }
    }

    const handleIceCandidateEvent = (e) => {
        console.log("Found Ice Candidate")
        if (e.candidate) {
            console.log(e.candidate)
            webSocketRef.current.send(
                JSON.stringify({ iceCandidate: e.candidate })
            )
        }
    }

    const handleTrackEvent = (e) => {
        console.log("Received Tracks")
        partnerVideo.current.srcObject = e.streams[0]
    }

    return (
        <div>
            <video autoPlay controls={true} ref={userVideo}></video>
            <video autoPlay controls={true} ref={partnerVideo}></video>
        </div>
    )
}
export default Room