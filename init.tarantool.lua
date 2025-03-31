box.cfg{
    listen = 3301,
    memtx_memory = 1024 * 1024 * 128,
}

box.schema.user.grant('guest', 'read,write,execute', 'universe')

local user = box.space._user.index.name:get('guest')
if user ~= nil then
    local privileges = box.schema.user.info('guest').universe
    if not privileges or #privileges == 0 then
        box.schema.user.grant('guest', 'read,write,execute', 'universe')
    end
end