package sdp

import (
	"strconv"
	"strings"
	"time"

	"github.com/notedit/sdp/transform"
)

type SDPInfo struct {
	version      int
	streams      map[string]*StreamInfo
	orderStreams []*StreamInfo
	medias       []*MediaInfo     // as we need to keep order
	candidates   []*CandidateInfo // keep order
	ice          *ICEInfo
	dtls         *DTLSInfo
	crypto       *CryptoInfo
}

func NewSDPInfo() *SDPInfo {

	sdp := &SDPInfo{
		version:      1,
		streams:      map[string]*StreamInfo{},
		orderStreams: []*StreamInfo{},
		medias:       []*MediaInfo{},
		candidates:   []*CandidateInfo{},
	}

	return sdp
}

func (s *SDPInfo) SetVersion(version int) {

	s.version = version
}

func (s *SDPInfo) AddMedia(media *MediaInfo) {

	s.medias = append(s.medias, media)
}

func (s *SDPInfo) GetMedia(mtype string) *MediaInfo {

	for _, media := range s.medias {
		if strings.ToLower(media.GetType()) == strings.ToLower(mtype) {
			return media
		}
	}
	return nil
}

func (s *SDPInfo) GetAudioMedia() *MediaInfo {
	for _, media := range s.medias {
		if strings.ToLower(media.GetType()) == "audio" {
			return media
		}
	}
	return nil
}

func (s *SDPInfo) GetVideoMedia() *MediaInfo {
	for _, media := range s.medias {
		if strings.ToLower(media.GetType()) == "video" {
			return media
		}
	}
	return nil
}

func (s *SDPInfo) GetMediasByType(mtype string) []*MediaInfo {

	medias := []*MediaInfo{}
	for _, media := range s.medias {
		if strings.ToLower(media.GetType()) == strings.ToLower(mtype) {
			medias = append(medias, media)
		}
	}
	return medias
}

func (s *SDPInfo) GetMediaByID(mid string) *MediaInfo {

	for _, media := range s.medias {
		if strings.ToLower(media.GetID()) == strings.ToLower(mid) {
			return media
		}
	}
	return nil
}

func (s *SDPInfo) ReplaceMedia(media *MediaInfo) bool {

	for i, rmedia := range s.medias {
		if rmedia.GetID() == media.GetID() {
			s.medias[i] = media
			return true
		}
	}
	return false
}

func (s *SDPInfo) GetMedias() []*MediaInfo {

	return s.medias
}

func (s *SDPInfo) GetVersion() int {

	return s.version
}

func (s *SDPInfo) GetDTLS() *DTLSInfo {

	return s.dtls
}

func (s *SDPInfo) SetDTLS(dtls *DTLSInfo) {

	s.dtls = dtls
}

func (s *SDPInfo) GetCrypto() *CryptoInfo {

	return s.crypto
}

func (s *SDPInfo) SetCrypto(crypto *CryptoInfo) {

	s.crypto = crypto
}

func (s *SDPInfo) GetICE() *ICEInfo {

	return s.ice
}

func (s *SDPInfo) SetICE(ice *ICEInfo) {

	s.ice = ice
}

func (s *SDPInfo) AddCandidate(candidate *CandidateInfo) {
	if candidate != nil {
		s.candidates = append(s.candidates, candidate)
	}
}

func (s *SDPInfo) AddCandidates(candidates []*CandidateInfo) {

	for _, candidate := range candidates {
		s.AddCandidate(candidate)
	}
}

func (s *SDPInfo) GetCandidates() []*CandidateInfo {

	return s.candidates
}

func (s *SDPInfo) GetStream(id string) *StreamInfo {

	return s.streams[id]
}

func (s *SDPInfo) GetStreams() map[string]*StreamInfo {

	return s.streams
}

func (s *SDPInfo) GetOrderStreams() []*StreamInfo {

	return s.orderStreams
}

func (s *SDPInfo) GetFirstStream() *StreamInfo {

	for _, stream := range s.orderStreams {
		return stream
	}
	return nil
}

func (s *SDPInfo) AddStream(stream *StreamInfo) {
	_, ok := s.streams[stream.GetID()]
	if ok { // exist
		s.streams[stream.GetID()] = stream
		for idx, v := range s.orderStreams {
			if stream.GetID() == v.GetID() {
				s.orderStreams[idx] = stream
				return
			}
		}
		return
	}

	s.streams[stream.GetID()] = stream
	s.orderStreams = append(s.orderStreams, stream)
}

func (s *SDPInfo) RemoveStream(stream *StreamInfo) {
	_, ok := s.streams[stream.GetID()]
	if ok { // exist
		for idx, v := range s.orderStreams {
			if stream.GetID() == v.GetID() {
				s.orderStreams = append(s.orderStreams[:idx], s.orderStreams[idx+1:]...)
				return
			}
		}
	}
	delete(s.streams, stream.GetID())

}

func (s *SDPInfo) RemoveAllStreams() {
	s.streams = make(map[string]*StreamInfo)
	s.orderStreams = s.orderStreams[:0]
}

func (s *SDPInfo) GetTrackByMediaID(mid string) *TrackInfo {
	for _, stream := range s.streams {
		for _, track := range stream.GetTracks() {
			if track.GetMID() == mid {
				return track
			}
		}
	}
	return nil
}

func (s *SDPInfo) GetStreamByMediaID(mid string) *StreamInfo {

	for _, stream := range s.streams {
		for _, track := range stream.GetTracks() {
			if track.GetMID() == mid {
				return stream
			}
		}
	}
	return nil
}

func (s *SDPInfo) GetVideoTracks() []*TrackInfo {

	tracks := []*TrackInfo{}
	for _, stream := range s.orderStreams {
		for _, track := range stream.GetTracks() {
			if strings.ToLower(track.GetMediaType()) == "video" {
				tracks = append(tracks, track)
			}
		}
	}
	return tracks
}

func (s *SDPInfo) GetAudioTracks() []*TrackInfo {

	tracks := []*TrackInfo{}
	for _, stream := range s.orderStreams {
		for _, track := range stream.GetTracks() {
			if strings.ToLower(track.GetMediaType()) == "audio" {
				tracks = append(tracks, track)
			}
		}
	}
	return tracks
}

func (s *SDPInfo) Answer(ice *ICEInfo, dtls *DTLSInfo, candidates []*CandidateInfo, medias map[string]*Capability) *SDPInfo {

	sdpInfo := NewSDPInfo()

	if ice != nil {
		sdpInfo.SetICE(ice.Clone())
	}

	if dtls != nil {
		sdpInfo.SetDTLS(dtls)
	}

	for _, candidate := range candidates {
		sdpInfo.AddCandidate(candidate)
	}

	for _, media := range s.medias {
		supported := medias[media.GetType()]
		if supported != nil {
			answer := media.AnswerCapability(supported)
			answer.Payloads = media.Payloads
			answer.Protocal = media.Protocal
			sdpInfo.AddMedia(answer)
		}
	}

	return sdpInfo
}

func (s *SDPInfo) String() string {

	sdpMap := &transform.SdpStruct{
		Version: 0,
		Media:   []*transform.MediaStruct{},
		Groups:  []*transform.GroupStruct{},
	}

	sdpMap.Origin = &transform.OriginStruct{
		Username:       "-",
		SessionId:      strconv.FormatInt(time.Now().UnixNano(), 10),
		SessionVersion: s.version,
		NetType:        "IN",
		IpVer:          4,
		Address:        "127.0.0.1",
	}

	sdpMap.Connection = &transform.ConnectionStruct{
		Version: 4,
		Ip:      "0.0.0.0",
	}

	sdpMap.Name = "media"

	if s.GetICE().IsLite() {
		sdpMap.Icelite = "ice-lite"
	}

	sdpMap.Timing = &transform.TimingStruct{
		Start: 0,
		Stop:  0,
	}

	sdpMap.MsidSemantic = &transform.MsidSemanticStruct{
		Semantic: "WMS",
		Token:    "*",
	}

	bundleType := "BUNDLE"
	bundleMids := []string{}

	for _, media := range s.medias {

		mediaMap := &transform.MediaStruct{
			Type:       media.GetType(),
			Port:       media.GetPort(),
			Protocal:   media.Protocal, // "UDP/TLS/RTP/SAVP",
			Fmtp:       []*transform.FmtpStruct{},
			Rtp:        []*transform.RtpStruct{},
			RtcpFb:     []*transform.RtcpFbStruct{},
			Ext:        []*transform.ExtStruct{},
			Bandwidth:  []*transform.BandwithStruct{},
			Candidates: []*transform.CandidateStruct{},
			SsrcGroups: []*transform.SsrcGroupStruct{},
			Ssrcs:      []*transform.SsrcStruct{},
			Rids:       []*transform.RidStruct{},
		}

		mediaMap.Direction = media.GetDirection().String()

		mediaMap.RtcpMux = "rtcp-mux"

		mediaMap.RtcpRsize = "rtcp-rsize"

		mediaMap.Mid = media.GetID()

		if media.GetDirection() == INACTIVE {
			mediaMap.Port = 0
		}

		if media.GetPort() != 0 {
			bundleMids = append(bundleMids, media.GetID())
		}

		if media.GetBitrate() > 0 {
			mediaMap.Bandwidth = append(mediaMap.Bandwidth, &transform.BandwithStruct{
				Type:  "AS",
				Limit: media.GetBitrate(),
			})
		}

		for _, candidate := range s.GetCandidates() {

			mediaMap.Candidates = append(mediaMap.Candidates, &transform.CandidateStruct{
				Foundation: candidate.GetFoundation(),
				Component:  candidate.GetComponentID(),
				Transport:  candidate.GetTransport(),
				Priority:   candidate.GetPriority(),
				Ip:         candidate.GetAddress(),
				Port:       candidate.GetPort(),
				Type:       candidate.GetType(),
				Raddr:      candidate.GetRelAddr(),
				Rport:      candidate.GetRelPort(),
			})
		}

		mediaMap.IceUfrag = s.GetICE().GetUfrag()
		mediaMap.IcePwd = s.GetICE().GetPassword()

		mediaMap.Fingerprint = &transform.FingerprintStruct{
			Type: s.GetDTLS().GetHash(),
			Hash: s.GetDTLS().GetFingerprint(),
		}

		mediaMap.Setup = s.GetDTLS().GetSetup().String()

		for _, codec := range media.GetCodecs() {

			if "video" == strings.ToLower(media.GetType()) {
				mediaMap.Rtp = append(mediaMap.Rtp, &transform.RtpStruct{
					Payload: codec.GetPayload(),
					Codec:   strings.ToUpper(codec.GetCodec()),
					Rate:    90000,
				})
			} else {
				if "opus" == strings.ToLower(codec.GetCodec()) {
					mediaMap.Rtp = append(mediaMap.Rtp, &transform.RtpStruct{
						Payload:  codec.GetPayload(),
						Codec:    codec.GetCodec(),
						Rate:     codec.GetRate(),
						Encoding: 2,
					})
				} else {
					mediaMap.Rtp = append(mediaMap.Rtp, &transform.RtpStruct{
						Payload: codec.GetPayload(),
						Codec:   codec.GetCodec(),
						Rate:    codec.GetRate(),
					})
				}
			}

			for _, rtcpfb := range codec.GetRTCPFeedbacks() {
				mediaMap.RtcpFb = append(mediaMap.RtcpFb, &transform.RtcpFbStruct{
					Payload: codec.GetPayload(),
					Type:    rtcpfb.GetID(),
					Subtype: strings.Join(rtcpfb.GetParams(), " "),
				})
			}

			if codec.HasRTX() {
				mediaMap.Rtp = append(mediaMap.Rtp, &transform.RtpStruct{
					Payload: codec.GetRTX(),
					Codec:   "rtx",
					Rate:    90000,
				})
				mediaMap.Fmtp = append(mediaMap.Fmtp, &transform.FmtpStruct{
					Payload: codec.GetRTX(),
					Config:  "apt=" + strconv.Itoa(codec.GetPayload()),
				})
			}

			params := codec.GetParams()

			if params != nil && len(params) > 0 {

				fmtp := &transform.FmtpStruct{
					Payload: codec.GetPayload(),
					Config:  "",
				}

				for k, v := range params {

					if fmtp.Config != "" {
						fmtp.Config = fmtp.Config + ";"
					}

					// k and value
					if v != "" {
						fmtp.Config = fmtp.Config + k + "=" + v
					} else {
						fmtp.Config = fmtp.Config + k
					}
				}

				mediaMap.Fmtp = append(mediaMap.Fmtp, fmtp)
			}
		}

		payloads := []int{}

		for _, rtp := range mediaMap.Rtp {
			payloads = append(payloads, rtp.Payload)
		}

		if strings.ToLower(media.mtype) == "application" {
			mediaMap.SctpMap = &transform.SctpMapStuct{
				5000,
				"webrtc-datachannel",
				1024,
			}
			mediaMap.Payloads = media.Payloads
			// mediaMap.SctpPort = 5000
			mediaMap.SctpMaxSize = 256 * 1024
			mediaMap.RtcpMux = ""
			mediaMap.RtcpRsize = ""
		} else {
			mediaMap.Payloads = intArrayToString(payloads, " ")
		}

		for id, uri := range media.GetExtensions() {

			mediaMap.Ext = append(mediaMap.Ext, &transform.ExtStruct{
				Value: id,
				Uri:   uri,
			})
		}

		for _, ridInfo := range media.GetRIDS() {

			rid := &transform.RidStruct{
				Id:        ridInfo.GetID(),
				Direction: ridInfo.GetDirection().String(),
				Params:    "",
			}

			if len(ridInfo.GetFormats()) > 0 {
				//rid.Params = "pt=" + strings.Join(ridInfo.GetFormats(), ",")
				rid.Params = "pt=" + intArrayToString(ridInfo.GetFormats(), ",")
			}

			for key, val := range ridInfo.GetParams() {
				if rid.Params == "" {
					rid.Params = key + "=" + val
				} else {
					rid.Params = rid.Params + ";" + key + "=" + val
				}
			}

			mediaMap.Rids = append(mediaMap.Rids, rid)
		}

		if media.GetSimulcastInfo() != nil {

			simulcast := media.GetSimulcastInfo()

			index := 1

			mediaMap.Simulcast = &transform.SimulcastStruct{}

			sendStreams := simulcast.GetSimulcastStreams(SEND)
			recvStreams := simulcast.GetSimulcastStreams(RECV)

			if sendStreams != nil && len(sendStreams) > 0 {
				list := ""
				for _, stream := range sendStreams {
					alternatives := ""
					for _, item := range stream {
						if alternatives == "" {
							if item.IsPaused() {
								alternatives = alternatives + "~" + item.GetID()
							} else {
								alternatives = alternatives + item.GetID()
							}
						} else {
							if item.IsPaused() {
								alternatives = alternatives + "," + "~" + item.GetID()
							} else {
								alternatives = alternatives + "," + item.GetID()
							}
						}
					}
					if list == "" {
						list = list + alternatives
					} else {
						list = list + ";" + alternatives
					}
				}
				mediaMap.Simulcast.Dir1 = "send"
				mediaMap.Simulcast.List1 = list
				index = index + 1
			}

			if recvStreams != nil && len(recvStreams) > 0 {
				list := ""
				for _, stream := range recvStreams {
					alternatives := ""
					for _, item := range stream {
						if alternatives == "" {
							if item.IsPaused() {
								alternatives = alternatives + "~" + item.GetID()
							} else {
								alternatives = alternatives + item.GetID()
							}
						} else {
							if item.IsPaused() {
								alternatives = alternatives + "," + "~" + item.GetID()
							} else {
								alternatives = alternatives + "," + item.GetID()
							}
						}
					}
					if list == "" {
						list = list + alternatives
					} else {
						list = list + ";" + alternatives
					}
				}
				if index == 1 {
					mediaMap.Simulcast.Dir1 = "recv"
					mediaMap.Simulcast.List1 = list
				}
				if index == 2 {
					mediaMap.Simulcast.Dir2 = "recv"
					mediaMap.Simulcast.List2 = list
				}
			}
		}

		sdpMap.Media = append(sdpMap.Media, mediaMap)
	}

	// streams
	for _, stream := range s.GetOrderStreams() {
		for _, track := range stream.GetTracks() {
			for _, md := range sdpMap.Media {
				// check if it is unified or plan b
				if track.GetMID() != "" {
					if track.GetMID() == md.Mid {
						groups := track.GetSourceGroupS()
						for _, group := range groups {
							md.SsrcGroups = append(md.SsrcGroups, &transform.SsrcGroupStruct{
								Semantics: group.GetSemantics(),
								Ssrcs:     uint32ArrayToString(group.GetSSRCs(), " "),
							})
						}
						ssrcs := track.GetSSRCS()
						for _, ssrc := range ssrcs {
							md.Ssrcs = append(md.Ssrcs, &transform.SsrcStruct{
								Id:        ssrc,
								Attribute: "cname",
								Value:     stream.GetID(),
							})
							md.Ssrcs = append(md.Ssrcs, &transform.SsrcStruct{
								Id:        ssrc,
								Attribute: "msid",
								Value:     stream.GetID() + " " + track.GetID(),
							})
						}
						md.Msid = stream.GetID() + " " + track.GetID()
						break
					}
				} else if strings.ToLower(md.Type) == strings.ToLower(track.GetMediaType()) {

					groups := track.GetSourceGroupS()
					for _, group := range groups {
						md.SsrcGroups = append(md.SsrcGroups, &transform.SsrcGroupStruct{
							Semantics: group.GetSemantics(),
							Ssrcs:     uint32ArrayToString(group.GetSSRCs(), " "),
						})
					}
					ssrcs := track.GetSSRCS()
					for _, ssrc := range ssrcs {
						md.Ssrcs = append(md.Ssrcs, &transform.SsrcStruct{
							Id:        ssrc,
							Attribute: "cname",
							Value:     stream.GetID(),
						})
						md.Ssrcs = append(md.Ssrcs, &transform.SsrcStruct{
							Id:        ssrc,
							Attribute: "msid",
							Value:     stream.GetID() + " " + track.GetID(),
						})
					}
					break
				}
			}
		}
	}
	sdpMap.Groups = append(sdpMap.Groups, &transform.GroupStruct{
		Mids: strings.Join(bundleMids, " "),
		Type: bundleType,
	})

	sdpStr, err := transform.Write(sdpMap)
	if err != nil {
		println(err)
	}

	return sdpStr
}

func (s *SDPInfo) Clone() *SDPInfo {

	cloned := NewSDPInfo()
	cloned.SetVersion(s.GetVersion())
	for _, media := range s.GetMedias() {
		cloned.AddMedia(media.Clone())
	}
	for _, stream := range s.GetOrderStreams() {
		cloned.AddStream(stream.Clone())
	}
	for _, candidate := range s.GetCandidates() {
		cloned.AddCandidate(candidate)
	}
	cloned.SetICE(s.GetICE().Clone())

	if s.GetDTLS() != nil {
		cloned.SetDTLS(s.GetDTLS().Clone())
	}
	if s.GetCrypto() != nil {
		cloned.SetCrypto(s.GetCrypto().Clone())
	}
	return cloned
}

// Unify return an unified plan version of the SDP info
func (s *SDPInfo) Unify() *SDPInfo {
	cloned := NewSDPInfo()

	cloned.version = s.version

	for _, media := range s.medias {
		cloned.AddMedia(media.Clone())
	}

	medias := map[string][]*MediaInfo{
		"audio": cloned.GetMediasByType("audio"),
		"video": cloned.GetMediasByType("video"),
	}

	for _, stream := range s.orderStreams {
		clonedStream := stream.Clone()
		for _, clonedTrack := range clonedStream.GetTracks() {
			var clonedMedia *MediaInfo
			if len(medias[clonedTrack.GetMediaType()]) == 0 {
				media := s.GetMedia(clonedTrack.GetMediaType())
				clonedMedia = media.Clone()
				clonedMedia.SetID(clonedTrack.GetID())
				cloned.AddMedia(clonedMedia)
			} else {
				mediaList := medias[clonedTrack.GetMediaType()]
				clonedMedia = mediaList[len(mediaList)-1]
				mediaList = mediaList[:len(mediaList)-1]
				medias[clonedTrack.GetMediaType()] = mediaList
			}
			clonedTrack.SetMID(clonedMedia.GetID())
		}
		cloned.AddStream(clonedStream)
	}

	for _, candidate := range s.GetCandidates() {
		cloned.AddCandidate(candidate.Clone())
	}

	cloned.SetICE(s.GetICE().Clone())

	if s.GetDTLS() != nil {
		cloned.SetDTLS(s.GetDTLS().Clone())
	}
	if s.GetCrypto() != nil {
		cloned.SetCrypto(s.GetCrypto().Clone())
	}

	return cloned
}

func Create(ice *ICEInfo, dtls *DTLSInfo, candidates []*CandidateInfo, capabilities map[string]*Capability) *SDPInfo {

	sdpInfo := NewSDPInfo()

	if ice != nil {
		sdpInfo.SetICE(ice.Clone())
	}

	if dtls != nil {
		sdpInfo.SetDTLS(dtls)
	}

	for _, candidate := range candidates {
		sdpInfo.AddCandidate(candidate)
	}

	dyn := 96

	for mType, capability := range capabilities {
		media := MediaInfoCreate(mType, capability)
		for _, codec := range media.GetCodecs() {
			if codec.GetPayload() >= 96 {
				dyn++
				codec.SetPayload(dyn)
			}
			if codec.GetRTX() > 0 {
				dyn++
				codec.SetRTX(dyn)
			}
		}
		sdpInfo.AddMedia(media)
	}

	return sdpInfo
}

func Create2(capabilities map[string]*Capability) *SDPInfo {

	sdpInfo := NewSDPInfo()
	dyn := 96
	for mType, capability := range capabilities {
		media := MediaInfoCreate(mType, capability)
		for _, codec := range media.GetCodecs() {
			if codec.GetPayload() >= 96 {
				dyn++
				codec.SetPayload(dyn)
			}
			if codec.GetRTX() > 0 {
				dyn++
				codec.SetRTX(dyn)
			}
		}
		sdpInfo.AddMedia(media)
	}

	return sdpInfo
}

func Parse(sdp string) (*SDPInfo, error) {

	sdpMap, err := transform.Parse(sdp)

	if err != nil {
		return nil, err
	}

	sdpInfo := NewSDPInfo()

	sdpInfo.SetVersion(sdpMap.Version)

	for _, md := range sdpMap.Media {

		media := md.Type
		mid := md.Mid

		if len(mid) == 0 {
			continue
		}

		mediaInfo := NewMediaInfo(mid, media)
		mediaInfo.Payloads = md.Payloads
		mediaInfo.Protocal = md.Protocal

		ufrag := md.IceUfrag
		pwd := md.IcePwd

		sdpInfo.SetICE(NewICEInfo(ufrag, pwd))

		for _, candiate := range md.Candidates {

			candidateInfo := NewCandidateInfo(
				candiate.Foundation,
				candiate.Component,
				candiate.Transport,
				candiate.Priority,
				candiate.Ip,
				candiate.Port,
				candiate.Type,
				candiate.Raddr,
				candiate.Rport)

			sdpInfo.AddCandidate(candidateInfo)
		}

		var remoteHash, remoteFingerprint string

		if sdpMap.Fingerprint != nil {
			remoteHash = sdpMap.Fingerprint.Type
			remoteFingerprint = sdpMap.Fingerprint.Hash
		}

		if md.Fingerprint != nil {
			remoteHash = md.Fingerprint.Type
			remoteFingerprint = md.Fingerprint.Hash
		}

		setup := SETUPACTPASS

		if md.Setup != "" {
			setup = SetupByValue(md.Setup)
		}

		sdpInfo.SetDTLS(NewDTLSInfo(setup, remoteHash, remoteFingerprint))

		if md.BundleOnly != "" && md.Port == 0 {
			md.Port = 9
		}

		direction := SENDRECV

		if md.Direction != "" {
			direction = DirectionbyValue(md.Direction)
		}

		mediaInfo.SetDirection(direction)
		mediaInfo.SetPort(md.Port)

		apts := map[int]int{}

		for _, fmt := range md.Rtp {

			payload := fmt.Payload
			codec := fmt.Codec
			rate := fmt.Rate

			if "RED" == strings.ToUpper(codec) || "ULPFEC" == strings.ToUpper(codec) {
				continue
			}

			params := map[string]string{}

			for _, fmtp := range md.Fmtp {

				if fmtp.Payload == payload {
					list := strings.Split(fmtp.Config, ";")

					for _, kv := range list {
						param := strings.Split(kv, "=")
						if len(param) < 2 {
							continue
						}
						params[param[0]] = param[1]
					}
				}
			}

			if "RTX" == strings.ToUpper(codec) {
				if apt, ok := params["apt"]; ok {
					aptint, _ := strconv.Atoi(apt)
					apts[aptint] = payload
				}
			} else {
				codecInfo := NewCodecInfo(codec, payload, rate)
				codecInfo.AddParams(params)
				mediaInfo.AddCodec(codecInfo)
			}
		}

		// rtx
		for pt1, pt2 := range apts {
			codecInfo := mediaInfo.GetCodecForType(pt1)
			if codecInfo != nil {
				codecInfo.SetRTX(pt2)
			}
		}

		// rtcpFb
		if md.RtcpFb != nil {
			for _, rtcfb := range md.RtcpFb {
				codecInfo := mediaInfo.GetCodecForType(rtcfb.Payload)
				if codecInfo != nil {
					id := rtcfb.Type
					params := []string{}
					if rtcfb.Subtype != "" {
						params = strings.Split(rtcfb.Subtype, " ")
					}
					codecInfo.AddRTCPFeedback(NewRTCPFeedbackInfo(id, params))
				}
			}
		}

		// extmap
		for _, extmap := range md.Ext {
			mediaInfo.AddExtension(extmap.Value, extmap.Uri)
		}

		for _, rid := range md.Rids {
			direction := DirectionWaybyValue(rid.Direction)
			ridInfo := NewRIDInfo(rid.Id, direction)

			formats := []string{}
			params := map[string]string{}

			if rid.Params != "" {
				list := transform.ParseParams(rid.Params)
				for k, v := range list {
					if k == "pt" {
						formats = strings.Split(v, ",")
					} else {
						params[k] = v
					}
				}
				ridInfo.SetFormats(formats)
				ridInfo.SetParams(params)
			}

			mediaInfo.AddRID(ridInfo)
		}

		encodings := [][]*TrackEncodingInfo{}

		if md.Simulcast != nil {

			simulcast := NewSimulcastInfo()

			if md.Simulcast.Dir1 != "" {
				direction := DirectionWaybyValue(md.Simulcast.Dir1)
				streamList := transform.ParseSimulcastStreamList(md.Simulcast.List1)
				for _, streams := range streamList {
					alternatives := []*SimulcastStreamInfo{}
					for _, stream := range streams {
						simulcastStreamInfo := NewSimulcastStreamInfo(stream.Scid, stream.Paused)
						alternatives = append(alternatives, simulcastStreamInfo)
					}
					simulcast.AddSimulcastAlternativeStreams(direction, alternatives)
				}
			}

			if md.Simulcast.Dir2 != "" {
				direction := DirectionWaybyValue(md.Simulcast.Dir2)
				streamList := transform.ParseSimulcastStreamList(md.Simulcast.List2)
				for _, streams := range streamList {
					alternatives := []*SimulcastStreamInfo{}
					for _, stream := range streams {
						simulcastStreamInfo := NewSimulcastStreamInfo(stream.Scid, stream.Paused)
						alternatives = append(alternatives, simulcastStreamInfo)
					}
					simulcast.AddSimulcastAlternativeStreams(direction, alternatives)
				}
			}

			// For all sending encodings
			for _, streams := range simulcast.GetSimulcastStreams(SEND) {
				alternatives := []*TrackEncodingInfo{}
				for _, stream := range streams {
					encoding := NewTrackEncodingInfo(stream.GetID(), stream.IsPaused())
					ridInfo := mediaInfo.GetRID(encoding.GetID())
					if ridInfo != nil {
						//Get associated payloads
						formats := ridInfo.GetFormats()
						for _, format := range formats {
							codecInfo := mediaInfo.GetCodecForType(format)
							if codecInfo != nil {
								encoding.AddCodec(codecInfo)
							}
						}
						encoding.SetParams(ridInfo.GetParams())
						alternatives = append(alternatives, encoding)
					}
				}

				if len(alternatives) > 0 {
					encodings = append(encodings, alternatives)
				}
			}

			mediaInfo.SetSimulcastInfo(simulcast)
		}

		sources := []*SourceInfo{}

		if md.Ssrcs != nil {
			for _, ssrcAttr := range md.Ssrcs {
				ssrc := ssrcAttr.Id
				key := ssrcAttr.Attribute
				value := ssrcAttr.Value

				var source *SourceInfo
				for _, sourceInfo := range sources {
					if sourceInfo.ssrc == ssrc {
						source = sourceInfo
						break
					}
				}
				if source == nil {
					source = NewSourceInfo(ssrc)
					sources = append(sources, source)
				}

				if strings.ToLower(key) == "cname" {
					source.SetCName(value)
				} else if strings.ToLower(key) == "msid" {
					ids := strings.Split(value, " ")
					// get stream id and track id
					streamId := ids[0]
					trackId := ids[1]

					source.SetStreamID(streamId)
					source.SetTrackID(trackId)

					stream := sdpInfo.GetStream(streamId)

					if stream == nil {
						stream = NewStreamInfo(streamId)
						sdpInfo.AddStream(stream)
					}

					track := stream.GetTrack(trackId)

					if track == nil {
						track = NewTrackInfo(trackId, media)
						track.SetMID(mid)
						track.SetEncodings(encodings)
						stream.AddTrack(track)
					}
					// Add ssrc
					track.AddSSRC(ssrc)

				}

			}
		}

		// Check if ther is a global msid
		// Why this?
		if md.Msid != "" {
			ids := strings.Split(md.Msid, " ")
			streamId := ids[0]
			trackId := ids[1]

			stream := sdpInfo.GetStream(streamId)

			if stream == nil {
				stream = NewStreamInfo(streamId)
				sdpInfo.AddStream(stream)
			}

			track := stream.GetTrack(trackId)

			if track == nil {
				track = NewTrackInfo(trackId, media)
				track.SetMID(mid)
				track.SetEncodings(encodings)
				stream.AddTrack(track)
			}

			for _, source := range sources {

				if source.GetStreamID() == "" {
					source.SetStreamID(streamId)
					source.SetTrackID(trackId)
					track.AddSSRC(source.ssrc)
				}
			}
		}

		for _, source := range sources {

			if source.GetStreamID() == "" {
				streamId := source.GetCName()
				trackId := mid

				source.SetStreamID(streamId)
				source.SetTrackID(trackId)

				stream := sdpInfo.GetStream(streamId)

				if stream == nil {
					stream = NewStreamInfo(streamId)
					sdpInfo.AddStream(stream)
				}

				track := stream.GetTrack(trackId)

				if track == nil {
					track = NewTrackInfo(trackId, media)
					track.SetMID(mid)
					track.SetEncodings(encodings)
					stream.AddTrack(track)
				}

				track.AddSSRC(source.ssrc)
			}
		}

		if md.SsrcGroups != nil {
			for _, ssrcGroupAttr := range md.SsrcGroups {
				ssrcs := strings.Split(ssrcGroupAttr.Ssrcs, " ")
				ssrcsint := []uint{}
				for _, ssrcstr := range ssrcs {
					ssrcint, _ := strconv.ParseUint(ssrcstr, 10, 32)
					ssrcsint = append(ssrcsint, uint(ssrcint))
				}
				group := NewSourceGroupInfo(ssrcGroupAttr.Semantics, ssrcsint)
				ssrc := ssrcsint[0]
				for _, source := range sources {
					if source.ssrc == ssrc {
						streamInfo := sdpInfo.GetStream(source.GetStreamID())
						if streamInfo != nil && streamInfo.GetTrack(source.GetTrackID()) != nil {
							streamInfo.GetTrack(source.GetTrackID()).AddSourceGroup(group)
						}
						break
					}
				}
			}
		}

		sdpInfo.AddMedia(mediaInfo)

	}

	return sdpInfo, nil
}

