box.cfg{
    listen = 3301
}

box.schema.user.grant('guest', 'read,write,execute', 'universe', nil, {if_not_exists = true})

local space = box.schema.space.create('votes', { if_not_exists = true })
space:format({
    {name = 'id', type = 'unsigned'},
    {name = 'user', type = 'string'},
    {name = 'option', type = 'string'}
})
space:create_index('primary', {
    parts = {'id'},
    if_not_exists = true
})
