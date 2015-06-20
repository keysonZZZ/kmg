package MSCHAPV2

import (
	. "github.com/bronze1man/kmgTest"
	"testing"
)

// rfc2759 Page 15 section 9.2 "9.2.  Hash Example"
func TestGenerateAuthenticatorResponse(ot *testing.T) {
	//Password:=[]byte("clientPass")
	Password := []byte("clientPass") //utf16 of "clientPass"
	UserName := []byte("User")
	AuthenticatorChallenge := [16]byte{0x5B, 0x5D, 0x7C, 0x7D, 0x7B, 0x3F, 0x2F, 0x3E, 0x3C, 0x2C, 0x60, 0x21, 0x32, 0x26, 0x26, 0x28}
	PeerChallenge := [16]byte{0x21, 0x40, 0x23, 0x24, 0x25, 0x5E, 0x26, 0x2A, 0x28, 0x29, 0x5F, 0x2B, 0x3A, 0x33, 0x7C, 0x7E}
	NTResponse := [24]byte{0x82, 0x30, 0x9E, 0xCD, 0x8D, 0x70, 0x8B, 0x5E, 0xA0, 0x8F, 0xAA, 0x39, 0x81, 0xCD, 0x83, 0x54, 0x42, 0x33, 0x11,
		0x4A, 0x3D, 0x85, 0xD6, 0xDF}
	AuthenticatorResponse := [20]byte{0x40, 0x7A, 0x55, 0x89, 0x11, 0x5F, 0xD0, 0xD6, 0x20, 0x9F, 0x51, 0x0F, 0xE9, 0xC0, 0x45, 0x66, 0x93,
		0x2C, 0xDA, 0x56}
	Challenge := []byte{0xD0, 0x2E, 0x43, 0x86, 0xBC, 0xE9, 0x12, 0x26}
	out0 := challengeHash(PeerChallenge, AuthenticatorChallenge, UserName)
	Equal(out0, Challenge)

	//PasswordHash := []byte{0x44, 0xEB, 0xBA, 0x8D, 0x53, 0x12, 0xB8, 0xD6, 0x11, 0x47, 0x44, 0x11, 0xF5, 0x69, 0x89, 0xAE}
	//out01 := md4(Password)
	//Equal(out01, PasswordHash)

	out1 := GenerateAuthenticatorResponse(Password, NTResponse, PeerChallenge, AuthenticatorChallenge, UserName)
	Equal(out1, AuthenticatorResponse)
}

func TestMsCHAPV2GetSendAndRecvKey(ot *testing.T) {
	Password := []byte("clientPass")
	NTResponse := [24]byte{0x82, 0x30, 0x9E, 0xCD, 0x8D, 0x70, 0x8B, 0x5E, 0xA0, 0x8F, 0xAA, 0x39, 0x81, 0xCD, 0x83, 0x54, 0x42, 0x33, 0x11,
		0x4A, 0x3D, 0x85, 0xD6, 0xDF}
	sendKey, recvKey := MsCHAPV2GetSendAndRecvKey(Password, NTResponse)
	Equal(sendKey, []byte{0x8B, 0x7C, 0xDC, 0x14, 0x9B, 0x99, 0x3A, 0x1B, 0xA1, 0x18, 0xCB, 0x15, 0x3F, 0x56, 0xDC, 0xCB})
	// 官方文档上没有写recvkey的值.
	Equal(recvKey, []byte{0xd5, 0xf0, 0xe9, 0x52, 0x1e, 0x3e, 0xa9, 0x58, 0x96, 0x45, 0xe8, 0x60, 0x51, 0xc8, 0x22, 0x26})
}
