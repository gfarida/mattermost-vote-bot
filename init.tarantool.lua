box.cfg{
    listen = 3301,
    memtx_memory = 1024 * 1024 * 128,
}

box.schema.user.grant('guest', 'read,write,execute', 'universe')

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