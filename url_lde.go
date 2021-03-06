/*
 This file was autogenerated via
 -------------------------------------------
 ldetool generate --package chinsert insertURL.lde
 -------------------------------------------
 do not touch it with bare hands!
*/

package chinsert

import (
	"bytes"
	"fmt"
	"strconv"
	"unsafe"
)

// URL ...
type URL struct {
	rest []byte
	Auth struct {
		Valid bool
		Data  []byte
	}
	Host   []byte
	Port   uint16
	DBName []byte
}

// Extract ...
func (p *URL) Extract(line []byte) (bool, error) {
	p.rest = line
	var err error
	var pos int
	var rest1 []byte
	var tmp []byte
	var tmpUint uint64
	rest1 = p.rest

	// Take until '@' as Data(string)
	if len(rest1) < 1 {
		p.Auth.Valid = false
		goto urlAuthLabel
	}
	pos = bytes.IndexByte(rest1[1:], '@')
	if pos >= 0 {
		p.Auth.Data = rest1[:pos+1]
		rest1 = rest1[pos+1+1:]
	} else {
		p.Auth.Valid = false
		goto urlAuthLabel
	}
	p.Auth.Valid = true
	p.rest = rest1
urlAuthLabel:

	// Take until ':' as Host(string)
	if len(p.rest) < 1 {
		return false, fmt.Errorf("Cannot slice from %d as only %d characters left in the rest (`\033[1m%s\033[0m`)", 1, len(p.rest), string(p.rest))
	}
	pos = bytes.IndexByte(p.rest[1:], ':')
	if pos >= 0 {
		p.Host = p.rest[:pos+1]
		p.rest = p.rest[pos+1+1:]
	} else {
		return false, fmt.Errorf("Cannot find `\033[1m%c\033[0m` in `\033[1m%s\033[0m` to bound data for field Host", ':', string(p.rest[1:]))
	}

	// Take until '/' (or all the rest if not found) as Port(uint16)
	pos = bytes.IndexByte(p.rest, '/')
	if pos >= 0 {
		tmp = p.rest[:pos]
		p.rest = p.rest[pos+1:]
	} else {
		tmp = p.rest
		p.rest = p.rest[len(p.rest):]
	}
	if tmpUint, err = strconv.ParseUint(*(*string)(unsafe.Pointer(&tmp)), 10, 16); err != nil {
		return false, fmt.Errorf("Cannot parse `%s`: %s", string(tmp), err)
	}
	p.Port = uint16(tmpUint)

	// Take the rest as DBName(string)
	p.DBName = p.rest
	p.rest = p.rest[len(p.rest):]
	return true, nil
}

// GetAuthData ...
func (p *URL) GetAuthData() (res []byte) {
	if p.Auth.Valid {
		res = p.Auth.Data
	}
	return
}

// Auth ...
type Auth struct {
	rest     []byte
	User     []byte
	Password []byte
}

// Extract ...
func (p *Auth) Extract(line []byte) (bool, error) {
	p.rest = line
	var pos int

	// Take until ':' (or all the rest if not found) as User(string)
	pos = bytes.IndexByte(p.rest, ':')
	if pos >= 0 {
		p.User = p.rest[:pos]
		p.rest = p.rest[pos+1:]
	} else {
		p.User = p.rest
		p.rest = p.rest[len(p.rest):]
	}

	// Take the rest as Password(string)
	p.Password = p.rest
	p.rest = p.rest[len(p.rest):]
	return true, nil
}
