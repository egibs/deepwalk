package deepwalk

import (
	"crypto/rand"
	"math/big"
	"strings"
)

var maxDepth = 10

// RandomInt generate a cryptographically secure random integer
func RandomInt(min, max int64) (int64, error) {
	if max-min <= 0 {
		return 0, nil
	}
	diff := max - min
	bigRand, err := rand.Int(rand.Reader, big.NewInt(diff))
	if err != nil {
		return 0, err
	}
	return min + bigRand.Int64(), nil
}

// RandomString generate random string based on a length argument
func RandomString(length int) (string, error) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, length)
	for i := range s {
		randomIndex, err := RandomInt(1, int64(len(letters)))
		if err != nil {
			return "", err
		}
		s[i] = letters[randomIndex]
	}
	return string(s), nil
}

// RandomKVPair generate a random key-value pair for use in a map
func RandomKVPair() (map[string]interface{}, error) {
	kvPair := make(map[string]interface{})
	j, err := RandomInt(1, 10)
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(j); i++ {
		key := strings.Builder{}
		k, err := RandomInt(1, 10)
		if err != nil {
			return nil, err
		}
		for j := 0; j < int(k); j++ {
			randomSeq, err := RandomString(12)
			if err != nil {
				return nil, err
			}
			key.WriteString(randomSeq)
		}
		choice, err := RandomInt(0, 2)
		if err != nil {
			return nil, err
		}
		if int(choice) == 0 {
			value := strings.Builder{}
			l, err := RandomInt(1, 10)
			if err != nil {
				return nil, err
			}
			for j := 0; j < int(l); j++ {
				randomSeq, err := RandomString(12)
				if err != nil {
					return nil, err
				}
				value.WriteString(randomSeq)
			}
			kvPair[key.String()] = value.String()
		} else {
			randomInt, err := RandomInt(1, 12)
			if err != nil {
				return nil, err
			}
			randomInt2, err := RandomInt(1, randomInt)
			if err != nil {
				return nil, err
			}
			kvPair[key.String()] = randomInt2
		}
	}
	return kvPair, nil
}

// SucessCases generate a random map, list of keys, and a expected value
func SucessCases(keys []string, depth int) (map[string]interface{}, []string, interface{}, error) {
	kvPair, err := RandomKVPair()
	if err != nil {
		return nil, nil, nil, err
	}
	key := ""
	for k := range kvPair {
		key = k
	}
	if depth == maxDepth {
		return kvPair, append(keys, key), kvPair[key], nil
	}
	nested, keys, expected, err := SucessCases(append(keys, key), depth+1)
	if err != nil {
		return nil, nil, nil, err
	}
	kvPair[key] = nested
	choice, err := RandomInt(0, 2)
	if err != nil {
		return nil, nil, nil, err
	}
	if int(choice) == 0 {
		listLength, err := RandomInt(1, 10)
		if err != nil {
			return nil, nil, nil, err
		}
		list := make([]interface{}, listLength)
		for i := range list {
			list[i] = nested
		}
		kvPair[key] = list
	}
	return kvPair, keys, expected, nil
}

// DefaultCases generate a random map or key-value pair for later use when generating test cases
func DefaultCases(depth int, maxDepth int) (interface{}, []interface{}, error) {
	kvPair, err := RandomKVPair()
	if err != nil {
		return nil, nil, err
	}
	key := ""
	for k := range kvPair {
		key = k
	}
	choice, err := RandomInt(0, 2)
	if err != nil {
		return nil, nil, err
	}
	if int(choice) == 0 {
		if depth != maxDepth {
			nested, _, err := DefaultCases(depth+1, maxDepth)
			if err != nil {
				return nil, nil, err
			}
			return []interface{}{map[string]interface{}{key: nested}}, nil, nil
		}
		return []interface{}{kvPair}, nil, nil
	}
	if depth != maxDepth {
		nested, _, err := DefaultCases(depth+1, maxDepth)
		if err != nil {
			return nil, nil, err
		}
		return map[string]interface{}{key: nested}, nil, nil
	}
	return kvPair, nil, nil
}
