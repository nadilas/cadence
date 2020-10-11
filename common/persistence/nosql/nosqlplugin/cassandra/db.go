// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cassandra

import (
	"errors"

	"github.com/gocql/gocql"

	"github.com/uber/cadence/common/service/config"

	"github.com/uber/cadence/common/persistence/nosql/nosqlplugin"
)

var (
	errConditionFailed = errors.New("internal condition fail error")
)

// cdb represents a logical connection to Cassandra database
type cdb struct {
	session *gocql.Session
}

var _ nosqlplugin.DB = (*cdb)(nil)

// NewCassandraDBFromSession returns a DB from a session
func NewCassandraDBFromSession(session *gocql.Session) nosqlplugin.DB {
	return &cdb{
		session: session,
	}
}

// NewCassandraDB return a new DB
func NewCassandraDB(cfg config.Cassandra) (nosqlplugin.DB, error) {
	session, err := CreateSession(cfg)
	if err != nil {
		return nil, err
	}
	return &cdb{
		session: session,
	}, nil
}

func (db *cdb) Close() {
	if db.session != nil {
		db.session.Close()
	}
}

func (db *cdb) PluginName() string {
	return PluginName
}

func (db *cdb) IsNotFoundError(err error) bool {
	return err == gocql.ErrNotFound
}

func (db *cdb) IsTimeoutError(err error) bool {
	if err == gocql.ErrTimeoutNoResponse {
		return true
	}
	if err == gocql.ErrConnectionClosed {
		return true
	}
	_, ok := err.(*gocql.RequestErrWriteTimeout)
	return ok
}

func (db *cdb) IsThrottlingError(err error) bool {
	if req, ok := err.(gocql.RequestError); ok {
		// gocql does not expose the constant errOverloaded = 0x1001
		return req.Code() == 0x1001
	}
	return false
}

func (db *cdb) IsConditionFailedError(err error) bool {
	if err == errConditionFailed {
		return true
	}
	return false
}
