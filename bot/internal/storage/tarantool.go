package storage

import (
	"github.com/tarantool/go-tarantool"
)

func InitTarantool(uri string) (*tarantool.Connection, error) {
	conn, err := tarantool.Connect(uri, tarantool.Opts{
		User: "guest",
		Pass: "guest",
	})
	if err != nil {
		return nil, err
	}

	_, err = conn.Eval(`
        if not box.space.polls then
            box.schema.create_space('polls')
            box.space.polls:format({
                {name = 'id', type = 'string'},
                {name = 'question', type = 'string'},
                {name = 'options', type = 'map'},
                {name = 'creator', type = 'string'},
                {name = 'active', type = 'boolean'}
            })
            box.space.polls:create_index('primary', {parts = {'id'}})
        end
    `, []interface{}{})

	return conn, err
}
