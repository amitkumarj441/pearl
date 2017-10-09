package pearl

import (
	"testing"

	"github.com/mmcloughlin/pearl/torcrypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMethodString(t *testing.T) {
	assert.Equal(t, "AUTH0003", AuthMethodEd25519SHA256RFC5705.String())
}

func TestAuthenticateCellRoundTrip(t *testing.T) {
	data := []byte{
		0x00, 0x00, 0x00, 0x00, 0x83, 0x01, 0x64, 0x00, 0x01, 0x01, 0x60, 0x41,
		0x55, 0x54, 0x48, 0x30, 0x30, 0x30, 0x31, 0x81, 0xe4, 0x71, 0x36, 0x1d,
		0x86, 0x96, 0x47, 0x49, 0x72, 0x0f, 0x6e, 0x79, 0x00, 0x0d, 0xfa, 0xa8,
		0x8f, 0x83, 0x4a, 0x67, 0x41, 0x5c, 0xc0, 0x34, 0xfe, 0xa2, 0xc1, 0x35,
		0xe7, 0x84, 0xcb, 0x87, 0x2b, 0x11, 0x3e, 0x05, 0x85, 0x31, 0x4a, 0x25,
		0x5c, 0x66, 0x94, 0x11, 0x2a, 0x18, 0xff, 0x1c, 0xcb, 0x2c, 0x5b, 0x40,
		0x26, 0xfc, 0x03, 0x2e, 0x8a, 0xa4, 0x01, 0xad, 0x92, 0xb4, 0x74, 0xb4,
		0xa0, 0xcf, 0xad, 0x1b, 0x16, 0xd5, 0x10, 0xbf, 0x67, 0x23, 0xfb, 0x3e,
		0x7a, 0x88, 0xea, 0x5b, 0x27, 0x13, 0x00, 0x65, 0x73, 0x8c, 0x14, 0xe2,
		0xca, 0x50, 0xc5, 0x6c, 0x3c, 0xa6, 0xa8, 0xbc, 0xc2, 0x41, 0x02, 0x2a,
		0xe3, 0x97, 0x32, 0x51, 0x99, 0xaa, 0xb7, 0x5d, 0x86, 0xd5, 0xc7, 0xe8,
		0x5e, 0x24, 0x22, 0xab, 0x5c, 0xaf, 0xe0, 0x1e, 0x30, 0x96, 0xe3, 0x0f,
		0x27, 0xd6, 0x5b, 0xef, 0x8f, 0x62, 0x83, 0x0a, 0x48, 0x45, 0x57, 0x45,
		0x1c, 0x3b, 0x28, 0x1f, 0x06, 0x6a, 0xbc, 0xf2, 0x38, 0xae, 0x86, 0xde,
		0xbe, 0xeb, 0x04, 0x29, 0x4c, 0xb9, 0x6b, 0x30, 0xe6, 0xad, 0x30, 0x25,
		0x5c, 0x5b, 0xec, 0xa4, 0xc4, 0x72, 0xed, 0xca, 0xbb, 0x65, 0xa8, 0x67,
		0x2c, 0x8a, 0x9d, 0xfa, 0xf0, 0x63, 0x4d, 0xd6, 0xc7, 0x8b, 0x9c, 0xba,
		0xd7, 0x80, 0xc5, 0x1c, 0xc6, 0x0b, 0xbc, 0x8b, 0x91, 0x3d, 0xf8, 0xb0,
		0x81, 0x14, 0x68, 0x3a, 0x56, 0x1a, 0xd2, 0xef, 0xfd, 0xe1, 0xea, 0x4b,
		0xcc, 0xe4, 0xf5, 0xc0, 0x01, 0x5e, 0x75, 0x3e, 0xe2, 0xdc, 0x16, 0x0f,
		0x59, 0x29, 0x30, 0x5d, 0x48, 0x3d, 0x52, 0x0b, 0x99, 0x9b, 0x26, 0x7a,
		0x19, 0x43, 0x7b, 0x73, 0x42, 0xd2, 0x0c, 0x61, 0x65, 0x39, 0xa6, 0x9d,
		0x60, 0x0d, 0xc8, 0x92, 0x40, 0x2d, 0x88, 0xbc, 0x69, 0xfd, 0xad, 0x39,
		0xcf, 0xb1, 0xa8, 0xde, 0x4d, 0x80, 0xb9, 0x5e, 0xd6, 0x2e, 0xd1, 0xb9,
		0x63, 0x97, 0x87, 0x2e, 0xe9, 0xb2, 0xeb, 0x7e, 0x99, 0xaf, 0xf5, 0xe5,
		0xe6, 0x8f, 0xb8, 0xac, 0x8b, 0xd6, 0x03, 0x6f, 0x5e, 0x8b, 0x5b, 0xa7,
		0x8e, 0x7c, 0x86, 0x33, 0x9c, 0xa1, 0x37, 0x8a, 0x5f, 0x0b, 0x0d, 0xbf,
		0x75, 0xff, 0xdc, 0xa2, 0xfd, 0x9d, 0x32, 0x55, 0x65, 0x0c, 0x8a, 0xe6,
		0x6c, 0x4e, 0xfc, 0xfa, 0xaa, 0x29, 0x80, 0x0a, 0x12, 0x93, 0xd6, 0xaa,
		0x89, 0xff, 0xdb, 0x06, 0x65, 0xf0, 0xee, 0xb2, 0x55, 0xfa, 0x84, 0x7a,
		0xdf, 0x65, 0xa7,
	}
	f := CircID4Format{}
	c := NewCellFromBuffer(f, data)
	a, err := ParseAuthenticateCell(c)
	require.NoError(t, err)
	c2, err := a.Cell(f)
	require.NoError(t, err)
	assert.Equal(t, c.Bytes(), c2.Bytes())
}

func TestAuthRSASHA256TLSSecret(t *testing.T) {
	a := AuthRSASHA256TLSSecret{
		AuthKey:           torcrypto.MustRSAPrivateKey(torcrypto.ParseRSAPrivateKeyPKCS1PEM([]byte("-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQDHtMM+7VEvWllFC7xoW96CaSIkgCOJiNtCKylUV86iD3qziLzE\nXQWgEecDmM5urbu+3tcpLVMqPbCp3gxzkdNozql1eydV0+JUw2AI3Nhbv89cppBA\n3W+MhckQ1VmMlaiJLg9xTOWClAuy4jQzdVnj5QKIi7W3ZT/UvSzvDkP9WwIDAQAB\nAoGBAKAr38jRqCKVkTGqlwMQY+cukT67M0V06X4phe1qu4UJaz0hd1z6yq82jJU6\n8p6cYw9URTd2bdRcRBwJxuzOUcK8AvRUUA7TXU8dG0/6pF5ScI+E2VKvBHgGIXQM\ni+Meogk2Fkt4RoVQRPobFxgXfsp8d6/pCX+MBMxE7F1VYHrZAkEA4oyTEr05UwHC\nMh7xWO6RZtzGvnmuux1FhtWqbNHLcgcggzv6UcvyH0s+R1hjpjaiT/dXk/PO9UaD\nJlFNQ/MNRwJBAOGq3jGXjQ4Y3dTqeOrlH/MYOUuDHlcFzY5HIpB8ptT4Al11R4B/\nqdElTI5Ej/EAdmebf29vOeL0yvHvaMKCiU0CQG4yPp/Q1v9fTZyfnHnLoYJNRYcF\nHU760ATkDX/dFH6kpNXw6LO85kr+iI6fmekRjiYjg7/9yd9YqxaKWXEB2qUCQGyq\nYNA0kAHHy5opRgymRFpEweIwwz1YWAE5E9XLkHJg8pKaVNH1p4pEkba4ITAF7v45\nDIZWYuN8yPTzOdjgDskCQBqkqe1wupf7InCHtRq9UwnB3s3nsbcgmJ80igWfjrGa\nHr3hF+LrpR3nWVwuZcsAcDb4xAI6KvEuFDZ1l+no5m0=\n-----END RSA PRIVATE KEY-----\n"))),
		ClientIdentityKey: torcrypto.MustRSAPublicKey(torcrypto.ParseRSAPublicKeyPKCS1PEM([]byte("-----BEGIN RSA PUBLIC KEY-----\nMIGJAoGBALaKBJ/sK8zr+0j7ih0YWk7jHLDYnZSBvseoRmUfTOuxkj8LOce8X/GG\nLPYMFJUTNL0ToQApC6TqbEuShzQyQLk9IHWRhVsmSDKYjLZepzdsvJx8gL5QaHea\nf5Ge3nmo+oUKdeX3rDQd07us/nLja3VUL2xKdd+hE81KMxhTjG4RAgMBAAE=\n-----END RSA PUBLIC KEY-----\n"))),
		ServerIdentityKey: torcrypto.MustRSAPublicKey(torcrypto.ParseRSAPublicKeyPKCS1PEM([]byte("-----BEGIN RSA PUBLIC KEY-----\nMIGJAoGBALMlpknZ4yhwp7TcjAAZjIcgyjqjSd4BJqbLWvhEFWvM5rhO+DWkLfuM\nssdS6FimnN5oItUYVx0W4RPKyuVeqdUK0F2gj+yVtgA5cUXAhhrJUQp4o4JBFrH3\ntivLapYfvNvhpT/Xo6kBeu29LwxYWgVYrKAK/d9RRVE9lJ1SOxuHAgMBAAE=\n-----END RSA PUBLIC KEY-----\n"))),
		ServerLogHash: []byte{
			0x3d, 0xff, 0xfa, 0x23, 0x44, 0x11, 0x38, 0x21, 0x67, 0x22, 0xb5, 0xe3, 0xdb, 0xb8, 0x66, 0x44,
			0x87, 0x0e, 0x41, 0x15, 0x72, 0x96, 0xde, 0x70, 0x14, 0xe8, 0xc4, 0x72, 0x99, 0x96, 0x8c, 0xa9,
		},
		ClientLogHash: []byte{
			0x36, 0x3f, 0xce, 0x3f, 0x6b, 0xd7, 0x2c, 0x9c, 0x25, 0x25, 0xd1, 0x07, 0x47, 0x3a, 0x97, 0xb5,
			0xbc, 0x1a, 0xfd, 0xba, 0xce, 0xe7, 0xb3, 0xde, 0x3c, 0xdf, 0x01, 0xa0, 0xd1, 0x4c, 0x70, 0xaf,
		},
		ServerLinkCert: []byte{
			0x30, 0x82, 0x02, 0x45, 0x30, 0x82, 0x01, 0xae, 0xa0, 0x03, 0x02, 0x01, 0x02, 0x02, 0x08, 0x4e,
			0xbe, 0xe3, 0xc7, 0xa4, 0xcc, 0x0a, 0x73, 0x30, 0x0d, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7,
			0x0d, 0x01, 0x01, 0x0b, 0x05, 0x00, 0x30, 0x22, 0x31, 0x20, 0x30, 0x1e, 0x06, 0x03, 0x55, 0x04,
			0x03, 0x0c, 0x17, 0x77, 0x77, 0x77, 0x2e, 0x64, 0x73, 0x68, 0x7a, 0x6c, 0x78, 0x69, 0x67, 0x6f,
			0x71, 0x63, 0x36, 0x7a, 0x7a, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x30, 0x1e, 0x17, 0x0d, 0x31, 0x37,
			0x30, 0x35, 0x31, 0x35, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x5a, 0x17, 0x0d, 0x31, 0x38, 0x30,
			0x35, 0x30, 0x36, 0x32, 0x33, 0x35, 0x39, 0x35, 0x39, 0x5a, 0x30, 0x24, 0x31, 0x22, 0x30, 0x20,
			0x06, 0x03, 0x55, 0x04, 0x03, 0x0c, 0x19, 0x77, 0x77, 0x77, 0x2e, 0x67, 0x79, 0x36, 0x65, 0x34,
			0x35, 0x6e, 0x68, 0x75, 0x35, 0x76, 0x62, 0x7a, 0x61, 0x33, 0x37, 0x68, 0x2e, 0x6e, 0x65, 0x74,
			0x30, 0x82, 0x01, 0x22, 0x30, 0x0d, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01,
			0x01, 0x05, 0x00, 0x03, 0x82, 0x01, 0x0f, 0x00, 0x30, 0x82, 0x01, 0x0a, 0x02, 0x82, 0x01, 0x01,
			0x00, 0xa6, 0x0f, 0x46, 0xf4, 0xff, 0xd5, 0xe1, 0x22, 0xdc, 0x5e, 0x8e, 0x7c, 0x3c, 0x54, 0x61,
			0xf8, 0xdc, 0x56, 0xda, 0x93, 0x55, 0xf3, 0x40, 0x14, 0xf8, 0xc7, 0xca, 0x6b, 0xea, 0x26, 0xf5,
			0x24, 0x55, 0x9e, 0x34, 0x1d, 0x7f, 0x86, 0xd8, 0x0c, 0xb0, 0x01, 0xdb, 0xb6, 0x12, 0x99, 0xe1,
			0xc4, 0x2e, 0x03, 0xb7, 0x32, 0x59, 0x49, 0xc5, 0xb0, 0x02, 0x0d, 0x51, 0x2f, 0xf1, 0xf3, 0x15,
			0x5d, 0xc0, 0x5d, 0x49, 0x76, 0xab, 0xad, 0xf9, 0xbb, 0x2a, 0x53, 0xb2, 0x58, 0x24, 0xf4, 0x90,
			0x22, 0xea, 0xff, 0xa3, 0x53, 0xae, 0x41, 0x18, 0xee, 0x82, 0x99, 0x3b, 0x0d, 0x12, 0x67, 0x90,
			0x25, 0x25, 0x04, 0x55, 0x2f, 0x72, 0xca, 0x21, 0x7f, 0xc5, 0x58, 0xab, 0x58, 0x66, 0x16, 0x11,
			0x54, 0x24, 0xc9, 0x24, 0xf5, 0x0a, 0x86, 0xef, 0x12, 0x43, 0xad, 0x88, 0x71, 0x21, 0x81, 0xe8,
			0x6f, 0x1f, 0x95, 0x1a, 0x4a, 0x1f, 0x57, 0x7c, 0x0e, 0x4b, 0x99, 0x7b, 0x18, 0x0c, 0xe4, 0x87,
			0xbd, 0xee, 0x6c, 0x60, 0x37, 0xd3, 0x02, 0x71, 0x32, 0x9f, 0x7e, 0x88, 0xac, 0x22, 0x36, 0x9e,
			0xad, 0x29, 0x2d, 0xcc, 0xf5, 0xd0, 0x8f, 0xf8, 0x26, 0x7f, 0x3b, 0x43, 0xbc, 0x30, 0x54, 0x15,
			0xab, 0x46, 0x73, 0x8c, 0x85, 0x4d, 0xa4, 0x49, 0xb0, 0x21, 0x4a, 0xde, 0xe8, 0xf0, 0x5c, 0x14,
			0x4e, 0x5f, 0xae, 0x3d, 0x28, 0xd6, 0xc1, 0x1c, 0x02, 0x2d, 0x3c, 0xed, 0xf5, 0xbb, 0x1e, 0x69,
			0x21, 0x23, 0x21, 0x13, 0xdf, 0x98, 0x96, 0xbf, 0x64, 0x33, 0x7e, 0x04, 0xa9, 0x6c, 0xea, 0x08,
			0x25, 0x05, 0x07, 0x0f, 0x1f, 0xb9, 0x60, 0x47, 0xbb, 0x32, 0x3d, 0x7a, 0x41, 0x62, 0x40, 0x11,
			0x26, 0xdf, 0x0e, 0xfe, 0x3b, 0x62, 0x17, 0x6d, 0x9a, 0xaf, 0x2b, 0xc3, 0xaa, 0x66, 0xc5, 0x3c,
			0x0d, 0x02, 0x03, 0x01, 0x00, 0x01, 0x30, 0x0d, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d,
			0x01, 0x01, 0x0b, 0x05, 0x00, 0x03, 0x81, 0x81, 0x00, 0x3f, 0x2d, 0x88, 0x9d, 0x1a, 0x6f, 0xec,
			0x3c, 0x8a, 0x91, 0x93, 0x48, 0xe3, 0xd2, 0x3c, 0x69, 0xd3, 0x90, 0x62, 0x44, 0x19, 0xe8, 0xc6,
			0x5f, 0x6f, 0x0e, 0x61, 0xe4, 0xbb, 0x68, 0x25, 0xaa, 0x71, 0xa2, 0x1d, 0x55, 0x09, 0x75, 0x94,
			0x90, 0xd0, 0xb6, 0xcf, 0xe3, 0xc8, 0xbd, 0xaf, 0x97, 0x68, 0xd8, 0x5b, 0xd7, 0xbc, 0xd3, 0x40,
			0xcc, 0x13, 0x80, 0xcd, 0xa9, 0xe7, 0xfc, 0x98, 0xe7, 0x7d, 0xc2, 0xdf, 0x17, 0xcc, 0x66, 0x7e,
			0xa6, 0xe4, 0x31, 0xd7, 0xe9, 0xfb, 0x30, 0x4a, 0x0c, 0x19, 0x4f, 0x7a, 0x21, 0x6f, 0x77, 0x06,
			0xd5, 0xa2, 0x79, 0x62, 0x1d, 0xde, 0xdd, 0x82, 0xfc, 0x2c, 0x0b, 0x07, 0x3c, 0xcc, 0x71, 0xc4,
			0x29, 0x58, 0x64, 0x90, 0x9f, 0x4a, 0x63, 0x69, 0x7d, 0x3e, 0x3b, 0x58, 0x51, 0x41, 0x7e, 0x1d,
			0x91, 0x4f, 0xde, 0x54, 0xb4, 0xb1, 0x24, 0x54, 0x2c,
		},
		TLSMasterSecret: []byte{
			0xb3, 0x2a, 0x03, 0xf4, 0x8d, 0xdd, 0x8a, 0x27, 0x29, 0x21, 0xb2, 0x54, 0x44, 0xa4, 0x99, 0x24,
			0x17, 0x64, 0x72, 0x31, 0xf7, 0x93, 0x3f, 0x0c, 0x1e, 0x3c, 0x28, 0xc1, 0x53, 0xc9, 0x40, 0xfc,
			0x4e, 0x92, 0x47, 0x81, 0x1a, 0x8d, 0x88, 0x97, 0xd4, 0xd1, 0x65, 0x36, 0x38, 0x4d, 0xfd, 0x13,
		},
		TLSClientRandom: []byte{
			0xf0, 0xed, 0x85, 0x5c, 0x8e, 0xe0, 0x6d, 0x22, 0xa6, 0xdf, 0x4d, 0x36, 0x7e, 0xfa, 0x21, 0x8b,
			0x15, 0x93, 0x71, 0xb9, 0x4a, 0xa7, 0xf6, 0x6c, 0x73, 0x16, 0x09, 0x38, 0x33, 0x43, 0xb9, 0xb4,
		},
		TLSServerRandom: []byte{
			0xca, 0xc6, 0xa4, 0xe5, 0xeb, 0x1e, 0xaf, 0x35, 0xec, 0xed, 0x76, 0xf2, 0x8b, 0x0b, 0xa0, 0x87,
			0x70, 0x8d, 0xe1, 0x00, 0xa7, 0x60, 0x7f, 0xe5, 0x6a, 0x49, 0xe1, 0x08, 0xe6, 0x7a, 0x36, 0x77,
		},
	}

	expect := []byte{
		0x41, 0x55, 0x54, 0x48, 0x30, 0x30, 0x30, 0x31, 0x2a, 0x90, 0x7f, 0x75,
		0x1f, 0x6a, 0xd4, 0x1d, 0xac, 0xb6, 0x8a, 0x23, 0x1d, 0xcd, 0x27, 0x86,
		0xd5, 0x32, 0xac, 0xc7, 0x1c, 0x4b, 0x25, 0xfe, 0x86, 0x42, 0xe0, 0x2d,
		0x25, 0x29, 0x8b, 0x71, 0x91, 0x3b, 0x62, 0x54, 0xb6, 0x9f, 0xfc, 0xad,
		0xc3, 0x40, 0x69, 0x9c, 0xd5, 0x4a, 0xd7, 0x1d, 0xfc, 0xcd, 0x26, 0x8e,
		0xa8, 0xed, 0xa2, 0x3a, 0x81, 0x68, 0x1f, 0x00, 0x10, 0x37, 0x47, 0xf3,
		0x3d, 0xff, 0xfa, 0x23, 0x44, 0x11, 0x38, 0x21, 0x67, 0x22, 0xb5, 0xe3,
		0xdb, 0xb8, 0x66, 0x44, 0x87, 0x0e, 0x41, 0x15, 0x72, 0x96, 0xde, 0x70,
		0x14, 0xe8, 0xc4, 0x72, 0x99, 0x96, 0x8c, 0xa9, 0x36, 0x3f, 0xce, 0x3f,
		0x6b, 0xd7, 0x2c, 0x9c, 0x25, 0x25, 0xd1, 0x07, 0x47, 0x3a, 0x97, 0xb5,
		0xbc, 0x1a, 0xfd, 0xba, 0xce, 0xe7, 0xb3, 0xde, 0x3c, 0xdf, 0x01, 0xa0,
		0xd1, 0x4c, 0x70, 0xaf, 0xc4, 0xd7, 0x72, 0xe8, 0x60, 0xcc, 0x74, 0xde,
		0xa7, 0x93, 0xfb, 0xd5, 0xd6, 0x39, 0xe0, 0x16, 0x1f, 0x4c, 0x0d, 0x42,
		0x23, 0x6c, 0xb8, 0x17, 0x9d, 0xd0, 0xdd, 0x29, 0x88, 0x94, 0x7e, 0x7c,
		0x1b, 0xbc, 0xae, 0xdf, 0xb8, 0x38, 0xc7, 0x25, 0x9a, 0x39, 0x46, 0xc6,
		0x42, 0x3c, 0x3d, 0x96, 0xe1, 0xd3, 0xb8, 0x21, 0xbb, 0x5b, 0xef, 0x18,
		0xdd, 0x21, 0xe2, 0x41, 0xdc, 0x5f, 0x6e, 0xd1, 0x5b, 0x9c, 0xd0, 0xc4,
		0x8a, 0x90, 0xc6, 0x1a, 0x60, 0x8d, 0x94, 0x6d, 0xe0, 0xfa, 0xc5, 0x90,
		0x7c, 0x25, 0x74, 0x34, 0x8d, 0xb7, 0x0e, 0x99,
	}

	body, err := a.Body()
	require.NoError(t, err)
	// assert equal up to the last 24 bytes of random
	assert.Equal(t, expect[:len(expect)-24], body[:len(body)-24])
}
