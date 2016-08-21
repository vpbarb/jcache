package client

import (
	"bufio"
	"errors"
	"log"
	"regexp"
	"strconv"
)

const endResponse = "END"

type dataFormatFunc func(*bufio.Reader) error

var (
	invalidDataFormatError = errors.New("Invalid data format")

	keyRegexp       = regexp.MustCompile("^KEY ([a-zA-Z0-9_]+)$")
	valueRegexp     = regexp.MustCompile("^VALUE ([0-9]+)$")
	hashFieldRegexp = regexp.MustCompile("^FIELD ([a-zA-Z0-9_]+) ([0-9]+)$")
	lenRegexp       = regexp.MustCompile("^LEN ([0-9]+)$")
	ttlRegexp       = regexp.MustCompile("^TTL ([0-9]+)$")
)

func emptyDataFormatter() dataFormatFunc {
	return func(response *bufio.Reader) (err error) {
		line, _, err := response.ReadLine()
		if err != nil {
			return err
		}
		str := string(line)
		if str == endResponse {
			return nil
		}
		return invalidDataFormatError
	}
}

func keysDataFormatter(keys *[]string) dataFormatFunc {
	return func(response *bufio.Reader) (err error) {
		var result []string
		for {
			line, _, err := response.ReadLine()
			if err != nil {
				return err
			}
			str := string(line)
			if str == endResponse {
				*keys = result
				return nil
			}
			matches := keyRegexp.FindStringSubmatch(str)
			log.Printf("Key string: %s, %+v", str, matches)
			if len(matches) == 0 {
				return invalidDataFormatError
			}
			result = append(result, matches[1])
		}
		return nil
	}
}

func valueDataFormatter(value *string) dataFormatFunc {
	return func(response *bufio.Reader) (err error) {
		var result string
		for {
			line, _, err := response.ReadLine()
			if err != nil {
				return err
			}
			str := string(line)
			if str == endResponse {
				*value = result
				return nil
			}
			matches := valueRegexp.FindStringSubmatch(str)
			if len(matches) == 0 {
				return invalidDataFormatError
			}
			length, err := strconv.Atoi(matches[1])
			if err != nil {
				return err
			}
			data := make([]byte, length, length)
			n, err := response.Read(data)
			if err != nil || n != length {
				return err
			}
			result = string(data)
			response.ReadLine()
		}
		return nil
	}
}

func valuesDataFormatter(values *[]string) dataFormatFunc {
	return func(response *bufio.Reader) (err error) {
		var result []string
		for {
			line, _, err := response.ReadLine()
			if err != nil {
				return err
			}
			str := string(line)
			if str == endResponse {
				*values = result
				return nil
			}
			matches := valueRegexp.FindStringSubmatch(str)
			if len(matches) == 0 {
				return invalidDataFormatError
			}
			length, err := strconv.Atoi(matches[1])
			if err != nil {
				return err
			}
			data := make([]byte, length, length)
			n, err := response.Read(data)
			if err != nil || n != length {
				return err
			}
			result = append(result, string(data))
			response.ReadLine()
		}
		return nil
	}
}

func ttlDataFormatter(ttl *uint64) dataFormatFunc {
	return func(response *bufio.Reader) (err error) {
		var result uint64
		for {
			line, _, err := response.ReadLine()
			if err != nil {
				return err
			}
			str := string(line)
			if str == endResponse {
				*ttl = result
				return nil
			}
			matches := ttlRegexp.FindStringSubmatch(str)
			if len(matches) == 0 {
				return invalidDataFormatError
			}
			data, err := strconv.Atoi(matches[1])
			if err != nil {
				return err
			}
			result = uint64(data)
		}
		return nil
	}
}

func hashDataFormatter(hash map[string]string) dataFormatFunc {
	return func(response *bufio.Reader) (err error) {
		result := make(map[string]string)
		for {
			line, _, err := response.ReadLine()
			if err != nil {
				return err
			}
			str := string(line)
			if str == endResponse {
				hash = result
				return nil
			}
			matches := hashFieldRegexp.FindStringSubmatch(str)
			if len(matches) == 0 {
				return invalidDataFormatError
			}
			length, err := strconv.Atoi(matches[2])
			if err != nil {
				return err
			}
			data := make([]byte, length, length)
			n, err := response.Read(data)
			if err != nil || n != length {
				return err
			}
			result[matches[1]] = string(data)
			response.ReadLine()
		}
		return nil
	}
}

func lenDataFormatter(length *uint64) dataFormatFunc {
	return func(response *bufio.Reader) (err error) {
		var result uint64
		for {
			line, _, err := response.ReadLine()
			if err != nil {
				return err
			}
			str := string(line)
			if str == endResponse {
				*length = result
				return nil
			}
			matches := lenRegexp.FindStringSubmatch(str)
			if len(matches) == 0 {
				return invalidDataFormatError
			}
			data, err := strconv.Atoi(matches[1])
			if err != nil {
				return err
			}
			result = uint64(data)
		}
		return nil
	}
}
