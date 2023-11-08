package ebcm

import "strings"

const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength = 20

	AddressZERO = "0x0000000000000000000000000000000000000000"
)

func IsAddress(address string) bool {
	if !strings.HasPrefix(address, "0x") {
		return false
	}
	return EIP55(address) != ""
}
func IsAddressP(address *string) bool {
	if IsAddress(*address) {
		*address = strings.TrimSpace(strings.ToLower(*address))
		return true
	}
	return false
}

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// isHex validates whether each byte is valid hexadecimal string.
func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

// IsHexAddress verifies whether a string can represent a valid hex-encoded
// Ethereum address or not.
func _IsHexAddress(s string) bool {
	if has0xPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && isHex(s)
}

//IsHashHex : check hash format (check 32bytes)
func IsHashHex(hash string) bool {
	//0x6f660bdca8a542968afbc25574a3caf65e6509cefaf259561c18904ab7851c73
	if !strings.HasPrefix(hash, "0x") {
		return false
	}

	for _, c := range []byte(hash[2:]) {
		if !isHexCharacter(c) {
			return false
		}
	} //for

	return len(hash) == 66
}
